<script>
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { setUrlParam } from "$lib/index.js";

	let repos = [];
	let err = "";
	let selected = "";

	onMount(() => {
		getRepos();

		selected = $page.url.searchParams.get("repo");

		getDefaultRepo();
	});

	async function getRepos() {
		try {
			repos = [];
			err = "";

			const response = await fetch("api/config/get_repos?active=y");

			if (response.ok) {
				repos = await response.json();
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		}
	}

	async function getDefaultRepo() {
		if (
			(!selected || selected === "") &&
			repos !== undefined &&
			repos.length > 0
		) {
			try {
				err = "";
				selected = "";

				const response = await fetch("api/config/get_default_repository");

				if (response.ok) {
					selected = await response.text();
				} else {
					throw new Error(await response.text());
				}
			} catch (error) {
				err = error.message;
			}
		}
	}

	function handleChange() {
		setUrlParam("repo", selected);
	}

	$: repos, getDefaultRepo();
</script>

{#if repos !== undefined && repos.length > 0}
	<div class="container">
		<p><slot></slot></p>

		<select
			bind:value={selected}
			on:change={handleChange}
			on:input={handleChange}>
			{#each repos as repo}
				<option value={repo.name}>{repo.name}</option>
			{/each}
		</select>
	</div>
{/if}

<style>
	div.container {
		display: inline-block;
	}

	select {
		background-color: var(--content-bg-alt);
		color: var(--text);
		padding: 8px;
		border-radius: 8px;
		border: 1px solid var(--border);
	}

	select:hover,
	select:focus {
		cursor: pointer;
	}

</style>
