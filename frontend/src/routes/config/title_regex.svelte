<script>
	import { onMount } from "svelte";
	import Button from "../../components/button.svelte";

	let temp_id = -100;
	let resultList = [];
	let result = "";
	let err = "";
	let loading = false;

	onMount(() => {
		get();
	});

	async function get() {
		try {
			resultList = {};
			err = "";

			const response = await fetch("api/config/get_title_regex_list");

			if (response.ok) {
				resultList = await response.json();

				resultList.push({
					title_regex_id: temp_id,
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

	async function set() {
		const data = resultList.map((titleRegex, index) => ({
			title_regex_id: index,
			regex_pattern: titleRegex.regex_pattern,
			link: titleRegex.link,
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

	async function del(id, idx) {
		try {
			err = "";
			result = "";

			if (id > 0) {
				const response = await fetch(
					"api/config/delete_regex?titleRegexId=" + id,
					{
						method: "POST",
					},
				);
				result = await response.text();
			}
			resultList.splice(idx, 1);
			resultList = resultList;
		} catch (error) {
			err = error.message;
		}
	}

	function handle_change(idx) {
		let entry = resultList[idx];

		if (
			entry.regex_pattern != "" &&
			entry.link != "" &&
			resultList[resultList.length - 1].title_regex_id === entry.title_regex_id
		) {
			temp_id--;
			resultList.push({
				title_regex_id: temp_id,
				regex_pattern: "",
				link: "",
				repository_name: "",
			});
		}
	}
</script>

<div class="container">
	<h3>Title Regex</h3>
	<p>Enter regex patterns and a link to insert with the pattern</p>
	<p>
		E.G. <code>"[Aa][Bb]#(\d+)"</code> with
		<code>"https://example.com/&lbrace;1&rbrace;/test"</code>
		would make a pull request title with <code>AB#123</code> in the title make AB#123 into a link
		to <code>https://example.com/123/test</code>
	</p>
	{#if err !== ""}
		<p class="bad">
			{err}
		</p>
	{:else if loading}
		<Loading size="64px">Loading title regex list...</Loading>
	{:else if resultList !== null}
		<table>
			<thead>
				<tr>
					<th>Regex Pattern</th>
					<th>Link</th>
				</tr>
			</thead>
			<tbody>
				{#each resultList as entry, idx}
					<tr>
						<td>
							<input
								type="text"
								placeholder="[Aa][Bb]#(\d+)"
								bind:value={entry.regex_pattern}
								on:change={() => handle_change(idx)}
								on:input={() => handle_change(idx)} />
						</td>
						<td>
							<input
								type="text"
								placeholder="example.com/"
								bind:value={entry.link}
								on:change={() => handle_change(idx)}
								on:input={() => handle_change(idx)} />
						</td>

						<td>
							<Button
								color="red"
								on:click={() => del(entry.title_regex_id, idx)}>
								delete
							</Button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
		{#if result !== ""}
			<p class="good">
				{result}
			</p>
		{/if}
		<div>
			<Button color="green" on:click={() => set()}>Save</Button>
		</div>
	{/if}
</div>

<style>
	div.container {
		flex: 1;
	}

	.good {
		color: var(--green);
	}

	.bad {
		color: var(--red);
	}

	input {
		background-color: var(--blue-bg);
		outline: 1px solid var(--border);
		border-radius: 4px;
	}
</style>
