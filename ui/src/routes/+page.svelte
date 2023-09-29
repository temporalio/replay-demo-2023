<script lang="ts">
    import { goto } from '$app/navigation';
    import { to, from, amount } from '$lib/stores/transfer';
    import { APIRoutes } from '$lib/utilities/url';
		import Icon from '@temporalio/ui/holocene/icon/icon.svelte';
    import { onMount } from 'svelte';

	const initiateWorkflow = async () => {
		const res = await fetch(APIRoutes.initiate, {
				method: "GET",
				headers: {
						"Content-Type": "application/json",
				},
		});
		const result = await res.json();
		const { workflowID, runID } = result;
		goto(`/${workflowID}/${runID}/from`);
	}

	const startSchedule = async () => {
		const res = await fetch(APIRoutes.schedule, {
				method: "GET",
		});
		goto('/schedules');
	}
	
	onMount(() => {
		$from = 'Checking';
		$to = 'Savings';
		$amount = '0.00'
	})
</script>

	<div class="flex flex-col gap-8 items-start w-full md:max-w-xl px-8 py-4">
	<div class="flex gap-4 items-center">
		<div class="flex flex-col">
			<Icon name="arrow-left" class="text-green-400 scale-90" />
			<Icon name="arrow-right" class="text-green-400 -mt-2 scale-90" />	
		</div>
		<h1 class="text-4xl">
			Bank Transfer
		</h1>
	</div>
	<div class="flex gap-2 items-center w-full">
		<button on:click={initiateWorkflow} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white py-4 rounded-xl">Initiate Transfer</button>
		<button on:click={startSchedule} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white py-4 rounded-xl">Schedule Transfer</button>
	</div>
</div>

