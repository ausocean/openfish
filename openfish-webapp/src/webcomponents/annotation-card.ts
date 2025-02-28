import { TailwindElement } from './tailwind-element'
import { css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { formatVideoTime } from '../utils/datetime.ts'
import type { Annotation } from '../api/annotation.ts'
import type { SeekEvent } from './playback-controls.ts'

@customElement('annotation-card')
export class AnnotationCard extends TailwindElement {
  @property({ type: Object })
  annotation: Annotation | undefined

  @property({ type: Boolean })
  outline = false

  @property({ type: Boolean })
  glow = false

  dispatchSeekEvent(time: number) {
    this.dispatchEvent(
      new CustomEvent('seek', { detail: time, bubbles: true, composed: true }) as SeekEvent
    )
  }

  render() {
    if (this.annotation === undefined) {
      return html`<div class="card"></div>`
    }

    const start = formatVideoTime(this.annotation.start)
    const end = formatVideoTime(this.annotation.end)

    const { common_name, species, ...rest } = this.annotation.observation

    const rows = repeat(
      Object.entries(rest),
      ([key, val]) => html`
      <tr>
      <td>${key}</td>
      <td>${val}</td>
      </tr>
    `
    )

    const table = html`
      <table class="cols-2">
        <thead>
            <tr>
                <th>Property</th>
                <th>Value</th>
            </tr>
        </thead>
        <tbody>
        ${rows}
        </tbody>
      </table>
    `

    return html`
    <article class="card grid gap-2 ${this.glow ? 'shadow-md shadow-sky-500/50 border-sky-500' : ''}">
      <span class="font-bold">${common_name}</span>
      <span class="bg-slate-200 rounded-sm text-sm py-0.5 px-2 place-self-start">${this.annotation.observer}</span>
      <span class="text-sm">${species}</span>
      <span class="text-sm text-nowrap place-self-end">
        <button class="link cursor-pointer" @click=${() => this.dispatchSeekEvent(this.annotation.start)}>${start}</button>
        -
        <button class="link cursor-pointer" @click=${() => this.dispatchSeekEvent(this.annotation.end)}>${end}</button>
      </span>
      
      ${Object.entries(rest).length > 0 ? table : html``}
    </article>
    `
  }

  static styles = [
    TailwindElement.styles!,
    css`
      article {
        grid-template-rows: fit-content fit-content fit-content;
        grid-template-columns: 1fr min-content;
      }
  
      table {
        font-size: 0.8rem;
        width: 100%;
        border-spacing: 0;
        grid-column: 1 / 3;
      }
      table th:nth-child(1) {
        width: 40%;
      }
      table th {
        text-align: left;
        border-bottom: 1px solid var(--gray-200);
      }
      table th, table td {
        padding: 0.25rem;
      }
      table tbody tr:hover {
        background-color: var(--gray-50)
      }
  `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'annotation-card': AnnotationCard
  }
}
