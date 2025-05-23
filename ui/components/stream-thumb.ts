import { TailwindElement } from './tailwind-element.ts'
import { html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { datetimeDifference, formatAsDate, formatDuration } from '../utils/datetime.ts'
import { extractVideoID } from '../utils/youtube.ts'
import type { VideoStreamWithJoins } from '@openfish/client'

@customElement('stream-thumb')
export class StreamThumb extends TailwindElement {
  @property()
  stream: VideoStreamWithJoins

  render() {
    return html`
      <article
        class="flex flex-col cursor-pointer overflow-clip rounded-md transition-shadow bg-slate-100 border border-blue-300 hover:border-blue-400 shadow-sm hover:shadow-lg hover:shadow-sky-500/25"
        onclick="window.location = '/watch?id=${this.stream.id}'"
      >
        <div class="relative">
          <img
            class="block w-full"
            src="https://i.ytimg.com/vi/${extractVideoID(this.stream.stream_url)}/maxresdefault.jpg"
          />
          <span class="absolute bottom-2 right-2 text-sm bg-slate-900/80 text-slate-50 px-2 rounded-full"
            >${formatDuration(datetimeDifference(this.stream.endTime, this.stream.startTime))}</span
          >
        </div>
        <footer class="flex flex-wrap text-sm p-2 *:w-full">
          <h3>${this.stream.capturesource.name}</h3>
          <p>${formatAsDate(this.stream.startTime)}</p>
        </footer>
      </article>
    `
  }

  static styles = TailwindElement.styles
}

declare global {
  interface HTMLElementTagNameMap {
    'stream-thumb': StreamThumb
  }
}
