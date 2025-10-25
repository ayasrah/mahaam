# Maintainability

### Overview

This page discusses code maintainability.

**Code maintainability** is how easily code can be understood, updated, fixed and extended over time.

**Code readability** is part of maintainability and its how the code is clear and self-explanatory.

**Naming** is part of maintainability, and it refers to choosing clear, consistent, breaf and descriptive names for the **`app, modules, classes, functions, variables and db tables`**.

### Importance

Maintainability is very important to cosider as it:

- **Reduces Costs**: Lots of **time and efforts are wasted** because of unmaintainable code, which means wasted money.
- **Increase Developer happiness**: Working on a clean and readable code is comfortable. From other side, when anyone's time and efforts are wasted, for sure will not be happy.
- **Healthy Environment**: Clean, maintainable code reflects a good engineering culture and vise versa.
- **Long-term Investment**: Code that's easy to maintain adapts better for people or technologies changes.
- **Reduced technical debt**: Maintainability prevents accumilation of quick hacks.
- **Faster feature delivery**: A well-structured code makes adding/editing features easier.
- **Faster onboarding**: New team members can understand the codebase quickly.

### Steps to maintainability

#### 1. Well Defined user stories

Understanding **what end users needs** is the first and last part of the system and is part of the big picture of the app. It may not be clear at starting new project, but its very important to focus on and align with during the journey.

Mahaam highly recommend to interact with users, product and domain experts during building the app in order to get solid business requirements.

#### 2. Well Designed Data model

If database design is mess, its very difficult to fix anything else. This is part of the big picture of the app.

Mahaam highly recommend to consult DB team when creating and reviewing the database schema, and reviewing the sqls as well.

#### 3. Simple Structure

Unfortunately, some online content about **clean architecture** just maps to over-engineered structure. Keep it simple and practical.

Break your application into logical, independent modules that can be developed, tested, and maintained separately.

I saw apps that did part 1 and 2 very well (have well-defined user requirments and robust database design), but messed here. Unfortunately, some smart developers don't pay attention to code maintainability. Which leads sometimes to overengineered, complex and unmaintanable code.

- Avoid multi-project structure in favor to modular monolith.
- Avoid deep nested folders, 2 to 3 levels should be good.
- Avoid too much folders which mostly contains one or two files, try grouping them.

#### 4. Naming

This is a main door that messes the readability, like long, vague, or unrelated names.

- Choose Clear, consistent, breaf and descriptive names.
- Choose simple, short, meaningful names.
- Start naming from domain, data models.
- One word is better than two, two is better than three. Eg: App name (Mahaam): one descriptive word that means jobs or tasks. Module names: User, Plan, Task.
- Context is considered: UpdateTitle(id, title) in TaskRepo means UpdateTaskTitle, id is TaskId, title is task title.
- Choose suffix and stick with it:
  - Controller or Handler or Router
  - Service
  - Repo, Repository, Provider or Dao
- Infrastructure abbreviations are acceptable (repo, infra, db) as they're repeated and directly translated. Domain abbreviations are not (pln for plan, tsk for task, u for user)

The app building journey starts from user needs, then **data model**, so once these 2 steps are in good state, you can start by building the database schema, choose good names for the db tables, then start to name classes and modules based on. Eg:

- Plans, Tasks tables
- Plans, Tasks modules
- PlanController, PlanService, PlanRepo, PlanModel classes

#### 5. Small files

- Each file should do one thing (**single responsibility**), eg: `TaskRepo` file only defines Task database operations.
- 50-300 lines of code is good per file.
- 300-500 lines of code is ok per file.
- Have you worked on a 3000+ line Java/C# file? Pain guaranteed.

#### 6. Small functions

- Each function should do one thing (**single responsibility**), eg: `TaskRepo.UpdateTitle` function only updates a task title.
- 5-15 lines of code is good.
- 15-30 lines of code is ok.
- 300+ line Java/C# function, if you’ve been there, you know the suffering.
- Break large functions into smaller, focused ones.
- Parameter count should be reasonable (up to 5 is acceptable).

#### 7. What goes where

- Place things in their proper place (e.g., no SQL in controllers!)
- DB interactions in repos, business logic in services, API definition in controllers
- Big door for spaghitti code

#### 8. Let reader easily know `what` before `how`

- In mahaam, once reader open `Feat` folder, he find 3 folders and knows **what are modules of the project**:

```bash
src/
└── feat/
	├── plans/
    ├── tasks/
    └── users/
```

- Once reader open any feat, eg , `Task`, he read the interface and know **what are the functionalities**

```C#
public interface ITaskController
{
	IActionResult Create(Guid planId, string title);
	IActionResult Delete(Guid planId, Guid id);
	IActionResult UpdateDone(Guid planId, Guid id, bool done);
	IActionResult UpdateTitle(Guid id, string title);
	IActionResult ReOrder(Guid planId, int oldIndex, int newIndex);
	IActionResult GetMany(Guid planId);
}
```

#### 8. Consistency

Be consistent in:

- Naming
- Code format
- Patterns used in the app

#### 9. Enhance/Refactor/Reengineer

Enhancements and refactoring are continous processes. In some cases, reengineering is needed, for the full app, a component, or even the DB. Pay this debt as soon as possible, to avoid building on broken foundations.

### Final word

Code maintainability needs to pay attention to, especially in the AI era, things will easily be unmanged without rules, template, standards and patterns to follow.
