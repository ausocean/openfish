import { defineConfig } from "vitepress";
import type { DefaultTheme } from "vitepress";
import { useOpenapi } from "vitepress-openapi/client";
import { useSidebar } from "vitepress-openapi";
import spec from '@openfish/client/swagger.json' with { type: 'json' }
import { vitepressDemoPlugin } from 'vitepress-demo-plugin'

useOpenapi({ spec});

const apiItems: DefaultTheme.SidebarItem[] = useSidebar({
	spec
})
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
  markdown: {
    config(md) {
      md.use(vitepressDemoPlugin, { locale: { 'en-US': 'en-US' } })
    },
  },
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
				{ text: 'Openfish Player Integration', link: '/developer-docs/player-integration' },
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
