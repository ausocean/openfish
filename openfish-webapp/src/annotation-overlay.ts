import { LitElement, css, html, svg } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { Annotation } from './api.types.ts'
import { repeat } from 'lit/directives/repeat.js'
import { resetcss } from './reset.css.ts'

export type MouseoverAnnotationEvent = CustomEvent<number | null>

@customElement('annotation-overlay')
export class AnnotationOverlay extends LitElement {
  @property({ type: Array })
  annotations: Annotation[] = []

  @property({ type: Number })
  activeAnnotation: number | null = null

  hoverAnnotation(id: number | null) {
    this.dispatchEvent(new CustomEvent('mouseover-annotation', { detail: id }))
  }

  render() {
    const rects = repeat(this.annotations, (annotation: Annotation) => {
      const x1 = annotation.boundingBox?.x1 ?? 0
      const y1 = annotation.boundingBox?.y1 ?? 0

      const x2 = annotation.boundingBox?.x2 ?? 100
      const y2 = annotation.boundingBox?.y2 ?? 100

      return svg`
        <g 
          @mouseover="${() => this.hoverAnnotation(annotation.id)}" 
          @mouseout=${() => this.hoverAnnotation(null)}
        > 
          <rect 
            class="annotation-rect ${annotation.id === this.activeAnnotation ? 'active' : ''}" 
            x="${x1}%" y="${y1}%" width="${x2 - x1}%" height="${y2 - y1}%" 
            stroke-width="4" fill="#00000000" 
          />
          <foreignobject class="node" x="${x1}%" y="${y1}%" width="${x2 - x1}%" height="${y2 - y1}%" >
            <span class="annotation-label">Ann. #${annotation.id}</span>              
          </foreignobject>
        </g>`
    })

    return html`
          <svg width="100%" height="100%">
            ${rects}
          </svg>`
  }

  static styles = css`
    ${resetcss}

    svg {
      pointer-events: none;
    }
    .annotation-rect {
      transition: fill 0.25s;
      stroke: var(--bright-blue-400);
    }
    .annotation-rect.active {
      fill: #CCEEEE20;
    }
    .annotation-label {
      font-size: 0.75rem;
      padding: 0.4rem 0.5rem;
      background-color: var(--bright-blue-400);
      color: var(--content);
    }`
}

declare global {
  interface HTMLElementTagNameMap {
    'annotation-overlay': AnnotationOverlay
  }
}
