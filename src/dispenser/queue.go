package main

type job struct {
	Priority int    `json:"priority"`
	Message  string `json:"message"`
}

var queue []job
