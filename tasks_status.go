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

// statusTasks will show status of all the Tasks in tasks.json
func statusTasks() {
	counter := 0

	// Create AWS session
	s, err := session.NewSession()
	if err != nil {
		log.Fatal("Couldn't create AWS Session.")
	}

	// Create the AWS Service
	svc := databasemigrationservice.New(s, &aws.Config{Region: aws.String("eu-west-1")})

	// Describe all tasks on AWS account
	params := &databasemigrationservice.DescribeReplicationTasksInput{}

	// Create the response
	resp, err := svc.DescribeReplicationTasks(params)
	if err != nil {
		log.Fatal("Couldn't describe Replication Task", err)
	}

	// Marshal the output and unmarshal it to golang
	output := new([]ReplicationTasks)
	stringMarshaled, _ := json.Marshal(resp)
	err = json.Unmarshal(stringMarshaled, output)
	if err != nil {
		log.Fatal("Couldn't JSON Unmarshal Output from Replication Task", err)
	}

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
		status := ""

		counter++
		fmt.Println("Task " + task.ReplicationTaskIdentifier + " is " + status)
	}

	fmt.Println("\nDONE! Listed status for", counter, "tasks.")
}
