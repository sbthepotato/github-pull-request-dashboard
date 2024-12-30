<script>
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { set_url_param } from "$lib/index.js";

	export let useDefault = false;
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
		set_url_param("repo", selected);
	}

	$: repos, getDefaultRepo();
</script>

<select bind:value={selected} on:change={handleChange} on:input={handleChange}>
	{#each repos as repo}
		<option value={repo.name}>{repo.name}</option>
	{/each}
</select>
