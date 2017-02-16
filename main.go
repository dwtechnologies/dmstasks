package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// Get the action flag
	action := flag.String("action", "", "What action to excute [create/start/stop/resume/delete/status]")
	flagRegion := flag.String("region", "", "What AWS Region to execute commands in, ex. eu-west-1")
	flag.Parse()

	// Set region
	region = *flagRegion
	if region == "" {
		printUsage()
	}

	switch *action {
	case "create":
		createTasks()

	case "start":
		startTasks(start)

	case "resume":
		startTasks(resume)

	case "reload":
		startTasks(reload)

	case "stop":
		stopTasks()

	case "delete":
		deleteTasks()

	case "status":
		statusTasks()

	default:
		printUsage()
	}
}

// printUsage prints the programs usage and exits.
func printUsage() {
	fmt.Println("Usage: ./dmstasks -action [create/start/stop/resume/delete/status] -region [aws-region]")
	fmt.Println("Both -action and -region are mandatory. Also nees AWS credentials to be set in ENV or in .aws/credentials file.")
	os.Exit(1)
}

// readConfig reads the config file and sets the config variable
func readConfig() {
	// Read the defaults file
	readSettings, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		log.Fatal("Couldn't read file " + settingsFile)
	}

	// Unmarshal the default settings into the Task object
	err = json.Unmarshal(readSettings, settings)
	if err != nil {
		log.Fatal("Couldn't unmarshal settings file " + settingsFile)
	}
}
