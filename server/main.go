package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
	"go.temporal.io/sdk/client"
	demo "replay-demo/client"
	"replay-demo/schedule"
	"replay-demo/workflows"
)

func returnWorkflowIds(workflowID, runID string, w http.ResponseWriter) {
	resp := make(map[string]string)
	resp["workflowID"] = workflowID
	resp["runID"] = runID
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func returnError(err error, w http.ResponseWriter) {
	resp := make(map[string]string)
	resp["success"] = ""
	resp["error"] = err.Error()
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func returnSuccess(w http.ResponseWriter) {
	resp := make(map[string]string)
	resp["success"] = "Update successful"
	resp["error"] = ""
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

type TransferRequestWithIDs struct {
	FromAccount string
	ToAccount   string
	Amount      float64
	WorkflowID  string
	RunID       string
}

func main() {
	mux := http.NewServeMux()

	c := demo.NewClient()
	defer c.Close()

	mux.HandleFunc("/initiate", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now().Unix()

		we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
			ID:        "transfer-" + fmt.Sprint(t),
			TaskQueue: "demo-tq",
		}, workflows.TransferWorkflow)

		if err != nil {
			log.Printf("error start wf: %v", err)
			returnError(err, w)
			return
		}

		returnWorkflowIds(we.GetID(), we.GetRunID(), w)
	})

	handleFunc := func(w http.ResponseWriter, r *http.Request, updateName string) {
		decoder := json.NewDecoder(r.Body)
		var t TransferRequestWithIDs
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		if err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		var updateHandle client.WorkflowUpdateHandle
		switch updateName {
		case workflows.SetFromAccountUpdateName:
			updateHandle, err = c.UpdateWorkflow(context.Background(), t.WorkflowID, t.RunID, updateName, t.FromAccount)
		case workflows.SetToAccountUpdateName:
			updateHandle, err = c.UpdateWorkflow(context.Background(), t.WorkflowID, t.RunID, updateName, t.ToAccount)
		case workflows.TransferAmountUpdateName:
			updateHandle, err = c.UpdateWorkflow(context.Background(), t.WorkflowID, t.RunID, updateName, t.Amount)
		}

		if err != nil {
			log.Printf("error update workflow for %v: %v", updateName, err)
			returnError(err, w)
			return
		}

		err = updateHandle.Get(context.Background(), nil)
		if err != nil {
			log.Printf("error get update result for %v: %v", updateName, err)
			returnError(err, w)
			return
		}

		returnSuccess(w)
	}

	mux.HandleFunc("/from-account", func(w http.ResponseWriter, r *http.Request) {
		handleFunc(w, r, workflows.SetFromAccountUpdateName)
	})

	mux.HandleFunc("/to-account", func(w http.ResponseWriter, r *http.Request) {
		handleFunc(w, r, workflows.SetToAccountUpdateName)
	})

	mux.HandleFunc("/amount", func(w http.ResponseWriter, r *http.Request) {
		handleFunc(w, r, workflows.TransferAmountUpdateName)
	})

	mux.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		// Start a schedule of payment workflows
		sClient := c.ScheduleClient()
		schedule.CreateSchedule(sClient, "schedule_every_5s", "payment_every_5s", schedule.MakeSpecEvery5Seconds(), false)
		schedule.CreateSchedule(sClient, "schedule_business_hourly", "payment_hourly", schedule.MakeSpecBusinessHoursHourly(), true)
		schedule.CreateSchedule(sClient, "schedule_custom", "payment_custom_schedule", schedule.MakeSpecCustomSchedule(), false)
	})

	handler := cors.Default().Handler(mux)

	// Start the HTTP server on port 7654
	fmt.Println("Starting server on :7654")
	if err := http.ListenAndServe(":7654", handler); err != nil {
		panic(err)
	}
}
