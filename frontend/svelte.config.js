import adapter from "@sveltejs/adapter-static";

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter({
			pages: "../backend/static",
			assets: "../backend/static",
			fallback: "index.html",
		}),
		paths: {
			base: "",
		},
	},
};

export default config;
