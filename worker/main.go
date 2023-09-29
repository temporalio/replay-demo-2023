package main

import (
	"context"
	"log"
	"strings"

	"replay-demo/client"
	"replay-demo/workflows"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/worker"
)

// BuildID is used for versioning. Update this whenever there is non-backward compatible workflow logic change.
// Once deployed with new BuildID, we need to update TaskQueue to set this new buildID as default so new workflow
// can be route to worker wit this build ID.
// Example: `temporal task-queue update-build-ids add-new-default --task-queue demo-tq --build-id 2.0`
const BuildID = "1.0"

func main() {
	c := client.NewClient()
	defer c.Close()

	// Set current worker as default. This is for demo convenience so this worker will always be the default version.
	// WARNING: DO NOT DO THIS IN PROD. Should set the default BuildID as part of deployment flow.
	// Doing this in worker code in prod will cause issue because older version worker may restart and would cause
	// the older version to be set as default again.
	SetCurrentWorkerAsDefault()

	w := worker.New(c, "demo-tq", worker.Options{BuildID: BuildID, UseBuildIDForVersioning: true})
	w.RegisterWorkflow(workflows.TransferWorkflow)
	w.RegisterWorkflow(workflows.BatchTransferWorkflow)
	a := &workflows.TransferActivity{
		TemporalClient: c,
	}
	w.RegisterActivity(a)
	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func SetCurrentWorkerAsDefault() {
	c := client.NewClient()
	defer c.Close()
	request := &workflowservice.UpdateWorkerBuildIdCompatibilityRequest{
		Namespace: client.GetNamespace(),
		TaskQueue: "demo-tq",
		Operation: &workflowservice.UpdateWorkerBuildIdCompatibilityRequest_AddNewBuildIdInNewDefaultSet{
			AddNewBuildIdInNewDefaultSet: BuildID,
		},
	}

	_, err := c.WorkflowService().UpdateWorkerBuildIdCompatibility(context.Background(), request)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		request := &workflowservice.UpdateWorkerBuildIdCompatibilityRequest{
			Namespace: client.GetNamespace(),
			TaskQueue: "demo-tq",
			Operation: &workflowservice.UpdateWorkerBuildIdCompatibilityRequest_PromoteSetByBuildId{
				PromoteSetByBuildId: BuildID,
			},
		}
		_, err = c.WorkflowService().UpdateWorkerBuildIdCompatibility(context.Background(), request)
	}
	log.Printf("Set default buildID: %v, Err: %v\n", BuildID, err)
}
