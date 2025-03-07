import { TailwindElement } from './tailwind-element'
import { type TemplateResult, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { provide, consume, createContext } from '@lit/context'
import { formatAsDatetimeRange } from '../utils/datetime'

export type ClickRowEvent<T> = CustomEvent<T>

export type HoverItemEvent = CustomEvent<string | undefined>

export const dataContext = createContext(Symbol('table'))
export const pkeyContext = createContext(Symbol('pkey'))
export const hoverContext = createContext(Symbol('hover-row'))

@customElement('data-table')
export class DataTable<T extends Record<string, any>> extends TailwindElement {
  @state()
  protected _page = 1

  @state()
  @provide({ context: dataContext })
  protected _items: T[] = []

  @state()
  @provide({ context: hoverContext })
  protected _hover: string | undefined

  @state()
  protected _totalPages = 0

  @property()
  src: string

  @property()
  colwidths = ''

  @property()
  @provide({ context: pkeyContext })
  pkey = 'id'

  onHoverItem(e: HoverItemEvent) {
    this._hover = e.detail
  }

  async deleteItem(item: T) {
    const url = new URL(`./${item[this.pkey]}`, `${document.location.origin}${this.src}`)
    url.searchParams.forEach((_, key) => url.searchParams.delete(key))
    await fetch(url, { method: 'DELETE' })
    await this.fetchData()
  }

  async connectedCallback() {
    super.connectedCallback()
    await this.fetchData()
  }

  async fetchData() {
    const perPage = 10

    const url = new URL(this.src, document.location.origin)
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
    <span class="text-slate-800">Page ${this._page} of ${this._totalPages}</span>
    <span class="flex gap-1">
      <button class="btn variant-slate" @click="${this.next}" .disabled=${this._page === 1}>Prev</button>
      <button class="btn variant-slate" @click="${this.prev}" .disabled=${this._page === this._totalPages}>Next</button>
    </span>
    `

    return html`
    <div class="table border border-slate-300 rounded-md overflow-clip" style="--colwidths: ${this.colwidths}" @hoverItem=${this.onHoverItem}>
      <slot></slot>
      <footer class="flex px-4 gap-1 justify-between items-center h-12 bg-slate-100">
        ${pagination}
      </footer>
    </div>
    `
  }

  static styles = [
    TailwindElement.styles!,
    css`
      .table {
        display: grid;
        grid-template-columns: var(--colwidths, "");
      }
      
      footer {
        grid-column: 1/-1;
      }
    `,
  ]
}

abstract class DataTableColumn<T extends Record<string, any>> extends TailwindElement {
  @consume({ context: dataContext, subscribe: true })
  @state()
  protected _items: T[] = []

  @consume({ context: hoverContext, subscribe: true })
  @state()
  protected _hover: string | undefined

  @consume({ context: pkeyContext, subscribe: true })
  @state()
  protected _pkey: string

  @property()
  align: 'left' | 'right' | 'center' = 'left'

  hoverItem(key: string | undefined) {
    this.dispatchEvent(
      new CustomEvent('hoverItem', {
        detail: key,
        bubbles: true,
        composed: true,
      })
    )
  }

  clickItem(item: T) {
    this.dispatchEvent(
      new CustomEvent('clickitem', {
        detail: item,
        bubbles: true,
        composed: true,
      })
    )
  }

  abstract renderTitle(): TemplateResult
  abstract renderCell(item: T): TemplateResult

  render() {
    return html`
    <div class="flex items-center px-4 h-12 bg-slate-200 border-b border-b-slate-300 cursor-pointer font-bold text-slate-700" style="justify-content: ${this.align}">
      ${this.renderTitle()}
    </div>
    ${repeat(
      this._items,
      (item) => html`
      <div 
        class="td flex items-center px-4 h-12 border-b border-b-slate-300 cursor-pointer transition-colors ${this._hover === item[this._pkey].toString() ? 'hover' : ''}" 
        style="justify-content: ${this.align}" 
        @click=${() => this.clickItem(item)} 
        @mouseenter=${() => this.hoverItem(item[this._pkey].toString())}
        @mouseleave=${() => this.hoverItem(undefined)}
      >
        ${this.renderCell(item)}
      </div>
      `
    )}
    `
  }

  static styles = [
    TailwindElement.styles!,
    css`
      :host {
        width: 1fr
      }

      .td.hover {
        background-color: var(--color-blue-100);
        color: var(--color-blue-800);
      }
    `,
  ]
}

@customElement('dt-col')
export class DataTableTextColumn<T extends Record<string, any>> extends DataTableColumn<T> {
  @property()
  title: string

  @property()
  key: keyof T

  @property()
  align: 'left' | 'right' | 'center' = 'left'

  renderTitle(): TemplateResult {
    return html`${this.title}`
  }

  renderCell(item: T): TemplateResult {
    return html`${item[this.key]}`
  }
}

@customElement('dt-daterange-col')
export class DataTableDateColumn<T extends Record<string, any>> extends DataTableColumn<T> {
  @property()
  title: string

  @property()
  startKey: keyof T

  @property()
  endKey: keyof T

  @property()
  align: 'left' | 'right' | 'center' = 'left'

  renderTitle(): TemplateResult {
    return html`${this.title}`
  }

  renderCell(item: T): TemplateResult {
    return html`${formatAsDatetimeRange(item[this.startKey], item[this.endKey])}`
  }
}

@customElement('dt-btn')
export class DataTableButton<T extends Record<string, any>> extends DataTableColumn<T> {
  @property()
  text: string

  @property()
  action: string

  clickButton(item: T) {
    this.dispatchEvent(
      new CustomEvent(this.action, {
        detail: item,
        bubbles: true,
        composed: true,
      })
    )
  }

  renderTitle(): TemplateResult {
    return html``
  }

  renderCell(item: T): TemplateResult {
    return html`
    <button type="button" class="btn size-sm variant-slate" @click=${() => this.clickButton(item)}>${this.text}</button>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'data-table': DataTable<Record<string, any>>
    'dt-col': DataTableColumn<Record<string, any>>
    'dt-btn': DataTableButton<Record<string, any>>
  }
}
