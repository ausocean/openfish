import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { Result } from '../utils/api.types.ts'
import { repeat } from 'lit/directives/repeat.js'

@customElement('data-select')
export class DataSelect<T extends Record<string, any>> extends LitElement {
  static formAssociated = true

  @property()
  name: string

  @property()
  src: string

  @property()
  pkey = 'id'

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
  }

  render() {
    return html`
    <select @input=${this.onInput} .name=${this.name}>
    <option .value=${null}>Any</option>
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
    'data-select': DataSelect<{ id: number }>
  }
}
