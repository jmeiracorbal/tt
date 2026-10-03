package api

import "time"

type ViewResponse struct {
	View View `json:"view"`
}

type View struct {
	Resources []Resource `json:"resources"`
	LogList   LogList    `json:"logList"`
}

type Resource struct {
	Name              string  `json:"name"`
	RuntimeStatus     string  `json:"runtimeStatus"`
	UpdateStatus      string  `json:"updateStatus"`
	CurrentBuild      Build   `json:"currentBuild"`
	BuildHistory      []Build `json:"buildHistory"`
	HasPendingChanges bool    `json:"hasPendingChanges"`
}

type Build struct {
	StartTime  time.Time `json:"startTime"`
	FinishTime time.Time `json:"finishTime"`
	Error      string    `json:"error"`
	SpanID     string    `json:"spanId"`
}

type LogList struct {
	FromCheckpoint int64        `json:"fromCheckpoint"`
	ToCheckpoint   int64        `json:"toCheckpoint"`
	Segments       []LogSegment `json:"segments"`
}

type LogSegment struct {
	SpanID string            `json:"spanId"`
	Time   string            `json:"time"`
	Text   string            `json:"text"`
	Level  string            `json:"level"`
	Fields map[string]string `json:"fields,omitempty"`
}
