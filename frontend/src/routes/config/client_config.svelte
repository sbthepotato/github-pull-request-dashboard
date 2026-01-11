<script>
	import Checkbox from "../../components/checkbox.svelte";
	import { onMount } from "svelte";

	let tv_mode = false;

	// // Load from localStorage on startup
	// onMount(() => {
	// 	const stored = localStorage.getItem("tv_mode");
	// 	if (stored !== null) {
	// 		tv_mode = stored === "true";
	// 	}
	// });

	// // Persist whenever it changes
	// $: localStorage.setItem("tv_mode", tv_mode);
</script>

<Checkbox bind:checked={tv_mode}>TV Mode</Checkbox>
