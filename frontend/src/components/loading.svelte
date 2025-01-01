<script>
	import Icon from "./icon.svelte";

	export let size = "128px";
</script>

<div>
	<p><slot></slot></p>
	<Icon name="mark-github-24" color="rainbow" {size} />
</div>

<style>
</style>
