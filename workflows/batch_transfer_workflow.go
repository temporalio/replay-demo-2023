package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func BatchTransferWorkflow(ctx workflow.Context) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
	})
	ctx = workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{StartToCloseTimeout: time.Second})

	var batchTransfers []TransferRequest
	var a *TransferActivity
	err := workflow.ExecuteActivity(ctx, a.GetBatchTransferRequest).Get(ctx, &batchTransfers)
	if err != nil {
		return err
	}

	for _, req := range batchTransfers {
		//var amount float64
		//err := workflow.ExecuteLocalActivity(ctx, a.GetPaymentAmount, req).Get(ctx, &amount)
		//if err != nil {
		//	return err
		//}
		//req.Amount = amount

		err = workflow.ExecuteActivity(ctx, a.Transfer, req).Get(ctx, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
