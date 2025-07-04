#!/usr/bin/env python3
"""
Load Test Visualization Script
Generates time-series graphs for load test metrics showing convexity patterns.
"""

import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np
import os
import argparse
from datetime import datetime
import warnings

# Suppress warnings for cleaner output
warnings.filterwarnings('ignore')

# Set style
plt.style.use('seaborn-v0_8')
sns.set_palette("husl")

def load_timeseries_data(csv_file):
    """Load time-series data from CSV file."""
    try:
        df = pd.read_csv(csv_file)
        df['Timestamp'] = pd.to_datetime(df['Timestamp'])
        return df
    except FileNotFoundError:
        print(f"Warning: {csv_file} not found")
        return None
    except Exception as e:
        print(f"Error loading {csv_file}: {e}")
        return None

def create_single_scenario_plot(df, scenario_name, output_dir):
    """Create plots for a single scenario."""
    fig, axes = plt.subplots(2, 2, figsize=(15, 12))
    fig.suptitle(f'Load Test Metrics - {scenario_name.title()} Scenario', fontsize=16, fontweight='bold')
    
    # Plot 1: Mean and P95 Latency over time
    ax1 = axes[0, 0]
    ax1.plot(df['Elapsed_Seconds'], df['Mean_Latency_Ms'], 
             label='Mean Latency', linewidth=2, marker='o', markersize=3)
    ax1.plot(df['Elapsed_Seconds'], df['P95_Latency_Ms'], 
             label='P95 Latency', linewidth=2, marker='s', markersize=3)
    ax1.set_xlabel('Time (seconds)')
    ax1.set_ylabel('Latency (ms)')
    ax1.set_title('Response Time Evolution')
    ax1.legend()
    ax1.grid(True, alpha=0.3)
    
    # Plot 2: P99 Latency over time
    ax2 = axes[0, 1]
    ax2.plot(df['Elapsed_Seconds'], df['P99_Latency_Ms'], 
             label='P99 Latency', linewidth=2, marker='^', markersize=3, color='red')
    ax2.set_xlabel('Time (seconds)')
    ax2.set_ylabel('Latency (ms)')
    ax2.set_title('P99 Latency Evolution')
    ax2.legend()
    ax2.grid(True, alpha=0.3)
    
    # Plot 3: Throughput over time
    ax3 = axes[1, 0]
    ax3.plot(df['Elapsed_Seconds'], df['Throughput_RPS'], 
             label='Throughput', linewidth=2, marker='D', markersize=3, color='green')
    ax3.set_xlabel('Time (seconds)')
    ax3.set_ylabel('Requests/Second')
    ax3.set_title('Throughput Evolution')
    ax3.legend()
    ax3.grid(True, alpha=0.3)
    
    # Plot 4: Success Rate over time
    ax4 = axes[1, 1]
    ax4.plot(df['Elapsed_Seconds'], df['Success_Rate_Percent'], 
             label='Success Rate', linewidth=2, marker='*', markersize=4, color='orange')
    ax4.set_xlabel('Time (seconds)')
    ax4.set_ylabel('Success Rate (%)')
    ax4.set_title('Success Rate Evolution')
    ax4.set_ylim(0, 105)
    ax4.legend()
    ax4.grid(True, alpha=0.3)
    
    plt.tight_layout()
    
    # Save plot
    output_file = os.path.join(output_dir, f'{scenario_name}_metrics.png')
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"Created {output_file}")
    plt.close()

def create_comparison_plot(dfs, scenarios, output_dir):
    """Create comparison plots across all scenarios."""
    fig, axes = plt.subplots(2, 2, figsize=(16, 12))
    fig.suptitle('Load Test Comparison - All Scenarios', fontsize=16, fontweight='bold')
    
    colors = ['blue', 'orange', 'red']
    
    # Plot 1: Mean Latency comparison
    ax1 = axes[0, 0]
    for i, (scenario, df) in enumerate(zip(scenarios, dfs)):
        if df is not None:
            ax1.plot(df['Elapsed_Seconds'], df['Mean_Latency_Ms'], 
                    label=f'{scenario.title()} ({len(df)} points)', 
                    linewidth=2, color=colors[i % len(colors)])
    ax1.set_xlabel('Time (seconds)')
    ax1.set_ylabel('Mean Latency (ms)')
    ax1.set_title('Mean Response Time Comparison')
    ax1.legend()
    ax1.grid(True, alpha=0.3)
    
    # Plot 2: P95 Latency comparison
    ax2 = axes[0, 1]
    for i, (scenario, df) in enumerate(zip(scenarios, dfs)):
        if df is not None:
            ax2.plot(df['Elapsed_Seconds'], df['P95_Latency_Ms'], 
                    label=f'{scenario.title()} ({len(df)} points)', 
                    linewidth=2, color=colors[i % len(colors)])
    ax2.set_xlabel('Time (seconds)')
    ax2.set_ylabel('P95 Latency (ms)')
    ax2.set_title('P95 Response Time Comparison')
    ax2.legend()
    ax2.grid(True, alpha=0.3)
    
    # Plot 3: Throughput comparison
    ax3 = axes[1, 0]
    for i, (scenario, df) in enumerate(zip(scenarios, dfs)):
        if df is not None:
            ax3.plot(df['Elapsed_Seconds'], df['Throughput_RPS'], 
                    label=f'{scenario.title()} ({len(df)} points)', 
                    linewidth=2, color=colors[i % len(colors)])
    ax3.set_xlabel('Time (seconds)')
    ax3.set_ylabel('Throughput (RPS)')
    ax3.set_title('Throughput Comparison')
    ax3.legend()
    ax3.grid(True, alpha=0.3)
    
    # Plot 4: Convexity Analysis (Rate of Change)
    ax4 = axes[1, 1]
    for i, (scenario, df) in enumerate(zip(scenarios, dfs)):
        if df is not None and len(df) > 1:
            # Calculate rate of change (derivative) for mean latency
            time_diff = df['Elapsed_Seconds'].diff()
            latency_diff = df['Mean_Latency_Ms'].diff()
            rate_of_change = latency_diff / time_diff
            
            ax4.plot(df['Elapsed_Seconds'][1:], rate_of_change[1:], 
                    label=f'{scenario.title()} Rate of Change', 
                    linewidth=2, color=colors[i % len(colors)])
    ax4.set_xlabel('Time (seconds)')
    ax4.set_ylabel('Latency Rate of Change (ms/s)')
    ax4.set_title('Convexity Analysis - Rate of Latency Change')
    ax4.legend()
    ax4.grid(True, alpha=0.3)
    ax4.axhline(y=0, color='black', linestyle='--', alpha=0.5)
    
    plt.tight_layout()
    
    # Save plot
    output_file = os.path.join(output_dir, 'comparison_all_scenarios.png')
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"Created {output_file}")
    plt.close()

def analyze_convexity(df, scenario_name):
    """Analyze convexity patterns in the data."""
    if df is None or len(df) < 3:
        return
    
    print(f"\n=== Convexity Analysis for {scenario_name.title()} ===")
    
    # Calculate first and second derivatives
    time_diff = df['Elapsed_Seconds'].diff()
    latency_diff = df['Mean_Latency_Ms'].diff()
    first_derivative = latency_diff / time_diff
    second_derivative = first_derivative.diff() / time_diff[1:]
    
    # Statistics
    mean_latency_start = df['Mean_Latency_Ms'].iloc[0]
    mean_latency_end = df['Mean_Latency_Ms'].iloc[-1]
    p95_latency_start = df['P95_Latency_Ms'].iloc[0]
    p95_latency_end = df['P95_Latency_Ms'].iloc[-1]
    
    print(f"Mean Latency: {mean_latency_start:.2f} → {mean_latency_end:.2f} ms ({((mean_latency_end/mean_latency_start-1)*100):+.1f}%)")
    print(f"P95 Latency:  {p95_latency_start:.2f} → {p95_latency_end:.2f} ms ({((p95_latency_end/p95_latency_start-1)*100):+.1f}%)")
    
    # Convexity indicators
    positive_acceleration = (second_derivative > 0).sum()
    negative_acceleration = (second_derivative < 0).sum()
    
    print(f"Acceleration points: {positive_acceleration} positive, {negative_acceleration} negative")
    
    if positive_acceleration > negative_acceleration:
        print("Pattern: CONVEX (accelerating degradation)")
    elif negative_acceleration > positive_acceleration:
        print("Pattern: CONCAVE (decelerating degradation)")
    else:
        print("Pattern: MIXED (variable acceleration)")

def main():
    parser = argparse.ArgumentParser(description='Visualize load test metrics')
    parser.add_argument('--reports-dir', default='reports', help='Directory containing CSV files')
    parser.add_argument('--output-dir', default='reports', help='Directory to save plots')
    parser.add_argument('--scenarios', nargs='+', default=['baseline', 'stress', 'breaking'], 
                       help='Scenarios to analyze')
    
    args = parser.parse_args()
    
    # Create output directory if it doesn't exist
    os.makedirs(args.output_dir, exist_ok=True)
    
    # Load data for each scenario
    dfs = []
    scenarios = args.scenarios
    
    print("Loading time-series data...")
    for scenario in scenarios:
        csv_file = os.path.join(args.reports_dir, f'timeseries_{scenario}.csv')
        df = load_timeseries_data(csv_file)
        dfs.append(df)
        
        if df is not None:
            print(f"Loaded {len(df)} data points for {scenario}")
            # Create individual scenario plot
            create_single_scenario_plot(df, scenario, args.output_dir)
            # Analyze convexity
            analyze_convexity(df, scenario)
        else:
            print(f"No data available for {scenario}")
    
    # Create comparison plot if we have data
    valid_data = [(scenario, df) for scenario, df in zip(scenarios, dfs) if df is not None]
    if len(valid_data) > 1:
        valid_scenarios, valid_dfs = zip(*valid_data)
        create_comparison_plot(valid_dfs, valid_scenarios, args.output_dir)
    
    print(f"\nVisualization complete! Check {args.output_dir} for generated plots.")

if __name__ == '__main__':
    main()