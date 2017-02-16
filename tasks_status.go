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
	readConfig()
	counter := 0

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

	// Create filter to use for gathering tasks from AWS
	taskFilter := []*string{}
	for _, task := range *tasks {
		taskArn := task.ReplicationTaskArn
		taskFilter = append(taskFilter, &taskArn)
	}

	// Create AWS session
	s, err := session.NewSession()
	if err != nil {
		log.Fatal("Couldn't create AWS Session.")
	}

	// Create the AWS Service
	svc := databasemigrationservice.New(s, &aws.Config{Region: &region})

	// Describe all tasks on AWS account
	params := &databasemigrationservice.DescribeReplicationTasksInput{
		Filters: []*databasemigrationservice.Filter{
			{
				Name:   aws.String("replication-task-arn"),
				Values: taskFilter,
			},
		},
	}

	// Create the response
	resp, err := svc.DescribeReplicationTasks(params)
	if err != nil {
		log.Fatal("Couldn't describe Replication Task", err)
	}

	// Marshal the output and unmarshal it to golang
	output := new(ReplicationTasks)
	outputMarshaled, _ := json.Marshal(resp)
	err = json.Unmarshal(outputMarshaled, output)
	if err != nil {
		log.Fatal("Couldn't JSON Unmarshal Output from Replication Task", err)
	}

	for _, dmsTask := range output.ReplicationTasks {
		fmt.Println("Task " + dmsTask.ReplicationTaskIdentifier + " is " + dmsTask.Status)
		counter++
	}

	fmt.Println("\nDONE! Listed status for", counter, "tasks.")
}
