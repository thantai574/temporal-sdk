package temporal_sdk

import (
	"context"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type Context struct {
	context.Context
	TraceID   string
	WfID      string
	RunningID string
}

func ContextActivityNotRetry(parent_ctx workflow.Context) workflow.Context {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	return workflow.WithActivityOptions(parent_ctx, ao)
}
