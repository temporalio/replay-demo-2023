package workflows

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

const totalAccountNumber = 10

type TransferActivity struct {
	TemporalClient client.Client
}

func (a *TransferActivity) Deposit(ctx context.Context, accountID string, amount float64) error {
	if strings.Contains(strings.ToLower(accountID), "piggy") {
		return temporal.NewNonRetryableApplicationError("deposit failed: piggy bank account is frozen", "account-frozen", nil)
	}
	// make bank API call to deposit the amount.
	return nil
}

func (a *TransferActivity) Withdraw(ctx context.Context, accountID string, amount float64) error {
	if strings.Contains(strings.ToLower(accountID), "piggy") {
		return temporal.NewNonRetryableApplicationError("withdraw failed: piggy bank account is frozen", "account-frozen", nil)
	}
	// make bank API call to withdraw the amount.
	return nil
}

func (a *TransferActivity) RevertDeposit(ctx context.Context, accountID string, amount float64) error {
	// make bank API call to revert deposit.
	return nil
}

func (a *TransferActivity) RevertWithdraw(ctx context.Context, accountID string, amount float64) error {
	// make bank API call to revert withdraw
	return nil
}

func (a *TransferActivity) GetBatchTransferRequest(ctx context.Context) ([]TransferRequest, error) {
	var requests []TransferRequest
	for i := 1; i <= totalAccountNumber; i++ {
		requests = append(requests, TransferRequest{
			FromAccount: fmt.Sprintf("from-account-%v", 1+rand.Intn(50)),
			ToAccount:   fmt.Sprintf("to-account-%v", 51+rand.Intn(50)),
			Amount:      10, // hard code amount
		})
	}
	return requests, nil
}

func (a *TransferActivity) Transfer(ctx context.Context, req TransferRequest) (string, error) {
	workflowID := fmt.Sprintf("%s_%s_%s_$%.2f", activity.GetInfo(ctx).WorkflowExecution.ID, req.FromAccount, req.ToAccount, req.Amount)
	_, err := a.TemporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "demo-tq",
	}, TransferWorkflow)
	if err != nil {
		return "", err
	}
	handle, err := a.TemporalClient.UpdateWorkflowWithOptions(ctx, &client.UpdateWorkflowWithOptionsRequest{
		UpdateID:   "batch-transfer-update",
		WorkflowID: workflowID,
		UpdateName: TransferUpdateName,
		Args:       []interface{}{req},
	})
	if err != nil {
		return "", err
	}

	err = handle.Get(ctx, nil)
	// Sleep 1s to slow down
	time.Sleep(time.Second)
	return workflowID, err
}

func (a *TransferActivity) GetPaymentAmount(req TransferRequest) (float64, error) {
	// return random number [1, 100)
	return 1 + rand.Float64()*99, nil
}
