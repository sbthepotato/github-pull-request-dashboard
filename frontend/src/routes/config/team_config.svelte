<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";

	let repository = "";
	let teams = [];
	let result = "";
	let err = "";

	onMount(() => {
		get_teams(false, repository);
	});

	async function get_teams(refresh, repository) {
		try {
			teams = [];
			err = "";

			let url = "api/config/get_teams";

			if (refresh) {
				url = url + "?refresh=y";
			}

			if (repository !== undefined && repository !== "") {
				url = url + "?repo=" + repository;
			}

			const response = await fetch(url);

			if (response.ok) {
				teams = await response.json();
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		}
	}

	async function set_teams() {
		const data = teams.map((team) => ({
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
</script>

<h2>Team Configuration</h2>

{#if err !== ""}
	<p>
		{err}
	</p>
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
	<p>{teams.length} teams found</p>
	{#if result !== ""}
		<p>
			{result}
		</p>
	{/if}
{:else}
	<p>No teams found</p>
{/if}

<Button color="green" on_click={() => get_teams(true)}>
	hard refresh team list
</Button>
<Button color="green" on_click={() => set_teams()}>Save Teams</Button>

<style>
</style>
