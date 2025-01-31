import DefaultTheme from 'vitepress/theme'
import type { Theme } from 'vitepress'

import { theme, useOpenapi } from 'vitepress-openapi/client'
import 'vitepress-openapi/dist/style.css'

import spec from '../../developer-docs/api/swagger.json' with { type: 'json' }

export default {
    ...DefaultTheme,
    async enhanceApp({app, router, siteData}) {
        // Set the OpenAPI specification.
        const openapi = useOpenapi({
            spec
        })

        // Use the theme.
        theme.enhanceApp({app, openapi})
    }
} satisfies Theme