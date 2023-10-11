import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { Result, VideoStream } from '../utils/api.types'
import { repeat } from 'lit/directives/repeat.js'
import { resetcss } from '../css/reset.css'

@customElement('stream-list')
export class StreamList extends LitElement {
  constructor() {
    super()
    this.fetch()
  }

  async fetch() {
    const res = await fetch(
      `http://localhost:3000/api/v1/videostreams?offset=${this.offset}&limit=${this.limit}`
    )
    this.page = (await res.json()) as Result<VideoStream>
  }

  @property()
  page: Result<VideoStream> = {
    results: [],
    offset: 0,
    limit: 20,
    total: 0,
  }

  @property()
  offset: number = 0

  limit: number = 20

  next() {
    this.offset += this.limit
  }

  prev() {
    this.offset -= this.limit
  }

  row(videostream: VideoStream) {
    return html`
    <tr @click=${() => (window.location.href = `/watch.html?id=${videostream.id}`)}>
      <td>${videostream.startTime}</td>
      <td>${videostream.endTime}</td>
      <td>${videostream.stream_url}</td>
    </tr>
    `
  }

  paginate() {
    return html`
    <footer>
    <div>Page ${this.page.offset / this.page.limit + 1}</div>
    <div>
    <button @click=${this.prev()} ?disabled=${this.page.offset == 0}>Previous</button>
    <button @click=${this.next()} ?disabled=${this.page.results.length < this.page.limit}>Next</button>
    </div>
    </footer>
    `
  }

  render() {
    return html`
    <main>
      <h1>Video Streams</h1>
      <table>
        <thead>
          <tr>
            <th>Start Time</th>
            <th>End Time</th>
            <th>Stream URL</th>
          </tr>
        </thead>
        <tbody>
          ${repeat(this.page?.results ?? [], this.row)}
        </tbody>
      </table>
      ${this.paginate()}
    </main>
    `
  }

  static styles = css`
    ${resetcss}

    table {
      display: grid;
      min-width: 100%;
      grid-template-columns: 
      minmax(10rem, 1fr)
      minmax(10rem, 1fr)
      minmax(20rem, 2fr);
      border: 1px solid var(--gray2);
    }

    thead,
    tbody,
    thead tr {
      display: contents;
    }

    tbody tr {
      grid-column: 1 / 4;
      display: grid;
      grid-template-columns: subgrid;
      grid-template-rows: subgrid;
    }

    th {
      position: sticky;
      top: 0;
    }

    th, td {
      padding: 1rem;
    }

    tr {
      background-color: var(--gray0);
      border: 1px solid transparent;
      transition: background-color 200ms ease-in-out, border-color 200ms ease-in-out;
    }
    
    tr:hover {
      background-color: var(--gray1);
      border: 1px solid var(--secondary);
    }

    tr:nth-child(even) {
      background-color: var(--bg);
    }

    tr:nth-child(even):hover {
      background-color: var(--gray0);
    }


    th {
      background-color: var(--gray1);
      border-bottom: 1px solid var(--gray2);
    }

    footer {
      padding: 1rem;
      display: flex;
      width: 100%;
      justify-content: space-between;
    }

    main {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 0 12rem;
    }

    h1 {
      width: 100%;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'stream-list': StreamList
  }
}
