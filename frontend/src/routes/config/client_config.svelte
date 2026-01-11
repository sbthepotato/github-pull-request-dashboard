<script>
	import Checkbox from "../../components/checkbox.svelte";
	import { onMount } from "svelte";
	import { browser } from "$app/environment";

	let tv_mode = false;
	let total_count = false;

	onMount(() => {
		if (browser) {
			if (localStorage.getItem("tv_mode") !== null) {
				tv_mode = true;
			}
			if (localStorage.getItem("total_count") !== null) {
				total_count = true;
			}
		}
	});

	function setTVMode(value) {
		tv_mode = value;

		if (browser) {
			if (tv_mode) {
				localStorage.setItem("tv_mode", String(value));
			} else {
				localStorage.removeItem("tv_mode");
			}
		}
	}

	function setTotalCount(value) {
		total_count = value;

		if (browser) {
			if (total_count) {
				localStorage.setItem("total_count", String(value));
			} else {
				localStorage.removeItem("total_count");
			}
		}
	}
</script>

<div>
	<Checkbox bind:checked={tv_mode} on:change={() => setTVMode(tv_mode)}>
		Enable TV Mode (will always hide buttons out of scroll view)
	</Checkbox>
	<Checkbox
		bind:checked={total_count}
		on:change={() => setTotalCount(total_count)}>
		Count all Pull Requests instead of only filtered
	</Checkbox>
</div>
