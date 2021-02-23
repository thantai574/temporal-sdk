package temporal_sdk

import (
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"testing"
	"time"
)

func TestSdk_TemporalWFCallSync(t *testing.T) {
	type fields struct {
		Client  client.Client
		Context Context
		Opts    Opts
	}
	type args struct {
		workFlow interface{}
		response interface{}
		agrs     []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{args: args{
			workFlow: SampleParallelWorkflow,
			response: nil,
			agrs:     nil,
		}},
	}
	sdk, _ := NewClient(Opts{
		ClientOptions: client.Options{
			HostPort: "localhost:7233",
		},
		TaskQueue: "test1",
		Name:      "test2",
	})
	res := []string{}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if err := sdk.TemporalWFCallSync(tt.args.workFlow, &res); (err != nil) != tt.wantErr {
				t.Errorf("TemporalWFCallSync() error = %v, wantErr %v", err, tt.wantErr)
				fmt.Println(err.Error())
			}
		})
	}
	fmt.Println(res)
	defer sdk.CloseConnection()

}

// SampleParallelWorkflow workflow definition
func SampleParallelWorkflow(ctx workflow.Context) ([]string, error) {
	logger := workflow.GetLogger(ctx)
	defer logger.Info("Workflow completed.")

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	result := ""
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, SampleActivity, "branch1.1").Get(ctx, &result)

	if err != nil {
		return []string{}, err
	}

	result2 := ""
	err = workflow.ExecuteActivity(ctx, SampleActivity, "branch1.2").Get(ctx, &result2)

	if err != nil {
		return []string{}, err
	}
	result3 := ""
	err = workflow.ExecuteActivity(ctx, SampleActivity, "branch1.4").Get(ctx, &result3)

	if err != nil {
		return []string{}, err
	}

	return []string{result, result2, result3}, nil
}

func SampleActivity(input string) (string, error) {
	name := "sampleActivity"
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + input, nil
}

func SampleActivityE(input string) (string, error) {
	name := "sampleActivity"
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + input, nil
}
