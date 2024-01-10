import { LitElement, css, html } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { Result, CaptureSource } from './api.types.ts'
import { repeat } from 'lit/directives/repeat.js'

export type SelectCaptureSourceEvent = CustomEvent<number | null>

@customElement('capture-source-dropdown')
export class CaptureSourceDropdown extends LitElement {
  @state()
  _items: CaptureSource[] = []

  connectedCallback() {
    super.connectedCallback()
    this.fetchData()
  }

  async fetchData() {
    try {
      const res = await fetch('http://localhost:3000/api/v1/capturesources?limit=999')
      const data = (await res.json()) as Result<CaptureSource>
      this._items = data.results
    } catch (error) {
      console.error(error)
    }
  }

  onSelectCaptureSource(event: InputEvent & { target: HTMLSelectElement }) {
    const options = {
      detail: event.target.value === 'null' ? null : Number(event.target.value),
      bubbles: true,
      composed: true,
    }
    this.dispatchEvent(new CustomEvent('selectcapturesource', options))
  }

  render() {
    return html`
    <select @input=${this.onSelectCaptureSource}>
    <option .value=${null}>Any</option>
    ${repeat(
      this._items,
      (item: CaptureSource) => html`<option .value=${item.id}>${item.name}</option>`
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
    'capture-source-dropdown': CaptureSourceDropdown
  }
}
