import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import resetcss from '../styles/reset.css?raw'
import btncss from '../styles/buttons.css?raw'
import { repeat } from 'lit/directives/repeat.js'
import { provide, consume, createContext } from '@lit/context'

export type ClickRowEvent<T> = CustomEvent<T>

export type HoverItemEvent = CustomEvent<number | undefined>

export const dataContext = createContext(Symbol('table'))
export const hoverContext = createContext(Symbol('hover-row'))

@customElement('data-table')
export class DataTable<T extends { id: number }> extends LitElement {
  @state()
  protected _page = 1

  @state()
  @provide({ context: dataContext })
  protected _items: T[] = []

  @state()
  @provide({ context: hoverContext })
  protected _hover: number | undefined

  @state()
  protected _totalPages = 0

  @property()
  src: string

  @property()
  colwidths = ''

  onHoverItem(e: HoverItemEvent) {
    this._hover = e.detail
  }

  async connectedCallback() {
    super.connectedCallback()
    await this.fetchData()
  }

  async fetchData() {
    const perPage = 10

    const url = new URL(`${import.meta.env.VITE_API_HOST}${this.src}`)
    url.searchParams.set('limit', String(perPage))
    url.searchParams.set('offset', String((this._page - 1) * perPage))

    try {
      const res = await fetch(url)
      const data = await res.json()
      this._items = data.results
      this._totalPages = Math.floor(data.total / perPage) + 1
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

  render() {
    const pagination = html`   
    <span class="mr-1">Page ${this._page} of ${this._totalPages}</span>
    <button @click="${this.next}" .disabled=${this._page === 1}>Prev</button>
    <button @click="${this.prev}" .disabled=${this._page === this._totalPages}>Next</button>
    `

    return html`
    <div class="table" style="--colwidths: ${this.colwidths}" @hoverItem=${this.onHoverItem}>
      <slot></slot>
      <footer>
        ${pagination}
      </footer>
    </div>
    `
  }

  static styles = css`
    ${unsafeCSS(resetcss)}
    ${unsafeCSS(btncss)}

    .table {
      display: grid;
      grid-template-columns: var(--colwidths, "");
      border-radius: 0.25rem;
      border: 1px solid var(--gray-100);
    }

    footer {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 1rem;
      height: 3rem;
      background-color: var(--gray-50);
      grid-column: 1/-1;
      gap: 0.25rem
    }
    `
}

@customElement('dt-col')
export class DataTableColumn<T extends { id: number }> extends LitElement {
  @consume({ context: dataContext, subscribe: true })
  @state()
  protected _items: T[] = []

  @consume({ context: hoverContext, subscribe: true })
  @state()
  protected _hover: number | undefined

  @property()
  title: string

  @property()
  key: keyof T

  @property()
  align: 'left' | 'right' | 'center' = 'left'

  hoverItem(id: number | undefined) {
    this.dispatchEvent(
      new CustomEvent('hoverItem', {
        detail: id,
        bubbles: true,

        composed: true,
      })
    )
  }

  render() {
    return html`
    <div class="th" style="text-align: ${this.align}">
        ${this.title}
    </div>
    ${repeat(
      this._items,
      (item) => html`
      <div class="td ${this._hover === item.id ? 'hover' : ''}" style="text-align: ${this.align}" @mouseenter=${() => this.hoverItem(item.id)} @mouseleave=${() => this.hoverItem(undefined)}>${item[this.key]}</div>
        `
    )}
    `
  }

  static styles = css`
    ${unsafeCSS(resetcss)}
    
    :host {
      width: 1fr
    }

    .th {
      display: flex;
      align-items: center;
      padding: 0 1rem;
      height: 3rem;
      font-weight: bold;
      background-color: var(--gray-50);
      border-bottom: 1px solid var(--gray-100);
    }

    .td {
      display: flex;
      align-items: center;
      padding: 0 1rem;
      height: 3rem;
      border-bottom: 1px solid var(--gray-100);
      cursor: pointer;

      &.hover {
        background-color: var(--gray-50);
        color: var(--blue-700);
      }
    }
    `
}

declare global {
  interface HTMLElementTagNameMap {
    'data-table': DataTable<{ id: number }>
    'dt-col': DataTableColumn<{ id: number }>
  }
}
