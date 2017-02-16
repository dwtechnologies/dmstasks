# dmstasks

Quick and dirty tasks creator for Amazon AWS DMS (Data Migration Service). 

This project was created because of how DMS works (it needs small task with around 1-20 tablers each ...)  
Basically we tried to migrate a large database (30GB) with around 570 tables with ongoing replication.  
It Would always go *well* until around 85% when we would start seeing errors.

After talking to AWS Support it seems that the only way to get DMS to work is to split the task to many smaller task with around 1-20 tables maxiumum per task.  
It should be split based on transactions/table names.

So as the lazy guy I am I created this quick and dirty program for generating, creating, starting, stoping, deleting and showing status.  

What you need:  
`Golang` - Build the program by issuing `go build` inside the dmstasks directory.  
`tables.txt` - Should contain a list of all the tables you will want to migrate split up into different task seperated by [task-name] (have a look in tables_example.txt)  
`settings.json` - Should contain a JSON file with the settings and parameters you will want the task to have. (Tip, create a task in the Console then describe it in aws cli) (see settings_example.json)

If you don't specify Target Schema in the settings.json the program will assume that you want the same Schema name on target.


Results from the commands below will be stored in a file called tasks.json. So please don't remove this.  
If you happen to remove a task in the console and not by the commands below, please find it in the JSON-file and remove it manually to skip seeing errors about that task.  

The program takes two arguments `-action` and `-region`. Both are mandatory.


### Then simply run: 

Create tasks  
`./dmstasks -action create -region eu-west-1`

Start tasks  
`./dmstasks -action start -region eu-west-1`

Stop tasks  
`./dmstasks -action stop -region eu-west-1`

Resume tasks  
`./dmstasks -action resume -region eu-west-1`

Delete tasks  
`./dmstasks -action delete -region eu-west-1`

Status (max 100 tasks atm :))  
`./dmstasks -action status -region eu-west-1`


----------

This is a quite quick and dirty program. Please create Pull Requests for new features and optimizations to the code.