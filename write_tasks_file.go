package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func writeTaskFile(t *[]ReplicationTask) {
	// Write output to file so that we can start/stop/resume at a later point without asking AWS for info
	tasksJSON, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		fmt.Println("Couldn't JSON marshal list of tasks", err)
		return
	}
	ioutil.WriteFile(tasksFile, tasksJSON, 0644)
}
