<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";
	import Loading from "../../components/loading.svelte";
	import { boolToString } from "$lib/index.js";

	export let repository = "";
	let teams = [];
	let result = "";
	let err = "";
	let loading = false;

	onMount(() => {
		get(false, repository);
	});

	async function get(refresh, repository) {
		try {
			loading = true;
			teams = [];
			err = "";

			const response = await fetch(
				"api/config/get_teams?refresh=" +
					boolToString(refresh) +
					"&repo=" +
					repository ?? null,
			);

			if (response.ok) {
				teams = await response.json();
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
		const data = teams.map((team) => ({
			slug: team.slug,
			name: team.name,
			repository_name: team.repository_name,
			review_order: team.review_order,
		}));

		try {
			err = "";
			result = "";

			const response = await fetch("api/config/set_teams", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(data),
			});

			result = await response.text();
		} catch (error) {
			err = error.message;
		}
	}

	$: repository, get(false, repository);
</script>

<div class="container">
	<h3>Team Configuration</h3>

	{#if err !== ""}
		<p>
			{err}
		</p>
	{:else if loading}
		<Loading size="64px">Loading teams...</Loading>
	{:else if teams.length > 0}
		<table>
			<thead>
				<tr>
					<th></th>
					<th>Review Order</th>
				</tr>
			</thead>
			<tbody>
				{#each teams as team}
					<tr>
						<td>
							{team.name}
						</td>
						<td>
							<input
								type="number"
								min="0"
								max={teams.length}
								bind:value={team.review_order} />
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
		{#if result !== ""}
			<p>
				{result}
			</p>
		{/if}
	{:else}
		<p>No teams found</p>
	{/if}

	<div class="button-container">
		<Button color="blue" on:click={() => get(true)}>
			Sync teams with GitHub
		</Button>
		<Button color="green" on:click={() => set()}>Save Teams</Button>
	</div>
</div>

<style>
	div.container {
		margin: 8px;
		flex: 1;
	}
</style>
