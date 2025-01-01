<script>
	import { onDestroy, onMount } from "svelte";
	import { page } from "$app/stores";
	import {
		set_many_url_params,
		set_url_param,
		string_to_bool,
	} from "$lib/index.js";

	import Button from "../components/button.svelte";
	import Checkbox from "../components/checkbox.svelte";
	import Searchbar from "../components/searchbar.svelte";
	import Loading from "../components/loading.svelte";
	import RepositorySelect from "../components/repositorySelect.svelte";
	import PRTable from "./pr_table.svelte";
	import PRAgg from "./pr_aggregation.svelte";

	let err = "";
	let result = {};
	let pr_list = {};
	let repository = "";

	let loading = false;

	let show_search = false;
	let auto_reload = false;
	let reload_interval;

	let user_filter = "";
	let user_filter_object = {};
	let search_query = "";

	onMount(() => {
		repository = $page.url.searchParams.get("repo");

		get_pr_list(false, repository);

		user_filter = $page.url.searchParams.get("user");

		auto_reload = string_to_bool(
			$page.url.searchParams.get("auto_reload"),
			false,
		);

		show_search = string_to_bool(
			$page.url.searchParams.get("show_search"),
			false,
		);
	});

	onDestroy(() => {
		clearInterval(reload_interval);
	});

	async function get_pr_list(refresh, repository) {
		try {
			loading = true;
			err = "";
			result = {};
			pr_list = {};

			let url = "api/dashboard/get_pr_list";

			if (refresh) {
				url = url + "?refresh=y";
			}

			if (
				repository !== undefined &&
				repository !== null &&
				repository !== ""
			) {
				url = url + "?repo=" + repository;
			}

			const response = await fetch(url);

			if (response.ok) {
				result = await response.json();

				pr_list = result.pull_requests;
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		} finally {
			loading = false;
		}
	}

	function handle_searchbar_change(event) {
		search_query = event.detail.value.toLowerCase();
		get_filter();
	}

	function handle_params() {
		user_filter = $page.url.searchParams.get("user");
		show_search = string_to_bool(
			$page.url.searchParams.get("show_search"),
			false,
		);

		const newRepository = $page.url.searchParams.get("repo") ?? "";

		if (
			newRepository !== "" &&
			repository !== "" &&
			newRepository !== repository
		) {
			repository = newRepository;
			get_pr_list(false, repository);
		}
	}

	function get_filter() {
		if (
			(user_filter !== null || search_query !== "") &&
			result.pull_requests !== undefined
		) {
			pr_list = result.pull_requests.filter(
				(pr) =>
					(user_filter === null ||
						pr.created_by.login === user_filter ||
						pr.review_overview?.some(
							(review) =>
								review.user?.login === user_filter &&
								review.state === "REVIEW_REQUESTED" &&
								pr.awaiting !== "Changes Requested",
						) ||
						(pr.unassigned === true &&
							pr.created_by.login != user_filter &&
							pr.awaiting === user_filter_object.team?.name &&
							pr.awaiting !== "Changes Requested")) &&
					(pr.title.toLowerCase().includes(search_query) ||
						pr.awaiting?.toLowerCase().includes(search_query) ||
						pr.created_by.login.toLowerCase().includes(search_query) ||
						pr.created_by.name.toLowerCase().includes(search_query) ||
						pr.base.label.toLowerCase().includes(search_query) ||
						pr.number.toString().includes(search_query) ||
						pr.review_overview?.some(
							(review) =>
								review.state === "REVIEW_REQUESTED" &&
								(review.user?.login.toLowerCase().includes(search_query) ||
									review.user?.name.toLowerCase().includes(search_query)),
						) ||
						pr.labels?.some((label) =>
							label.name.toLowerCase().includes(search_query),
						)),
			);

			if (user_filter !== null) {
			}
		} else {
			pr_list = result.pull_requests;
		}
	}

	function get_current_user() {
		if (user_filter == null) {
			user_filter_object = {};
		} else {
			if (result.users !== undefined) {
				result.users.forEach((user) => {
					if (user?.login === user_filter) {
						user_filter_object = user;
						return true;
					}
				});
			}
		}
	}

	function clear_filters() {
		set_many_url_params({ user: null });
		user_filter = null;
		search_query = "";
		get_filter();
	}

	$: $page.url.search, handle_params();
	$: result, get_current_user(), get_filter();
	$: user_filter, get_current_user(), get_filter();
	$: if (show_search) {
		set_url_param("show_search", "y");
	} else {
		set_url_param("show_search");
	}
	$: if (auto_reload) {
		set_url_param("auto_reload", "y");
	} else {
		set_url_param("auto_reload");
	}
</script>

<section class="pr-table">
	{#if err !== ""}
		{err}
	{:else if loading}
		<Loading>Loading PR list...</Loading>
	{:else}
		<PRAgg {pr_list} review_teams={result.review_teams} />
		{#if show_search}
			<Searchbar
				value={search_query}
				placeholder="Search Pull Requests..."
				on:change={handle_searchbar_change}
				on:input={handle_searchbar_change} />
		{/if}
		{#if user_filter === null}
			<PRTable {pr_list} />
		{:else if user_filter !== null}
			<PRTable
				title="Created by {user_filter}"
				pr_list={pr_list?.filter(
					(pr) => pr.created_by.login === user_filter,
				)} />

			<PRTable
				title="{user_filter} requested reviewer"
				pr_list={pr_list?.filter((pr) =>
					pr.review_overview?.some(
						(review) =>
							review.user?.login === user_filter &&
							review.state === "REVIEW_REQUESTED",
					),
				)} />

			{#if user_filter_object.team}
				<PRTable
					show_empty={false}
					title="Waiting on {user_filter_object.team
						.name} - Not assigned to anyone else"
					pr_list={pr_list?.filter(
						(pr) =>
							pr.unassigned === true &&
							pr.created_by.login != user_filter &&
							pr.awaiting === user_filter_object.team.name,
					)} />
			{/if}
		{/if}
	{/if}
</section>

<section class="buttons">
	<Button color="grey" to="/config">Config</Button>
	<Button color="blue" on_click={() => get_pr_list(true)}>
		Refresh PR List
	</Button>
	<RepositorySelect />
	<Checkbox id="auto_reload" bind:checked={auto_reload}>Auto Refresh</Checkbox>
	<Checkbox id="show_search" bind:checked={show_search}>Show Search</Checkbox>
	{#if user_filter !== null || search_query !== ""}
		<Button color="blue" on_click={() => clear_filters()}>Clear Filters</Button>
	{/if}
</section>

<style>
	section.pr-table {
		margin-bottom: 32px;
	}
</style>
