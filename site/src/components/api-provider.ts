import { useApiProvider } from '@openfish/ui/components/api-provider.ts'
import { createClient } from '@openfish/client'

export const client = createClient()
useApiProvider(client)
