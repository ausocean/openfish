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

    try {
      const res = await fetch('/api/v1/auth/me')
      if (res.ok) {
        const json = await res.json()
        this.user = plainToInstance(User, json)
      }
      if (res.status === 404) {
        window.location.href = '/welcome.html'
      }
    } catch (error) {
      console.error(error)
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
