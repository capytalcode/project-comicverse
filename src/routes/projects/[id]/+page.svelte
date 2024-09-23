<script lang="ts">
	import type { PageData } from './$types';

	export let data: PageData;
	const pages = data.project.pages;

	let modal = false;

	let reader, scroll;
	let color = pages[0].background;
</script>

<pre style="position: fixed; bottom: 0; font-size: 0.6rem;">
<code
		>{JSON.stringify(
			{
				color: color,
				scroll: scroll
			},
			null,
			2
		)}
</code>
</pre>

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
	<article
		class="reader"
		style={`--bg-color: ${color}`}
		bind:this={reader}
		on:scroll={() => (scroll = reader.scrollTop)}
	>
		<div class="pages">
			{#each pages as page}
				<div class="page" style={`background-color:${page.background}`}>
					<img width="1080" height="1920" src={`/files/${data.project.id}/${page.src}`} />
					<form method="POST" action="?/delete-file" class="delete-file">
						<fieldset role="group">
							<input type="text" disabled value={`${page.src}`} name="file" />
							<input type="submit" value="Delete page" class="pico-background-red-500" />
						</fieldset>
					</form>
				</div>
			{/each}
		</div>
	</article>
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
