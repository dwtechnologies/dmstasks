package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

func createTasks() {
	readConfig()

	// Check that the minimum amount of settings is set
	if settings.ReplicationTaskIdentifier == "" || settings.SourceEndpointArn == "" || settings.TargetEndpointArn == "" || settings.ReplicationInstanceArn == "" || settings.SourceSchema == "" || settings.MigrationType == "" {
		fmt.Printf("usage: dmstasks\n")
		fmt.Printf("Atleast ReplicationTaskIdentifier, SourceEndpointArn, TargetEndpointArn, ReplicationInstanceArn, SourceSchema and MigrationType must be set! If TargetSchema is left empty it will be set to the same value as SourceSchema.\n\n")
		flag.PrintDefaults()
		fmt.Printf("\n")
		os.Exit(0)
	}

	// If TargetSchema is emtpy set to same as SourceSchema
	if settings.TargetSchema == "" {
		settings.TargetSchema = settings.SourceSchema
	}

	// Generate and create the tasks
	tasks := new(Tasks)
	tasks.Tasks = addTasks()

	// Generate task JSON files
	createTasksOnAws(tasks)
}

// addTasks takes r *int (ruleid number) and generates Mappings to tables based on tables file. Returns *[]Rules
func addTasks() *[]Task {
	tasks := new([]Task)

	// Read the tables from file
	tblsFile, err := ioutil.ReadFile(tablesFile)
	if err != nil {
		log.Fatal("Couldn't read file " + tablesFile)
	}

	// Split on [, to get name and all tables
	tbls := strings.Split(string(tblsFile), "[")

	for _, row := range tbls {
		ruleid = 1
		task := genTask(&row)
		if task == nil {
			continue
		}
		*tasks = append(*tasks, *task)
	}

	return tasks
}

func genTask(r *string) *Task {
	task := new(Task)
	*task = *settings

	// Split by name
	split := strings.SplitN(*r, "]", 2)
	if len(split) != 2 {
		return nil
	}

	// Set the name of the task
	name := split[0]
	task.ReplicationTaskIdentifier += name

	// Set CloudWatch LogGroup
	if task.ReplicationTaskSettings.Logging.EnableLogging == true {
		task.ReplicationTaskSettings.Logging.CloudWatchLogGroup = "dms-tasks-" + task.ReplicationTaskIdentifier
		task.ReplicationTaskSettings.Logging.CloudWatchLogStream = task.ReplicationTaskSettings.Logging.CloudWatchLogGroup
	}

	// Get the mappings for the current job
	params := strings.Split(split[1], "\n")
	for _, param := range params {
		param := strings.TrimSpace(param)
		if param == "" {
			continue
		}
		rule := genRule(&param)
		task.Mappings.TableMappings = append(task.Mappings.TableMappings, *rule)
	}

	// Add default exclude all
	// exclude := defaultExclude()
	// task.Mappings.TableMappings = append(task.Mappings.TableMappings, *exclude)

	// Add rename schema mapping - As last ID
	if settings.SourceSchema != settings.TargetSchema {
		rename := rename()
		task.Mappings.TableMappings = append(task.Mappings.TableMappings, *rename)
	}

	return task
}

// genRule take r *string (row) and creates a include Rule out of it
func genRule(r *string) *Rules {
	rulestr := strconv.Itoa(ruleid)
	ruleid++
	return &Rules{
		RuleType: "selection",
		RuleID:   rulestr,
		RuleName: rulestr,
		ObjectLocator: ObjectLocator{
			SchemaName: settings.SourceSchema,
			TableName:  *r,
		},
		RuleAction: "include",
	}
}

// rename creates rule for renaming schema from source name to target name and returns *Rules
func rename() *Rules {
	rulestr := strconv.Itoa(ruleid)
	ruleid++
	return &Rules{
		RuleType:   "transformation",
		RuleID:     rulestr,
		RuleName:   rulestr,
		RuleTarget: "schema",
		ObjectLocator: ObjectLocator{
			SchemaName: settings.SourceSchema,
		},
		RuleAction: "rename",
		Value:      settings.TargetSchema,
	}
}

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
	svc := databasemigrationservice.New(s, &aws.Config{Region: &region})

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
			switch {
			case strings.Contains(err.Error(), alreadyExists):
				fmt.Println("Task", task.ReplicationTaskIdentifier, "already exists")
				continue
			}

			fmt.Println("Couldn't create Replication Task", err)
			continue
		}

		// Marshal the output and unmarshal it to golang
		output := new(Reply)
		stringMarshaled, _ := json.Marshal(resp)
		err = json.Unmarshal(stringMarshaled, output)
		if err != nil {
			fmt.Println("Couldn't JSON Unmarshal Output from Replication Task", err)
			continue
		}

		counter++
		fmt.Println("Task created: " + output.Reply.ReplicationTaskIdentifier)
		tasksCreated = append(tasksCreated, output.Reply)
	}

	// Load whats currenty in the file
	readTasks, err := ioutil.ReadFile(tasksFile)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			log.Fatal("Couldn't read file "+tasksFile, err)
		}
	}

	// Create tasksFromFile and unmarshal the JSON
	tasksFromFile := new([]ReplicationTask)
	err = json.Unmarshal(readTasks, tasksFromFile)
	if err != nil {
		if string(readTasks) != "" {
			log.Fatal("Couldn't JSON unmarshal file "+tasksFile, err)
		}
	}

	if (len(tasksCreated)) > 0 {
		tasksCreated = append(tasksCreated, *tasksFromFile...)
		// Write tasks.jsons file
		writeTaskFile(&tasksCreated)
	}

	fmt.Println("\nDONE! Created", counter, "tasks.")
}
