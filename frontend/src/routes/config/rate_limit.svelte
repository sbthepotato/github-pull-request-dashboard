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
	<h2>rate limit check</h2>

	<div>
		<ul>
			{#each Object.entries(answer) as [key, value]}
				<li>
					<strong>{key}</strong>
					<ul>
						{#each Object.entries(value) as [childKey, childValue]}
							<li>{childKey} - {childValue}</li>
						{/each}
					</ul>
				</li>
			{/each}
		</ul>

		<Button color="blue" on_click={get}>check rate limit</Button>
	</div>
</div>

<style>
	div.container {
		flex: 1;
	}

	li {
		list-style: none;
	}
</style>
