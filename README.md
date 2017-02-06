# dispenser
Modern Job Queue for Modern Developers

Dispenser is a modern job queue designed to be: fast, reliable, feature rich, and tailored towards the needs of developers.

### Primary Features
- No compromise on speed or reliability
- Simple JSON interface
- Queue fully in memory with optional persistence
- Blocking Operation (receive a job immediately or block until ready)
- Non-blocking Operation (receive a job immediately or be notified none are available)

### API Reference

#### Status: `/`
##### Request Body
```
empty
```
##### Response Body
```json
{
  "name": "dispenser",
  "version": "0.0.1",
  "timestamp": "2017-02-05T23:42:51-06:00",
  "status": "ok",
  "payload": {
    "system": {
      "pid": 10169,
      "cpu_count": 8
    },
    "stats": {
      "idle_workers": 0,
      "queued_jobs": 0
    }
  }
}
```

#### List Jobs: `/jobs`
##### Request Body
```
empty
```
##### Response Body (queue empty)
```json
[]
```
##### Response Body (queue not empty)
```json
[
  {
    "job_num": 1,
    "hash": "5066400f81e3ba8e6160279b4fad9d6ed5598584",
    "timestamp": "2017-02-05T23:48:14-06:00",
    "priority": 1,
    "message": "job message here"
  }
]
```
