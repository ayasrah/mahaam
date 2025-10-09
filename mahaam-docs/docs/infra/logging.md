# Logging

### Overview

Logging is saving **events** done in the app to storage.

### Example

```log
2025-08-01 03:00:10,187 INFO mahaam-api-v1.0 started on nodeIP=1.2.3.4 with healthID=8eeb192c-ffef-4830-8954-12c69d3d6b5d
2025-08-01 03:10:58,187 INFO user created with id=2f7d68cd-0216-4b80-bdc2-de12d2a37099
2025-08-01 03:20:58,187 ERROR max allowed plans reached for userID=5d3c0e60-5030-4dad-973b-0ba9b97f7d5b
```

### Purpose

- Monitor the app
- Discover and fix errors

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
