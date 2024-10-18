import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { CaptureSource, Result, VideoStream } from '../utils/api.types.ts'
import { repeat } from 'lit/directives/repeat.js'
import resetcss from '../styles/reset.css'
import { datetimeDifference, formatAsDate, formatDuration } from '../utils/datetime.ts'
import type { Filter } from './stream-filter.ts'
import { extractVideoID } from '../utils/youtube.ts'

type VideoStreamItem = Omit<VideoStream, 'capturesource'> & {
  capturesource: CaptureSource
  first?: boolean
}

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

    try {
      const params = new URLSearchParams()
      if (this._page === 1) {
        params.set('limit', String(perPageFirst))
        params.set('offset', String(0))
      } else {
        params.set('limit', String(perPage))
        params.set('offset', String(perPageFirst + (this._page - 2) * perPage))
      }

      for (const key in this._filter) {
        params.set(key, String(this._filter[key as keyof Filter]))
      }

      // Get videostreams and join on their capture source.
      const res = await fetch(`/api/v1/videostreams?${params.toString()}`)
      const data = (await res.json()) as Result<VideoStream>

      const promises: Promise<VideoStreamItem>[] = []
      for (const stream of data.results) {
        promises.push(
          (async (stream: VideoStream) => {
            const res = await fetch(`/api/v1/capturesources/${stream.capturesource}`)
            return {
              ...stream,
              capturesource: await res.json(),
            }
          })(stream)
        )
      }
      const items = await Promise.all(promises)
      if (this._page === 1 && items.length > 0) {
        items[0].first = true
      }
      this._items = items
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
    const streams = (stream: VideoStreamItem) => html`
    <article onclick="window.location = '/watch.html?id=${stream.id}'" class="${stream.first ? 'first-item' : ''}">
      <div class="img-container">
        <img src="https://i.ytimg.com/vi/${extractVideoID(stream.stream_url)}/maxresdefault.jpg">
        <span class="duration">${formatDuration(datetimeDifference(stream.endTime, stream.startTime))}</span>
      </div>
      <footer>
        <h3>${stream.capturesource.name}</h3>
        <p>${formatAsDate(stream.startTime)}</p>
      </footer>
    </article>
    `

    const pagination = html`   
    <span class="mr-1">Page ${this._page} of ${this._totalPages}</span>
    <button @click="${this.next}" .disabled=${this._page === 1}>Prev</button>
    <button @click="${this.prev}" .disabled=${this._page === this._totalPages}>Next</button>
    `

    return html`
    <main>
      ${repeat(this._items, streams)}
    </main>
    <footer class="pagination">
      ${pagination}
    </footer>
    `
  }

  static styles = css`
    ${resetcss}

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
      background-color: var(--gray-100);
      border: 1px solid var(--blue-400);
      overflow: clip;
      box-shadow: var(--shadow-soft-sm);
      transition: box-shadow;
      transition-duration: 200ms;

      &:hover {
        box-shadow: var(--shadow-soft-md);
      }
      
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

      .img-container {
        position: relative;
      }

      .duration {
          position: absolute;
          bottom: 0.5em;
          right: 0.5em;

          font-size: 0.8rem;
          background-color: rgba(0, 0, 0, 0.8);
          color: var(--bg);
          padding: 0 .5rem;
          border-radius: 999999px;
          white-space: pre-line;
      }
    }

    footer.pagination {
      display: flex;
      justify-content: center;
      gap: 0.25rem;
      padding: 1rem 0;
    }
    `
}

declare global {
  interface HTMLElementTagNameMap {
    'stream-list': StreamList
  }
}
