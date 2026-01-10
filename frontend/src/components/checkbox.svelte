<script>
	export let checked;
	export let id = "";
	export let disabled = false;
</script>

<input type="checkbox" {id} name={id} {disabled} bind:checked class:disabled />
<label for={id} class:disabled>
	<slot></slot>
</label>

<style>
	label,
	input {
		cursor: pointer;
	}

	.disabled {
		cursor: not-allowed;
	}
</style>
