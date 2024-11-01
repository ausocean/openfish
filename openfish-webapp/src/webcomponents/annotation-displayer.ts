import { LitElement, css, html, svg } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import resetcss from '../styles/reset.css?lit'
import { createRef, ref, type Ref } from 'lit/directives/ref.js'
import { findClosestKeypointPair, interpolateKeypoints } from '../utils/keypoints.ts'
import type { Annotation } from '../api/annotation.ts'

export type MouseoverAnnotationEvent = CustomEvent<number | null>

abstract class AnnotationDisplayer extends LitElement {
  @property({ type: Array })
  annotations: Annotation[] = []

  @property({ type: Number })
  activeAnnotation: number | null = null

  @property({ type: Number })
  currentTime = 0

  dispatchMouseOverAnnotation(id: number | null) {
    this.dispatchEvent(new CustomEvent('mouseover-annotation', { detail: id }))
  }
}

@customElement('annotation-overlay')
export class AnnotationOverlay extends AnnotationDisplayer {
  render() {
    const rects = repeat(this.annotations, (annotation: Annotation) => {
      // Interpolate between keypoints.
      const kpPair = findClosestKeypointPair(annotation.keypoints, this.currentTime)
      const box = interpolateKeypoints(kpPair, this.currentTime)

      return svg`
          <g 
            @mouseover="${() => this.dispatchMouseOverAnnotation(annotation.id)}" 
            @mouseout=${() => this.dispatchMouseOverAnnotation(null)}
          > 
            <rect 
              class="annotation-rect ${annotation.id === this.activeAnnotation ? 'active' : ''}" 
              x="${box.xmin}%" y="${box.ymin}%" width="${box.w}%" height="${box.h}%" 
              stroke-width="3px" fill="#00000000" 
            />
            <foreignobject x="${box.xmin}%" y="${box.ymin}%" width="${box.w}%" height="${box.h}%" >
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
  listContainer: Ref<HTMLElement> = createRef()

  render() {
    const currentIdx = this.annotations.findIndex((a) => this.currentTime < a.end)

    const items = repeat(
      this.annotations,
      (a) => a.id,
      (a: Annotation, i: number) => html`
      ${currentIdx === i ? html`<hr class="playback-hr"/>` : html``}
    <annotation-card
      .annotation=${a}
      .glow=${this.activeAnnotation === a.id}
      .outline=${a.start <= this.currentTime && this.currentTime <= a.end}
      @mouseover=${() => this.dispatchMouseOverAnnotation(a.id)}
      @mouseout=${() => this.dispatchMouseOverAnnotation(null)}
    />`
    )

    return html`
    <div ${ref(this.listContainer)}>
        ${items}
    </div>
    `
  }

  static styles = css`
    ${resetcss}
    div {
      height: 100%;
      overflow-y: scroll;
      display: flex;
      flex-direction: column;
      gap: 1rem;
      padding: 1rem 1rem;
      overflow: visible; 
    }

    annotation-card.active {
      border: 1px solid var(--bright-blue-400)
    }

    .playback-hr {
      border: 1px solid var(--red-500);
      width: 100%;
      height: 0;
      margin-top: -0.5rem;
      margin-bottom: -0.5rem;
    }
    `
}

declare global {
  interface HTMLElementTagNameMap {
    'annotation-list': AnnotationList
    'annotation-overlay': AnnotationOverlay
  }
}
