import { TailwindElement } from './tailwind-element'
import { css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { datetimeDifference, formatAsDate, formatDuration } from '../utils/datetime.ts'
import type { Filter } from './stream-filter.ts'
import { extractVideoID } from '../utils/youtube.ts'
import { client } from '../api'
import type { components } from '@openfish/client/schema.ts'

type VideoStreamItem = components['schemas']['services.VideoStreamWithJoins'] & {
  first?: boolean
}

@customElement('stream-list')
export class StreamList extends TailwindElement {
  @property({ type: Object })
  set filter(val: Filter) {
    this._filter = val
    this.fetchData()
  }

  @state()
  protected _page = 1

  @state()
  protected _items: VideoStreamItem[] = []

  @state()
  protected _totalPages = 0

  @state()
  protected _filter: Filter = {}

  connectedCallback() {
    super.connectedCallback()
    this.fetchData()
  }

  async fetchData() {
    // First page has 12 items, the rest have 15 because of this particular layout.
    const perPage = 15
    const perPageFirst = 12

    let limit = perPageFirst
    let offset = 0
    if (this._page !== 1) {
      limit = perPage
      offset = (this._page - 2) * perPage + perPageFirst
    }

    const { data, error } = await client.GET('/api/v1/videostreams', {
      params: {
        query: {
          limit,
          offset,
          capturesource: this._filter.captureSource,
        },
      },
    })

    if (error !== undefined) {
      console.error(error)
    }

    if (data !== undefined) {
      this._items = data.results
      if (this._page === 1 && data.results.length > 0) {
        this._items[0].first = true
      }
      this._totalPages = Math.floor(data.total / perPage)
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
    const streams = (stream: VideoStreamItem) => html`
      <article
        class="border border-blue-300 hover:border-blue-400 shadow-sm hover:shadow-lg hover:shadow-sky-500/25 ${
          stream.first ? 'first-item' : ''
        }"
        onclick="window.location = '/watch.html?id=${stream.id}'"
      >
        <div class="img-contain">
          <img
            src="https://i.ytimg.com/vi/${extractVideoID(stream.stream_url)}/maxresdefault.jpg"
          />
          <span class="duration"
            >${formatDuration(datetimeDifference(stream.endTime, stream.startTime))}</span
          >
        </div>
        <footer>
          <h3>${stream.capturesource.name}</h3>
          <p>${formatAsDate(stream.startTime)}</p>
        </footer>
      </article>
    `

    const pagination = html`
      <span>Page ${this._page} of ${this._totalPages}</span>
      <span class="flex gap-1">
        <button
          class="btn variant-slate"
          @click="${this.next}"
          .disabled=${this._page === 1}
        >
          Prev
        </button>
        <button
          class="btn variant-slate"
          @click="${this.prev}"
          .disabled=${this._page === this._totalPages}
        >
          Next
        </button>
      </span>
    `

    return html`
      <main>${repeat(this._items, streams)}</main>
      <footer
        class="flex gap-1 pt-4 justify-between items-baseline border-t border-slate-300"
      >
        ${pagination}
      </footer>
    `
  }

  static styles = [
    TailwindElement.styles!,
    css`
      main {
        display: grid;
        grid-template-columns: repeat(5, 1fr);
        grid-template-rows: repeat(3, 1fr);
        gap: 1rem;
      }

      .first-item {
        grid-row: span 2;
        grid-column: span 2;

        footer {
          padding: 1rem;
        }
        h3 {
          font-size: 2rem;
        }
      }

      article {
        display: flex;
        flex-direction: column;
        cursor: pointer;
        border-radius: 0.5rem;
        background-color: var(--color-slate-100);
        overflow: clip;
        transition: box-shadow;
        transition-duration: 200ms;

        footer {
          padding: 0.5rem;
          font-size: 0.8rem;
          display: flex;
          flex-wrap: wrap;

          * {
            width: 100%;
          }
        }
        img {
          display: block;
          width: 100%;
        }

        .img-contain {
          position: relative;
        }

        .duration {
          position: absolute;
          bottom: 0.5em;
          right: 0.5em;

          font-size: 0.8rem;
          background-color: rgba(0, 0, 0, 0.8);
          color: var(--bg);
          padding: 0 0.5rem;
          border-radius: 999999px;
          white-space: pre-line;
        }
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'stream-list': StreamList
  }
}
