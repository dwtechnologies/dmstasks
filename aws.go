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

// createTasksOnAws tasks t *Tasks (Tasks to create) and creates them on AWS.
func createTasksOnAws(t *Tasks) {
	counter := 0

	// Slice to store created info on
	tasksCreated := []ReplicationTask{}

	// Create AWS session
	s, err := session.NewSession()
	if err != nil {
		log.Fatal("Couldn't create AWS Session.")
	}

	// Create the AWS Service
	svc := databasemigrationservice.New(s, &aws.Config{Region: aws.String("eu-west-1")})

	// Loop through the Tasks
	for _, task := range *t.Tasks {
		// Convert TableMappings and ReplicationTaskSettings to JSON
		tableMappings, err := json.Marshal(task.Mappings)
		if err != nil {
			log.Fatal("Couldn't convert TableMappings to JSON", err)
		}

		replicationSettings, err := json.Marshal(task.ReplicationTaskSettings)
		if err != nil {
			log.Fatal("Couldn't convert ReplicationSettings to JSON", err)
		}

		params := &databasemigrationservice.CreateReplicationTaskInput{
			MigrationType:             &task.MigrationType,
			ReplicationInstanceArn:    &task.ReplicationInstanceArn,
			ReplicationTaskIdentifier: &task.ReplicationTaskIdentifier,
			SourceEndpointArn:         &task.SourceEndpointArn,
			TargetEndpointArn:         &task.TargetEndpointArn,
			TableMappings:             aws.String(string(tableMappings)),
			ReplicationTaskSettings:   aws.String(string(replicationSettings)),
			Tags: []*databasemigrationservice.Tag{
				{
					Key:   aws.String("Name"),
					Value: &task.ReplicationTaskIdentifier,
				},
			},
		}

		resp, err := svc.CreateReplicationTask(params)
		if err != nil {
			log.Fatal("Couldn't create Replication Task", err)
		}

		// Marshal the output and unmarshal it to golang
		output := new(Reply)
		stringMarshaled, _ := json.Marshal(resp)
		_ = json.Unmarshal(stringMarshaled, output)

		fmt.Println("Task created: " + output.Reply.ReplicationTaskIdentifier)
		tasksCreated = append(tasksCreated, output.Reply)

		counter++
	}

	// Write output to file so that we can start/stop/resume at a later point without asking AWS for info
	tasksJSON, _ := json.MarshalIndent(tasksCreated, "", "\t")
	ioutil.WriteFile(tasksFile, tasksJSON, 0644)

	fmt.Println("DONE! Created", counter, "tasks.")
}

func startTasksOnAws(t *Tasks) {

}
