<script lang="ts">
	import { onMount } from 'svelte';

	type Page = {
		title: string;
		src: string;
		background: string;
		iteraction: Iteraction[];
	};

	type Iteraction = {
		x: number;
		y: number;
		link: string;
	};

	export let page: Page;
	export let projectId: string;

	let image: Element;
	let width: number;
	let height: number;
	let browser = false;

	function setCoords() {
		let rect = image.getBoundingClientRect();
		width = rect.width;
		height = rect.height;
	}

	onMount(() => {
		setCoords();
		browser = true;
	});
</script>

<div style="position: relative;" on:resize={() => setCoords()}>
	{#if page.iteraction !== undefined && browser}
		{#each page.iteraction as i}
			<a
				class="iteraction-box"
				href={i.link}
				target="_blank"
				style={[
					`margin-left:${(i.x / 100) * width}px;`,
					`margin-top:${(i.y / 100) * height}px;`
				].join('')}
			></a>
		{/each}
	{/if}
	<img bind:this={image} width="1080" height="1920" src={`/files/${projectId}/${page.src}`} />
</div>

<style>
	.iteraction-box {
		width: 30px;
		height: 30px;
		display: block;
		background-color: #ff0000;
		opacity: 0.3;
		position: absolute;
		z-index: 100;
	}
</style>
