<script>
	import Button from "../../components/button.svelte";

	let answer = "";
	let err = "";

	async function helloGo() {
		try {
			answer = "";
			err = "";

			const response = await fetch("api/config/hello_go");

			if (response.ok) {
				answer = await response.text();
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		}
	}
</script>

<div class="container">
	<h2>Connection Test</h2>

	<p>
		{#if err}
			<br />
			<span class="bad">{err}</span>
		{:else if answer}
			<br />
			<span class="good">{answer}</span>
		{/if}
	</p>

	<Button color="blue" on:click={() => helloGo}
		>Say hello to the backend</Button>
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
