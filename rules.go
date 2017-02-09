package main

import "strconv"

// // genRule take r *string (row) and creates a include Rule out of it
// func genRule(r *string) *Rules {
// 	return &Rules{
// 		Type:         "Include",
// 		SourceSchema: settings.SourceSchema,
// 		SourceTable:  *r,
// 	}
// }

// // defaultExclude creates the default exclude all tables in schema and returns *Rules
// func defaultExclude() *Rules {
// 	return &Rules{
// 		Type:         "Exclude",
// 		SourceSchema: settings.SourceSchema,
// 		SourceTable:  "%",
// 	}
// }

// // defaultRename creates rule for renaming schema from source name to target name and returns *Rules
// func defaultRename() *Rules {
// 	return &Rules{
// 		Type:         "Transformation",
// 		SourceSchema: settings.SourceSchema,
// 		TargetSchema: settings.SourceSchema,
// 	}
// }

// genRule take r *string (row) and creates a include Rule out of it
func genRule(r *string) *Rules {
	rulestr := strconv.Itoa(ruleid)
	ruleid++
	return &Rules{
		RuleType: "selection",
		RuleID:   rulestr,
		RuleName: rulestr,
		ObjectLocator: ObjectLocator{
			SchemaName: settings.SourceSchema,
			TableName:  *r,
		},
		RuleAction: "include",
	}
}

// defaultExclude creates the default exclude all tables in schema and returns *Rules
func defaultExclude() *Rules {
	rulestr := strconv.Itoa(ruleid)
	ruleid++
	return &Rules{
		RuleType: "selection",
		RuleID:   rulestr,
		RuleName: rulestr,
		ObjectLocator: ObjectLocator{
			SchemaName: settings.SourceSchema,
			TableName:  "%",
		},
		RuleAction: "exclude",
	}
}

// defaultRename creates rule for renaming schema from source name to target name and returns *Rules
func defaultRename() *Rules {
	rulestr := strconv.Itoa(ruleid)
	ruleid++
	return &Rules{
		RuleType:   "transformation",
		RuleID:     rulestr,
		RuleName:   rulestr,
		RuleTarget: "schema",
		ObjectLocator: ObjectLocator{
			SchemaName: settings.SourceSchema,
		},
		RuleAction: "rename",
		Value:      settings.TargetSchema,
	}
}
