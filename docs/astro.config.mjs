// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightThemeRapide from 'starlight-theme-rapide'
import starlightOpenAPI, { openAPISidebarGroups } from 'starlight-openapi'
import starlightImageZoom from 'starlight-image-zoom'

export default defineConfig({
	base: '/openfish',
	integrations: [
		starlight({
			title: "OpenFish Documentation",
			description: "OpenFish is an open-source system written in Golang for classifying marine species. Tasks involve importing video or image data, classifying and annotating data (both manually and automatically), searching, and more. ",
			plugins: [
				starlightThemeRapide(),
				starlightOpenAPI([
				  {
				   base: 'api',
				   label: 'API documentation',
				   schema: './src/assets/swagger.json',
				  },
				]),
				starlightImageZoom(),
			  ],
			social: {
				github: "https://github.com/ausocean/openfish",
			},
			logo: {
				src: '/src/assets/logo.webp',
			},
			
			customCss: [
				'/src/styles/custom.css',
			],
			sidebar: [
				{
					slug: 'project-overview'
				},
				{
					label: 'User Guide',
					autogenerate: { directory: 'user-guide' },
				},
				{
					label: 'Developer Documentation',
					items: [
						{slug: 'developer-docs/getting-started'},
						{slug: 'developer-docs/openfish-cli'},
						{slug: 'developer-docs/contributing'},
						...openAPISidebarGroups,
					]
				},
			],
		}),
	],
});
