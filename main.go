package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Get the action flag
	action := flag.String("action", "", "What action to excute --action [create/start/stop/resume/remove]")
	flag.Parse()

	switch *action {
	case "create":
		// Create tasks
		createTasks()

	case "start":
		// Start the tasks stored in tasks.json
		startTasks(start)

	case "resume":
		// Resuma stopped tasks
		startTasks(resume)

	case "reload":
		// Reload target
		startTasks(reload)

	case "stop":
		// Stop the tasks stored in tasks.json
		stopTasks()

	case "delete":
		// Stop the tasks stored in tasks.json
		deleteTasks()

	case "status":
		// Show status of the tasks in tasks.json

	default:
		fmt.Println("Usage: ./dmstasks --action [create/start/stop/resume/remove]")
		fmt.Println("Needs AWS ENV vars to be set or in .aws/credentials file.")
		os.Exit(0)
	}
}
