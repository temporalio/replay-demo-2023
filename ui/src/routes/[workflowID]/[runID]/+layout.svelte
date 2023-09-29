<script lang="ts">
  import { page } from '$app/stores';
  import Icon from '@temporalio/ui/holocene/icon/icon.svelte';

  $: ({ workflowID, runID } = $page.params)
</script>

	<div class="flex flex-col gap-8 items-start w-full md:max-w-xl px-8 py-4">
		<div class="flex flex-col gap-2 items-start w-full">
		<div class="flex gap-4 items-center">
			<div class="flex flex-col">
				<Icon name="arrow-left" class="text-green-400 scale-90" />
				<Icon name="arrow-right" class="text-green-400 -mt-2 scale-90" />	
			</div>
			<h1 class="text-4xl">
				Bank Transfer
			</h1>
		</div>
		<div class="flex gap-2">
			<p class="text-sm text-gray-400">{workflowID}</p>
			<p class="text-sm text-gray-400">{runID}</p>
		</div>
	</div>
  <slot />
</div>

