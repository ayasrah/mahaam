# Logging

### Overview

Logging is saving **events** done in the app to storage.

### Example

```log
2025-08-01 03:00:10,187 INFO mahaam-api-v1.0 started on nodeIP=1.2.3.4 with healthID=8eeb192c-ffef-4830-8954-12c69d3d6b5d
2025-08-01 03:10:58,187 INFO TrafficId: f0aed2a0-9cc9-4a5e-85aa-f15ad9337035, user created with id=2f7d68cd-0216-4b80-bdc2-de12d2a37099
2025-08-01 03:20:58,187 ERROR TrafficId: 4e4d04d0-7bb3-47af-9e19-32aff5fda5b7, max allowed plans reached for userID=5d3c0e60-5030-4dad-973b-0ba9b97f7d5b
2025-08-01 16:37:53.762 INFO TrafficId: dc1f92ac-7c96-4c02-aa44-31fe5707a3e9, OTP sent to abc@example.com

```

### Purpose

- Monitor the app
- Debug and fix errors

### Saving Options

- Files
- External log service
- Database

### Log to file vs DB

- File:
  - **faster**: in Microseconds, no network, just I/O to local disk.
  - Distributed, so needs to be aggregated from multiple servers
  - Searching needs extra work and setup like ELK(Elasticsearch/Logstash/Kibana).
- DB:
  - slower Network call (in Milliseconds)
  - More **scalable** as its centralized
  - **Searching** is easy using sql queries.

Mahaam took a hybrid approach, first log to file(performance), then to DB asynchronously.

### Log config

- format:`timestamp level trafficID message %d{yyyy-MM-dd HH:mm:ss,SSS} %-5p %s%e%n`
- `%d{yyyy-MM-dd HH:mm:ss,SSS}`: timestamp with milliseconds
- `%-5p`: log level (left-aligned, 5 chars)
- `%s`: traffic ID for request correlation
- `%e%n`: exception stack trace and newline
- file size: 50M (rolls over when reached)
- number of files to keep: 10 (oldest deleted automatically)
- file names: mahaam.log, mahaam.log.1, mahaam.log.2, etc.
- location: /var/log/mahaam/

### Log Model

- timestamp: created_at
- level: type
- context: traffic_id, node_ip
- message

### Log Levels

- DEBUG: Detailed diagnostic info (used in dev).
- INFO: General app events (e.g., startup, shutdown).
- WARN: Unexpected but non-breaking situations.
- ERROR: Recoverable failures.
- FATAL / CRITICAL: Unrecoverable system errors.

Mahaam uses only two levels, **INFO** and **ERROR**.

### Log Retention & Rotation

- Define how long logs are kept.
- Use tools to:
  - Rotate log files (daily, size-based)
  - Compress or archive old logs
  - Delete expired logs (to save space)

### See

- Logging implementation in: [C#](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-cs/Src/Infra/Monitor/Log.cs), [Java](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-java/src/main/java/mahaam/infra/Log.java), [Go](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-go/utils/log/logs.go), [TypeScript](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-ts/src/infra/log.ts), [Python](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-py/infra/log.py)
