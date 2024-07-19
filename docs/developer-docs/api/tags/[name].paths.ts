import spec from '../../../public/swagger.json' assert { type: 'json' }
import { useOpenapi } from 'vitepress-theme-openapi'

export default {
    paths() {
        const openapi = useOpenapi({ spec })

        const json = openapi.json as typeof spec

        return json.tags.map(tag => ({params: tag}))

    },
}