package temporal_sdk

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"testing"
)

func TestSdk_Worker(t *testing.T) {
	sdk, _ := NewClient(Opts{
		ClientOptions: client.Options{
			HostPort: "localhost:7233",
		},
		TaskQueue: "test1",
		Name:      "test2",
	})

	err := sdk.RegisterWorker(SampleParallelWorkflow, SampleActivity, SampleActivityE).Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
	defer sdk.CloseConnection()
	select {}

}
