<script>
	import { onMount } from "svelte";
	import User from "../../components/user.svelte";
	import Button from "../../components/button.svelte";
	import Loading from "../../components/loading.svelte";
	import { boolToString } from "$lib/index.js";

	export let repository = "";
	let err = "";
	let result = {};
	let loading = false;

	onMount(() => {
		getUsers(false, "members", repository);
	});

	async function get(refresh, type, repository) {
		try {
			loading = true;
			err = "";
			result = {};

			const response = await fetch(
				"api/config/get_users?refresh=" +
					boolToString(refresh) +
					"&type=" +
					type +
					"&repo=" +
					repository ?? null,
			);

			if (response.ok) {
				result = await response.json();
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		} finally {
			loading = false;
		}
	}

	$: repository, get(false, "members", repository);
</script>

<div class="container">
	<h3>Member Configuration</h3>

	{#if err !== ""}
		<p>
			{err}
		</p>
	{:else if loading}
		<Loading size="64px">Loading Members...</Loading>
	{:else if result !== null}
		<table>
			<thead>
				<tr>
					<th>Team Name</th>
					<th>Members</th>
				</tr>
			</thead>
			<tbody>
				{#each Object.entries(result) as [teamName, users]}
					<tr>
						<td class="team-name">{teamName}</td>
						<td>
							{#each users as user}
								<div class="user-container">
									<User {user} />
								</div>
							{/each}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{:else}
		<p>No members found</p>
	{/if}

	<div class="button-container">
		<Button color="blue" on_click={() => get(true, "users", repository)}>
			Sync all users with GitHub
		</Button>
		<Button color="blue" on_click={() => get(true, "members", repository)}>
			Sync repository team members with GitHub
		</Button>
	</div>
</div>

<style>
	div.container {
		margin: 8px;
		flex: 2;
	}

	.team-name {
		font-weight: bold;
	}

	.user-container {
		display: inline-flex;
		align-items: center;
		margin: 8px;
	}
</style>
