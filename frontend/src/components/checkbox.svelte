<script>
	export let checked;
	export let id = "";
	export let name = "";
	export let disabled = false;

	function handleChange(event) {
		checked = event.target.checked;
		const changeEvent = new CustomEvent("change", {
			detail: { id, checked },
			bubbles: true,
		});
		dispatchEvent(changeEvent);
	}
</script>

<label class:disabled>
	<input
		type="checkbox"
		{id}
		{name}
		{disabled}
		bind:checked
		on:change={handleChange} />
	<slot></slot>
</label>

<style>
	label {
		cursor: pointer;
	}

	label.disabled {
		cursor: not-allowed;
	}
</style>
