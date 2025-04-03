import { createClient } from '@openfish/client'

export const client = createClient()

export type Result<T> = {
  results: T[]
  offset: number
  limit: number
  total: number
}
