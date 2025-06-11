---
applyTo: '**'
---
# Project Memory System

## Memory Management Instructions

You are equipped with a persistent memory system to maintain and evolve your understanding of this project. Follow these guidelines:

### 1. Memory File Location
- **File**: `.memory.md` in the project root
- **Purpose**: Store cumulative knowledge about the project
- **Format**: Structured Markdown with clear sections

### 2. Memory Reading Protocol
At the start of each conversation:
- Read the current `.memory.md` file if it exists
- Use this information to inform your responses
- Reference stored knowledge when relevant

### 3. Memory Update Process
After each significant interaction:
- Update `.memory.md` with new insights about:
  - Project architecture and structure
  - Key patterns and conventions used
  - Important business logic or domain knowledge
  - Common issues and their solutions
  - User preferences and coding style
  - Dependencies and external integrations
  - Testing strategies and patterns

### 4. Memory Structure

Maintain the following sections in `.memory.md`:

```markdown
# Project Memory

## Last Updated
[Current date and interaction summary]

## Project Overview
- **Name**: [Project name]
- **Type**: [Application type]
- **Primary Language**: [Main language]
- **Framework**: [Key frameworks]

## Architecture & Structure
- [Key architectural patterns]
- [Important modules/components]
- [Data flow and relationships]

## Key Patterns & Conventions
- [Coding patterns observed]
- [Naming conventions]
- [File organization principles]

## Business Domain Knowledge
- [Domain-specific concepts]
- [Business rules and logic]
- [Key entities and relationships]

## Technical Insights
- [Important technical decisions]
- [Performance considerations]
- [Security patterns]
- [Error handling approaches]

## Dependencies & Integrations
- [External libraries and versions]
- [API integrations]
- [Database connections]
- [Third-party services]

## User Preferences
- [Coding style preferences]
- [Preferred solutions and approaches]
- [Communication style notes]

## Common Issues & Solutions
- [Recurring problems and fixes]
- [Debugging approaches]
- [Deployment considerations]

## Recent Learnings
- [New insights from recent conversations]
- [Evolving understanding]
- [Questions for future exploration]
```

### 5. Memory Refinement Guidelines

- **Accuracy**: Only store verified information and solution. If in doubt, please verify before updating.
- **Relevance**: Focus on information that helps future interactions
- **Evolution**: Update existing knowledge when new information contradicts or refines it
- **Conciseness**: Keep entries clear and actionable
- **Dating**: Track when information was learned or updated
- Do not store too much detail; focus on high-level insights, store Reference links to detailed documentation or code when necessary.
- Do not store changes from refactoring or minor code changes unless they significantly alter understanding of the project.

### 6. Memory Usage in Responses

When answering questions:
- Reference relevant stored knowledge
- Build upon previous learnings
- Identify gaps in understanding
- Propose memory updates for significant new insights

### 7. Memory Initialization

If `.memory.md` doesn't exist, create it with:
- Initial project analysis based on current workspace
- Observed patterns from available code
- Preliminary understanding of the project structure

## Example Memory Update Process

After helping with a significant task:

1. **Analyze** what was learned about the project
2. **Identify** which memory sections need updates  
3. **Update** `.memory.md` with new or refined information
4. **Confirm** the update was successful
5. **Summarize** what was added to memory

---

**Note**: This memory system helps me provide increasingly relevant and context-aware assistance as our collaboration evolves. The memory file serves as a knowledge base that grows with each interaction.