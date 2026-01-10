<script>
	import Checkbox from "../../components/checkbox.svelte";

	function tvMode(checked) {
		console.log(checked);
	}
</script>

<Checkbox on:change={() => tvMode(checked)}>TV Mode</Checkbox>
