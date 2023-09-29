package workflows

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/workflow"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	SetFromAccountUpdateName = "set-from-account"
	SetToAccountUpdateName   = "set-to-account"
	TransferAmountUpdateName = "transfer-amount"

	TransferUpdateName = "transfer"

	DailyAmountLimit = 100000.0
)

type TransferRequest struct {
	FromAccount string
	ToAccount   string
	Amount      float64
}

func TransferWorkflow(ctx workflow.Context) error {
	log := workflow.GetLogger(ctx)

	var a *TransferActivity
	var pendingCompensations []func(workflow.Context) error
	var transferErr error
	var transferAttempted, transferDone bool
	transferHandlerFunc := func(ctx workflow.Context, req TransferRequest) error {
		transferAttempted = true
		defer func() { transferDone = true }()

		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Second * 10,
		})

		// add compensation for withdraw in case it fails.
		pendingCompensations = append(pendingCompensations, func(ctx workflow.Context) error {
			return workflow.ExecuteActivity(ctx, a.RevertWithdraw, req.FromAccount, req.Amount).Get(ctx, nil)
		})
		transferErr = workflow.ExecuteActivity(ctx, a.Withdraw, req.FromAccount, req.Amount).Get(ctx, nil)

		if transferErr != nil {
			return transferErr
		}

		// Deposit (adding compensation just in case it fails)
		pendingCompensations = append(pendingCompensations, func(ctx workflow.Context) error {
			return workflow.ExecuteActivity(ctx, a.RevertDeposit, req.ToAccount, req.Amount).Get(ctx, nil)
		})
		transferErr = workflow.ExecuteActivity(ctx, a.Deposit, req.ToAccount, req.Amount).Get(ctx, nil)
		return transferErr
	}
	transferValidator := func(ctx workflow.Context, fromAccount, toAccount string, amount float64) error {
		if transferAttempted {
			log.Debug("Rejecting transfer request", "transferAttempted", transferAttempted)
			return fmt.Errorf("transfer already attempted")
		}
		if fromAccount == "" {
			log.Debug("Rejecting transfer request", "from-account", fromAccount)
			return fmt.Errorf("from account is not set (%v)", fromAccount)
		}
		if toAccount == "" {
			log.Debug("Rejecting transfer request", "to-account", toAccount)
			return fmt.Errorf("to account is not set (%v)", toAccount)
		}
		if amount <= 0 {
			log.Debug("Rejecting transfer request", "transfer-amount", amount)
			return fmt.Errorf("invalid transfer amount (%v)", amount)
		}
		if amount > DailyAmountLimit {
			log.Debug("Rejecting transfer request", "transfer-amount", amount)
			return fmt.Errorf("transfer amount ($%s) exceeds daily limit ($%s)", formatMoney(amount), formatMoney(DailyAmountLimit))
		}

		return nil
	}

	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		TransferUpdateName,
		transferHandlerFunc,
		workflow.UpdateHandlerOptions{Validator: func(ctx workflow.Context, req TransferRequest) error {
			return transferValidator(ctx, req.FromAccount, req.ToAccount, req.Amount)
		}},
	); err != nil {
		return err
	}

	// below 3 updates are for page flow
	var fromAccountID, toAccountID string
	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		SetFromAccountUpdateName,
		func(ctx workflow.Context, accountID string) error {
			fromAccountID = accountID
			return nil
		},
		workflow.UpdateHandlerOptions{Validator: func(ctx workflow.Context, accountID string) error {
			if strings.Contains(strings.ToLower(accountID), "crypto") {
				log.Debug("Rejecting from account", "from-account", accountID)
				return fmt.Errorf("crypto account is not supported (%v)", accountID)
			}
			return nil
		}},
	); err != nil {
		return err
	}
	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		SetToAccountUpdateName,
		func(ctx workflow.Context, accountID string) error {
			toAccountID = accountID
			return nil
		},
		workflow.UpdateHandlerOptions{Validator: func(ctx workflow.Context, accountID string) error {
			if strings.Contains(strings.ToLower(accountID), "crypto") {
				log.Debug("Rejecting from account", "to-account", accountID)
				return fmt.Errorf("crypto account is not supported (%v)", accountID)
			}
			return nil
		}},
	); err != nil {
		return err
	}
	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		TransferAmountUpdateName,
		func(ctx workflow.Context, amount float64) error {
			return transferHandlerFunc(ctx, TransferRequest{
				FromAccount: fromAccountID,
				ToAccount:   toAccountID,
				Amount:      amount,
			})
		},
		workflow.UpdateHandlerOptions{Validator: func(ctx workflow.Context, amount float64) error {
			return transferValidator(ctx, fromAccountID, toAccountID, amount)
		}},
	); err != nil {
		return err
	}

	// block until transfer is done.
	workflow.Await(ctx, func() bool { return transferDone })

	if transferErr != nil {
		// execute saga compensations
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Second * 10,
		})
		var compensationErrs []error
		for i := len(pendingCompensations) - 1; i >= 0; i-- {
			compensationErrs = append(compensationErrs, pendingCompensations[i](ctx))
		}
		return errors.Join(compensationErrs...)
	}

	return nil
}

func formatMoney(amount float64) string {
	em := message.NewPrinter(language.English)
	return em.Sprintf("%.2f", amount)
}
