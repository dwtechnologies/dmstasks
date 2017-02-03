package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func createTasks() *Tasks {
	// Create Tasks
	tasks := new(Tasks)
	tasks.Tasks = addTasks()

	return tasks
}

// addTasks takes r *int (ruleid number) and generates Mappings to tables based on tables file. Returns *[]Rules
func addTasks() *[]Task {
	tasks := new([]Task)

	// Read the tables from file
	tblsFile, err := ioutil.ReadFile(tables)
	if err != nil {
		log.Fatal("Couldn't read file " + tables)
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
	*task = *deftask

	// Split by name
	split := strings.SplitN(*r, "]", 2)
	if len(split) != 2 {
		return nil
	}

	// Set the name of the task
	name := split[0]
	task.ReplicationTaskIdentifier = *sourceschema + "-" + *targetschema + "-" + name

	// Set ARNs
	task.SourceEndpointArn = *source
	task.TargetEndpointArn = *target
	task.ReplicationInstanceArn = *replication

	// Set CloudWatch LogGroup
	if task.ReplicationTaskSettings.Logging.EnableLogging == true {
		task.ReplicationTaskSettings.Logging.CloudWatchLogGroup = "dms-tasks-" + task.ReplicationTaskIdentifier
		task.ReplicationTaskSettings.Logging.CloudWatchLogStream = task.ReplicationTaskSettings.Logging.CloudWatchLogGroup
	}

	// Add default exclude all
	exclude := defaultExclude()
	task.TableMappings = append(task.TableMappings, *exclude)

	// Get the mappings for the current job
	params := strings.Split(split[1], "\n")
	for _, param := range params {
		param := strings.TrimSpace(param)
		if param == "" {
			continue
		}
		rule := genRule(&param)
		task.TableMappings = append(task.TableMappings, *rule)
	}

	// Add rename schema mapping - As last ID
	if sourceschema != targetschema {
		rename := defaultRename()
		task.TableMappings = append(task.TableMappings, *rename)
	}

	return task
}
