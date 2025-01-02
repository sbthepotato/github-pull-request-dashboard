<script>
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";

	export let type = "button";
	export let on_click = () => {};
	export let to = null;
	export let color = "grey";

	function click_handler(event) {
		if (to) {
			const url_prefix = import.meta.env.VITE_URL_PATH;
			const params = $page.url.searchParams.toString();
			const fullUrl = params ? `${url_prefix+to}?${params}` : url_prefix+to;
			goto(fullUrl);
		} else {
			on_click(event);
		}
	}
</script>

<button class={color} {type} on:click={click_handler}>
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

	button:hover {
		cursor: pointer;
	}
</style>
