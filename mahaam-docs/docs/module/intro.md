# Microservices Modules

### Overview

A **module** is a set of related classes, functions and data models. For example, in Mahaam plans is a module, in ecommerce, orders, products, catalogs, and reviews are all modules.

This section discusses module internal parts.

### Mahaam Modules

- Defining app modules, and each module boundaries are important step.
- Folding by modules highly enhances readability and maintainability.
- Mahaam has 3 business modules: Plan, Task and User, and one infra module which is monitoring.

### Module Layers

these are the layers (vertical slice) that a module cosists of:

- Controller Layer: Module APIs.
- Service Layer: Module business logic.
- Repo Layer: Database access layer.
- Models: Data models and DTOs.
