import { createContext } from '@lit/context'
import type { OpenfishClient, User } from '@openfish/client'

export const userContext = createContext<User | null>(Symbol('current-user'))
export const clientContext = createContext<OpenfishClient>(Symbol('openfish-client'))
