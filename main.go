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
	defaults = "defaults.json"
	tables   = "tables.txt"
)

var (
	source       = new(string)
	target       = new(string)
	replication  = new(string)
	sourceschema = new(string)
	targetschema = new(string)
	ruleid       = 1
	deftask      = &Task{}
)

func main() {
	source := flag.String("source", "", "ARN to the Source Endpoint")
	target = flag.String("target", "", "ARN to the Target Endpoint")
	replication = flag.String("replication", "", "ARN to the Replication Instance")
	sourceschema = flag.String("sourceschema", "", "Schema name on Source Endpoint")
	targetschema = flag.String("targetschema", "", "Schema name on Target Endpoint")
	flag.Parse()

	// Check that the minimum amount of settings is set
	if *source == "" || *target == "" || *replication == "" || *sourceschema == "" {
		fmt.Printf("usage: dmstasks -[replication/source/sourceschema/target] string (optinal -targetschema string)\n")
		fmt.Printf("Will generate task files in ./tasks directory based on tables in tables.txt with settings from defaults.json.\n\n")
		flag.PrintDefaults()
		fmt.Printf("\n")
		os.Exit(0)
	}

	// If target schema is not set, assume it's the same schema name as source
	if *targetschema == "" {
		*targetschema = *sourceschema
	}

	// Read the defaults file
	defs, err := ioutil.ReadFile(defaults)
	if err != nil {
		log.Fatal("Couldn't read file " + defaults)
	}

	// Unmarshal the default settings into the Task object
	err = json.Unmarshal(defs, deftask)

	// Generate the tasks
	tasks := createTasks()

	// Generate task JSON files
	genJSONFiles(tasks)
}

func genJSONFiles(t *Tasks) {
	for _, task := range *t.Tasks {
		// Convert object to JSON
		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Fatal("Couldn't convert Task to JSON")
		}

		// Create the tasks directory
		os.Mkdir("tasks", 0755)

		// Write the task to a JSON-file
		ioutil.WriteFile("tasks/"+task.ReplicationTaskIdentifier+".json", []byte(taskJSON), 0644)
	}
}
