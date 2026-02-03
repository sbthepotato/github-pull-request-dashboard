<script>
	export let type = "button";
	export let color = "grey";
</script>

<button class={color} {type} on:click>
	<slot></slot>
</button>

<style>
	button {
		color: var(--text);
		display: inline;
		border: none;
		padding: 8px;
		border-radius: 8px;
		font-weight: bold;
		margin: 4px;
	}

	button.grey {
		background-color: var(--button-grey);
	}

	button.green {
		background-color: var(--button-green);
	}

	button.blue {
		background-color: var(--border-blue);
	}

	button.red {
		background-color: var(--red);
	}

	button:hover {
		cursor: pointer;
	}
</style>
