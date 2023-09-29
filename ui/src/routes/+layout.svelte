<script>
	import Icon from '@temporalio/ui/holocene/icon/icon.svelte';
	import "@fontsource/cooper-hewitt";
	import '../app.css';
</script>

<svelte:head>
	<title>TemporalCash</title>
	<meta name="description" content="Temporal feature demo app" />
</svelte:head>

<header class="fixed w-full px-8 py-2 top-0">
	<h1><a href="/" class="flex gap-2"><Icon name="temporal-logo" class="text-green-400" />TemporalCash</a></h1>
</header>
<main class="w-full">
	<section class="flex flex-col items-center justify-center grow h-[100vh] w-[100vw] bg-black text-white">
		<slot />
	</section>
</main>

<footer class="fixed bottom-0 w-full px-8 py-2 text-[12px] flex items-center justify-center">
	<p><a href="https://temporal.io">temporal.io</a> to learn more</p>
</footer>

<style lang="postcss">
	header, main, footer {
		@apply text-white;
	}

	.background-animate {
		@apply bg-gradient-to-r from-green-300 via-green-400 to-green-500;

    background-size: 400%;

    -webkit-animation: AnimationName 3s ease infinite;
    -moz-animation: AnimationName 3s ease infinite;
    animation: AnimationName 3s ease infinite;
  }

  @keyframes AnimationName {
    0%,
    100% {
      background-position: 0% 50%;
    }
    50% {
      background-position: 100% 50%;
    }
  }

</style>
