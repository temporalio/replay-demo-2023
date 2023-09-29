<script>
	import { goto } from '$app/navigation';
	import { amount, to, from } from '$lib/stores/transfer';
</script>

<svelte:head>
	<title>Transfer</title>
	<meta name="description" content="Transfer funds" />
</svelte:head>

<div class="flex flex-col gap-2">
	<h1 class="text-2xl">Funds transferred</h1>
	<p>From: {$from}</p>
	<p>To: {$to}</p>
	<p>Amount: ${$amount}</p>
</div>
<div class="flex gap-2 items-center w-full">
  <button on:click={() => goto('/')} class="w-full bg-gray-900 hover:bg-green-400 border-2 hover:border-green-400 hover:text-white disabled:bg-red-400 py-4 rounded-xl">Start Over</button>
</div>

