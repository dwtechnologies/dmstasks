package main

const (
	settingsFile = "settings.json"
	tablesFile   = "tables.txt"
	tasksFile    = "tasks.json"
	start        = "start-replication"
	resume       = "resume-processing"
	reload       = "reload-target"
)

var (
	ruleid   = 1
	settings = &Task{}
)
