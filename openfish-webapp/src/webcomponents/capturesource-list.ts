import { LitElement, css, html } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import type { CaptureSource, Result } from '../utils/api.types.ts'
import { repeat } from 'lit/directives/repeat.js'
import resetcss from '../styles/reset.css'
import btncss from '../styles/buttons.css'

export type DeleteItemEvent = CustomEvent<number>

@customElement('capturesource-list')
export class CaptureSourceList extends LitElement {
  @state()
  protected _page = 1

  @state()
  protected _items: CaptureSource[] = []

  @state()
  protected _totalPages = 0

  connectedCallback() {
    super.connectedCallback()
    this.fetchData()
  }

  async fetchData() {
    const perPage = 10

    try {
      const params = new URLSearchParams()
      params.set('limit', String(10))
      params.set('offset', String((this._page - 1) * perPage))

      const res = await fetch(`/api/v1/capturesources?${params.toString()}`)
      const data = (await res.json()) as Result<CaptureSource>
      this._items = data.results
      this._totalPages = Math.floor(data.total / perPage)
    } catch (error) {
      console.error(error)
    }
  }

  prev() {
    this._page += 1
    this.fetchData()
  }

  next() {
    this._page -= 1
    this.fetchData()
  }

  dispatchDeleteItemEvent(id: number) {
    this.dispatchEvent(new CustomEvent('deleteitem', { detail: id }) as DeleteItemEvent)
  }

  render() {
    const header = html`
    <tr>
      <th>Name</th>
      <th>Camera Hardware</th>
      <th>Site ID</th>
      <th>Location</th>
      <th></th>
    </tr>
    `

    const rows = (source: CaptureSource) => html`
    <tr>
      <td>${source.name}</td>
      <td>${source.camera_hardware}</td>
      <td>${source.site_id ?? '-'}</td>
      <td>${source.location}</td>
      <td><button class="btn btn-sm btn-secondary" @click=${() =>
        this.dispatchDeleteItemEvent(source.id)}>Delete</button></td>
    </tr>
    `

    const pagination = html`   
    <span class="mr-1">Page ${this._page} of ${this._totalPages}</span>
    <button @click="${this.next}" .disabled=${this._page === 1}>Prev</button>
    <button @click="${this.prev}" .disabled=${this._page === this._totalPages}>Next</button>
    `

    return html`
    <table>
    <thead>
      ${header}
    </thead>
    <tbody>
      ${repeat(this._items, rows)}
    </tbody>
    <tfoot>
      ${pagination}
    </tfoot>
    </table>
    `
  }

  static styles = css`
    ${resetcss}
    ${btncss}

    table {
      display: grid;  
      grid-template-columns: 1fr 1fr 1fr 1fr min-content;
      border-radius: 0.25rem;
      border: 1px solid var(--gray-100);
    }

    thead, tbody, tr {
      display: contents;
    }

    tbody td {
      padding: 0.5rem 1rem;
      border-bottom: 1px solid var(--gray-100);
    }
    tr {
      cursor: pointer
    }
    tr:hover td {
      background-color: var(--gray-50);
      color: var(--blue-700)
    }

    thead th {
      background-color: var(--gray-50);
      border-bottom: 1px solid var(--gray-100);
      padding: 0.5em 0;
    }

    tfoot {
      background-color: var(--gray-50);
      padding: 0.5em 0;
      grid-column: 1/-1;
      display: flex;
      justify-content: center;
      gap: 0.25rem
    }

    th {
      padding: 0.5rem
    }
    
    td {
      text-align: center
    }`
}

declare global {
  interface HTMLElementTagNameMap {
    'capturesource-list': CaptureSourceList
  }
}
