import spec from '../swagger.json' with { type: 'json' }
import { useOpenapi } from 'vitepress-openapi'

export default {
    paths() {
        const openapi = useOpenapi({ spec })

        const json = openapi.json as typeof spec

        return json.tags.map(tag => ({params: tag}))

    },
}