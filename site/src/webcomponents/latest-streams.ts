import { TailwindElement } from './tailwind-element'
import { html } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { client } from '../api'
import type { VideoStreamWithJoins } from '@openfish/client'
import './stream-thumb'

@customElement('latest-streams')
export class LatestStreams extends TailwindElement {
  @state()
  protected _items: VideoStreamWithJoins[] = []

  connectedCallback() {
    super.connectedCallback()
    this.fetchData()
  }

  async fetchData() {
    const { data, error } = await client.GET('/api/v1/videostreams', {
      params: {
        query: {
          limit: 7,
          offset: 0,
        },
      },
    })

    if (error !== undefined) {
      console.error(error)
    }

    if (data !== undefined) {
      this._items = data.results
    }
  }

  render() {
    return html`
      <div class="h-full grid grid-cols-3 lg:grid-cols-4 gap-4 justify-start items-center">
        <div class="min-w-1/4 space-y-4">
            <h2 class="text-3xl font-bold text-blue-50 mt-8">Latest streams</h2>
            <a class="btn variant-slate" href="/streams">View More</a>
        </div>
        ${repeat(this._items, (item) => html`<stream-thumb class="aspect-[4/3] max-h-64 transition-transform duration-300 ease-in-out hover:-translate-y-1" .stream=${item}></stream-thumb>`)}
      </div>
    `
  }

  static styles = TailwindElement.styles
}

declare global {
  interface HTMLElementTagNameMap {
    'latest-streams': LatestStreams
  }
}
