import { TailwindElement } from './tailwind-element'
import { html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { client } from '../api'
import { repeat } from 'lit/directives/repeat.js'
import type { PaginatedPath } from '@openfish/client'

type NamedItem = { name: string } & Record<string, any>

// <data-select> is a select form element that will fetch data from
// the API and present them as options.

@customElement('data-select')
export class DataSelect extends TailwindElement {
  static formAssociated = true

  @property()
  name: string

  @property()
  src: PaginatedPath

  @property()
  pkey = 'id'

  @property()
  value: string

  @property()
  defaultText = 'Please select'

  @state()
  private _items: NamedItem[] = []

  private _internals: ElementInternals

  constructor() {
    super()

    this._internals = this.attachInternals()
  }

  async connectedCallback() {
    super.connectedCallback()

    const { data, error } = await client.GET(this.src, {
      params: { query: { limit: 100 } },
    })

    if (error !== undefined) {
      console.error(error)
    }
    if (data !== undefined) {
      this._items = data.results as NamedItem[]
    }
  }

  onInput(e: InputEvent & { target: HTMLSelectElement }) {
    this._internals.setFormValue(e.target.value)
    this.value = e.target.value
    this.dispatchEvent(new Event('input'))
  }

  render() {
    return html`
      <select
        @input=${this.onInput}
        .name=${this.name}
        .value=${this.value}
        class="w-full"
      >
        <option value="">${this.defaultText}</option>
        ${repeat(
          this._items,
          (item: NamedItem) => html`<option .value=${item[this.pkey]}>${item.name}</option>`
        )}
      </select>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'data-select': DataSelect
  }
}
