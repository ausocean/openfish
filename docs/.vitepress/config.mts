import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
	base: "/openfish/",
	title: "OpenFish Documentation",
	description:
		"OpenFish is an open-source system written in Golang for classifying marine species. Tasks involve importing video or image data, classifying and annotating data (both manually and automatically), searching, and more. ",
	ignoreDeadLinks: [/^https?:\/\/localhost:5173/],
	themeConfig: {
		// https://vitepress.dev/reference/default-theme-config
		nav: [
			{ text: "Project Overview", link: "/project-overview" },
			{ text: "User Guide", link: "/user-guide/todo" },
			{ text: "Developer Documentation", link: "/developer-docs/getting-started" },
			{ text: "AusOcean", link: "https://www.ausocean.org/" },
		],

		sidebar: {
			"/developer-docs/": [
				{ text: "Getting Started", link: "/developer-docs/getting-started" },
				{ text: "Contributing", link: "/developer-docs/contributing" },
				{ text: "OpenFish CLI", link: "/developer-docs/openfish-cli" },
				{
					text: "API Documentation",
					base: "/developer-docs/api",
					items: [
						{ text: "Introduction", link: "/" },
						{ text: "Authentication", link: "/authentication" },
						{ text: "Roles and Permissions", link: "/roles-and-permissions" },
						{
							text: "API Usage",
							link: "/api-usage",
							items: [
								{ text: "Capture Sources", link: "/capture-sources" },
								{ text: "Video Streams", link: "/video-streams" },
								{ text: "Annotations", link: "/annotations" },
								{ text: "Species", link: "/species" },
								{ text: "Users", link: "/users" },
							],
						},
					],
				},
			],
			"/user-guide/": [],
		},

		socialLinks: [
			{ icon: "github", link: "https://github.com/ausocean/openfish" },
		],

		search: {
			provider: "local",
		},
	},
});
