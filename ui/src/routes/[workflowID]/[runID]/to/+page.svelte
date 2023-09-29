<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { APIRoutes } from '$lib/utilities/url';
  import Icon from '@temporalio/ui/holocene/icon/icon.svelte';
  import Alert from '@temporalio/ui/holocene/alert.svelte';
  import { to } from '$lib/stores/transfer';

  $: ({ workflowID, runID } = $page.params)

  let errorMessage = '';

  const back = () => {
    goto(`/${workflowID}/${runID}/from`);
  }

	const updateWorkflow = async () => {
		const res = await fetch(APIRoutes.toAccount, {
				method: "POST",
				headers: {
						"Content-Type": "application/json",
				},
				body: JSON.stringify({
						WorkflowID: workflowID,
						RunId: runID,
						ToAccount: $to,
				}),
		});
		const { success, error } = await res.json();
    if (success) {
      goto(`/${workflowID}/${runID}/amount`);
    } else {
      errorMessage = error
    }
	}
</script>

<div class="flex flex-col gap-4 w-full">
  {#if errorMessage}
    <Alert intent="error">{errorMessage}</Alert>
  {/if}
  <div class="flex flex-col gap-2 bg-gray-900 p-4">
    <p class="text-sm flex gap-1 items-center">Transfer to<Icon name="arrow-right" class="scale-90" /></p>
    <select
      class="bg-gray-900 text-gray-200 focus:outline-none text-xl h-12 border border-white px-2 rounded-xl" 
      id="transfer-to-filter"
      bind:value={$to}
  >
    {#each ['Checking', 'Savings', 'Crypto', 'Piggy Bank', 'Wallet'] as value}
      <option {value}>{value}</option>
    {/each}
  </select>
  </div>
</div>
<div class="flex gap-2 items-center w-full">
  <button  on:click={back} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white disabled:bg-red-400 py-4 rounded-xl">Back</button>
  <button disabled={!$to} on:click={updateWorkflow} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white disabled:bg-red-400 py-4 rounded-xl">Next</button>
</div>

