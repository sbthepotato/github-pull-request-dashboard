// place files you want to import through the `$lib` alias in this folder.
import { goto } from "$app/navigation";

const goto_options = {
	keepFocus: true,
	noScroll: false,
	replaceState: true,
};

const date_options = {
	month: "short",
	day: "numeric",
	hour: "numeric",
	minute: "numeric",
	hour12: false,
};

export function setUrlParam(param_name, param_value) {
	if (typeof window !== "undefined") {
		const url = new URL(window.location.href);
		if (
			param_value !== null &&
			param_value !== undefined &&
			param_value !== ""
		) {
			url.searchParams.set(param_name, param_value);
		} else {
			url.searchParams.delete(param_name);
		}

		goto(url.pathname + url.search, goto_options);
	}
}

export function setManyUrlParams(params) {
	if (typeof window !== "undefined") {
		const url = new URL(window.location.href);

		Object.entries(params).forEach(([param_name, param_value]) => {
			if (
				param_value !== null &&
				param_value !== undefined &&
				param_value !== ""
			) {
				url.searchParams.set(param_name, param_value);
			} else {
				url.searchParams.delete(param_name);
			}
		});

		goto(url.pathname + url.search, goto_options);
	}
}

export function getTextLuminance(hexColor) {
	// Remove the "#" if it exists
	const color = hexColor.replace("#", "");

	// Convert hex to RGB
	const r = parseInt(color.substring(0, 2), 16);
	const g = parseInt(color.substring(2, 4), 16);
	const b = parseInt(color.substring(4, 6), 16);

	// Calculate the luminance of the color
	const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;

	// Return white for dark colors and black for light colors
	return luminance > 0.5 ? "#0d1117" : "#f0f6fc";
}

export function stringToBool(string, otherwise = false) {
	if (string === "y") {
		return true;
	} else if (string === "n") {
		return false;
	} else {
		return otherwise;
	}
}

export function boolToString(bool) {
	if (bool) {
		return "y";
	}
	return "n";
}

export function getPrettyDate(date_str) {
	const date = new Date(date_str);
	return date.toLocaleString("en-us", date_options);
}

export function addLinks(text, replaceList) {
	let result = text;

	replaceList.forEach((item) => {
		result = result.replace(new RegExp(item.pattern, "g"), (match, p1) => {
			return `<a href="${`${item.link}${p1}`}" target="_blank">${match}</a>`;
		});
	});

	return result;
}
