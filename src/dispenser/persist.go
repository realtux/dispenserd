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
    for {
        time.Sleep(time.Duration(config.PersistInterval) * time.Second)

        if !running {
            break
        }

        WriteQueue()
    }
}

func LoadQueue() {
    json_queue, err := ioutil.ReadFile(ROOT + "/config/queue.json")

    if err == nil {
        ju_err := json.Unmarshal(json_queue, &queue)

        if ju_err == nil {
            var highest_job_num uint64 = 0

            for _, v := range queue {
                if v.JobNum > highest_job_num {
                    highest_job_num = v.JobNum
                }
            }

            total_jobs = highest_job_num
        }
    }
}
