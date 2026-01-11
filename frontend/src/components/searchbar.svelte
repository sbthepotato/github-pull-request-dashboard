<script>
	export let placeholder = "Search...";
	export let value = "";
</script>

<div class="search-container">
	<input
		type="search"
		{placeholder}
		class={value ? "no-rounding" : ""}
		bind:value
		on:change
		on:input />
	{#if value}
		<button on:click={clearSearch}>X</button>
	{/if}
</div>

<style>
	.search-container {
		margin: 8px auto;
		position: relative;
	}

	input {
		display: inline-block;
		box-sizing: border-box;
		padding: 8px 16px;
		border-radius: 4px;
		width: 35vw;
		min-width: 256px;
	}

	.no-rounding {
		border-radius: 4px 0px 0px 4px;
	}

	button {
		position: absolute;
		height: 100%;
		width: 32px;
		border: none;
		background-color: var(--border-blue);
		border-radius: 0px 4px 4px 0px;
	}

	button:hover {
		cursor: pointer;
	}

	input::-webkit-search-decoration,
	input::-webkit-search-cancel-button,
	input::-webkit-search-results-button,
	input::-webkit-search-results-decoration {
		-webkit-appearance: none;
	}
</style>
