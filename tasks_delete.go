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
func deleteTasks() {
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
		params := &databasemigrationservice.DeleteReplicationTaskInput{
			ReplicationTaskArn: aws.String(task.ReplicationTaskArn),
		}

		_, err := svc.DeleteReplicationTask(params)
		if err != nil {
			// If the task doesn't exists we shouldn't keep it in the tasks.json file - just continue
			if strings.Contains(err.Error(), notFound) {
				continue
			}

			// If the task errored and not because it doesn't exists, keep it in the tasks.json file
			*remainingTasks = append(*remainingTasks, task)

			// Go through the different statuses that might have made the action failed
			switch {
			case strings.Contains(err.Error(), isRunning):
				fmt.Println("Please stop task", task.ReplicationTaskIdentifier, stopBeforeDeleting)
				continue
			case strings.Contains(err.Error(), "is currently being stopped"):
				fmt.Println("Please wait until task", task.ReplicationTaskIdentifier, waitForStop)
				continue
			case strings.Contains(err.Error(), "is already being deleted"):
				fmt.Println("Task", task.ReplicationTaskIdentifier, beingDeleted)
				continue
			}

			fmt.Println("Couldn't delete Replication Task", err)
			continue
		}

		counter++
		fmt.Println("Task deleted: " + task.ReplicationTaskIdentifier)
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

	fmt.Println("\nDONE! Deleted", counter, "tasks.")
}
