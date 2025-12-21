<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";
	import Checkbox from "../../components/checkbox.svelte";
	import Loading from "../../components/loading.svelte";

	let repos = [];
	let setResult = "";
	let err = "";
	let loading = false;

	onMount(() => {
		getRepos(false);
	});

	async function getRepos(refresh) {
		try {
			loading = true;
			repos = [];
			err = "";

			let url = "api/config/get_repos";

			if (refresh) {
				url = url + "?refresh=y";
			}

			const response = await fetch(url);

			if (response.ok) {
				repos = await response.json();
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		} finally {
			loading = false;
		}
	}

	async function setRepos() {
		const data = repos.map((repo) => ({
			name: repo.name,
			enabled: repo.enabled,
		}));

		try {
			err = "";
			setResult = "";

			const response = await fetch("api/config/set_repos", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(data),
			});

			setResult = await response.text();
		} catch (error) {
			err = error.message;
		}
	}
</script>

<div class="container">
	<h3>Repository Configuration</h3>
	<h4>Set the active repositories</h4>

	<div>
		{#if err !== ""}
			<p>
				{err}
			</p>
		{:else if loading}
			<Loading size="64px">Loading repositories...</Loading>
		{:else if repos.length > 0}
			<ul>
				{#each repos as repo}
					<li>
						<Checkbox
							id={repo.name}
							name={repo.name}
							bind:checked={repo.enabled}>
							{repo.name}
						</Checkbox>
					</li>
				{/each}
			</ul>

			<p>{repos.length} repositories found</p>
		{/if}

		<Button color="blue" on_click={() => getRepos(true)}>
			Sync repositories with GitHub
		</Button>
		<Button color="green" on_click={() => setRepos()}>Save Repositories</Button>

		{#if setResult !== ""}
			<p>
				{setResult}
			</p>
		{/if}
	</div>
</div>

<style>
	div.container {
		flex: 1;
	}

	ul {
		list-style: none;
		columns: 3;
	}

	li {
		padding: 12px;
	}
</style>
