package main

import (
	"context"
	"log"
	"os"

	demo "replay-demo/client"
	"replay-demo/schedule"
	"replay-demo/workflows"

	"go.temporal.io/sdk/client"
)

func main() {
	mode := "schedule"
	if len(os.Args) == 2 {
		mode = os.Args[1]
	}
	switch mode {
	case "schedule":
		createSchedules()
	case "update":
		runDemoUpdate()
	}
}

func runDemoUpdate() {
	c := demo.NewClient()
	defer c.Close()
	_, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        "transfer-1",
		TaskQueue: "demo-tq",
	}, workflows.TransferWorkflow)

	if err != nil {
		log.Fatalf("error start wf: %v", err)
	}

	updateHandle, err := c.UpdateWorkflow(context.Background(), "transfer-1", "", "transfer", workflows.TransferRequest{
		FromAccount: "from-account-id",
		ToAccount:   "to-account-id-piggy-bank",
		Amount:      10,
	})
	if err != nil {
		log.Fatalf("error update wf: %v", err)
	}
	err = updateHandle.Get(context.Background(), nil)
	log.Printf("Update failed: %v\n", err)

	_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        "transfer-2",
		TaskQueue: "demo-tq",
	}, workflows.TransferWorkflow)

	if err != nil {
		log.Fatalf("error start wf: %v", err)
	}
	_, err = c.UpdateWorkflow(context.Background(), "transfer-2", "", "transfer", workflows.TransferRequest{
		FromAccount: "from-account-id",
		ToAccount:   "to-account-id",
		Amount:      10,
	})
}

func createSchedules() {
	c := demo.NewClient()
	defer c.Close()
	sClient := c.ScheduleClient()
	schedule.CreateSchedule(sClient, "schedule_every_5s", "payment_every_5s", schedule.MakeSpecEvery5Seconds(), false)
	schedule.CreateSchedule(sClient, "schedule_business_hourly", "payment_hourly", schedule.MakeSpecBusinessHoursHourly(), true)
	schedule.CreateSchedule(sClient, "schedule_custom", "payment_custom_schedule", schedule.MakeSpecCustomSchedule(), false)
}
