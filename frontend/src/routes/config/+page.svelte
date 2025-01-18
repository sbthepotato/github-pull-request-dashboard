<script>
	import { page } from "$app/stores";

	import Button from "../../components/button.svelte";
	import HelloGo from "./hello_go.svelte";
	import RepoConfig from "./repo_config.svelte";
	import TeamConfig from "./team_config.svelte";
	import MemberConfig from "./member_config.svelte";
	import RepositorySelect from "../../components/repositorySelect.svelte";

	let repository = "";

	function handleParams() {
		repository = $page.url.searchParams.get("repo");
	}

	$: $page.url.search, handleParams();
</script>

<Button to="/">Back to home</Button>

<RepositorySelect>Select repository to configure</RepositorySelect>

<div class="container">
	<TeamConfig {repository} />
	<MemberConfig {repository} />
</div>

<div class="container">
	<RepoConfig />
</div>

<div class="container">
	<HelloGo />
</div>

<footer>
	made by sbthepotato | <a
		href="https://github.com/sbthepotato/gh-pull-request-dashboard"
		target="_blank">GitHub Repository</a>
</footer>

<style>
	div.container {
		display: flex;
		margin: 8px 4px 48px;
	}

	footer {
		margin: 48px 4px 16px;
	}
</style>
