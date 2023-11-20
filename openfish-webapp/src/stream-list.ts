import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { Result, VideoStream } from './api.types.ts'
import { repeat } from 'lit/directives/repeat.js'
import { resetcss } from './reset.css.ts'
import { datetimeDifference, formatAsDatetime, formatDuration } from './datetime.ts'

@customElement('stream-list')
export class StreamList extends LitElement {
  @property({ type: Number })
  page = 1

  @property({ type: Array })
  items: VideoStream[] = []

  @property({ type: Number })
  totalPages = 0

  attributeChangedCallback() {
    this.fetchData(this.page)
  }

  connectedCallback() {
    super.connectedCallback()
    this.fetchData(this.page)
  }

  async fetchData(page: number) {
    const perPage = 10

    try {
      const res = await fetch(
        `http://localhost:3000/api/v1/videostreams?limit=10&offset=${(page - 1) * perPage}`
      )
      const data = (await res.json()) as Result<VideoStream>
      this.items = data.results
      this.totalPages = Math.floor(data.total / perPage)
      console.log(this.totalPages)
    } catch (error) {
      console.error(error)
    }
  }

  prev() {
    this.page += 1
    this.fetchData(this.page)
  }

  next() {
    this.page -= 1
    this.fetchData(this.page)
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
    <span class="mr-1">Page ${this.page} of ${this.totalPages}</span>
    <button @click="${this.next}" .disabled=${this.page === 1}>Prev</button>
    <button @click="${this.prev}" .disabled=${this.page === this.totalPages}>Next</button>
    `

    return html`
    <table>
    <thead>
      ${header}
    </thead>
    <tbody>
      ${repeat(this.items, rows)}
    </tbody>
    <tfoot>
      ${pagination}
    </tfoot>
    </table>
    `
  }

  static styles = css`
    ${resetcss}

    :host {
      width: min(100vw, 80rem);
    }

    table {
      display: grid;  
      grid-template-columns: 1fr 20rem 10rem;
      border-radius: 0.25rem;
      border: 1px solid var(--gray1);
    }

    thead, tbody, tr {
      display: contents;
    }

    tbody td {
      padding: 0.5rem 1rem;
      border-bottom: 1px solid var(--gray1);
    }
    tr {
      cursor: pointer
    }
    tr:hover td {
      background-color: var(--gray0);
      color: var(--primary)
    }

    thead th {
      background-color: var(--gray0);
      border-bottom: 1px solid var(--gray1);
      padding: 0.5em 0;
    }

     tfoot {
      background-color: var(--gray0);
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
