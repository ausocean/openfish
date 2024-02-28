import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { buttonStyles, resetcss } from './reset.css'
import { zip } from './array-utils'

type Species = { species: string; common_name: string; images?: Image[] }
type Image = { src: string; attribution: string }

export type ObservationEvent = CustomEvent<Record<string, string>>

abstract class ObservationEditor extends LitElement {
  @state()
  protected _keys: string[] = []

  @state()
  protected _vals: string[] = []

  @property({ type: Object })
  set observation(val: Record<string, string>) {
    this._keys = Object.keys(val)
    this._vals = Object.values(val)
  }
  get observation() {
    return Object.fromEntries(zip(this._keys, this._vals))
  }

  dispatchObservationEvent() {
    this.dispatchEvent(
      new CustomEvent('observation', { detail: this.observation }) as ObservationEvent
    )
  }
}

@customElement('species-selection')
export class SpeciesSelection extends ObservationEditor {
  @state()
  private _speciesList: Species[] = []

  private selectSpecies(val: string) {
    const { species, common_name } = this._speciesList.find((s) => s.species === val)!
    this.observation = { species, common_name }
    this.dispatchObservationEvent()
  }

  connectedCallback() {
    super.connectedCallback()
    this.fetchGuide()
  }

  private async fetchGuide() {
    try {
      // TODO: fetch list from an API.
      const res = await fetch(
        `${import.meta.env.VITE_API_HOST}/api/v1/species/recommended?limit=999`
      )
      this._speciesList = (await res.json()).results
    } catch (error) {
      console.error(error)
    }
  }

  render() {
    const items = repeat(
      this._speciesList,
      ({ species, common_name, images }) => html`
        <li class="card" @click=${() => this.selectSpecies(species)} data-selected=${
          this.observation.species === species
        }>
          <span class="attribution">${images?.at(0)?.attribution ?? 'No image available'}</span>
          <img src=${images?.at(0)?.src ?? 'placeholder.svg'} />
          <div class="common_name">${common_name}</div>
          <div class="species">${species}</div>
        </li>
    `
    )

    return html`
    <ul>
        ${items}
    </ul>
    `
  }

  static styles = css`
    ${resetcss}
    ul {
        list-style-type: none;
        margin:0;
        padding: 0;
        overflow-y: scroll;
        
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        grid-template-rows: auto;
        gap: 1rem
    }

    .card {
        background-color: var(--gray-50);
        border: 2px solid var(--blue-300);
        border-radius: .5rem;
        box-shadow:  var(--shadow-sm);
        overflow: clip;
        cursor: pointer;
        color: var(--gray-900);
        position: relative;

        & img {
            width: 100%;
            aspect-ratio:  4 / 3;
            object-fit: cover;
        }

        & .common_name {
            padding: 0 .5rem;
            font-weight: bold;
        }

        & .species {
            padding: 0 .5rem;
            padding-bottom: 0.5rem;
            font-size: 0.8rem;
        }

        & .attribution {  
          font-size: 0.5rem;
          background-color: rgba(0, 0, 0, 0.5);
          color: var(--bg);
          padding: 0.5em 1em;
          border-radius: 999999px;
          position: absolute;
          top: 0.5em;
          right: 0.5em;
          white-space: pre-line;
        }
    }

    .card[data-selected="true"] {
      background-color: var(--blue-200);
      color: black;
      border-color: var(--bright-blue-400);
      box-shadow:  var(--shadow-lg), 0px 0px 10px 2px color-mix(in srgb, var(--bright-blue-400) 80%, transparent);

    }

  `
}

type TextInputEvent = InputEvent & { target: HTMLInputElement }

@customElement('advanced-editor')
export class AdvancedEditor extends ObservationEditor {
  // Mutating an array doesn't trigger an update.
  // https://lit.dev/docs/components/properties/#mutating-properties
  addRow() {
    this._keys.push('')
    this._vals.push('')
    this.requestUpdate()
    this.dispatchObservationEvent()
  }

  updateKey(idx: number, key: string) {
    this._keys[idx] = key
    this.requestUpdate()
    this.dispatchObservationEvent()
  }

  updateVal(idx: number, val: string) {
    this._vals[idx] = val
    this.requestUpdate()
    this.dispatchObservationEvent()
  }

  render() {
    const rows = repeat(
      zip(this._keys, this._vals),
      ([key, val], idx) => html`
      <tr>
      <td><input type="text" .value=${key} @input=${(ev: TextInputEvent) =>
        this.updateKey(idx, ev.target.value)}></input></td>
      <td><input type="text" .value=${val} @input=${(ev: TextInputEvent) =>
        this.updateVal(idx, ev.target.value)}></input></td>
      </tr>
    `
    )

    return html`
    <div class="card">
    <table>
    <thead>
        <tr>
            <th>Property</th>
            <th>Value</th>
        </tr>
    </thead>
    <tbody>
    ${rows}
    <tr>
    <td colspan="2"><button class="btn-sm btn-orange btn-fullwidth" @click=${this.addRow}>+ Add information</button></td>
    </tr>
    </tbody>
</table>
</div>
    `
  }

  static styles = css`
    ${resetcss}
    ${buttonStyles}
    .card {
      background-color: var(--gray-50);
      padding: 1rem;
      border-radius: .5rem;
      box-shadow:  var(--shadow-sm);
      width: calc(var(--aside-width) - 3rem);
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
      border-bottom: 1px solid var(--gray-200);
    }
    table th, table td {
      padding: 0.25rem;
    }
    table tbody tr:hover {
      background-color: var(--gray-50)
    }

    tbody:last-child td {

      padding: 0.5rem;
    }
    .btn-fullwidth {
      text-align: center;
      width: 100%;  
    }

    input {
      width: 100%;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'advanced-editor': AdvancedEditor
    'species-selection': SpeciesSelection
  }
}
