import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { Annotation } from './api.types.ts'
import { repeat } from 'lit/directives/repeat.js'

/**
 * TODO: write component documentation
 *
 */
@customElement('annotation-card')
export class AnnotationCard extends LitElement {
  @property({ type: Object })
  annotation: Annotation | undefined

  @property({ type: Boolean })
  outline = false

  hoverAnnotation(id: number | null) {
    this.dispatchEvent(new CustomEvent('mouseover-annotation', { detail: id }))
  }

  render() {
    if (this.annotation === undefined) {
      return html`<div class="card"></div>`
    }

    const start = new Date(this.annotation.timespan.start)
    const end = new Date(this.annotation.timespan.end)
    const duration = (end.getTime() - start.getTime()) / 1000
    const tz = new Intl.DateTimeFormat('en-AU', { day: '2-digit', timeZoneName: 'short' })
      .format(start)
      .slice(4)

    const startDate = new Intl.DateTimeFormat('en-AU', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: '2-digit',
    }).format(start)
    const startTime = new Intl.DateTimeFormat('en-AU', {
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
    }).format(start)
    const endDate = new Intl.DateTimeFormat('en-AU', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: '2-digit',
    }).format(end)
    const endTime = new Intl.DateTimeFormat('en-AU', {
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
    }).format(end)

    const rangeFormatted =
      startDate === endDate
        ? html`<em>${startDate}</em>, from <em>${startTime}</em> until <em>${endTime}</em>`
        : html`From <em>${startDate}, ${startTime}</em> until <em>${endDate}, ${endTime}</em>`

    const rows = repeat(
      Object.entries(this.annotation.observation),
      ([key, val]) => html`
      <tr>
      <td>${key}</td>
      <td>${val}</td>
      </tr>
    `
    )

    return html`
    <div class="card ${this.outline ? 'outline' : ''}" @mouseover="${() =>
      this.hoverAnnotation(this.annotation?.id ?? null)}" @mouseout=${() =>
      this.hoverAnnotation(null)}>
    <div class="header">
      <span class="title">Annotation #${this.annotation.id}</span>
      <span class="observer">${this.annotation.observer}</span>
    </div>
    <div class="timestamps">
      <div>
        <span>Time: </span>
        <span>${rangeFormatted} [${tz}]</span>
      </div>
      <div>
        <span>Duration: </span>
        <span><em>${duration} seconds</em></span>
      </div>
    </div>
    <table>
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
    
    </div>
    `
  }

  static styles = css`
  .card {
    background-color: var(--gray0);
    border: 1px solid var(--gray2);
    padding: 1rem;
    border-radius: .5rem;
    box-shadow:  0 0.25rem 0.25rem -0.25rem #00000040;
    transition: box-shadow 0.25s;
  }
  .card.outline {
    border: 1px solid var(--secondary);
    box-shadow:  0 1rem 3rem -1rem #00000080;
  }
  .header {
    display: flex; 
    justify-content: space-between;
    align-items: baseline;
    width: 100%;
    border-bottom: 1px solid var(--gray1);
    padding-bottom: 0.5rem;
  }

  .header>span {
    margin-right: 1rem;
  }
  .header>span:last-child {
    margin-right: 0;
  }

  .title {
    font-weight: bold;
    font-size: 1rem;
  }

  .observer {
    background-color: var(--gray1);
    font-size: 0.8rem;
    padding: 0.125rem 0.75rem;
    border-radius: 0.5rem;
  }

  .timestamps {
    padding: 0.5rem 0;
    width: 100%;
    font-size: 0.8rem;
    color: var(--gray3);
  }
  .timestamps>div>:nth-child(1) {
    display: inline-block;
    width: 4rem;
  }
  .timestamps em {
    color: var(--content);
  }

  table {
    font-size: 0.8rem;
    width: 100%;
    border-spacing: 0;
  }
  table th:nth-child(1) {
    width: 40%;
  }
  table th {
    text-align: left;
    border-bottom: 1px solid var(--gray1);
  }
  table th, table td {
    padding: 0.25rem;
  }
  table tbody tr:hover {
    background-color: var(--gray0)
  }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'annotation-card': AnnotationCard
  }
}
