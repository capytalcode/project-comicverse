<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageData } from './$types';

	export let data: PageData;
	const pages = data.project.pages;

	let modal = false;

	let reader: Element;
	let scroll: number;
	let color = pages[0]?.background;
	let currentPage = 0;
	let colorPerc = 0;
	let currentColor = color;
	let nextColor = pages[1]?.background ?? pages[0]?.background;
	let currentChunk = 0;
	let nextChunk = 0;

	let browser = false;
	let maxScroll = 0;
	let chunk = 0;
	let chunks: number[] = [];
	onMount(() => {
		browser = true;
		maxScroll = Math.max(reader.scrollHeight - reader.clientHeight);
	});

	function blendHexColors(c1: string, c2: string, ratio: number): string {
		// Convert hex colors to RGB
		const rgb1 = [
			parseInt(c1.substring(1, 3), 16) / 255,
			parseInt(c1.substring(3, 5), 16) / 255,
			parseInt(c1.substring(5, 7), 16) / 255
		];
		const rgb2 = [
			parseInt(c2.substring(1, 3), 16) / 255,
			parseInt(c2.substring(3, 5), 16) / 255,
			parseInt(c2.substring(5, 7), 16) / 255
		];

		// Blend RGB values
		const blendedRgb = [
			rgb1[0] * (1 - ratio) + rgb2[0] * ratio,
			rgb1[1] * (1 - ratio) + rgb2[1] * ratio,
			rgb1[2] * (1 - ratio) + rgb2[2] * ratio
		];

		// Convert blended RGB to hex
		const blendedHex = `#${
			Math.round(blendedRgb[0] * 255)
				.toString(16)
				.padStart(2, '0') +
			Math.round(blendedRgb[1] * 255)
				.toString(16)
				.padStart(2, '0') +
			Math.round(blendedRgb[2] * 255)
				.toString(16)
				.padStart(2, '0')
		}`;

		return blendedHex;
	}
</script>

{#if browser}
	<pre style="position: fixed; bottom: 0; font-size: 0.6rem;">
<code
			>{JSON.stringify(
				{
					page: currentPage,
					color: {
						background: color,
						current: currentColor,
						next: nextColor,
						percentage: colorPerc
					},
					scroll: {
						current: scroll,
						max: maxScroll,
						chunks: chunks,
						currentChunk: currentChunk,
						nextChunk: nextChunk
					}
				},
				null,
				2
			)}
</code>
</pre>
{/if}

<dialog open={modal}>
	<article>
		<header>
			<button
				aria-label="Close"
				rel="prev"
				on:click={() => {
					modal = false;
				}}
			></button>
			<p>
				<strong>Add new page</strong>
			</p>
		</header>
		<form method="POST" action="?/addpage" enctype="multipart/form-data">
			<input type="text" required placeholder="Page title" name="title" />
			<input type="color" required placeholder="Background color" name="color" />
			<input type="file" required name="file" />
			<input type="submit" value="Add page" />
		</form>
	</article>
</dialog>
<section class="project">
	<aside>
		<section>
			<h1>{data.project.title}</h1>
			<p class="id">{data.project.id}</p>
			<button
				class="add"
				on:click={() => {
					modal = true;
				}}>Add page</button
			>
			<form method="POST">
				<input type="submit" formaction="?/delete" value="Delete" class="pico-background-red-500" />
			</form>
		</section>
	</aside>
	{#key maxScroll}
		{#if browser}
			<article
				class="reader"
				style={`--bg-color: ${color}`}
				bind:this={reader}
				on:scroll={() => {
					scroll = reader.scrollTop;
					if (maxScroll === 0) {
						maxScroll = Math.max(reader.scrollHeight - reader.clientHeight);
						chunk = Math.round(maxScroll / pages.length);
						for (let i = 0; i < pages.length; i++) {
							chunks = [...chunks, chunk * i];
						}
					}
					let i = chunks.findIndex((c) => c > scroll - chunk);
					currentColor = pages[i]?.background;
					nextColor = pages[i + 1]?.background ?? currentColor;

					currentChunk = chunks[i];
					nextChunk = chunks[i + 1] ?? maxScroll;

					colorPerc = ((scroll - currentChunk) / (nextChunk - currentChunk)) * 100;

					color = blendHexColors(currentColor, nextColor, colorPerc / 100);

					currentPage = i;
				}}
			>
				<div class="pages">
					{#each pages as page, key}
						{@const coord = key * chunk}
						<div class="page" style={`background-color:${page.background}`}>
							<img width="1080" height="1920" src={`/files/${data.project.id}/${page.src}`} />
							<form method="POST" action="?/delete-file" class="delete-file">
								<fieldset role="group">
									<input type="text" disabled value={`${page.src}`} name="file" />
									<input type="submit" value="Delete page" class="pico-background-red-500" />
								</fieldset>
							</form>
						</div>
						<code>{coord}</code>
					{/each}
				</div>
			</article>
		{/if}
	{/key}
</section>

<style>
	h1 {
		font-size: 1.5rem;
	}

	.id {
		font-size: 0.5rem;
		opacity: 0.3;
	}

	.reader {
		display: flex;
		width: 80vw;
		justify-content: center;
		padding-top: 5rem;
		padding-bottom: 5rem;
		margin-bottom: 0;

		background-color: var(--bg-color);
		height: 100vh;
		overflow-y: scroll;
	}

	.page {
		width: calc(1080px / 2.5);
		min-height: calc(1920px / 2.5);
		background-color: #fff;

		display: flex;
		flex-direction: column;
		justify-content: space-between;

		padding: 0;

		& form {
			margin: 1rem;
			margin-bottom: 0;
		}
	}
	.pages {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.project {
		display: flex;
		margin-bottom: 0;
	}

	.add {
		width: 100%;
		margin-bottom: 0.5rem;
	}

	aside {
		padding: 1rem;
		width: 20vw;
	}
</style>
