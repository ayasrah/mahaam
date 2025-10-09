# Mahaam System Design

### Overview

This section discusses some of Mahaam system design aspects.

### Purpose

Helps understanding the big picture of the app.

### Topics

This section has following topics:

1. **Mahaam functions**: Usecases, business domain, functionalities.
2. **Mahaam data model**: Database ERD.
3. **Mahaam API Design**: Functions exposed by APIs.
4. **Mahaam Structure**: Folding options: by features or by layers.
5. **Maintainability**: Code maintainability.

### Building App Steps

Mahaam believes that these steps can boast app creating productivity:

1. Understanding user needs and business needs.
2. Design the workflows and UX.
3. Design the data models and database schema.
4. Maintainability: structure, refactor,...
5. Start implementing features, start by simplest
6. Refine, repeat and validate from 1-5 (may needs to change the use case, data model, and database design accordingly).

##### 1. Understanding business domain

- Mahaam domain is TODO app,which has mainly `Plans` (Projects, Lists), each plan has `Tasks` (todos)
- Eg: Groceries is a plan, bread, yogurt is tasks
- Plan owner can share plan with other users

##### 2. Design the data model

Mahaam has these entities: `Users, Plans, Tasks, PlanMembers`

##### 3. Design the app structure

Mahaam has 2 main folders: **Feat**, and **Infra** (more about in next section). Importance is that every code block is known where should be placed, so once a developer starts a feature implementation he just fills in blanks as the blueprint already established.

##### 4. Feature Implementation

Start implementing features, beginning with the simplest ones. Mahaam follows this implementation order:

Each feature follows the three-layer pattern: Controller → Service → Repository

##### 5. Refinement and Validation

Continuously refine and validate the implementation. This phase may require changes to use cases, data models, and database design.
