import { sveltekit } from "@sveltejs/kit/vite";
import adapter from "@sveltejs/adapter-static";
import { readFileSync } from "node:fs";
import { defineConfig } from "vite";

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

const basePath = readEnv("base_path").replace(/\/+$/, "");

export default defineConfig({
	plugins: [
		sveltekit({
			alias: {
				$lib: "src/lib",
			},
			adapter: adapter({
				pages: "../backend/static",
				assets: "../backend/static",
				fallback: "index.html",
			}),
			paths: {
				// Set base_path in backend/.env when hosting under a subpath,
				// e.g. base_path=/pr for https://devops.example.com/pr/.
				// Use an empty value for root hosting.
				base: basePath,
			},
		}),
	],
	server: {
		port: 5173,
		proxy: {
			"/api": {
				target: "http://localhost:8080",
				changeOrigin: true,
			},
		},
	},
	preview: {
		port: 5173,
		proxy: {
			"/api": {
				target: "http://localhost:8080",
				changeOrigin: true,
			},
		},
	},
});
