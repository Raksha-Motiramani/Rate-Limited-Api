# Rate-Limited API with Redis (Sliding Window + Queue Worker)

## Overview

This project implements a **rate-limited API service in Go using Redis**. The system enforces a per-user sliding-window rate limit, queues excess requests for retry, and exposes a `/stats` endpoint to monitor usage.

The implementation is designed to be **concurrency-safe**, **observable**, and **production-considerate**, while remaining simple enough to run locally without deployment.

---

## Features

* Sliding-window rate limiter (per user)
* Redis-backed request tracking
* Retry queue for throttled requests
* Background worker to process queued requests
* `/stats` endpoint for monitoring usage
* Concurrency-safe design using atomic Redis operations

---

## Tech Stack

* Go (Golang)
* Redis
* REST API (net/http)

---

## Project Structure

```
handler/
service/
redis/
worker/
main.go
```

Responsibilities:

* **handler/** → HTTP request handling
* **service/** → rate limiting + stats logic
* **redis/** → Redis client initialization
* **worker/** → background retry queue processor

---

## How to Run the Project Locally

### 1. Start Redis

Make sure Redis is installed and running:

```
redis-server
```

Verify:

```
redis-cli ping
```

Expected output:

```
PONG
```

---

### 2. Install dependencies

Inside the project folder:

```
go mod tidy
```

---

### 3. Start the API server

Run:

```
go run main.go
```

Server starts at:

```
http://localhost:8080
```

---

### 4. Test the request endpoint

Example:

```
curl -X POST http://localhost:8080/request \\
-H "Content-Type: application/json" \\
-d '{"user_id":"userTest","payload":"hello"}'
```

Send multiple requests quickly to trigger rate limiting.

---

### 5. Check stats endpoint

```
curl http://localhost:8080/stats
```

Example response:

```
{
  "userTest": {
    "total_requests": 5,
    "last_minute_requests": 5,
    "queued_requests": 2
  }
}
```

---

## Rate Limiting Algorithm

This project uses a **Sliding Window Rate Limiting Algorithm** implemented with Redis Sorted Sets.

For each user:

```
rate_limit:<userID>
```

Stores timestamps of requests within the last 60 seconds.

Workflow:

1. Add request timestamp
2. Remove expired timestamps
3. Count active timestamps
4. Allow or reject request

This guarantees accurate rate limiting even under concurrent traffic.

---

## Retry Queue Design

When a request exceeds the rate limit:

```
LPUSH request_queue
```

A background worker continuously processes queued requests:

```
RPOP request_queue
```

Benefits:

* prevents request loss
* smooths burst traffic
* simulates production retry pipelines

---

## Stats Endpoint Design

The `/stats` endpoint returns:

```
{
  total_requests
  last_minute_requests
  queued_requests
}
```

Metrics are derived from:

| Metric               | Source             |
| -------------------- | ------------------ |
| total_requests       | Redis counter      |
| last_minute_requests | Sorted set window  |
| queued_requests      | Redis queue length |

This provides lightweight observability into system usage.

---

## Concurrency Handling

Concurrency safety is ensured using **atomic Redis operations**:

* ZADD
* ZREMRANGEBYSCORE
* ZCARD
* LPUSH
* RPOP
* LLEN

Redis guarantees atomic execution per command, ensuring accurate rate limiting under parallel requests.

This allows multiple clients to safely interact with the API simultaneously.

---

## Design Decisions

### 1. Redis for Rate Limiting

Redis was chosen because:

* supports sorted sets
* provides atomic operations
* extremely fast
* commonly used in production throttling systems

---

### 2. Sliding Window Instead of Fixed Window

Sliding window provides:

* smoother traffic control
* fewer edge-case bursts
* more accurate enforcement

compared to fixed window rate limiting.

---

### 3. Background Worker for Retry Queue

Instead of rejecting excess requests permanently:

queued retry processing improves:

* reliability
* fairness
* user experience

---

### 4. Stats Endpoint for Observability

Production systems require monitoring.

Providing `/stats` allows:

* debugging
* traffic visibility
* usage tracking

without external tooling.

---

## Limitations

Current implementation uses a single-node Redis instance.

If Redis becomes unavailable:

* rate limiting stops working
* queue processing pauses
* stats become unavailable

Production systems typically use:

* Redis Sentinel
* Redis Cluster

for high availability.

Another limitation:

Queue is currently global instead of partitioned per user.

Large-scale systems may shard queues across workers.

---

## Improvements With More Time

If extended further, the following enhancements would be added:

### 1. Redis Lua Script for Fully Atomic Rate Limiting

Combine:

```
ZADD
ZREMRANGEBYSCORE
ZCARD
```

into a single atomic Redis script.

Improves correctness under extreme concurrency.

---

### 2. Distributed Worker Scaling

Allow multiple worker instances to consume the retry queue.

Improves throughput under heavy traffic.

---

### 3. Per‑User Retry Queues

Replace:

```
request_queue
```

with:

```
request_queue:<userID>
```

Improves fairness across users.

---

### 4. Configurable Rate Limits

Allow dynamic limits via configuration instead of constants.

Example:

```
LIMIT=5
WINDOW=60
```

loaded from environment variables.

---

### 5. Metrics Integration

Expose Prometheus-compatible metrics such as:

* request throughput
* rejection rate
* queue size trends

---

## Conclusion

This project demonstrates a production-considerate implementation of a sliding-window rate-limited API using Redis.

It supports concurrent clients safely, prevents request loss using a retry queue, and provides runtime observability through a stats endpoint.

The system can be extended easily for distributed deployment scenarios with minimal architectural changes.
