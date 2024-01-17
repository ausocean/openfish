import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { Annotation } from './api.types.ts'
import { repeat } from 'lit/directives/repeat.js'
import { resetcss } from './reset.css.ts'

export type MouseoverAnnotationEvent = CustomEvent<number | null>

@customElement('annotation-list')
export class AnnotationList extends LitElement {
  @property({ type: Array })
  annotations: Annotation[] = []

  @property({ type: Number })
  activeAnnotation: number | null = null

  dispatchMouseOverAnnotation(id: number | null) {
    this.dispatchEvent(new CustomEvent('mouseover-annotation', { detail: id }))
  }

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
    ${resetcss}
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
  }
}
