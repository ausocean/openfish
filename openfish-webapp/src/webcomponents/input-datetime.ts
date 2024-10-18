import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'

// An input element that wraps an <input type="datetime-local">
// but returns a date in RFC 3339 format, using the client's default
// timezone.

@customElement('input-datetime')
export class InputDatetime extends LitElement {
  static formAssociated = true

  @property()
  name: string

  @property()
  value = ''

  @property()
  required: boolean

  private _internals: ElementInternals

  onInput(e: InputEvent & { target: HTMLInputElement }) {
    const date = new Date(e.target.value)
    this._internals.setFormValue(date.toISOString())
  }

  constructor() {
    super()

    this._internals = this.attachInternals()
  }

  render() {
    return html`<input type="datetime-local" @input=${this.onInput} .name=${this.name} .value=${this.value} .required=${this.required}>`
  }

  static styles = css`
  :host {
    display: contents;
  }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'input-timezone': InputDatetime
  }
}
