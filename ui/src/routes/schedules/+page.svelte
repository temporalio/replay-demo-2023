<script lang="ts">
    import { goto } from '$app/navigation';
  import Icon from '@temporalio/ui/holocene/icon/icon.svelte';
</script>

	<div class="flex flex-col gap-8 items-start w-full md:max-w-xl px-8 py-4">
		<div class="flex flex-col gap-2 items-start w-full">
		<div class="flex gap-4 items-center">
			<div class="flex flex-col">
				<Icon name="arrow-left" class="text-green-400 scale-90" />
				<Icon name="arrow-right" class="text-green-400 -mt-2 scale-90" />	
			</div>
			<h1 class="text-4xl">
				Transfers Scheduled
			</h1>
		</div>
	</div>
  <div class="flex gap-2 items-center w-full">
    <button on:click={() => goto('/')} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white disabled:bg-red-400 py-4 rounded-xl">Back</button>
  </div>  
</div>

