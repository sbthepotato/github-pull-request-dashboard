<script>
	import { onDestroy, onMount } from "svelte";
	import { page } from "$app/stores";
	import { browser } from "$app/environment";
	import {
		setUrlParam,
		stringToBool,
		boolToString,
		redirect,
		getPrettyDate,
	} from "$lib/index.js";
	import Button from "../components/button.svelte";
	import Checkbox from "../components/checkbox.svelte";
	import Searchbar from "../components/searchbar.svelte";
	import Loading from "../components/loading.svelte";
	import RepositorySelect from "../components/repositorySelect.svelte";
	import PRTable from "./pr_table.svelte";
	import PRAgg from "./pr_aggregation.svelte";

	// variables for showing result
	let err = "";
	let result = {};
	let pr_list = {};
	let repository = "";
	let loading = false;

	// settings
	let show_search = false;
	let tv_mode = false;
	let total_count = false;
	let auto_reload = false;
	let seamless_reload = false;
	let last_updated = false;

	// filters and other
	let reload_interval;
	let user_filter = "";
	let user_filter_object = {};
	let search_query = "";
	let updated_time;

	onMount(() => {
		repository = $page.url.searchParams.get("repo");

		getPullRequests(false, repository);

		// temporary warning because of the rename
		if ($page.url.searchParams.get("created_by")) {
			window.alert(
				"the created by filter has been renamed to user, please delete your current bookmark and make a new one.",
			);
		}

		user_filter = $page.url.searchParams.get("user");

		if (browser) {
			if (localStorage.getItem("tv_mode") !== null) {
				tv_mode = true;
			}
			if (localStorage.getItem("total_count") !== null) {
				total_count = true;
			}
			if (localStorage.getItem("auto_reload") !== null) {
				auto_reload = true;
			}
			if (localStorage.getItem("seamless_reload") !== null) {
				seamless_reload = true;
			}
			if (localStorage.getItem("last_updated") !== null) {
				last_updated = true;
			}
		}

		show_search = stringToBool(
			$page.url.searchParams.get("show_search"),
			false,
		);
	});

	onDestroy(() => {
		clearInterval(reload_interval);
	});

	async function getPullRequests(refresh, repository) {
		try {
			if (Object.keys(result).length === 0 || !seamless_reload) {
				loading = true;
			}
			err = "";
			result = {};
			pr_list = {};

			const response = await fetch(
				"api/dashboard/get_pr_list?refresh=" +
					boolToString(refresh) +
					"&repo=" +
					repository ?? null,
			);

			if (response.ok) {
				result = await response.json();

				pr_list = result.pull_requests;
				updated_time = result.updated;
			} else {
				throw new Error(await response.text());
			}
		} catch (error) {
			err = error.message;
		} finally {
			loading = false;
		}
	}

	function handleSearchbarChange(value) {
		search_query = value.toLowerCase();
		getFilter();
	}

	function handleParams() {
		user_filter = $page.url.searchParams.get("user");
		show_search = stringToBool(
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
			result = {}
			updated_time = null
			getPullRequests(false, repository);
		}
	}

	function getFilter() {
		if (
			(user_filter !== null || search_query !== "") &&
			result.pull_requests !== undefined
		) {
			let cleaned_query = "";
			let anti_search = false;

			if (search_query[0] === "!") {
				cleaned_query = search_query.substring(1);
				anti_search = true;
			} else {
				cleaned_query = search_query;
			}

			if (anti_search) {
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
						!pr.title.toLowerCase().includes(cleaned_query) &&
						!pr.awaiting?.toLowerCase().includes(cleaned_query) &&
						!pr.created_by.login.toLowerCase().includes(cleaned_query) &&
						!pr.created_by.name?.toLowerCase().includes(cleaned_query) &&
						!pr.base.label.toLowerCase().includes(cleaned_query) &&
						!pr.number.toString().includes(cleaned_query) &&
						!pr.review_overview?.some(
							(review) =>
								review.state === "REVIEW_REQUESTED" &&
								review.user?.login.toLowerCase().includes(cleaned_query) &&
								review.user?.name.toLowerCase().includes(cleaned_query),
						) &&
						!pr.labels?.some((label) =>
							label.name.toLowerCase().includes(cleaned_query),
						),
				);
			} else {
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
						(pr.title.toLowerCase().includes(cleaned_query) ||
							pr.awaiting?.toLowerCase().includes(cleaned_query) ||
							pr.created_by.login.toLowerCase().includes(cleaned_query) ||
							pr.created_by.name?.toLowerCase().includes(cleaned_query) ||
							pr.base.label.toLowerCase().includes(cleaned_query) ||
							pr.number.toString().includes(cleaned_query) ||
							pr.review_overview?.some(
								(review) =>
									review.state === "REVIEW_REQUESTED" &&
									(review.user?.login.toLowerCase().includes(cleaned_query) ||
										review.user?.name.toLowerCase().includes(cleaned_query)),
							) ||
							pr.labels?.some((label) =>
								label.name.toLowerCase().includes(cleaned_query),
							)),
				);
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

	function clearFilters() {
		setUrlParam("user", null);
		user_filter = null;
		search_query = "";
		getFilter();
	}

	$: $page.url.search, handleParams();
	$: result, get_current_user(), getFilter();
	$: user_filter, get_current_user(), getFilter();
	$: if (show_search) {
		setUrlParam("show_search", "y");
	} else {
		setUrlParam("show_search");
	}
	$: if (auto_reload) {
		reload_interval = setInterval(function () {
			getPullRequests(false, repository);
		}, 600000);
	} else {
		clearInterval(reload_interval);
	}
	$: handleSearchbarChange(search_query);
</script>

<section class="pr-table" class:tv_mode>
	{#if err !== ""}
		{err}
	{:else if loading}
		<Loading>Loading PR list...</Loading>
	{:else}
		<PRAgg
			pr_list={total_count ? result.pull_requests : pr_list}
			review_teams={result.review_teams} />
		{#if show_search}
			<Searchbar
				bind:value={search_query}
				placeholder="Search Pull Requests..." />
		{/if}
		{#if last_updated && updated_time}
			<p class="last-updated">
				Lasted updated {getPrettyDate(updated_time, true)}
			</p>
		{/if}
		{#if user_filter === null}
			<PRTable {pr_list} />
		{:else if user_filter !== null}
			<PRTable
				pr_list={pr_list?.filter((pr) => pr.created_by.login === user_filter)}>
				{#if user_filter_object.name}
					Created by {user_filter_object.name}
				{:else}
					Created by @{user_filter}
				{/if}
			</PRTable>

			<PRTable
				pr_list={pr_list?.filter((pr) =>
					pr.review_overview?.some(
						(review) =>
							review.user?.login === user_filter &&
							review.state === "REVIEW_REQUESTED",
					),
				)}>
				{#if user_filter_object.name}
					{user_filter_object.name} requested reviewer
				{:else}
					@{user_filter} requested reviewer
				{/if}
			</PRTable>

			{#if user_filter_object.team}
				<PRTable
					show_empty={false}
					pr_list={pr_list?.filter(
						(pr) =>
							pr.unassigned === true &&
							pr.created_by.login != user_filter &&
							pr.awaiting === user_filter_object.team.name,
					)}>
					Waiting on {user_filter_object.team.name} - Not assigned to anyone else
				</PRTable>
			{/if}
		{/if}
	{/if}
</section>

<section class="buttons">
	<Button color="grey" on:click={() => redirect("/config")}>Config</Button>
	<Button color="blue" on:click={() => getPullRequests(true, repository)}>
		Refresh PR List
	</Button>
	<RepositorySelect />
	<Checkbox id="show_search" bind:checked={show_search}>Show Search</Checkbox>
	{#if user_filter !== null || search_query !== ""}
		<Button color="blue" on:click={() => clearFilters()}>Clear Filters</Button>
	{/if}
</section>

<style>
	section.pr-table {
		margin-bottom: 32px;
	}

	.tv_mode {
		min-height: 100%;
	}

	.last-updated {
		z-index: 1;
		color: var(--text-alt);
		background-color: var(--content-bg-alt);
		border-radius: 8px 0 0 0;
		font-size: small;
		padding: 8px;
		margin: 0;
		position: fixed;
		bottom: 0px;
		right: 0px;
	}
</style>
