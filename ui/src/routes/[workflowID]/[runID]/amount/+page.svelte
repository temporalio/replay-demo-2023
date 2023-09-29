<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { amount } from '$lib/stores/transfer';
  import { APIRoutes } from '$lib/utilities/url';
  import Alert from '@temporalio/ui/holocene/alert.svelte';

  $: ({ workflowID, runID } = $page.params)

  let errorMessage = '';

	$: decimalAmount = parseFloat($amount)

  const back = () => {
    goto(`/${workflowID}/${runID}/to`);
  }
  
	const updateWorkflow = async () => {
		console.log('Update workflow');
		const res = await fetch(APIRoutes.amount, {
				method: "POST",
				headers: {
						"Content-Type": "application/json",
				},
				body: JSON.stringify({
						WorkflowID: workflowID,
						RunId: runID,
						Amount: decimalAmount,
				}),
		});
		const { success, error } = await res.json();
    if (success) {
      goto(`/${workflowID}/${runID}/transfer`);
    } else {
      errorMessage = error
    }
	}
</script>

<div class="flex flex-col gap-4 w-full">
  {#if errorMessage}
    <Alert intent="error">{errorMessage}</Alert>
  {/if}
	<div class="flex flex-col">
		<p class="text-sm text-gray-400">Amount</p>
		<div class="flex gap-2 items-center">
			<p class="text-3xl text-gray-400 font-thin">$</p>
			<input class="text-6xl bg-black border-none outline-none w-full" bind:value={$amount} />
		</div>
	</div>
</div>
<div class="flex gap-2 items-center w-full">
  <button on:click={back} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white disabled:bg-red-400 py-4 rounded-xl">Back</button>
	<button on:click={() => goto('/')} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white disabled:bg-red-400 py-4 rounded-xl">Start Over</button>
  <button on:click={updateWorkflow} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white disabled:bg-red-400 py-4 rounded-xl">Transfer</button>
</div>

