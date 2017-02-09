package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	settingsFile = "settings.json"
	tablesFile   = "tables.txt"
	tasksFile    = "tasks.json"
)

var (
	ruleid   = 1
	settings = &Task{}
)

func main() {
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

	// Generate the tasks
	tasks := createTasks()

	// Generate task JSON files
	createTasksOnAws(tasks)
}
