# Testing

### Overview

Testing is the process of ensuring the app works as expected.

### Types

- **Manual Testing**: Group of tests run by a human.
- **Automated Testing**: Group of written tests run using a tool or a script.

### Levels

- **Unit Testing**: A Test written for a small piece of code in isolation like a function.
- **Integration Testing**: A Test written for a functionality (group of units).
- **API Testing**: A subset of integration test, and its a test written for an API endpoint.
- **End-to-end Testing**: A Test written for a functionality from end-user perspective.

### Importance

Writing automated tests is not optional for production apps. The more you invest in, the more confident your deployments will be.

### Mahaam Testing

Mahaam has the minimal accepted testing suite, it has automated API integration tests using postman, and only the happy path scenarios.

```bash
mahaam-api-test/
├── mahaam_local.postman_environment.json   # Environment variables
└── mahaam.postman_collection.json  		# Tests collection
```

### Run

#### 1. Run with Postman

Open both `mahaam.postman_collection` and `mahaam_local.postman_environment` in postman then run the collection against the environment.

<video controls>
  <source src="../public/tests_postman.mp4" type="video/mp4">
  Your browser does not support the video tag.
</video>

#### 2. Run with Newman

Install newman cli

```bash
npm install -g newman
```

then run the tests

```bash
newman run ./mahaam.postman_collection.json -e ./mahaam_local.postman_environment.json
```

<video controls>
  <source src="../public/tests_newman.mp4" type="video/mp4">
  Your browser does not support the video tag.
</video>

### Sample Test Script

```JavaScript
pm.test(pm.info.requestName, function () {
  pm.expect(pm.response.code).to.eq(200);
});
```
