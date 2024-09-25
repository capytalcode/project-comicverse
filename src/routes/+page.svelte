<script lang="ts">
	import type { Project } from '$lib';
	import type { PageData } from './$types';

	export let data: PageData;

	if ((data.projects.length + 1) % 3 !== 0) {
		data.projects.push({ Name: '', ID: '' });
		data.projects.push({ Name: '', ID: '' });
	}
</script>

<section>
	{#each data.projects as p}
		<article>
			<h1><a data-sveltekit-reload href={`/projects/${p.ID}`}>{p.Name}</a></h1>
			<p class="id">{p.ID}</p>
		</article>
	{/each}
	<article>
		<form method="POST">
			<fieldset role="group">
				<input type="text" name="project-name" placeholder="Project Name" required />
				<input type="submit" value="Create" />
			</fieldset>
		</form>
	</article>
</section>

<style>
	section {
		display: grid;

		@media (min-width: 768px) {
			grid-template-columns: repeat(3, minmax(0%, 1fr));
		}
		grid-column-gap: var(--pico-grid-column-gap);
		grid-row-gap: var(--pico-grid-row-gap);
		padding: 1rem var(--pico-grid-row-gap);
	}

	.id {
		font-size: 0.7rem;
		opacity: 0.3;
	}
</style>
