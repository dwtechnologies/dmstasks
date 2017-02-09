package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		params := &databasemigrationservice.StopReplicationTaskInput{
			ReplicationTaskArn: aws.String(task.ReplicationTaskArn),
		}

		_, err := svc.StopReplicationTask(params)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "is not currently running"):
				fmt.Println("Task", task.ReplicationTaskIdentifier, "is currently not running")
				continue
			case strings.Contains(err.Error(), "is already being stopped"):
				fmt.Println("Task", task.ReplicationTaskIdentifier, "is already being stopped")
				continue
			case strings.Contains(err.Error(), "does not exist"):
				fmt.Println("Task", task.ReplicationTaskIdentifier, "doesn't exists")
				removeTask(tasks, id)
			}

			fmt.Println("Couldn't stop Replication Task", err)
			continue
		}

		counter++
		fmt.Println("Task stopped: " + task.ReplicationTaskIdentifier)
	}

	fmt.Println("\nDONE! Stopped", counter, "tasks.")
}
