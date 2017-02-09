package main

// Tasks is a slice of all the different Tasks
type Tasks struct {
	Tasks *[]Task
}

// Task contains all the parameters for the Task
type Task struct {
	ReplicationTaskIdentifier string                  `json:"ReplicationTaskIdentifier"`
	SourceEndpointArn         string                  `json:"SourceEndpointArn"`
	TargetEndpointArn         string                  `json:"TargetEndpointArn"`
	ReplicationInstanceArn    string                  `json:"ReplicationInstanceArn"`
	SourceSchema              string                  `json:"SourceSchema"`
	TargetSchema              string                  `json:"TargetSchema"`
	MigrationType             string                  `json:"MigrationType"`
	Mappings                  Mappings                `json:"Mappings"`
	ReplicationTaskSettings   ReplicationTaskSettings `json:"ReplicationTaskSettings"`
}

// ReplicationTaskSettings contains all the config parameters for the ReplicationTaskSettings
type ReplicationTaskSettings struct {
	TargetMetadata                    TargetMetadata                    `json:"TargetMetadata,omitempty"`
	FullLoadSettings                  FullLoadSettings                  `json:"FullLoadSettings,omitempty"`
	Logging                           Logging                           `json:"Logging,omitempty"`
	ControlTablesSettings             ControlTablesSettings             `json:"ControlTablesSettings,omitempty"`
	StreamBufferSettings              StreamBufferSettings              `json:"StreamBufferSettings,omitempty"`
	ChangeProcessingDdlHandlingPolicy ChangeProcessingDdlHandlingPolicy `json:"ChangeProcessingDdlHandlingPolicy,omitempty"`
	ErrorBehavior                     ErrorBehavior                     `json:"ErrorBehavior,omitempty"`
	ChangeProcessingTuning            ChangeProcessingTuning            `json:"ChangeProcessingTuning,omitempty"`
}

// TargetMetadata contains all the parameters for the Target
type TargetMetadata struct {
	TargetSchema        string `json:"TargetSchema,omitempty"`
	SupportLobs         bool   `json:"SupportLobs,omitempty"`
	FullLobMode         bool   `json:"FullLobMode,omitempty"`
	LobChunkSize        int    `json:"LobChunkSize,omitempty"`
	LimitedSizeLobMode  bool   `json:"LimitedSizeLobMode,omitempty"`
	LobMaxSize          int    `json:"LobMaxSize,omitempty"`
	LoadMaxFileSize     int    `json:"LoadMaxFileSize,omitempty"`
	ParallelLoadThreads int    `json:"ParallelLoadThreads,omitempty"`
	BatchApplyEnabled   bool   `json:"BatchApplyEnabled,omitempty"`
}

// FullLoadSettings contains all the parameters for the FullLoadSettings
type FullLoadSettings struct {
	FullLoadEnabled                 bool   `json:"FullLoadEnabled,omitempty"`
	ApplyChangesEnabled             bool   `json:"ApplyChangesEnabled,omitempty"`
	TargetTablePrepMode             string `json:"TargetTablePrepMode,omitempty"`
	CreatePkAfterFullLoad           bool   `json:"CreatePkAfterFullLoad,omitempty"`
	StopTaskCachedChangesApplied    bool   `json:"StopTaskCachedChangesApplied,omitempty"`
	StopTaskCachedChangesNotApplied bool   `json:"StopTaskCachedChangesNotApplied,omitempty"`
	ResumeEnabled                   bool   `json:"ResumeEnabled,omitempty"`
	ResumeMinTableSize              int    `json:"ResumeMinTableSize,omitempty"`
	ResumeOnlyClusteredPKTables     bool   `json:"ResumeOnlyClusteredPKTables,omitempty"`
	MaxFullLoadSubTasks             int    `json:"MaxFullLoadSubTasks,omitempty"`
	TransactionConsistencyTimeout   int    `json:"TransactionConsistencyTimeout,omitempty"`
	CommitRate                      int    `json:"CommitRate,omitempty"`
}

// Logging contains all the log entries
type Logging struct {
	EnableLogging       bool            `json:"EnableLogging,omitempty"`
	LogComponents       []LogComponents `json:"LogComponents,omitempty"`
	CloudWatchLogGroup  string          `json:"CloudWatchLogGroup,omitempty"`
	CloudWatchLogStream string          `json:"CloudWatchLogStream,omitempty"`
}

// LogComponents contains the ID and Severity level
type LogComponents struct {
	ID       string `json:"Id,omitempty"`
	Severity string `json:"Severity,omitempty"`
}

// ControlTablesSettings contains all the parameters for the ControlTablesSettings
type ControlTablesSettings struct {
	ControlSchema               string `json:"ControlSchema,omitempty"`
	HistoryTimeslotInMinutes    int    `json:"HistoryTimeslotInMinutes,omitempty"`
	HistoryTableEnabled         bool   `json:"HistoryTableEnabled,omitempty"`
	SuspendedTablesTableEnabled bool   `json:"SuspendedTablesTableEnabled,omitempty"`
	StatusTableEnabled          bool   `json:"StatusTableEnabled,omitempty"`
}

// StreamBufferSettings contains all the parameters for the StreamBufferSettings
type StreamBufferSettings struct {
	StreamBufferCount        int `json:"StreamBufferCount,omitempty"`
	StreamBufferSizeInMB     int `json:"StreamBufferSizeInMB,omitempty"`
	CtrlStreamBufferSizeInMB int `json:"CtrlStreamBufferSizeInMB,omitempty"`
}

// ChangeProcessingDdlHandlingPolicy contains all the parameters for ChangeProcessingDdlHandlingPolicy
type ChangeProcessingDdlHandlingPolicy struct {
	HandleSourceTableDropped   bool `json:"HandleSourceTableDropped,omitempty"`
	HandleSourceTableTruncated bool `json:"HandleSourceTableTruncated,omitempty"`
	HandleSourceTableAltered   bool `json:"HandleSourceTableAltered,omitempty"`
}

// ErrorBehavior contains all the parameters for the ErrorBehavior
type ErrorBehavior struct {
	DataErrorPolicy               string `json:"DataErrorPolicy,omitempty"`
	DataTruncationErrorPolicy     string `json:"DataTruncationErrorPolicy,omitempty"`
	DataErrorEscalationPolicy     string `json:"DataErrorEscalationPolicy,omitempty"`
	DataErrorEscalationCount      int    `json:"DataErrorEscalationCount,omitempty"`
	TableErrorPolicy              string `json:"TableErrorPolicy,omitempty"`
	TableErrorEscalationPolicy    string `json:"TableErrorEscalationPolicy,omitempty"`
	TableErrorEscalationCount     int    `json:"TableErrorEscalationCount,omitempty"`
	RecoverableErrorCount         int    `json:"RecoverableErrorCount,omitempty"`
	RecoverableErrorInterval      int    `json:"RecoverableErrorInterval,omitempty"`
	RecoverableErrorThrottling    bool   `json:"RecoverableErrorThrottling,omitempty"`
	RecoverableErrorThrottlingMax int    `json:"RecoverableErrorThrottlingMax,omitempty"`
	ApplyErrorDeletePolicy        string `json:"ApplyErrorDeletePolicy,omitempty"`
	ApplyErrorInsertPolicy        string `json:"ApplyErrorInsertPolicy,omitempty"`
	ApplyErrorUpdatePolicy        string `json:"ApplyErrorUpdatePolicy,omitempty"`
	ApplyErrorEscalationPolicy    string `json:"ApplyErrorEscalationPolicy,omitempty"`
	ApplyErrorEscalationCount     int    `json:"ApplyErrorEscalationCount,omitempty"`
	FullLoadIgnoreConflicts       bool   `json:"FullLoadIgnoreConflicts,omitempty"`
}

// ChangeProcessingTuning contains all the parameters for the ChangeProcessingTuning
type ChangeProcessingTuning struct {
	BatchApplyPreserveTransaction bool `json:"BatchApplyPreserveTransaction,omitempty"`
	BatchApplyTimeoutMin          int  `json:"BatchApplyTimeoutMin,omitempty"`
	BatchApplyTimeoutMax          int  `json:"BatchApplyTimeoutMax,omitempty"`
	BatchApplyMemoryLimit         int  `json:"BatchApplyMemoryLimit,omitempty"`
	BatchSplitSize                int  `json:"BatchSplitSize,omitempty"`
	MinTransactionSize            int  `json:"MinTransactionSize,omitempty"`
	CommitTimeout                 int  `json:"CommitTimeout,omitempty"`
	MemoryLimitTotal              int  `json:"MemoryLimitTotal,omitempty"`
	MemoryKeepTime                int  `json:"MemoryKeepTime,omitempty"`
	StatementCacheSize            int  `json:"StatementCacheSize,omitempty"`
}

// Mappings contains all the schema mappings
type Mappings struct {
	TableMappings []Rules `json:"rules"`
}

// Rules contains the different mapping rules
type Rules struct {
	RuleType      string        `json:"rule-type,omitempty"`
	RuleID        string        `json:"rule-id,omitempty"`
	RuleName      string        `json:"rule-name,omitempty"`
	RuleTarget    string        `json:"rule-target,omitempty"`
	ObjectLocator ObjectLocator `json:"object-locator,omitempty"`
	RuleAction    string        `json:"rule-action,omitempty"`
	Value         string        `json:"value,omitempty"`
}

// ObjectLocator contains all the ObjectLocator parameters
type ObjectLocator struct {
	SchemaName string `json:"schema-name,omitempty"`
	TableName  string `json:"table-name,omitempty"`
}

// Reply contains the answer from the AWS Create task function
type Reply struct {
	Reply ReplicationTask `json:"ReplicationTask"`
}

// ReplicationTask contains the ReplicationTask data from the AWS Reply
type ReplicationTask struct {
	ReplicationTaskIdentifier string `json:"ReplicationTaskIdentifier"`
	ReplicationTaskArn        string `json:"ReplicationTaskArn"`
}
