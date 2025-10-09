# Mahaam API Design

### Overview

This page explains Mahaam's architecture and API design.

### Mahaam Components

The Mahaam system follows a simple three-tier architecture:

<img src="/mahaam_arch.svg" alt="Mahaam System Architecture" width="600"  />

- Client: The Mobile App.
- Server: The Backend API Service.
- Data: The Database.

**Data Flow:**

1. Users interact with Mobile App.
2. The app sends HTTP requests to the backend service.
3. The service processes requests and interacts with the database.
4. Data flows back through the same path.

### API Design

**Base URL**: `/mahaam-api`

##### User APIs

```bash
POST 	/users/send-me-otp		# sends OTP to user email
POST 	/users/verify-otp		# verifies OTP
POST 	/users					# creates a new user
POST 	/users/refresh-token	# refreshes token
POST 	/users/logout			# logs out from device
PATCH 	/users/name				# updates user full name
DELETE 	/users					# deletes user account
GET 	/users/devices			# gets user devices
GET 	/users/suggested-emails	# gets suggested emails
DELETE 	/users/suggested-emails	# deletes a suggested email

# note: userId is available in JWT header
```

##### Plan APIs

```bash
POST 	/plans					# creates a plan
PUT 	/plans					# updates a plan
GET 	/plans					# gets user plans
GET 	/plans/{id}				# gets a plan details
DELETE 	/plans/{id}				# deletes a plan
PATCH 	/plans/{id}/share		# shares a plan with a user
PATCH 	/plans/{id}/unshare		# unshares a plan with a user
PATCH 	/plans/{id}/leave		# a user leaves a plan
PATCH 	/plans/{id}/type		# updates plan type
PATCH 	/plans/reorder			# reorders plans
```

##### Task APIs

```bash
POST 	/plans/{planId}/tasks				# creates a task in plan
GET 	/plans/{planId}/tasks				# gets all tasks in plan
DELETE 	/plans/{planId}/tasks/{id}			# deletes a task
PATCH 	/plans/{planId}/tasks/{id}/done		# marks task as done/undone
PATCH 	/plans/{planId}/tasks/{id}/title	# updates task title
PATCH 	/plans/{planId}/tasks/reorder		# reorders tasks in plan
```

##### Audit APIs

```bash
POST 	/audit/error		# creates an error message
POST 	/audit/info			# creates an info message
```

##### Health APIs

```bash
GET 	/health				# checks system health
```
