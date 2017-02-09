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
	svc := databasemigrationservice.New(s, &aws.Config{Region: aws.String("eu-west-1")})

	// Read the defaults file
	readTasks, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		log.Fatal("Couldn't read file "+tasksFile, err)
	}

	// Create tasks and unmarshal the JSON
	tasks := new([]ReplicationTask)
	err = json.Unmarshal(readTasks, tasks)
	if err != nil {
		log.Fatal("Couldn't JSON unmarshal file "+tasksFile, err)
	}

	// Start all the tasks stored in tasks
	for id, task := range *tasks {
		params := &databasemigrationservice.DeleteReplicationTaskInput{
			ReplicationTaskArn: aws.String(task.ReplicationTaskArn),
		}

		_, err := svc.DeleteReplicationTask(params)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "is running"):
				fmt.Println("Please stop task", task.ReplicationTaskIdentifier, "before trying to delete it")
				continue
			case strings.Contains(err.Error(), "is currently being stopped"):
				fmt.Println("Please wait until task", task.ReplicationTaskIdentifier, "has stopped before trying to delete it")
				continue
			case strings.Contains(err.Error(), "is already being deleted"):
				fmt.Println("Task", task.ReplicationTaskIdentifier, "is already being deleted")
				continue
			case strings.Contains(err.Error(), "not found"):
				removeTask(tasks, id)
				continue
			}

			fmt.Println("Couldn't delete Replication Task", err)
			continue
		}

		// Remove task from tasks.json
		removeTask(tasks, id)

		counter++
		fmt.Println("Task deleted: " + task.ReplicationTaskIdentifier)
	}

	// If we have no tasks left, delete the whole file
	switch {
	case len(*tasks) == 0:
		err := os.Remove(tasksFile)
		if err != nil {
			fmt.Println("Couldn't remove tasks files", err)
		}

	default:
		// Write remaining tasks to tasks-file
		writeTaskFile(tasks)
	}

	fmt.Println("\nDONE! Deleted", counter, "tasks.")
}

func removeTask(tasks *[]ReplicationTask, id int) {
	length := len(*tasks)
	switch {
	case length > id:
		*tasks = append((*tasks)[:id], (*tasks)[id+1:]...)
	case length == id:
		tasks = new([]ReplicationTask)
	}
}
