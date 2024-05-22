import { LitElement, css, html, svg, unsafeCSS } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import type { Annotation } from '../utils/api.types.ts'
import { repeat } from 'lit/directives/repeat.js'
import resetcss from '../styles/reset.css?raw'

export type MouseoverAnnotationEvent = CustomEvent<number | null>

abstract class AnnotationDisplayer extends LitElement {
  @property({ type: Array })
  annotations: Annotation[] = []

  @property({ type: Number })
  activeAnnotation: number | null = null

  dispatchMouseOverAnnotation(id: number | null) {
    this.dispatchEvent(new CustomEvent('mouseover-annotation', { detail: id }))
  }
}

@customElement('annotation-overlay')
export class AnnotationOverlay extends AnnotationDisplayer {
  render() {
    const rects = repeat(this.annotations, (annotation: Annotation) => {
      const x1 = annotation.boundingBox?.x1 ?? 0
      const y1 = annotation.boundingBox?.y1 ?? 0

      const x2 = annotation.boundingBox?.x2 ?? 100
      const y2 = annotation.boundingBox?.y2 ?? 100

      return svg`
        <g 
          @mouseover="${() => this.dispatchMouseOverAnnotation(annotation.id)}" 
          @mouseout=${() => this.dispatchMouseOverAnnotation(null)}
        > 
          <rect 
            class="annotation-rect ${annotation.id === this.activeAnnotation ? 'active' : ''}" 
            x="${x1}%" y="${y1}%" width="${x2 - x1}%" height="${y2 - y1}%" 
            stroke-width="3px" fill="#00000000" 
          />
          <foreignobject x="${x1}%" y="${y1}%" width="${x2 - x1}%" height="${y2 - y1}%" >
            <span class="annotation-label">
              <span>${annotation.observation.common_name}</span>
              <span>(${annotation.observation.species})</span>
            </span>              
          </foreignobject>
        </g>`
    })

    return html`
          <svg width="100%" height="100%">
            ${rects}
          </svg>`
  }

  static styles = css`
    ${unsafeCSS(resetcss)}

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
      padding: 0.1rem 0.2rem;
      background-color: var(--bright-blue-400);
      color: var(--content);
      width: 100%;
      top: -1px;
      position: absolute;
    }
    .annotation-label>span {
      text-wrap: nowrap;
    }
    `
}

@customElement('annotation-list')
export class AnnotationList extends AnnotationDisplayer {
  render() {
    const items = repeat(
      this.annotations,
      (annotation: Annotation) => html`
    <annotation-card
      .annotation=${annotation}
      .outline=${this.activeAnnotation === annotation.id}
      @mouseover=${() => this.dispatchMouseOverAnnotation(annotation.id)}
      @mouseout=${() => this.dispatchMouseOverAnnotation(null)}
    />`
    )

    return html`
    <div>
        ${items}
    </div>
    `
  }

  static styles = css`
    ${unsafeCSS(resetcss)}
    div {
      height: 100%;
      overflow-y: scroll;
      display: flex;
      flex-direction: column;
      gap: 1rem;
      padding: 1rem 0;
      overflow: visible; 
    }

    annotation-card.active {
      border: 1px solid var(--bright-blue-400)
    }
    `
}

declare global {
  interface HTMLElementTagNameMap {
    'annotation-list': AnnotationList
    'annotation-overlay': AnnotationOverlay
  }
}
