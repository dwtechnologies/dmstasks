# dmstasks

Quick and dirty tasks creator for Amazon AWS DMS (Data Migration Service). 

This project was created because of how DMS works (it needs small task with around 8-12 tablers each ...)  
Basically we tried to migrate a large database (30GB) with around 570 tables with ongoing replication.  
It Would always go *well* until around 85% when we would start seeing errors.

After talking to AWS Support it seems that the only way to get DMS to work is to split the task to many smaller task with around 10-20 tables maxiumum per task.  
It should be split based on transactions/table names.

So as the lazy guy I am I created this quick and dirty program for generating JSON task files (that can be imported using aws cli).  

What you need:  
`Golang` - Build the program by issuing `go build` inside the dmstasks directory.  
`tables.txt` - Should contain a list of all the tables you will want to migrate split up into different task seperated by [task-name]  
`defaults.json` - Should contain a JSON file with the default parameters you will want the task to have. (Tip, create a task in the Console then describe it in aws cli)

Please have a look in the supploed `tables_example.txt` and `defaults_example.json` ...

Then simply run:  
`./dmstask -source "arn-to-source" -target "arn-to-target" -replication "arn-to-replication-instance" -sourceschema "schema-name-on-source" -targetschema "schema-name-on-target"`

(-targetschema is optional, if left empty it will assume that the schema name is the same as on source)


You will now end up with a lot of JSON files in the `tasks` directory, ready for your pleasure.  
Have a look in the supplied `createtasks.sh` and `runtasks.sh` for automated task creation and activation.