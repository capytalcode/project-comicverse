<svelte:head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css" />
	<link
		rel="stylesheet"
		href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.colors.min.css"
	/>
	<meta name="color-scheme" content="dark" />
</svelte:head>

<main>
	<slot></slot>
</main>

<style>
	main {
		width: 100vw;
		height: 100vh;
	}
</style>
