import { createContext } from '@lit/context'
import type { User } from '../api/user'

export const userContext = createContext<User | null>(Symbol('current-user'))
