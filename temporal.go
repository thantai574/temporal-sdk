package temporal_sdk

import (
	"context"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Opts struct {
	ClientOptions client.Options
	TaskQueue     string
	Name          string
}

type Sdk struct {
	client.Client
	*Context
	Opts
}

type ISdk interface {
	TemporalWFCallSync(workFlow interface{}, response interface{}, agrs ...interface{}) (err error)
	TemporalWFCallASync(workFlow interface{}, agrs ...interface{}) (err error)
	RegisterWorker(wf interface{}, activities ...interface{}) worker.Worker
	CloseConnection()
}

func NewClient(opts Opts) (sdk ISdk, e error) {
	c, err := client.NewClient(opts.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Sdk{
		Client: c,
		Context: &Context{
			Context: context.Background(),
			TraceID: uuid.New(),
		},
		Opts: opts,
	}, err
}

func (sdk *Sdk) TemporalWFCallSync(workFlow interface{}, response interface{}, agrs ...interface{}) (err error) {
	options := client.StartWorkflowOptions{
		ID:        sdk.Opts.Name + uuid.New(),
		TaskQueue: sdk.Opts.TaskQueue,
	}

	wf, err := sdk.ExecuteWorkflow(sdk.Context, options, workFlow, agrs...)

	if err != nil {
		return
	}

	sdk.Context.RunningID = wf.GetRunID()
	sdk.Context.WfID = wf.GetID()
	err = wf.Get(sdk.Context, &response)

	return
}

func (sdk *Sdk) TemporalWFCallASync(workFlow interface{}, agrs ...interface{}) (err error) {
	options := client.StartWorkflowOptions{
		ID:        sdk.Opts.Name + uuid.New(),
		TaskQueue: sdk.Opts.TaskQueue,
	}

	wf, err := sdk.ExecuteWorkflow(sdk.Context, options, workFlow, agrs...)

	if err != nil {
		return
	}

	sdk.Context.RunningID = wf.GetRunID()
	sdk.Context.WfID = wf.GetID()
	return
}

func (sdk *Sdk) RegisterWorker(wf interface{}, activities ...interface{}) worker.Worker {
	w := worker.New(sdk.Client, sdk.TaskQueue, worker.Options{})

	w.RegisterWorkflow(wf)

	for _, v := range activities {
		w.RegisterActivity(v)
	}

	return w

}

func (sdk *Sdk) CloseConnection() {
	sdk.Close()
}
