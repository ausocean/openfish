import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { Result, VideoStream } from '../utils/api.types.ts'
import { repeat } from 'lit/directives/repeat.js'
import resetcss from '../styles/reset.css'
import { datetimeDifference, formatAsDatetime, formatDuration } from '../utils/datetime.ts'
import type { Filter } from './stream-filter.ts'

@customElement('stream-list')
export class StreamList extends LitElement {
  @property({ type: Object })
  set filter(val: Filter) {
    this._filter = val
    this.fetchData()
  }

  @state()
  protected _page = 1

  @state()
  protected _items: VideoStream[] = []

  @state()
  protected _totalPages = 0

  @state()
  protected _filter: Filter = {}

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

      for (const key in this._filter) {
        params.set(key, String(this._filter[key as keyof Filter]))
      }

      const res = await fetch(`/api/v1/videostreams?${params.toString()}`)
      const data = (await res.json()) as Result<VideoStream>
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

  render() {
    const header = html`
    <tr>
      <th>Stream URL</th>
      <th>Date & Time</th>
      <th>Duration</th>
    </tr>
    `

    const rows = (stream: VideoStream) => html`
    <tr onclick="window.location = '/watch.html?id=${stream.id}'">
      <td>${stream.stream_url}</td>
      <td>${formatAsDatetime(stream.startTime)}</td>
      <td>${formatDuration(datetimeDifference(stream.endTime, stream.startTime))}</td>
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

    table {
      display: grid;  
      grid-template-columns: 1fr 20rem 10rem;
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
      grid-column: 1/ span 3;
      display: flex;
      justify-content: center;
      gap: 0.25rem
    }

    th {
      padding: 0.5rem
    }
    
    .mr-1 {
      margin-right: 1rem;
    }`
}

declare global {
  interface HTMLElementTagNameMap {
    'stream-list': StreamList
  }
}
