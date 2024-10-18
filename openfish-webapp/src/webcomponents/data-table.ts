import { LitElement, type TemplateResult, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import resetcss from '../styles/reset.css'
import btncss from '../styles/buttons.css'
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

  async deleteItem(item: T) {
    await fetch(`${this.src}/${item.id}`, { method: 'DELETE' })
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
    ${resetcss}
    ${btncss}

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

abstract class DataTableColumn<T extends { id: number }> extends LitElement {
  @consume({ context: dataContext, subscribe: true })
  @state()
  protected _items: T[] = []

  @consume({ context: hoverContext, subscribe: true })
  @state()
  protected _hover: number | undefined

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
    <div class="th" style="text-align: ${this.align}">
      ${this.renderTitle()}
    </div>
    ${repeat(
      this._items,
      (item) => html`
      <div 
        class="td ${this._hover === item.id ? 'hover' : ''}" 
        style="text-align: ${this.align}" 
        @click=${() => this.clickItem(item)} 
        @mouseenter=${() => this.hoverItem(item.id)} 
        @mouseleave=${() => this.hoverItem(undefined)}
      >
        ${this.renderCell(item)}
      </div>
      `
    )}
    `
  }

  static styles = css`
    ${resetcss}
    
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

@customElement('dt-col')
export class DataTableTextColumn<T extends { id: number }> extends DataTableColumn<T> {
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

@customElement('dt-btn')
export class DataTableButton<T extends { id: number }> extends DataTableColumn<T> {
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
    <button type="button" class="btn btn-sm btn-secondary" @click=${() => this.clickButton(item)}>${this.text}</button>
    `
  }

  static styles = css`
    ${DataTableColumn.styles}
    ${btncss}
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'data-table': DataTable<{ id: number }>
    'dt-col': DataTableTextColumn<{ id: number }>
    'dt-btn': DataTableButton<{ id: number }>
  }
}
