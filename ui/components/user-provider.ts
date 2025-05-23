import { LitElement, css, html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { consume, provide } from '@lit/context'
import { clientContext, userContext } from '../utils/context'
import type { OpenfishClient, User } from '@openfish/client'

@customElement('user-provider')
export class UserProvider extends LitElement {
  @provide({ context: userContext })
  user: User | null = null

  @consume({ context: clientContext, subscribe: true })
  client!: OpenfishClient

  async connectedCallback() {
    super.connectedCallback()

    const { data, error, response } = await this.client.GET('/api/v1/auth/me')

    if (response.status === 404) {
      window.location.href = '/welcome'
    }

    if (error !== undefined) {
      console.error(error)
    }

    if (data !== undefined) {
      this.user = data
    }
  }

  render() {
    return html`<slot></slot>`
  }

  static styles = css`
    :host {
      display: contents;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'user-provider': UserProvider
  }
}
