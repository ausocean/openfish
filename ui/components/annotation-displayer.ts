import { TailwindElement } from './tailwind-element'
import { css, html, svg } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { createRef, ref, type Ref } from 'lit/directives/ref.js'
import { findClosestKeypointPair, interpolateKeypoints } from '../utils/keypoints.ts'
import type { AnnotationWithJoins } from '@openfish/client'
import { parseVideoTime } from '../utils/datetime.ts'
import { plainToInstance } from 'class-transformer'
import { Keypoint } from '../utils/keypoints.ts'

export type MouseoverAnnotationEvent = CustomEvent<number | null>

abstract class AnnotationDisplayer extends TailwindElement {
  @property({ type: Array })
  accessor annotations: AnnotationWithJoins[] = []

  @property({ type: Number })
  accessor activeAnnotation: number | null = null

  @property({ type: Number })
  accessor currentTime = 0

  dispatchMouseOverAnnotation(id: number | null) {
    this.dispatchEvent(new CustomEvent('mouseover-annotation', { detail: id }))
  }
}

@customElement('annotation-overlay')
export class AnnotationOverlay extends AnnotationDisplayer {
  render() {
    const rects = repeat(this.annotations, (annotation: AnnotationWithJoins) => {
      // Interpolate between keypoints.
      const keypoints = annotation.keypoints.map((k) => plainToInstance(Keypoint, k))
      const kpPair = findClosestKeypointPair(keypoints, this.currentTime)
      const box = interpolateKeypoints(kpPair, this.currentTime)

      return svg`
        <g
          @mouseover="${() => this.dispatchMouseOverAnnotation(annotation.id)}"
          @mouseout=${() => this.dispatchMouseOverAnnotation(null)}
        >

          <foreignobject x="${box.xmin}%" y="${box.ymin}%" width="${box.w}%" height="${box.h}%" class="relative">
            <span class="px-1 py-0.5 text-xs bg-slate-900/50 text-white w-full absolute top-0 right-0 text-center text-nowrap">
              ${annotation.identifications[0].species.common_name}
            </span>
          </foreignobject>
          <rect
            class="annotation-rect stroke-sky-400 data-active:fill-white/25 transition-colors"
            x="${box.xmin}%" y="${box.ymin}%" width="${box.w}%" height="${box.h}%"
            stroke-width="2px" fill="#00000000"
            ?data-active=${annotation.id === this.activeAnnotation}
          />
        </g>`
    })

    return html` <svg width="100%" height="100%">${rects}</svg>`
  }

  static styles = [
    TailwindElement.styles!,
    css`
      svg {
        pointer-events: none;
      }
    `,
  ]
}

@customElement('annotation-list')
export class AnnotationList extends AnnotationDisplayer {
  listcontain: Ref<HTMLElement> = createRef()

  render() {
    const currentIdx = this.annotations.findIndex((a) => this.currentTime < parseVideoTime(a.end))
    const renderItem = (a: AnnotationWithJoins, i: number) =>
      html` <li class="contents">
        <div class="pt-2 flex flex-col items-center gap-2 translate-y-2.5">
          <div
            class="rounded-full w-3 aspect-square border-2 border-blue-300 ${
              currentIdx === i ? 'bg-blue-300' : ''
            }"
          ></div>
          ${
            i === this.annotations.length - 1
              ? html``
              : html`<div class="bg-blue-600 w-0.5 h-full"></div>`
          }
        </div>

        <annotation-card
          class="pb-3"
          .annotation=${a}
          .active=${
            parseVideoTime(a.start) <= this.currentTime && this.currentTime <= parseVideoTime(a.end)
          }
          @mouseover=${() => this.dispatchMouseOverAnnotation(a.id)}
          @mouseout=${() => this.dispatchMouseOverAnnotation(null)}
        />
      </li>`

    return html`
      <ul ${ref(this.listcontain)} class="grid p-2 overflow-hidden">
        ${repeat(this.annotations, (a) => a.id, renderItem)}
      </ul>
    `
  }

  static styles = [
    TailwindElement.styles!,
    css`
      ul {
        grid-template-columns: 1.5rem 1fr;
        grid-auto-rows: auto;
        column-gap: 0.5rem;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'annotation-list': AnnotationList
    'annotation-overlay': AnnotationOverlay
  }
}
