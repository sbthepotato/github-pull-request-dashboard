<script>
	export let checked;
	export let id = "";
	export let name = "";

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
	<input type="checkbox" {id} {name} bind:checked on:change={handleChange} />
	<slot></slot>
</label>
