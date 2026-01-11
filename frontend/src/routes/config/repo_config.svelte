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
		get(false);
	});

	async function get(refresh) {
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

	async function set() {
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
					<li class:enabled={repo.enabled} class:archived={repo.archived}>
						<Checkbox
							id={repo.name}
							name={repo.name}
							disabled={repo.archived}
							bind:checked={repo.enabled}>
							{repo.name}
						</Checkbox>
					</li>
				{/each}
			</ul>

			<p>{repos.length} repositories found</p>
		{/if}

		<Button color="blue" on:click={() => get(true)}>
			Sync repositories with GitHub
		</Button>
		<Button color="green" on:click={() => set()}>Save Repositories</Button>

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
		margin: auto;
		padding: 0;
		list-style: none;
		display: flex;
		flex-wrap: wrap;
		justify-content: space-around;
	}

	li {
		margin: 8px;
		border-radius: 8px;
		cursor: pointer;
		min-width: 330px;
		white-space: nowrap;
		display: flex;
		align-items: center;
		cursor: inherit;
	}

	li.enabled {
		font-weight: bold;
	}

	li.archived {
		color: var(--yellow);
		opacity: 50%;
		cursor: not-allowed;
	}
</style>
