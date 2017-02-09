package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

// startTasks will start all tasks stored in tasksFile
func startTasks(t string) {
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
		params := &databasemigrationservice.StartReplicationTaskInput{
			ReplicationTaskArn:       aws.String(task.ReplicationTaskArn),
			StartReplicationTaskType: &t,
		}

		_, err := svc.StartReplicationTask(params)
		if err != nil {
			fmt.Println("Couldn't start Replication Task", err)
			continue
		}

		fmt.Println("Task started: " + task.ReplicationTaskIdentifier)
	}

	counter++
	fmt.Println("\nDONE! Started", counter, "tasks.")
}
