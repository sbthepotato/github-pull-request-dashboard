<script>
	import { base } from "$app/paths";
	import Button from "../../components/button.svelte";

	let answer = {};
	let err = "";

	function entries(value) {
		return Object.entries(value ?? {});
	}

	function rateClass(value) {
		const used = value?.used ?? 0;
		const limit = value?.limit ?? 0;

		return {
			good: used > 0 && used / limit < 0.55,
			ok: used > 0 && used / limit >= 0.55 && used / limit <= 0.9,
			bad: used > 0 && used / limit > 0.9,
		};
	}

	async function get() {
		try {
			answer = {};
			err = "";

			const response = await fetch(`${base}/api/config/rate_limit`);

			if (response.ok) {
				answer = await response.json();
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
		{#each entries(answer) as [key, value]}
			{#if value != null}
				<code
					class="flex-item"
					class:good={rateClass(value).good}
					class:ok={rateClass(value).ok}
					class:bad={rateClass(value).bad}>
					<strong>{key}</strong>

					{#each entries(value) as [childKey, childValue]}
						<br />{childKey}: {childValue}
					{/each}
				</code>
			{/if}
		{/each}
	</div>
	<Button color="blue" on:click={() => get()}>check rate limit</Button>
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

	code.good {
		color: var(--green);
	}

	code.ok {
		color: var(--yellow);
	}

	code.bad {
		color: var(--red);
	}
</style>
