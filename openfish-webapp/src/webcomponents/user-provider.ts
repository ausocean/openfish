import { LitElement, css, html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { User } from '../api/user'
import { plainToInstance } from 'class-transformer'
import { provide } from '@lit/context'
import { userContext } from '../utils/context'

@customElement('user-provider')
export class UserProvider extends LitElement {
  @provide({ context: userContext })
  user: User | null = null

  async connectedCallback() {
    super.connectedCallback()

    if (import.meta.env.DEV) {
      this.user = plainToInstance(User, { email: 'user@localhost', role: 'admin' })
    } else {
      try {
        const res = await fetch('/api/v1/auth/me')
        const json = await res.json()
        this.user = plainToInstance(User, json)
      } catch (error) {
        console.error(error)
      }
    }
  }

  render() {
    return html`<slot></slot>`
  }

  static styles = css`
  :host {
    display: contents
  }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'user-provider': UserProvider
  }
}
