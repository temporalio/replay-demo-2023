package schedule

import (
	"context"
	"log"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"replay-demo/workflows"
)

func CreateSchedule(c client.ScheduleClient, scheduleID, workflowID string, spec client.ScheduleSpec, triggerNow bool) {
	action := &client.ScheduleWorkflowAction{
		ID:        workflowID,
		Workflow:  workflows.BatchTransferWorkflow,
		TaskQueue: "demo-tq",

		// Set short timeout so we don't accumulate too many concurrent running workflows from 5s schedule if demo worker
		// is down while we allow all overlap runs.
		WorkflowRunTimeout: 30 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := c.Create(ctx, client.ScheduleOptions{
		ID:     scheduleID,
		Spec:   spec,
		Action: action,

		// This is for versioning demo purpose to show concurrent running workflows of 2 versions
		Overlap: enums.SCHEDULE_OVERLAP_POLICY_ALLOW_ALL,

		TriggerImmediately: triggerNow,
	})

	if err == temporal.ErrScheduleAlreadyRunning {
		log.Printf("Schedule %v already registered.", scheduleID)
	} else if err != nil {
		log.Printf("Failed to create schedule: %v", err)
	} else {
		log.Printf("Schedule %v created.", scheduleID)
	}
	return
}

func MakeSpecEvery5Seconds() client.ScheduleSpec {
	return client.ScheduleSpec{
		// Run the schedule every 5s
		Intervals: []client.ScheduleIntervalSpec{
			{
				Every: 5 * time.Second,
			},
		},
	}
}

func MakeSpecBusinessHoursHourly() client.ScheduleSpec {
	// Run hourly from 9am to 5pm, Monday to Friday on pacific time zone.
	// Equivalent to CRON: TZ=US/Pacific 0 9-17 * * 1-5
	spec := client.ScheduleSpec{
		Calendars: []client.ScheduleCalendarSpec{
			{
				Hour:      []client.ScheduleRange{{Start: 9, End: 17}}, // 9am ~ 5pm
				DayOfWeek: []client.ScheduleRange{{Start: 1, End: 5}},  // Monday ~ Friday
			},
		},
		TimeZoneName: "US/Pacific", // will automatically adjust between PST/PDT
		// to spread load for large number of schedules
		Jitter: 5 * time.Minute,
	}
	return spec
}

func MakeSpecCustomSchedule() client.ScheduleSpec {
	// Run every Thursday at 2pm, only in January, August, and December
	spec := client.ScheduleSpec{
		Calendars: []client.ScheduleCalendarSpec{
			{
				// 2pm
				Hour: []client.ScheduleRange{{Start: 14}},
				// Thursday
				DayOfWeek: []client.ScheduleRange{{Start: 4}},
				// September/December
				Month: []client.ScheduleRange{{Start: 9}, {Start: 12}},
			},
		},
		TimeZoneName: "US/Pacific", // will automatically adjust between PST/PDT
	}
	return spec
}
