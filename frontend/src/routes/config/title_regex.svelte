<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";

	let result = {};
	let err = "";
	let loading = false;

	onMount(() => {
		getTitles();
	});

	async function getTitles() {
		try {
			result = {};
			err = "";

			const response = await fetch("api/config/get_title_regex_list");

			if (response.ok) {
				result = await response.json();
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		}
	}
</script>

<div class="container">
	<h3>Title Regex</h3>
	<p>Enter regex patterns and a link to insert with the pattern</p>
	{#if err !== ""}
		<p>
			{err}
		</p>
	{:else if loading}
		<Loading size="64px">Loading title regex list...</Loading>
	{:else if result !== null}
		<table>
			<thead>
				<tr>
					<th>Regex Pattern</th>
					<th>Link</th>
				</tr>
			</thead>
			<tbody>
				{#each result as title}
					<tr>
						<td> {title.regex} </td>
						<td> {title.link} </td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<style>
	div.container {
		flex: 1;
	}

	span.good {
		color: var(--green);
	}

	span.bad {
		color: var(--red);
	}
</style>
