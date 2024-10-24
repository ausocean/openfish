import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { zip } from '../utils/array-utils'
import resetcss from '../styles/reset.css'
import btncss from '../styles/buttons.css'

type Species = { species: string; common_name: string; images?: Image[] }
type Image = { src: string; attribution: string }

export type ObservationEvent = CustomEvent<Record<string, string>>

abstract class AbstractObservationEditor extends LitElement {
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

@customElement('observation-editor')
export class ObservationEditor extends AbstractObservationEditor {
  @state()
  private _editorMode: 'simple' | 'advanced' = 'simple'

  _onObservation(ev: ObservationEvent) {
    this.observation = ev.detail
    this.dispatchObservationEvent()
  }

  render() {
    return html`
    <menu>
    <h4>Observation</h4>
    <button class="btn-sm ${
      this._editorMode === 'simple' ? ' btn-secondary' : 'btn-outline'
    }"  @click=${() => {
      this._editorMode = 'simple'
    }}>Select species</button>
    <button class="btn-sm ${
      this._editorMode === 'advanced' ? ' btn-secondary' : 'btn-outline'
    }" @click=${() => {
      this._editorMode = 'advanced'
    }}>Additional observations</button>
    </menu>

      ${
        this._editorMode === 'simple'
          ? html`<species-selection .observation=${this.observation} @observation=${this._onObservation} ></species-selection>`
          : html`<advanced-editor .observation=${this.observation} @observation=${this._onObservation}></advanced-editor>`
      }
    `
  }

  static styles = css`
  ${resetcss}
  ${btncss}
  menu {
    display: flex;
    justify-content: end;
    margin: 0;
    padding: 0.5rem 1rem;
    gap: 0.5rem;
    background-color: var(--blue-600);

    & > h4 {
      color: var(--gray-50);
      margin-right: auto;
    }  

    & button[data-active="true"] {
      background-color: var(--gray-50);
      color: var(--gray-900);
    }

    & button[data-active="false"] {
      background-color: transparent;
      color: var(--gray-50);
    }
  }
  `
}

@customElement('species-selection')
export class SpeciesSelection extends AbstractObservationEditor {
  @state()
  private _speciesList: Species[] = []

  @state()
  offset = 0

  @state()
  _search = ''

  private selectSpecies(val: string) {
    const { species, common_name } = this._speciesList.find((s) => s.species === val)!
    this.observation = { species, common_name }
    this.dispatchObservationEvent()
  }

  connectedCallback() {
    super.connectedCallback()
    this.fetchMore()
  }

  private async fetchMore() {
    try {
      const params = new URLSearchParams({
        limit: '20',
        offset: this.offset.toString(),
      })
      if (this._search.length > 0) {
        params.set('search', this._search)
      }
      const res = await fetch(`/api/v1/species/recommended?${params}`)
      this._speciesList.push(...(await res.json()).results)
      this.offset += 20
      this.requestUpdate()
    } catch (error) {
      console.error(error)
    }
  }

  private async search(e: TextInputEvent) {
    this._search = e.target.value
    this.offset = 0
    this._speciesList = []
    this.fetchMore()
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
    <header>
      <input type="text" placeholder="Search species" @input=${this.search}></input>
    </header>
    <div class="scrollable">
    <div>
      <ul>
          ${items}
      </ul>
      <footer>
        <button class="btn btn-blue" @click=${this.fetchMore}>Load more</button>
      </footer>
    </div>
    </div>
    `
  }

  static styles = css`
    ${resetcss}
    ${btncss}

    .scrollable {
      position: relative;
      height: calc(100% - 6rem);
      overflow-y: scroll;
    }
    .scrollable > * {
      position: absolute;
      left: 0;
      top: 0;
    }
    
    header {
      background-color: var(--blue-600);
      padding: 0.5rem 1rem;
      border-bottom: 1px solid var(--blue-500);
      box-shadow: var(--shadow-sm);
    }

    input {
      width: 100%;
      padding: 0.5rem;
      font-size: 1rem;
      background-color: var(--blue-700);
      border: 1px solid var(--blue-800);
      border-radius: 0.25rem;
      color: var(--blue-100);

      &:focus {
        outline: none;
        color: var(--blue-50);
        border-color: var(--blue-400);
      }
    }

    ul {
        list-style-type: none;
        margin:0;
        overflow-y: scroll;
        
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        grid-template-rows: auto;
        gap: 1rem;
        padding: 1rem;
    }

    footer {
        width: 100%;
        display: flex;
        justify-content: center;
        padding-bottom: 1rem;
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
export class AdvancedEditor extends AbstractObservationEditor {
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
    <div class="root">
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
    </div>
    `
  }

  static styles = css`
    ${resetcss}
    ${btncss}
    .root {
      padding: 1rem;
    }
    .card {
      background-color: var(--gray-50);
      padding: 1rem;
      border-radius: .5rem;
      box-shadow:  var(--shadow-sm);
      width: 100%;
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
