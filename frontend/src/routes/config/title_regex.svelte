<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";
	import { addLinks } from "$lib/index.js";

	let titleRegexList = [];
	let err = "";
	let loading = false;

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

				titleRegexList.push({ regex_pattern: "", link: "" });
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
		}));
	}
</script>

<div class="container">
	<h3>Title Regex</h3>
	<p>Enter regex patterns and a link to insert with the pattern</p>
	<p>
		E.G. <code>"[Aa][Bb]#(\\d+)"</code> with
		<code>"https://example.com/"</code>
		would become:
		{@html addLinks("AB#123", [
			{
				regex_pattern: "[Aa][Bb]#(\\d+)",
				link: "https://example.com/",
			},
		])}
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
				</tr>
			</thead>
			<tbody>
				{#each titleRegexList as entry}
					<tr>
						<td>
							<input
								type="text"
								placeholder="[Aa][Bb]#(\\d+)"
								value={entry.regex_pattern} />
						</td>
						<td>
							<input
								type="text"
								placeholder="example.com/"
								value={entry.link} />
						</td>
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
