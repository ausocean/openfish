import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import resetcss from '../styles/reset.css'
import { formatVideoTime } from '../utils/datetime.ts'
import type { Annotation } from '../api/annotation.ts'

@customElement('annotation-card')
export class AnnotationCard extends LitElement {
  @property({ type: Object })
  annotation: Annotation | undefined

  @property({ type: Boolean })
  outline = false

  @property({ type: Boolean })
  glow = false

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
    <article class="card ${this.glow ? 'glow' : ''} ${this.outline ? 'outline' : ''}">
      <span class="title">${common_name}</span>
      <span class="observer">${this.annotation.observer}</span>
      <span class="species">${species}</span>
      <span class="timestamps"><a href="#">${start}</a> - <a href="#">${end}</a></span>
      
      ${Object.entries(rest).length > 0 ? table : html``}
    </article>
    `
  }

  static styles = css`
  ${resetcss}

  a:not(.btn) {
    font-weight: 500;
    color: var(--content);
    text-decoration: underline;
  
    &:hover {
      color: var(--bright-blue-500);
    }
  }

  .card {
    background-color: var(--gray-50);
    border: 2px solid var(--blue-300);
    padding: 0.75rem;
    border-radius: .5rem;
    box-shadow:  var(--shadow-sm);
    transition: box-shadow 0.25s;

    display: grid;
    gap: 0.5rem;
    grid-template-rows: min-content min-content min-content;
    grid-template-columns: 1fr min-content;
    width: 100%;

  }
  .card.glow { 
    border: 2px solid var(--bright-blue-400);
    box-shadow:  var(--shadow-lg), 0px 0px 10px 2px color-mix(in srgb, var(--bright-blue-400) 80%, transparent);
  }
  .card.outline {
    border: 2px solid var(--bright-blue-400);
    box-shadow: var(--bright-blue-400) 0 0 0 2px inset;
  }
  .card.outline.glow {
    box-shadow: var(--bright-blue-400) 0 0 0 2px inset, var(--shadow-lg), 0px 0px 10px 2px color-mix(in srgb, var(--bright-blue-400) 80%, transparent);
  }


  .title {
    font-weight: bold;
    font-size: 1rem;
  }

  .observer {
    background-color: var(--gray-200);
    font-size: 0.8rem;
    padding: 0.125rem 0.75rem;
    border-radius: 0.5rem;
  }

  .timestamps {
    font-size: 0.8rem;
    text-wrap: nowrap;
    place-self: end;
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
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'annotation-card': AnnotationCard
  }
}
