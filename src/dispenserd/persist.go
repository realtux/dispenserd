package main

import (
    "encoding/json"
    //"fmt"
    "io/ioutil"
    "time"
)

func WriteQueue() {
//op_start := time.Now()
    mu.Lock()
    json_queue, _ := json.MarshalIndent(queue, "", "  ")
    ioutil.WriteFile(ROOT+"/config/queue.json", json_queue, 0644)
    mu.Unlock()
//fmt.Println("queue persist took:", time.Since(op_start))
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

    var tmp_queue = map[string]job_set{}

    if err != nil {
        return
    }

    ju_err := json.Unmarshal(json_queue, &tmp_queue)

    if ju_err != nil {
        return
    }

    var highest_job_num uint64 = 0
    var current_lane string = "main"

    for _, lane := range tmp_queue {
        // skip over empty saved queues
        if len(lane) == 0 {
            continue
        }

        var previous_priority uint = 0
        priority := -1

        for i, v := range lane {
            if priority == -1 {
                priority = int(*v.Priority)
            }

            // continue establish state of the current lane for later
            if v.Lane != nil {
                current_lane = *v.Lane

                // add lane to indexes
                _, ok := indexes[current_lane]

                if !ok {
                    indexes[current_lane] = make(map[uint]uint64)
                }
            }

            // ascertain the newest high job num
            if v.JobNum > highest_job_num {
                highest_job_num = v.JobNum
            }

            // ascertain if the lane exists and append if it doesn't
            if !UtilInArray(*v.Lane, lanes) {
                lanes = append(lanes, *v.Lane)
            }

            // check for priority change to establish index boundaries
            if int(*v.Priority) != priority && int(*v.Priority) > priority {
                indexes[current_lane][previous_priority] = uint64(i)
                priority = int(*v.Priority)
            }

            previous_priority = *v.Priority
        }

        // determine where the lane will go, queue[0] is established by default
        queue[current_lane] = lane
    }

    // set the total jobs so dispenserd knows what the next unique job id is
    total_jobs = highest_job_num
}
