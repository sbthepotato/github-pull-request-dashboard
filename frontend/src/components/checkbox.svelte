<script>
	export let checked;
	export let id = "";
	export let disabled = false;
</script>

<label class:disabled>
	<input
		type="checkbox"
		{id}
		{disabled}
		bind:checked
		on:change
		class:disabled />

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
