package main

import (
    "crypto/sha1"
    "fmt"
    "io"
    "math/rand"
    //"sort"
    "strconv"
    "time"
)

type job struct {
    JobNum    uint64  `json:"job_num"`
    Hash      string  `json:"hash"`
    Timestamp string  `json:"timestamp"`
    Lane      *string `json:"lane"`
    Priority  *uint   `json:"priority"`
    Message   *string `json:"message"`
}

type job_set []job

var queue = []job_set{}
var lanes = []string{"main"}
var indexes = map[string]map[uint]uint64{"main": make(map[uint]uint64)}

var idle_workers = 0
var total_jobs uint64 = 0
var current_jobs uint64 = 0

func (a job_set) Len() int      { return len(a) }
func (a job_set) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a job_set) Less(i, j int) bool {
    if *a[i].Priority < *a[j].Priority {
        return true
    }
    if *a[i].Priority > *a[j].Priority {
        return false
    }
    return a[i].JobNum < a[j].JobNum
}

func InitJob() job {
    mu.Lock()
    total_jobs += 1
    var job_num uint64 = total_jobs
    mu.Unlock()

    rand_number := rand.Intn(999999999)
    timestamp := time.Now().Format(time.RFC3339)

    hash := sha1.New()
    io.WriteString(hash, strconv.Itoa(rand_number)+timestamp)

    var default_lane string = JOB_DEFAULT_LANE
    var default_priority uint = JOB_DEFAULT_PRIORITY

    return job{
        JobNum:    job_num,
        Hash:      fmt.Sprintf("%x", hash.Sum(nil)),
        Timestamp: timestamp,
        Lane:      &default_lane,
        Priority:  &default_priority,
        Message:   nil,
    }
}

func InitJobTemplate() job {
    return job{
        JobNum:    0,
        Hash:      "",
        Timestamp: "",
        Lane:      nil,
        Priority:  nil,
        Message:   nil,
    }
}

func InsertJob(job job) {
op_start := time.Now()
    mu.Lock()

    // get the queue index where this lane is stored
    li := LaneIndex(*job.Lane)
    lane := *job.Lane

    queue_size := uint64(len(queue[li]))

    // if the queue is empty, simply make the job the only entry
    if queue_size == 0 {
        queue[li] = append(queue[li], job)
        indexes[lane][*job.Priority] = 1
    } else {
        // check for a boundary index
        i, ok := indexes[lane][*job.Priority]

        if ok {
            queue[li] = append(queue[li][:i], append(job_set{job}, queue[li][i:]...)...)
            indexes[lane][*job.Priority] += 1
        } else {
            i = uint64(i)
            priority := -1

            for i_qjob, qjob := range queue[li] {
                i = uint64(i_qjob)

                // keep priority state
                if int(*qjob.Priority) > priority {
                    priority = int(*qjob.Priority)
                }

                // if job is higher priority than everything, insert first
                if *job.Priority < *qjob.Priority {
                    queue[li] = append(job_set{job}, queue[li]...)
                    i += 1
                    break
                }

                // if job is lower priority than everything, insert last
                if *job.Priority >= *qjob.Priority && queue_size == i + 1 {
                    queue[li] = append(queue[li], job)
                    i = queue_size + 1
                    break
                }

                // job in middle of priorities
                if queue_size > i + 1 && *job.Priority < *queue[li][i + 1].Priority {
                    queue[li] = append(queue[li][:i+1], append(job_set{job}, queue[li][i+1:]...)...)
                    i += 2
                    break
                }
            }

            indexes[lane][*job.Priority] = i
        }

        // increment all boundaries
        for k, _ := range indexes[lane] {
            if k > *job.Priority {
                indexes[lane][k] += 1
            }
        }
    }

    current_jobs += 1

    if idle_workers > 0 {
        ready <- 1
    }

    mu.Unlock()
fmt.Println("mutex on, insertion took:", time.Since(op_start))
}

func LaneIndex(search_lane string) int {
    index := -1

    if search_lane == "" || search_lane == "main" {
        return 0
    }

    // check if index exists
    _, ok := indexes[search_lane]

    if !ok {
        indexes[search_lane] = make(map[uint]uint64)
    }

    for i, lane := range lanes {
        if search_lane == lane {
            index = i
            break
        }
    }

    if index == -1 {
        lanes = append(lanes, search_lane)
        queue = append(queue, job_set{})

        new_index := len(lanes) - 1

        return new_index
    } else {
        return index
    }
}

func ResetIndexes(lane string) {
    indexes[lane] = make(map[uint]uint64)
}

func DecrementIndexes(lane string) {
    // decrement all indexes such that they aren't zero
    for k, v := range indexes[lane] {
        if v - 1 == 0 {
            delete(indexes[lane], k)
            continue
        }

        indexes[lane][k] -= 1
    }
}
