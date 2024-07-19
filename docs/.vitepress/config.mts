import { defineConfig } from "vitepress";

import { useSidebar, useOpenapi } from "vitepress-theme-openapi";
import spec from "../public/swagger.json" assert { type: "json" };
import type { DefaultTheme } from "vitepress";

const openapi = useOpenapi();
openapi.setSpec(spec);
const apiItems: DefaultTheme.SidebarItem[] = useSidebar()
	.generateSidebarGroups()
	.map((s: DefaultTheme.SidebarItem) => ({
		...s,
		collapsed: true,
		link: `/tags/${s.text}`,
	}));

// https://vitepress.dev/reference/site-config
export default defineConfig({
	base: "/openfish/",
	title: "OpenFish Documentation",
	description:
		"OpenFish is an open-source system written in Golang for classifying marine species. Tasks involve importing video or image data, classifying and annotating data (both manually and automatically), searching, and more. ",
	ignoreDeadLinks: [/^https?:\/\/localhost:5173/],
	themeConfig: {
		nav: [
			{ text: "Project Overview", link: "/project-overview" },
			{ text: "User Guide", link: "/user-guide/annotating-videos" },
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
					collapsed: true,
					link: "/",
					items: apiItems,
				},
			],
			"/user-guide/": [
				{ text: "Annotating Videos", link: "/user-guide/annotating-videos" },
			],
		},
		socialLinks: [
			{ icon: "github", link: "https://github.com/ausocean/openfish" },
		],
		search: {
			provider: "local",
		},
	},
});
