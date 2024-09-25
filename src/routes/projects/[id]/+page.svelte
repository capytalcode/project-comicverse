<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageData } from './$types';
	import { arrayBufferToBlob } from 'blob-util';
	import IImage from './IteractiveImage.svelte';

	export let data: PageData;
	const pages = data.project.pages;

	let modal = false;

	let reader: Element;
	let scroll: number;
	let color = hexToRgb(pages[0]?.background ?? '#181818');
	let currentPage = 0;
	let colorPerc = 0;
	let currentColor = color;
	let nextColor = hexToRgb(pages[1]?.background ?? pages[0]?.background ?? '#181818');
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

	function hexToRgb(color: string): number[] {
		return [
			parseInt(color.substring(1, 3), 16) / 255,
			parseInt(color.substring(3, 5), 16) / 255,
			parseInt(color.substring(5, 7), 16) / 255
		];
	}

	function rgbToHex(rgb: number[]): string {
		return `#${
			Math.round(rgb[0] * 255)
				.toString(16)
				.padStart(2, '0') +
			Math.round(rgb[1] * 255)
				.toString(16)
				.padStart(2, '0') +
			Math.round(rgb[2] * 255)
				.toString(16)
				.padStart(2, '0')
		}`;
	}

	function blendRgbColors(c1: number[], c2: number[], ratio: number): number[] {
		return [
			c1[0] * (1 - ratio) + c2[0] * ratio,
			c1[1] * (1 - ratio) + c2[1] * ratio,
			c1[2] * (1 - ratio) + c2[2] * ratio
		];
	}

	let fileInput: Element;
	let blobUrl: string | undefined = undefined;
	let currentIteraction: { x: number; y: number; link: string };
	let iteractionUrl = '';
	let iteractions: { x: number; y: number; link: string }[] = [];
	let imageElement: Element;
	let imageX = 0;
	let imageY = 0;
	let imageWidth = 0;
	let imageHeight = 0;
	function readFile(file: Blob) {
		let reader = new FileReader();
		reader.onloadend = function (e) {
			let buf = e.target?.result;
			let blob = arrayBufferToBlob(buf as ArrayBuffer, file.type);
			blobUrl = window.URL.createObjectURL(blob);
		};

		reader.readAsArrayBuffer(file);
	}

	let temp: any;

	let images: Map<string, { width: number; height: number }> = new Map();
</script>

{#if browser}
	<pre style="position: fixed; bottom: 0; font-size: 0.6rem;">
<code
			>{JSON.stringify(
				{
					page: currentPage,
					color: {
						background: rgbToHex(color),
						current: rgbToHex(currentColor),
						next: rgbToHex(nextColor),
						percentage: Math.round(colorPerc)
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

	<details style="position: fixed; bottom: 0; font-size: 0.6rem;">
		<pre>
<code>{JSON.stringify(pages)}</code>
</pre>
	</details>
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
			{#if blobUrl}
				<div class="blob-image">
					<div class="blob-image-image">
						<div
							class="iteraction-box"
							style={`${[
								`margin-left:${Math.round((imageX / 100) * imageWidth)}px;`,
								`margin-top:${Math.round((imageY / 100) * imageHeight)}px;`
							].join('')} pointer-events: none;`}
						></div>
						{#each iteractions as i}
							<a
								class="iteraction-box"
								href={i.link}
								style={[
									`margin-left:${Math.round((i.x / 100) * imageWidth)}px;`,
									`margin-top:${Math.round((i.y / 100) * imageHeight)}px;`
								].join('')}
							></a>
						{/each}
						<img
							style="margin: auto 0;"
							src={blobUrl}
							bind:this={imageElement}
							on:mousemove={(e) => {
								let rect = imageElement.getBoundingClientRect();
								imageX = Math.round(((e.clientX - rect.left) / rect.width) * 100);
								imageY = Math.round(((e.clientY - rect.top) / rect.height) * 100);
								imageWidth = rect.width;
								imageHeight = rect.height;
							}}
							on:click={() => {
								currentIteraction ||= { x: 0, y: 0, link: '' };
								currentIteraction.x = imageX;
								currentIteraction.y = imageY;
							}}
							alt=""
						/>
					</div>
					<fieldset role="group">
						<input
							type="url"
							placeholder="Iteraction url"
							bind:value={iteractionUrl}
							on:change={() => {
								currentIteraction ||= { x: 0, y: 0, link: '' };
								currentIteraction.link = iteractionUrl;
							}}
							on:mouseout={() => {
								currentIteraction ||= { x: 0, y: 0, link: '' };
								currentIteraction.link = iteractionUrl;
							}}
						/>
						<input
							type="button"
							value="Add"
							on:click|preventDefault={() => {
								iteractions.push({
									x: currentIteraction.x,
									y: currentIteraction.y,
									link: currentIteraction.link
								});
								iteractions = iteractions;
							}}
						/>
					</fieldset>
					<div>
						<code>{imageX} {imageY}</code>
						<code>
							Iteraction: {JSON.stringify(currentIteraction)}
						</code>
						<input
							style="display:hidden;"
							type="text"
							value={JSON.stringify(iteractions.filter(Boolean))}
							name="iteractions"
						/>
					</div>
				</div>
			{/if}
			<input
				type="file"
				required
				name="file"
				bind:this={fileInput}
				on:change={() => {
					// @ts-ignore
					readFile(fileInput.files[0]);
				}}
			/>
			<input type="submit" value="Add page" />
		</form>
	</article>
</dialog>
<section class="project">
	<aside>
		<a data-sveltekit-reload href="/" style="font-size: 0.5rem;">Return home</a>
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
				style={`--bg-color: rgba(${color.map((c) => c * 255).join(',')}, 0.8)`}
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
					currentColor = hexToRgb(pages[i]?.background);
					nextColor = pages[i + 1]?.background ? hexToRgb(pages[i + 1]?.background) : currentColor;

					currentChunk = chunks[i];
					nextChunk = chunks[i + 1] ?? maxScroll;

					colorPerc = ((scroll - currentChunk) / (nextChunk - currentChunk)) * 100;

					color = blendRgbColors(currentColor, nextColor, colorPerc / 100);

					currentPage = i;
				}}
			>
				<div class="pages">
					{#each pages as page, key}
						{@const coord = key * chunk}
						<div class="page" style={`background-color:${page.background}`}>
							<IImage {page} projectId={data.project.id} />
							<form method="POST" action="?/delete-file" class="delete-file">
								<fieldset role="group">
									<input type="text" value={`${page.src}`} name="file" />
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

	.iteraction-box {
		width: 30px;
		height: 30px;
		display: block;
		background-color: #ff0000;
		opacity: 0.3;
		position: absolute;
		z-index: 100;
	}

	.blob-image-image {
		position: relative;
	}

	.reader {
		display: flex;

		width: 80vw;
		height: 100vh;

		justify-content: center;
		padding-top: 5rem;
		padding-bottom: 5rem;
		margin-bottom: 0;

		background-color: var(--bg-color);
		overflow-y: scroll;
	}

	.page {
		width: calc(1080px / 3.5);
		min-height: calc(1920px / 3.5);
		@media (min-width: 1024px) {
			width: calc(1080px / 2.5);
			min-height: calc(1920px / 2.5);
		}
		background-color: #fff;

		display: flex;
		flex-direction: column;
		justify-content: space-between;

		padding: 0;

		box-shadow: 0rem 1rem 1rem 0rem rgba(0, 0, 0, 0.5);

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
		& * {
			font-size: 0.8rem !important;
			@media (min-width: 1024px) {
				font-size: unset;
			}
		}
	}
</style>
