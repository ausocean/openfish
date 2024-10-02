import { useOpenapi, httpVerbs } from 'vitepress-theme-openapi'
import spec from '../swagger.json' assert { type: 'json' }

export default {
    paths() {

        const openapi = useOpenapi({ spec })

        const json = openapi.json as typeof spec

        type Params = { operationId: string }
        const results: { params: Params }[] = []

        for (const path in json.paths) {
            for (const method in json.paths[path]) {
                const { operationId } = json.paths[path][method]
                    results.push({
                        params: {
                            operationId,
                        },
                    })
            }
        }

        return results
    },
}