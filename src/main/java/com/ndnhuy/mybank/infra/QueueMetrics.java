package com.ndnhuy.mybank.infra;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.Gauge;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import lombok.Builder;
import lombok.Getter;

import java.util.concurrent.BlockingQueue;
import java.util.function.Supplier;

@Getter
public class QueueMetrics {
  private final MeterRegistry registry;

  // Arrival metrics
  private Counter transfersSubmitted;
  private Timer submissionTime;

  // Service metrics
  private Timer serviceTime;
  private Counter transfersCompleted;

  // Queue metrics (optional, can be null)
  private Gauge queueLength;
  private Timer waitTime;

  // System metrics (optional, can be null)
  private Gauge systemUtilization;

  // Time tracking for rate calculations
  private long startTime;

  private Runnable resetMetrics;

  public QueueMetrics(MeterRegistry registry) {
    this.registry = registry;
    this.startTime = System.nanoTime();
    Runnable initFunc = () -> {
      this.transfersSubmitted = Counter.builder("transfers.submitted").register(registry);
      this.submissionTime = Timer.builder("transfers.submission.time").register(registry);
      this.serviceTime = Timer.builder("transfers.service.time").register(registry);
      this.transfersCompleted = Counter.builder("transfers.completed").register(registry);
      this.waitTime = Timer.builder("transfers.wait.time").register(registry);
    };

    this.resetMetrics = () -> {
      // Remove existing meters from registry
      registry.remove(transfersSubmitted);
      registry.remove(submissionTime);
      registry.remove(serviceTime);
      registry.remove(transfersCompleted);
      registry.remove(waitTime);

      if (queueLength != null) {
        registry.remove(queueLength);
        queueLength = null;
      }

      if (systemUtilization != null) {
        registry.remove(systemUtilization);
        systemUtilization = null;
      }

      // Reset start time for new observation period
      this.startTime = System.nanoTime();

      initFunc.run();
    };

    initFunc.run();
  }

  public void setQueueLengthGauge(BlockingQueue<?> queue, MeterRegistry registry) {
    this.queueLength = Gauge.builder("transfers.queue.length", queue, q -> q.size()).register(registry);
  }

  public void setSystemUtilizationGauge(Supplier<Double> utilizationSupplier, MeterRegistry registry) {
    // Gauge.builder expects Supplier<Number>
    this.systemUtilization = Gauge.builder("system.utilization", (Supplier<Number>) utilizationSupplier::get)
        .register(registry);
  }

  public void reset() {
    if (resetMetrics != null) {
      resetMetrics.run();
    } else {
      throw new IllegalStateException("Reset function is not initialized");
    }
  }

  public Report getReport() {
    return Report.builder()
        .transfersSubmitted(transfersSubmitted.count())
        .submissionTimeInMilliseconds(submissionTime.totalTime(java.util.concurrent.TimeUnit.MILLISECONDS))
        .serviceTimeInMilliseconds(serviceTime.totalTime(java.util.concurrent.TimeUnit.MILLISECONDS))
        .transfersCompleted(transfersCompleted.count())
        .queueLength(queueLength != null ? queueLength.value() : 0)
        .waitTimeInMilliseconds(waitTime.totalTime(java.util.concurrent.TimeUnit.MILLISECONDS))
        .systemUtilization(systemUtilization != null ? systemUtilization.value() : 0)
        .queueMetrics(this)  // Pass reference to parent QueueMetrics
        .build();
  }

  // ============================================
  // Queuing Theory Statistical Methods
  // ============================================

  /**
   * Get the observation duration in seconds since metrics collection started.
   * @return observation duration in seconds
   */
  public double getObservationDurationSeconds() {
    return (System.nanoTime() - startTime) / 1_000_000_000.0;
  }

  /**
   * Get the mean response time (total time from arrival to completion).
   * Response Time = Wait Time + Service Time
   * @return mean response time in milliseconds, or 0 if no requests completed
   */
  public double getMeanResponseTime() {
    double completed = transfersCompleted.count();
    if (completed == 0) return 0;
    
    double totalWaitTime = waitTime.totalTime(java.util.concurrent.TimeUnit.MILLISECONDS);
    double totalServiceTime = serviceTime.totalTime(java.util.concurrent.TimeUnit.MILLISECONDS);
    return (totalWaitTime + totalServiceTime) / completed;
  }

  /**
   * Get the mean wait time (average time spent waiting in queue).
   * @return mean wait time in milliseconds, or 0 if no requests completed
   */
  public double getMeanWaitTime() {
    double completed = transfersCompleted.count();
    if (completed == 0) return 0;
    
    return waitTime.totalTime(java.util.concurrent.TimeUnit.MILLISECONDS) / completed;
  }

  /**
   * Get the mean service time (average time spent being processed).
   * @return mean service time in milliseconds, or 0 if no requests completed
   */
  public double getMeanServiceTime() {
    double completed = transfersCompleted.count();
    if (completed == 0) return 0;
    
    return serviceTime.totalTime(java.util.concurrent.TimeUnit.MILLISECONDS) / completed;
  }

  /**
   * Get the arrival rate (λ) - requests submitted per second.
   * @return arrival rate in requests/second
   */
  public double getArrivalRate() {
    double duration = getObservationDurationSeconds();
    if (duration <= 0) return 0;
    
    return transfersSubmitted.count() / duration;
  }

  /**
   * Get the service rate (μ) - requests completed per second.
   * Also known as throughput.
   * @return service rate in requests/second
   */
  public double getServiceRate() {
    double duration = getObservationDurationSeconds();
    if (duration <= 0) return 0;
    
    return transfersCompleted.count() / duration;
  }

  /**
   * Get the throughput (same as service rate).
   * @return throughput in requests/second
   */
  public double getThroughput() {
    return getServiceRate();
  }

  /**
   * Get the traffic intensity (ρ = λ/μ).
   * Represents the system load factor.
   * - ρ < 1: System is stable
   * - ρ = 1: System at capacity
   * - ρ > 1: System overloaded
   * @return traffic intensity (dimensionless)
   */
  public double getTrafficIntensity() {
    double serviceRate = getServiceRate();
    if (serviceRate <= 0) return Double.POSITIVE_INFINITY;
    
    return getArrivalRate() / serviceRate;
  }

  /**
   * Get the average queue length using Little's Law (L = λW).
   * L = average number of customers in the queue
   * λ = arrival rate
   * W = average wait time
   * @return average queue length (number of requests)
   */
  public double getAverageQueueLength() {
    double arrivalRate = getArrivalRate();
    double meanWaitTimeSeconds = getMeanWaitTime() / 1000.0; // Convert ms to seconds
    
    return arrivalRate * meanWaitTimeSeconds;
  }

  /**
   * Get the current instantaneous queue length.
   * @return current queue length (snapshot)
   */
  public double getCurrentQueueLength() {
    return queueLength != null ? queueLength.value() : 0;
  }

  /**
   * Get the current system utilization percentage.
   * @return system utilization as percentage (0-100)
   */
  public double getSystemUtilizationPercentage() {
    if (systemUtilization == null) return 0;
    return systemUtilization.value() * 100;
  }

  @Builder
  public static class Report {
    private final double transfersSubmitted;
    private final double submissionTimeInMilliseconds;
    private final double serviceTimeInMilliseconds;
    private final double transfersCompleted;
    private final double queueLength;
    private final double waitTimeInMilliseconds;
    private final double systemUtilization;
    private final QueueMetrics queueMetrics;

    public void print() {
      if (queueMetrics == null) {
        printLegacyReport();
        return;
      }

      System.out.println("╔══════════════════════════════════════════════════════════════╗");
      System.out.println("║                    QUEUING SYSTEM ANALYSIS                  ║");
      System.out.println("╚══════════════════════════════════════════════════════════════╝");
      
      // System Overview
      System.out.println("\n📊 SYSTEM OVERVIEW:");
      System.out.printf("   Observation Duration: %.2f seconds%n", queueMetrics.getObservationDurationSeconds());
      System.out.printf("   Requests Submitted:   %.0f%n", transfersSubmitted);
      System.out.printf("   Requests Completed:   %.0f%n", transfersCompleted);
      
      double trafficIntensity = queueMetrics.getTrafficIntensity();
      String systemStatus = getSystemStatus(trafficIntensity);
      System.out.printf("   System Status:        %s%n", systemStatus);

      // Performance Metrics
      System.out.println("\n⚡ PERFORMANCE METRICS:");
      System.out.printf("   Mean Response Time:   %.2f ms%n", queueMetrics.getMeanResponseTime());
      System.out.printf("   Mean Wait Time:       %.2f ms%n", queueMetrics.getMeanWaitTime());
      System.out.printf("   Mean Service Time:    %.2f ms%n", queueMetrics.getMeanServiceTime());
      System.out.printf("   Throughput:           %.2f requests/sec%n", queueMetrics.getThroughput());

      // Queuing Theory Analysis
      System.out.println("\n🔬 QUEUING THEORY ANALYSIS:");
      System.out.printf("   Arrival Rate (λ):      %.2f requests/sec%n", queueMetrics.getArrivalRate());
      System.out.printf("   Service Rate (μ):      %.2f requests/sec%n", queueMetrics.getServiceRate());
      System.out.printf("   Traffic Intensity (ρ): %.3f %s%n", trafficIntensity, getTrafficIntensityWarning(trafficIntensity));
      System.out.printf("   Avg Queue Length (L):  %.2f requests%n", queueMetrics.getAverageQueueLength());
      System.out.printf("   Current Queue Length:  %.0f requests%n", queueMetrics.getCurrentQueueLength());
      System.out.printf("   System Utilization:    %.1f%%%n", queueMetrics.getSystemUtilizationPercentage());

      // Performance Assessment
      System.out.println("\n🎯 PERFORMANCE ASSESSMENT:");
      printPerformanceAssessment(queueMetrics);
      
      System.out.println("\n" + "═".repeat(66));
    }

    private void printLegacyReport() {
      System.out.println("Queue Metrics Report (Legacy):");
      System.out.printf("Transfers Submitted: %.0f%n", transfersSubmitted);
      System.out.printf("Submission Time (ms): %.2f%n", submissionTimeInMilliseconds);
      System.out.printf("Service Time (ms): %.2f%n", serviceTimeInMilliseconds);
      System.out.printf("Transfers Completed: %.0f%n", transfersCompleted);
      System.out.printf("Queue Length: %.0f%n", queueLength);
      System.out.printf("Wait Time (ms): %.2f%n", waitTimeInMilliseconds);
      System.out.printf("System Utilization: %.2f%%%n", systemUtilization * 100);
    }

    private String getSystemStatus(double trafficIntensity) {
      if (trafficIntensity >= 1.0) return "🔴 OVERLOADED";
      if (trafficIntensity >= 0.8) return "🟡 HIGH LOAD";
      if (trafficIntensity >= 0.5) return "🟢 MODERATE LOAD";
      return "🟢 LOW LOAD";
    }

    private String getTrafficIntensityWarning(double trafficIntensity) {
      if (trafficIntensity >= 1.0) return "⚠️  SYSTEM OVERLOADED!";
      if (trafficIntensity >= 0.8) return "⚠️  Approaching capacity";
      return "✅ Stable";
    }

    private void printPerformanceAssessment(QueueMetrics metrics) {
      double responseTime = metrics.getMeanResponseTime();
      double waitTime = metrics.getMeanWaitTime();
      double trafficIntensity = metrics.getTrafficIntensity();

      System.out.printf("   Response Time:  %s%n", assessResponseTime(responseTime));
      System.out.printf("   Queue Buildup:  %s%n", assessQueueBuildup(waitTime));
      System.out.printf("   System Health:  %s%n", assessSystemHealth(trafficIntensity));
    }

    private String assessResponseTime(double responseTime) {
      if (responseTime < 100) return "🟢 Excellent (< 100ms)";
      if (responseTime < 500) return "🟡 Good (< 500ms)";
      if (responseTime < 1000) return "🟠 Fair (< 1s)";
      return "🔴 Poor (> 1s)";
    }

    private String assessQueueBuildup(double waitTime) {
      if (waitTime < 10) return "🟢 Minimal queue buildup";
      if (waitTime < 100) return "🟡 Moderate queue buildup";
      if (waitTime < 500) return "🟠 Significant queue buildup";
      return "🔴 Severe queue buildup";
    }

    private String assessSystemHealth(double trafficIntensity) {
      if (trafficIntensity >= 1.0) return "🔴 System cannot keep up with demand";
      if (trafficIntensity >= 0.8) return "🟡 System near capacity, consider scaling";
      if (trafficIntensity >= 0.5) return "🟢 System handling load well";
      return "🟢 System has excess capacity";
    }
  }
}
