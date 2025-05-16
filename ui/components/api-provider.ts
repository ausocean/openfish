import { LitElement, css, html } from 'lit'
import { provide } from '@lit/context'
import { clientContext } from '../utils/context'
import type { OpenfishClient } from '@openfish/client'

export function useApiProvider(client: OpenfishClient) {
  class ApiProvider extends LitElement {
    @provide({ context: clientContext })
    client: OpenfishClient = client

    render() {
      return html`<slot></slot>`
    }

    static styles = css`
      :host {
        display: contents;
      }
    `
  }

  customElements.define('api-provider', ApiProvider)
}
