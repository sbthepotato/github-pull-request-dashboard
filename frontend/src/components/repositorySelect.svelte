<script>
	import { onMount } from "svelte";

	let repos = [];
	let err = "";

	onMount(() => {
		getRepos();
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
</script>

<select>
	{#each repos as repo}
		<option value={repo.name}>{repo.name}</option>
	{/each}
</select>
