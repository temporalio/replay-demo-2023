package workflows_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"replay-demo/workflows"
)

type updateCallback struct {
	accepted    bool
	rejectedErr error
	completeErr error
	result      interface{}
}

func (uc *updateCallback) Accept() {
	uc.accepted = true
}

func (uc *updateCallback) Reject(err error) {
	uc.rejectedErr = err
}

func (uc *updateCallback) Complete(success interface{}, err error) {
	uc.result = success
	uc.completeErr = err
}

func TestTransferWorkflow_Reject(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(workflows.TransferWorkflow)
	a := &workflows.TransferActivity{}
	env.RegisterActivity(a)

	cb1 := updateCallback{}
	cb2 := updateCallback{}

	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(workflows.TransferUpdateName, "transaction-id-1", &cb1, workflows.TransferRequest{
			FromAccount: "my-from-account",
			ToAccount:   "my-to-account",
			Amount:      -1, // invalid amount
		})
		env.UpdateWorkflow(workflows.TransferUpdateName, "transaction-id-2", &cb2, workflows.TransferRequest{
			FromAccount: "my-from-account",
			ToAccount:   "my-to-account",
			Amount:      1000000, // exceed daily limit
		})
	}, time.Nanosecond)

	// Run workflow
	env.ExecuteWorkflow(workflows.TransferWorkflow)

	require.False(t, cb1.accepted)
	require.Error(t, cb1.rejectedErr)
	require.Contains(t, cb1.rejectedErr.Error(), "invalid transfer amount")

	require.False(t, cb2.accepted)
	require.Error(t, cb2.rejectedErr)
	require.Contains(t, cb2.rejectedErr.Error(), "exceeds daily limit")

	// workflow eventually timeout
	err := env.GetWorkflowResult(nil)
	require.Error(t, err)
}

func TestTransferWorkflow_Succeed(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(workflows.TransferWorkflow)
	a := &workflows.TransferActivity{}
	env.RegisterActivity(a)

	cb1 := updateCallback{}

	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(workflows.TransferUpdateName, "transaction-id-1", &cb1, workflows.TransferRequest{
			FromAccount: "my-from-account",
			ToAccount:   "my-to-account",
			Amount:      10,
		})
	}, time.Nanosecond)

	// Run workflow
	env.ExecuteWorkflow(workflows.TransferWorkflow)

	require.True(t, cb1.accepted)
	require.NoError(t, cb1.completeErr)
	err := env.GetWorkflowResult(nil)
	require.NoError(t, err)
}

func TestTransferWorkflow_InvalidToAccount_Compensate(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(workflows.TransferWorkflow)
	a := &workflows.TransferActivity{}
	env.RegisterActivity(a)

	cb1 := updateCallback{}

	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(workflows.TransferUpdateName, "transaction-id-1", &cb1, workflows.TransferRequest{
			FromAccount: "my-from-account",
			ToAccount:   "my-to-account-piggy-bank",
			Amount:      10,
		})
	}, time.Nanosecond)

	// Run workflow
	env.ExecuteWorkflow(workflows.TransferWorkflow)

	require.True(t, cb1.accepted)
	require.Error(t, cb1.completeErr)
	require.Contains(t, cb1.completeErr.Error(), "piggy bank service down")
	err := env.GetWorkflowResult(nil)
	require.NoError(t, err)
}
