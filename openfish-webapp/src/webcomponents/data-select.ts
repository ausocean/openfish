import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { Result } from '../api'
import { repeat } from 'lit/directives/repeat.js'

type NamedItem = { name: string } & Record<string, any>

// <data-select> is a select form element that will fetch data from
// the API and present them as options.

@customElement('data-select')
export class DataSelect<T extends NamedItem> extends LitElement {
  static formAssociated = true

  @property()
  name: string

  @property()
  src: string

  @property()
  pkey = 'id'

  @property()
  value: string

  @property()
  defaultText = 'Please select'

  @state()
  private _items: T[] = []

  private _internals: ElementInternals

  constructor() {
    super()

    this._internals = this.attachInternals()
  }

  async connectedCallback() {
    super.connectedCallback()

    try {
      const url = new URL(this.src, document.location.origin)
      url.searchParams.set('limit', String(999))
      const res = await fetch(url)
      const data = (await res.json()) as Result<T>
      this._items = data.results
    } catch (error) {
      console.error(error)
    }
  }

  onInput(e: InputEvent & { target: HTMLSelectElement }) {
    this._internals.setFormValue(e.target.value)
    this.value = e.target.value
    this.dispatchEvent(new Event('input'))
  }

  render() {
    return html`
    <select @input=${this.onInput} .name=${this.name} .value=${this.value}>
    <option value="">${this.defaultText}</option>
    ${repeat(
      this._items,
      (item: T) => html`<option .value=${item[this.pkey]}>${item.name}</option>`
    )}
    </select>
    `
  }

  static styles = css`
    select {
      width: 100%;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'data-select': DataSelect<NamedItem>
  }
}
