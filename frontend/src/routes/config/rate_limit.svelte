<script>
	import Button from "../../components/button.svelte";

	let answer = "";
	let err = "";

	async function get() {
		try {
			answer = "";
			err = "";

			const response = await fetch("api/config/rate_limit");

			if (response.ok) {
				answer = await response.json();
				console.log(answer);
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		}
	}
</script>

<div class="container">
	<h2>Rate Limit Check</h2>

	<div class="flex-container">
		{#each Object.entries(answer) as [key, value]}
			<code class="flex-item" class:used={answer.used > 0}>
				<strong>{key}</strong>

				{#each Object.entries(value) as [childKey, childValue]}
					<br />{childKey}: {childValue}
				{/each}
			</code>
		{/each}
	</div>
	<Button color="blue" on_click={get}>check rate limit</Button>
</div>

<style>
	div.container {
		flex: 1;
	}

	div.flex-container {
		display: flex;
		flex-wrap: wrap;
		justify-content: space-around;
	}

	code.flex-item {
		box-sizing: border-box;
		margin: 8px 2px;
		flex: 1;
		flex-grow: 0;
		white-space: nowrap;
		text-align: left;
	}
</style>
