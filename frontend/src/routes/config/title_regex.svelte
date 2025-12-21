<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";

	let titleRegexList = [];
	let err = "";
	let loading = false;
	let result = "";

	onMount(() => {
		getTitles();
	});

	async function getTitles() {
		try {
			titleRegexList = {};
			err = "";

			const response = await fetch("api/config/get_title_regex_list");

			if (response.ok) {
				titleRegexList = await response.json();

				titleRegexList.push({
					regex_pattern: "",
					link: "",
					repository_name: "",
				});
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		}
	}

	async function setTitleRegex() {
		const data = titleRegexList.map((titleRegex, index) => ({
			title_regex_id: index,
			regex_pattern: titleRegex.regex_pattern,
			link: titleRegex.link,
			repository_name: titleRegex.repository_name,
		}));

		try {
			err = "";
			result = "";

			const response = await fetch("api/config/set_regex", {
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

<div class="container">
	<h3>Title Regex</h3>
	<p>Enter regex patterns and a link to insert with the pattern</p>
	<p>
		E.G. <code>"[Aa][Bb]#(\d+)"</code> with
		<code>"https://example.com/"</code>
		would make a pull request title with AB#123 in the title make AB#123 into a link
		to <code>https://example.com/123</code>
	</p>
	{#if err !== ""}
		<p>
			{err}
		</p>
	{:else if loading}
		<Loading size="64px">Loading title regex list...</Loading>
	{:else if titleRegexList !== null}
		<table>
			<thead>
				<tr>
					<th>Regex Pattern</th>
					<th>Link</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each titleRegexList as entry}
					<tr>
						<td>
							<input
								type="text"
								placeholder="[Aa][Bb]#(\d+)"
								bind:value={entry.regex_pattern} />
						</td>
						<td>
							<input
								type="text"
								placeholder="example.com/"
								bind:value={entry.link} />
						</td>
						<td>
							<input type="text" bind:value={entry.repository_name} />
						</td>
						<td>
							<button color="red">delete</button>
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
		<div>
			<Button color="green" on_click={() => setTitleRegex()}>Save</Button>
		</div>
	{/if}
</div>

<style>
	div.container {
		flex: 1;
	}
</style>
