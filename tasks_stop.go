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
	for _, task := range *tasks {
		params := &databasemigrationservice.StopReplicationTaskInput{
			ReplicationTaskArn: aws.String(task.ReplicationTaskArn),
		}

		_, err := svc.StopReplicationTask(params)
		if err != nil {
			// Check to see that the task actually is running
			if strings.Contains(err.Error(), "is not currently running") {
				continue
			}
			fmt.Println("Couldn't stop Replication Task", err)
			continue
		}

		counter++
		fmt.Println("Task stopped: " + task.ReplicationTaskIdentifier)
	}

	fmt.Println("\nDONE! Stopped", counter, "tasks.")
}
