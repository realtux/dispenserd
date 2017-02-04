package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

// config
type options struct {
    PersistQueue    bool `json:"persist_queue"`
    PersistInterval int  `json:"persist_interval"`
}

func init_options() options {
    return options{
        PersistQueue:    false,
        PersistInterval: 60,
    }
}

var config = init_options()

func ConfigLoad() {
    // read config
    data, err := ioutil.ReadFile("config/config.json")

    if err != nil {
        fmt.Println("could not open config.json, perhaps it doesn't exist?")
        os.Exit(1)
    }

    ju_err := json.Unmarshal(data, &config)

    if ju_err != nil {
        fmt.Println("error parsing config.json, likely invalid json")
        os.Exit(1)
    }

    if config.PersistQueue {
        go Persist()
    }
}
