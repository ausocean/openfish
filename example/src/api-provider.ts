import { useApiProvider } from '@openfish/ui/components/api-provider.ts'
import { createClient } from '@openfish/client'

const client = createClient({ baseUrl: 'http://localhost:8080/' })

useApiProvider(client)
