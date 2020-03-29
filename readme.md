# dispenserd [![Build Status](https://travis-ci.org/realtux/dispenserd.svg?branch=master)](https://travis-ci.org/realtux/dispenserd)
Yet Another Job Queue, But Better

dispenserd is a job queue designed to be: fast, reliable, feature rich, and tailored towards the needs of developers.

### Primary Features
- No compromise on stability or reliability
- Simple JSON interface
- Full queue in memory with optional persistence
- Blocking Operation (receive a job immediately or block until ready)
- Non-blocking Operation (receive a job immediately or be notified none are available)
- Job priorities
- Job lanes to separate classes of jobs for different consumers

---

### Contributing
If you'd like to contribute, fork the repo and beginning writing code. If you open a PR, tests will begin to run. If tests don't pass, please correct the part of your code causing tests to fail and resubmit your PR. A dispenderd roadmap will be published in the coming weeks.

---

### Dependencies
- Golang 1.6+ (tests run against 1.6.x, 1.7.x, 1.8.x)

### Installation

##### Run from source folder (Linux/macOS)
```bash
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

##### CentOS/RHEL package
```
coming soon
```

##### Ubuntu package
```
coming soon
```

##### Optionally run tests (Python 3+ required)
```bash
./run_tests
```

---

### Configuration

You can specify a configuration file located at `$ROOT/config/config.json` if you'd like, but it's completely optional. The config file that comes with the source represents all of the default values. dispenserd will run just fine with no modifications or no config file at all.

#### address (0-255,0-255,0-255,0-255:1-65536), default: 0.0.0.0:8282
Set this to the address you'd like dispenserd to bind to.

#### persist_queue [true/false], default: true
Setting `persist_queue` to true will cause the following three things to happen over the course of dispenserd running:
 1. Upon start-up, dispenserd will look for `$ROOT/data/queue.json`. If present, it will attempt to parse that file and load the jobs into the queue.
 2.  Periodically during program execution, every `persist_interval` seconds, dispenserd will write the contents of the queue to `$ROOT/data/queue.json`.
 3.  Upon receiving a SIGTERM, dispenserd will write the contents of the queue to `$ROOT/data/queue.json`.

#### persist_interval [0-2147483647], default: 3600
If `persist_queue` is set to true, `persist_interval` will be used by dispenserd to determine how often to write the contents of the queue to `$ROOT/data/queue.json`. A value of `0` will tell dispenserd to never persist the queue during program execution and only persist the queue when a SIGTERM is received.

If for some reason your use case involves queues that are abnormally large (perhaps > 2,000,000 pending jobs), the size of `$ROOT/data/queue.json` could be quite large and take time to write to disk, depending on your hardware. Unless your server infrastructure is prone to abrupt power loss, it's probably safe to leave this at the default value.

#### throttle_receive [true/false], default: false
If `throttle_receive` is set to true, a `25ms` delay will be enforced between consuming jobs via `/receive_block`. This effectively limits receive operations to 40/sec max.

#### throttle_schedule [true/false], default: false
If `throttle_schedule` is set to true, a `25ms` delay will be enforced between scheduling jobs. This effectively limits schedule operations to 40/sec max.

---

### Performance Considerations

#### Scalability with a high rate of job scheduling/receiving
Because dispenserd uses HTTP and a new connection for every operation, it's possible to run into a situation where a system runs out of sockets to use to connect with dispenserd. The time where this would most likely happen is during sustained high rate of job scheduling with or without a high rate of job receives. After a connection is closed, the connection sits in a `TIME_WAIT` status for somewhere between 30-120 seconds on average. By default, you only have about 29,000 sockets available. If you are scheduling/receiving jobs at a rate faster than 200/sec, you will likely hit this limit. There are a few ways of handling this situation:
 1. Throttle dispenserd with either the `throttle_schedule` and/or `throttle_receive` options.
 2. Tune your system to lower the `TIME_WAIT` status.
 3. Tune your system to reuse sockets in the `TIME_WAIT` status.
Any one of these will work, and each has their pros and cons.

---

### API

### Status: `/`
Required Request Parameters: None

Optional Request Parameters: None

##### Request Body
```
empty
```
##### Response Body
```json
{
    "name": "dispenserd",
    "version": "x.x.x",
    "timestamp": "2017-03-14T15:42:51-06:00",
    "status": "ok",
    "payload": {
        "system": {
            "pid": 12345,
            "cpu_count": 32
        },
        "queued_jobs": {
            "main": 0
        },
        "idle_workers": {
            "main": 0
        }
    }
}
```

### List Jobs: `/jobs`
Required Request Parameters: None

Optional Request Parameters: None

##### Request Body
```
empty
```
##### Response Body (single lane, queue empty)
```json
{
    "main": []
}
```
##### Response Body (multiple lanes, all queues empty)
```json
{
    "main": [],
    "other_lane": [],
    "other_lane2": []
}
```
##### Response Body (queue not empty)
```json
{
    "main": [
        {
            "job_num": 1,
            "hash": "5066400f81e3ba8e6160279b4fad9d6ed5598584",
            "timestamp": "2017-02-05T23:48:14-06:00",
            "priority": 1,
            "message": "job message here"
        }
    ]
}
```

### Schedule Job: `/schedule`
Required Request Parameters:
 - **message:** string

Optional Request Parameters:
 - **lane:** string, default: main
 - **priority:** 0-4294967295, default: 512

##### Request Body (main lane)
```json
{
    "priority": 0,
    "message": "msg here. stringified json, serialized objects, etc..."
}
```
##### Request Body (specific lane)
```json
{
    "priority": 0,
    "lane": "video",
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

### Receive Job (blocking): `/receive_block`
Required Request Parameters: None

Optional Request Parameters:
 - **lane:** string, default: main

##### Request Body (main lane)
```
empty
```
##### Request Body (specific lane)
```json
{
    "lane": "video"
}
```
##### Response Body (no job ready)
```
since this is a blocking response, the request will simply hang until a job is ready
```
##### Response Body (job ready)
```json
{
    "status": "ok",
    "code": 0,
    "message": "message from your job"
}
```

### Receive Job (non-blocking): `/receive_noblock`
Required Request Parameters: None

Optional Request Parameters:
 - **lane:** string, default: main

##### Request Body
```
empty
```
##### Request Body (specific lane)
```json
{
    "lane": "video"
}
```
##### Response Body (no job ready)
```json
{
    "status": "ok",
    "code": 2,
    "message": "empty queue"
}
```
##### Response Body (job ready)
```json
{
    "status": "ok",
    "code": 0,
    "message": "message from your job"
}
```
