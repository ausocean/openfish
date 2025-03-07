import { useOpenapi } from 'vitepress-openapi/client'
import spec from '@openfish/client/swagger.json' with { type: 'json' }

export default {
    paths() {
        const openapi = useOpenapi({ spec })

        const json = openapi.json as typeof spec

        return json.tags.map(tag => ({params: tag}))

    },
}