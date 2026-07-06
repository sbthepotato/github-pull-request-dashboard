import adapter from "@sveltejs/adapter-static";
import { readFileSync } from "node:fs";

// Read config from the backend .env so the base path has a single source of
// truth shared with the Go server. Falls back to the BASE_PATH env var, then
// to root ("").
function readEnv(key) {
	if (process.env[key]) {
		return process.env[key];
	}
	try {
		const env = readFileSync("../backend/.env", "utf-8");
		for (const line of env.split("\n")) {
			const trimmed = line.trim();
			if (!trimmed || trimmed.startsWith("#")) {
				continue;
			}
			const idx = trimmed.indexOf("=");
			if (idx !== -1 && trimmed.slice(0, idx).trim() === key) {
				return trimmed.slice(idx + 1).trim();
			}
		}
	} catch {
		// no .env file, fall through to default
	}
	return "";
}

// Must start with / and not end with / (empty string for root).
const basePath = readEnv("base_path").replace(/\/+$/, "");

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter({
			pages: "../backend/static",
			assets: "../backend/static",
			fallback: "index.html",
		}),
		paths: {
			// Set base_path in backend/.env when hosting under a subpath,
			// e.g. base_path=/pr for https://devops.kabal.com/pr/
			base: basePath,
		},
	},
};

export default config;
