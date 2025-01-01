<script>
	import { onMount } from "svelte";
	import User from "../../components/user.svelte";
	import Button from "../../components/button.svelte";
	import Loading from "../../components/loading.svelte";

	let err = "";
	let result = {};
	let loading = false;

	onMount(() => {
		getUsers(false, "members");
	});

	async function getUsers(refresh, type) {
		try {
			loading = true;
			err = "";
			result = {};

			let url = "api/config/";

			if (type === "users") {
				url = url + "get_users";
			} else {
				url = url + "get_members";
			}

			if (refresh) {
				url = url + "?refresh=y";
			}

			const response = await fetch(url);

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
</script>

<div class="container">
	<h2>Member Configuration</h2>

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
		<Button color="blue" on_click={() => getUsers(true, "users")}>
			Sync all users with GitHub
		</Button>
		<Button color="blue" on_click={() => getUsers(true, "members")}>
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
