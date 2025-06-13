---
applyTo: '**'
---
# GitHub Copilot Prompt Guide

## Core Design Principles

**YAGNI (You Aren't Gonna Need It):** Implement only what's needed right now. Avoid building for hypothetical future requirements.

**Simplicity First:** Choose the simplest solution that solves the current problem. Add complexity only when proven necessary by real requirements.

**Deep Modules, Simple Interfaces:** Create modules that hide complexity behind clean, minimal interfaces that reduce cognitive load.

## Approach Framework

When responding to my requests, follow this lean approach:

1. **Understand the Immediate Need**
   - Focus on the specific problem at hand
   - Ask clarifying questions only if the requirement is truly unclear
   - Resist the urge to solve adjacent or future problems

2. **Start Simple**
   - Choose the most straightforward solution that works
   - Avoid premature abstractions and design patterns
   - Prefer direct, obvious implementations over clever ones

3. **Implement Incrementally**
   - Write the minimal code that solves the problem
   - Keep interfaces simple and focused
   - Add features only when explicitly requested

4. **Explain Concisely**
   - Focus on what the code does, not potential enhancements
   - Highlight only the essential design decisions
   - Avoid discussing patterns or optimizations unless directly relevant

## Code Style Preferences

- Write straightforward code that solves the current problem
- Prefer obvious solutions over clever abstractions
- Use simple, descriptive names that reduce mental mapping
- Add comments only when the code itself cannot be made clearer
- Build upon existing patterns rather than introducing new complexity
- Avoid premature optimization and over-engineering

## Context Understanding

- Pay attention to the folder structure and file relationships
- Consider existing patterns in the codebase
- Reference relevant files or code sections in your responses
- Build upon existing functionality rather than reinventing

## Communication Style

- Be direct and to the point
- Focus on solving the immediate problem
- Avoid suggesting enhancements unless explicitly asked
- Use simple language that reduces cognitive load
- Show, don't just tell - provide working examples

## Decision Making

- Choose the solution that requires the least cognitive effort to understand
- Avoid creating abstractions until you have multiple concrete examples
- When unsure, pick the simpler option and iterate
- Resist the temptation to solve problems that don't exist yet

## Learning Approach

- Focus on understanding the current codebase and problem
- Explain decisions only when they're non-obvious
- Share knowledge that helps solve the immediate task
- Avoid theoretical discussions unless they directly apply

---

When responding to my coding tasks, prioritize simplicity and solving the actual problem at hand. I value working solutions that are easy to understand and maintain over elaborate architectures that anticipate future needs.