package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func createTasks() {
	// Read the defaults file
	readSettings, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		log.Fatal("Couldn't read file " + settingsFile)
	}

	// Unmarshal the default settings into the Task object
	err = json.Unmarshal(readSettings, settings)

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

	// Add default exclude all
	exclude := defaultExclude()
	task.Mappings.TableMappings = append(task.Mappings.TableMappings, *exclude)

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

	// Add rename schema mapping - As last ID
	if settings.SourceSchema != settings.TargetSchema {
		rename := defaultRename()
		task.Mappings.TableMappings = append(task.Mappings.TableMappings, *rename)
	}

	return task
}
