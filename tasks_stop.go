package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

// startTasks will start all tasks stored in tasksFile
func stopTasks() {
	counter := 0

	// Create AWS session
	s, err := session.NewSession()
	if err != nil {
		log.Fatal("Couldn't create AWS Session.")
	}

	// Create the AWS Service
	svc := databasemigrationservice.New(s, &aws.Config{Region: &region})

	// Read the defaults file
	readTasks, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		log.Fatal("Couldn't read file "+tasksFile, err)
	}

	// Create tasks and unmarshal the JSON
	tasks := new([]ReplicationTask)
	remainingTasks := new([]ReplicationTask) // Tasks that will be saved (if they couldn't be removed for example)
	err = json.Unmarshal(readTasks, tasks)
	if err != nil {
		log.Fatal("Couldn't JSON unmarshal file "+tasksFile, err)
	}

	// Start all the tasks stored in tasks
	for _, task := range *tasks {
		params := &databasemigrationservice.StopReplicationTaskInput{
			ReplicationTaskArn: aws.String(task.ReplicationTaskArn),
		}

		_, err := svc.StopReplicationTask(params)
		if err != nil {
			// If the task doesn't exists we shouldn't keep it in the tasks.json file - just continue
			if strings.Contains(err.Error(), doesntExists) {
				continue
			}

			// If the task errored and not because it doesn't exists, keep it in the tasks.json file
			*remainingTasks = append(*remainingTasks, task)

			switch {
			case strings.Contains(err.Error(), currentlyNotRunning):
				fmt.Println("Task", task.ReplicationTaskIdentifier, "is currently not running")
				continue
			case strings.Contains(err.Error(), allreadyBeingStopped):
				fmt.Println("Task", task.ReplicationTaskIdentifier, "is already being stopped")
				continue
			}

			fmt.Println("Couldn't stop Replication Task", err)
			continue
		}

		// If the task errored and not because it doesn't exists, keep it in the tasks.json file
		*remainingTasks = append(*remainingTasks, task)

		counter++
		fmt.Println("Task stopped: " + task.ReplicationTaskIdentifier)
	}

	// If we have no tasks left, delete the whole file
	switch {
	case len(*remainingTasks) == 0:
		err := os.Remove(tasksFile)
		if err != nil {
			fmt.Println("Couldn't remove tasks files", err)
		}

	default:
		// Write remaining tasks to tasks-file
		writeTaskFile(remainingTasks)
	}

	fmt.Println("\nDONE! Stopped", counter, "tasks.")
}
