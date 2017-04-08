package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "time"
)

func WriteQueue() {
    op_start := time.Now()

    mu.Lock()
    json_queue, _ := json.MarshalIndent(queue, "", "  ")
    ioutil.WriteFile(ROOT+"/config/queue.json", json_queue, 0644)
    mu.Unlock()

    fmt.Println("queue persist took:", time.Since(op_start))
}

func Persist() {
    // start a loop and write the queue at the prescribed time
    for {
        time.Sleep(time.Duration(config.PersistInterval) * time.Second)

        // don't persist if dispenserd is shutting down
        if !running {
            break
        }

        WriteQueue()
    }
}

func LoadQueue() {
    // load in the queue
    json_queue, err := ioutil.ReadFile(ROOT + "/config/queue.json")

    current_jobs = 0

    var tmp_queue = []job_set{}

    if err == nil {
        ju_err := json.Unmarshal(json_queue, &tmp_queue)

        if ju_err == nil {
            var highest_job_num uint64 = 0

            var current_lane string = "main"

            for _, lane := range tmp_queue {
                // skip over empty saved queues
                if len(lane) == 0 {
                    continue
                }

                // add the current jobs from the length of the lane
                current_jobs += uint64(len(lane))

                for _, v := range lane {
                    // ascertain the newest high job num
                    if v.JobNum > highest_job_num {
                        highest_job_num = v.JobNum
                    }

                    // ascertain if the lane exists and append if it doesn't
                    if !UtilInArray(*v.Lane, lanes) {
                        lanes = append(lanes, *v.Lane)
                    }

                    // continue establish state of the current lane for later
                    if v.Lane != nil {
                        current_lane = *v.Lane
                    }
                }

                // determine where the lane will go, queue[0] is established by default
                if current_lane == "main" {
                    queue[0] = lane
                } else {
                    queue = append(queue, lane)
                }
            }

            total_jobs = highest_job_num
        }
    }
}
