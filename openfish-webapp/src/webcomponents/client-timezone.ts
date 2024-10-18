import { LitElement, css, html } from 'lit'
import { customElement } from 'lit/decorators.js'

@customElement('client-timezone')
export class ClientTimezone extends LitElement {
  value = Intl.DateTimeFormat().resolvedOptions().timeZone

  render() {
    return html`(${this.value})`
  }

  static styles = css`
  :host {
    display: contents;
    color: var(--gray-800);
  }      
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'client-timezone': ClientTimezone
  }
}
