<script>
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { setUrlParam } from "$lib/index.js";

	export let useDefault = false;
	export let text = "";
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
			useDefault &&
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
					handleChange();
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
	{#if text !== ""}
		<p>{text}</p>
	{/if}

	<select
		bind:value={selected}
		on:change={handleChange}
		on:input={handleChange}>
		{#each repos as repo}
			<option value={repo.name}>{repo.name}</option>
		{/each}
	</select>
{/if}
