<script>
	import Checkbox from "../../components/checkbox.svelte";
	import { onMount } from "svelte";
	import { browser } from "$app/environment";

	let tv_mode = false;
	let total_count = false;
	let auto_refresh = false;
	let seamless_reload = false;
	let last_updated = false;
	let disable_regex = false;

	onMount(() => {
		if (browser) {
			if (localStorage.getItem("tv_mode") !== null) {
				tv_mode = true;
			}
			if (localStorage.getItem("total_count") !== null) {
				total_count = true;
			}
			if (localStorage.getItem("auto_refresh") !== null) {
				auto_refresh = true;
			}
			if (localStorage.getItem("seamless_reload") !== null) {
				seamless_reload = true;
			}
			if (localStorage.getItem("last_updated") !== null) {
				last_updated = true;
			}
			if (localStorage.getItem("disable_regex") !== null) {
				disable_regex = true;
			}
		}
	});

	function handleChange(name, value) {
		if (browser) {
			if (value) {
				localStorage.setItem(name, String(value));
			} else {
				localStorage.removeItem(name);
			}
		}
	}
</script>

<div>
	<ul>
		<li>
			<Checkbox
				bind:checked={tv_mode}
				on:change={() => handleChange("tv_mode", tv_mode)}>
				Enable TV Mode (will always hide buttons out of scroll view)
			</Checkbox>
		</li>
		<li>
			<Checkbox
				bind:checked={total_count}
				on:change={() => handleChange("total_count", total_count)}>
				Count all Pull Requests instead of only filtered
			</Checkbox>
		</li>
		<li>
			<Checkbox
				bind:checked={auto_refresh}
				on:change={() => handleChange("auto_refresh", auto_refresh)}>
				Automatically refresh dashboard (every 10 minutes)
			</Checkbox>
		</li>
		<li>
			<Checkbox
				bind:checked={seamless_reload}
				on:change={() => handleChange("seamless_reload", seamless_reload)}>
				Reload without showing loading animation on dashboard
			</Checkbox>
		</li>
		<li>
			<Checkbox
				bind:checked={last_updated}
				on:change={() => handleChange("last_updated", last_updated)}>
				Show last updated time in bottom right of dashboard
			</Checkbox>
		</li>
		<li>
			<Checkbox
				bind:checked={disable_regex}
				on:change={() => handleChange("disable_regex", disable_regex)}>
				Disable regex links in titles
			</Checkbox>
		</li>
	</ul>
</div>
