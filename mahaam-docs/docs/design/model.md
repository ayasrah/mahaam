# Mahaam Data Model

### Overview

This page explains Mahaam data models and the Entity-Relationship Diagram (ERD).

### Creation Steps

At data model step, Mahaam suggests to follow a UX-Driven approach by:

1. Understanding user and business needs.
2. Designing workflows.
3. Sketching high-level pages with functions.
4. Designing the data models accordingly.
5. Consult DB team to review the design.

This ensures the database structure directly supports the user needs rather than being designed in isolation.

### Data Models (Entities)

#### Business Models

- **Users**: App users.
- **Plans**: Projects or lists.
- **Tasks**: Plan todo items.
- **PlanMembers**: Junction table enabling plan sharing between users.
- **Devices**: Devices that user logged in using.
- **SuggestedEmails**: Suggested emails for plan sharing functionality.

#### Observability Models

- **Health**: Service instances info, like started, last pulsed, node,...etc.
- **Traffic**: API request/response info, like path, status code, elapsed time,...etc.
- **Logs**: App logs, like: events, errors, and audits.

<img src="/mahaam_erd.png" alt="ERD" width="430" style="border-radius:5px;"/>
