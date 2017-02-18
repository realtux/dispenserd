# dispenserd
Modern Job Queue for Modern Developers

Dispenserd is a modern job queue designed to be: fast, reliable, feature rich, and tailored towards the needs of developers.

### Primary Features
- No compromise on stability or reliability
- Simple JSON interface
- Queue fully in memory with optional persistence
- Blocking Operation (receive a job immediately or block until ready)
- Non-blocking Operation (receive a job immediately or be notified none are available)
- Job priorities

### Dependencies
- Golang 1.6+

### Installation
---

##### Run from source folder (Linux/macOS)
```
# get the code
git clone https://github.com/ebrian/dispenser
cd dispenser

# build
./build

# run in foreground
./dispenserd

# run in background
./dispenserd -d
```

---

### API Reference
---

#### Status: `/`
##### Request Body
```
empty
```
##### Response Body
```json
{
  "name": "dispenserd",
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

---

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

---

#### Schedule Job: `/schedule`
##### Request Body
```json
{
    "priority": 1,
    "message": "msg here. stringified json, serialized objects, etc..."
}
```
##### Response Body (valid request)
```json
{
  "status": "ok",
  "code": 0
}
```
##### Response Body (invalid request)
```json
{
  "status": "error",
  "code": 1,
  "message": "missing message"
}
```

---

#### Receive Job (blocking): `/receive_block`
##### Request Body
```
empty
```
##### Response Body (no job ready)
```
since this is a blocking response, the request will simply hang until a job is ready
```
##### Response Body (job ready)
```json
plain text job
```

---

#### Receive Job (non-blocking): `/receive_noblock`
##### Request Body
```
empty
```
##### Response Body (no job ready)
```json
{
  "status": "ok",
  "code": 0,
  "message": "empty queue"
}
```
##### Response Body (job ready)
```json
plain text job
```
