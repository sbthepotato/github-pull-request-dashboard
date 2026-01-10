<script>
	import { page } from "$app/stores";

	import Button from "../../components/button.svelte";
	import HelloGo from "./hello_go.svelte";
	import RepoConfig from "./repo_config.svelte";
	import TeamConfig from "./team_config.svelte";
	import TitleRegexConfig from "./title_regex.svelte";
	import MemberConfig from "./member_config.svelte";
	import RepositorySelect from "../../components/repositorySelect.svelte";
	import RateLimit from "./rate_limit.svelte";
	import ClientConfig from "./client_config.svelte";

	let repository = "";

	function handleParams() {
		repository = $page.url.searchParams.get("repo");
	}

	$effect(() => {
		$page.url.search;
		handleParams();
	});
</script>

<Button to="/">Back to home</Button>

<section>
	<h2>Client Configuration</h2>
	<p>
		Configure settings that will be saved in cookies and only apply to the
		current client
	</p>
	<ClientConfig></ClientConfig>
</section>

<section>
	<h2>Server Configuration</h2>
	<p>
		Configure settings that will be saved on the server and apply for everyone
	</p>
	<div>
		<RepositorySelect>Select repository to configure</RepositorySelect>
	</div>

	<div class="container">
		<TeamConfig {repository} />
		<MemberConfig {repository} />
	</div>

	<div class="container">
		<RepoConfig />
	</div>

	<div class="container">
		<TitleRegexConfig />
	</div>

	<div class="container">
		<RateLimit />
	</div>

	<div class="container">
		<HelloGo />
	</div>
</section>

<footer>
	made by sbthepotato | <a
		href="https://github.com/sbthepotato/github-pull-request-dashboard"
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
