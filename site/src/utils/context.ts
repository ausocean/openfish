import { createContext } from '@lit/context'
import type { User } from '@openfish/client'

export const userContext = createContext<User | null>(Symbol('current-user'))
