package main

const (
	settingsFile         = "settings.json"
	tablesFile           = "tables.txt"
	tasksFile            = "tasks.json"
	start                = "start-replication"
	resume               = "resume-processing"
	reload               = "reload-target"
	doesntExists         = "does not exist"
	notFound             = "not found"
	cantBeStarted        = "Task cannot be started, invalid state"
	currentlyNotRunning  = "is not currently running"
	allreadyBeingStopped = "is already being stopped"
	isRunning            = "is running"
	stopBeforeDeleting   = "before trying to delete it"
	waitForStop          = "has stopped before trying to delete it"
	beingDeleted         = "is already being deleted"
	alreadyExists        = "A task with this name already exists"
)

var (
	ruleid   = 1
	settings = &Task{}
	region   = "eu-west-1"
)
