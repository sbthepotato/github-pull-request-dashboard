<script>
	import { onMount } from "svelte";
	import { setUrlParam } from "$lib/index.js";

	export let user;
	export let size = "l";
	export let type = "a";

	let element = HTMLElement;

	onMount(() => {
		if (type === "a") {
			element.href = user.html_url;
			element.target = "_blank";
		} else {
			element.href = "";
			element.target = "";
		}
	});

	function click_handler() {
		if (type === "a") {
			return;
		} else {
			setUrlParam("user", user.login);
		}
	}
</script>

<svelte:element
	this={type}
	role="button"
	tabindex="0"
	on:click={() => {
		click_handler();
	}}
	bind:this={element}
	class="container">
	{#if size !== "xs"}
		<img src={user.avatar_url} alt="{user.login} avatar" class={size} />
	{/if}
	<div class="name-container">
		{#if user.name !== undefined}
			<span class="name">{user.name}</span>
			{#if size !== "xs"}
				<span class="login">@{user.login}</span>
			{/if}
		{:else}
			<span class="big-login">@{user.login}</span>
		{/if}
	</div>
</svelte:element>

<style>
	.container {
		display: inline-flex;
		align-items: center;
		color: var(--text);
		font-size: medium;
	}

	.container:hover {
		color: var(--text-links);
		cursor: pointer;
	}

	.container:hover span.name,
	.container:hover span.big-login {
		text-decoration: underline;
	}

	img {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		margin-right: 8px;
	}

	img.s {
		width: 20px;
		height: 20px;
	}

	div.name-container {
		display: flex;
		flex-direction: column;
	}

	div.name-container > span {
		display: block;
		text-align: left;
	}

	span.login {
		color: var(--text-alt);
	}
</style>
