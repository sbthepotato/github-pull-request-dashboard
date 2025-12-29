<script>
	export let checked;
	export let id = "";
	export let name = "";
	export let show_checkbox = true;

	function handleChange(event) {
		checked = event.target.checked;
		const changeEvent = new CustomEvent("change", {
			detail: { id, checked },
			bubbles: true,
		});
		dispatchEvent(changeEvent);
	}
</script>

<label>
	<input
		type="checkbox"
		{id}
		{name}
		bind:checked
		on:change={handleChange}
		class:hidden={!show_checkbox} />
	<slot></slot>
</label>

<style>
	input[type="checkbox"].hidden {
		opacity: 0;
	}
</style>
