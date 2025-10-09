# Monitoring

#### Overview

Monitoring tracks app health and performance in real-time.

#### Components

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

#### Data Models

**Traffic**

```sql
x_traffic (id, user_id, device_id, code, elapsed, method, path,
          payload, response, app_version, app_store, ip, created_at)
```

**Health**

```sql
x_health (id, node_ip, node_name, app_name, app_version, created_at,
          updated_at, stopped)
```

**Logs**

```sql
x_log (traffic_id, type, message, ip, created_at)
```

#### Implementation

- **Async logging**: Traffic and logs saved asynchronously
- **Error handling**: Failed monitoring doesn't break main flow
- **Performance**: Non-blocking operations
- **Privacy**: Sensitive data filtered out

#### Benefits

- Track app performance
- Monitor user activity
- Debug issues quickly
- Ensure app availability
