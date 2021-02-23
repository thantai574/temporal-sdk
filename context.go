package temporal_sdk

import "context"

type Context struct {
	context.Context
	TraceID   string
	WfID      string
	RunningID string
}
