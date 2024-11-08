import DefaultTheme from 'vitepress/theme'
import type { Theme } from 'vitepress'

import { theme, useOpenapi, useTheme } from 'vitepress-openapi'
import 'vitepress-openapi/dist/style.css'

import spec from '../../developer-docs/api/swagger.json' with { type: 'json' }

export default {
    ...DefaultTheme,
    async enhanceApp({app, router, siteData}) {
        // Set the OpenAPI specification.
        const openapi = useOpenapi()
        openapi.setSpec(spec)

        // Optionally, configure the theme.
        const themeConfig = useTheme()
        themeConfig.setLocale('en')

        // Use the theme.
        theme.enhanceApp({app})
    }
} satisfies Theme