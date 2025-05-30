import { TailwindElement } from './tailwind-element'
import { html, svg } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { parseVideoTime } from '../utils/datetime'
import type { AnnotationWithJoins } from '@openfish/client'

@customElement('timeline-heatmap')
export class TimelineHeatmap extends TailwindElement {
  @property({ type: Number })
  accessor duration = 0

  @property({ type: Array })
  accessor annotations: AnnotationWithJoins[] = []

  render() {
    const svgContents = repeat(this.annotations, (a) => {
      const x = (parseVideoTime(a.start) / this.duration) * 100
      const width = Math.max(((a.duration / this.duration) * 100) / 1000, 0.25) // Give them a min width of 0.25% so they are legible.

      return svg`<rect class="fill-green-500 opacity-50" x="${x}%" y="0%" width="${width}%" height="100%" />`
    })

    return html`
      <svg class="absolute inset z-10 w-full h-6">${svgContents}</svg>
    `
  }

  static styles = TailwindElement.styles
}

declare global {
  interface HTMLElementTagNameMap {
    'timeline-heatmap': TimelineHeatmap
  }
}
