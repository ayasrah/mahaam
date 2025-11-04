# Monitoring

### Overview

Monitoring tracks app health and performance in real-time.

### Components

**Health Check**

- Endpoint: `/health`
- Returns: app name, version, environment
- Used by load balancers and health checks

**Traffic Logging**

- Records all API requests automatically
- Excludes: `/swagger`, `/health`, `/audit` paths
- Captures: method, path, status code, response time, user/device IDs
- Privacy: omits response data for user endpoints

**Audit Endpoints**

- `/audit/error`: Log client-side errors
- `/audit/info`: Log client-side info messages

**Server Health**

- Records server start/stop events
- Sends periodic pulses to track uptime
- Stores node IP, name, app details

### Data Models

**Traffic**

```sql
traffic (id, user_id, device_id, code, elapsed, method, path,
          payload, response, app_version, app_store, ip, created_at)
```

**Health**

```sql
health (id, node_ip, node_name, app_name, app_version, created_at,
          updated_at, stopped)
```

**Logs**

```sql
logs (traffic_id, type, message, ip, created_at)
```

### Implementation

- **Async logging**: Traffic and logs saved asynchronously
- **Error handling**: Failed monitoring doesn't break main flow
- **Performance**: Non-blocking operations
- **Privacy**: Sensitive data filtered out

### Benefits

- Observe all events done
- Debug issues quickly
- Track app performance
- Ensure app availability

### See

- Monitor module in: [C#](https://github.com/ayasrah/mahaam/tree/main/mahaam-api-cs/Src/Infra/Monitor), [Java](https://github.com/ayasrah/mahaam/tree/main/mahaam-api-java/src/main/java/mahaam/infra/monitor), [TypeScript](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-ts/src/infra/monitor), [Python](https://github.com/ayasrah/mahaam/tree/main/mahaam-api-py/infra/monitor)
