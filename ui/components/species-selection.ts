import { TailwindElement } from '@openfish/ui/components/tailwind-element'
import { html, nothing, TemplateResult } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import type { OpenfishClient, Species } from '@openfish/client'

import '@openfish/ui/components/species-thumb'
import { consume } from '@lit/context'
import { clientContext } from '../utils/context'
import { debounce } from '../utils/debounce'

export type SpeciesSelectionEvent = CustomEvent<number | null>
type TextInputEvent = InputEvent & { target: HTMLInputElement }

@customElement('species-selection')
export class SpeciesSelection extends TailwindElement {
  @consume({ context: clientContext, subscribe: true })
  accessor client!: OpenfishClient

  @state()
  protected accessor _keys: string[] = []

  @state()
  protected accessor _vals: string[] = []

  @property({ type: Object })
  accessor selection: Species | null = null

  @state()
  private accessor _speciesList: Species[] = []

  @state()
  accessor offset = 0

  @state()
  accessor _search = ''

  @state()
  accessor _loading = true

  @state()
  accessor _isMore = false

  private selectSpecies(species: Species) {
    this.selection = species
    this.dispatchEvent(new CustomEvent('selection', { detail: this.selection?.id }))
  }

  connectedCallback() {
    super.connectedCallback()
    this.fetchMore()
  }

  private async fetchMore() {
    const { data, error } = await this.client.GET('/api/v1/species', {
      params: {
        query: {
          limit: 20,
          offset: this.offset,
          search: this._search.length > 0 ? this._search : undefined,
        },
      },
    })

    if (error !== undefined) {
      console.error(error)
    }

    if (data !== undefined) {
      this._speciesList.push(...data.results)
      this.offset += 20
      this._loading = false;
      this._isMore = data.total == data.limit;
      this.requestUpdate()
    }
  }

  private debouncedFetch = debounce(this.fetchMore, 300);

  private async search(e: TextInputEvent) {
    this._search = e.target.value
    this.offset = 0
    this._speciesList = []
    this._loading = true;
    this.debouncedFetch()
  }

  render() {
    let results : TemplateResult | typeof nothing = nothing;
    if (!this._loading) {
      if (this._speciesList.length > 0) {
        results = html`
          <ul
            class="grid overflow-hidden gap-4 p-4 grid-cols-2 auto-rows-auto"
          >
            ${repeat(
              this._speciesList,
              (species) => species.id,
              (species) => html`
                <species-thumb
                    class="hover:ring-2 hover:ring-sky-400/50 data-selected:ring-2 data-selected:ring-sky-400 data-selected:ring-offset-4 ring-offset-1 ring-offset-blue-700 transition-shadow rounded-md"
                    .species=${species}
                    @click=${() => this.selectSpecies(species)}
                    ?data-selected=${this.selection?.id === species.id}
                ></species-thumb>
            `
            )}
          </ul>
        `
      } else {
        results = html`
          <p class="text-2xl text-center text-white mt-5">No Results</p>
        `
      }
    }
    return html`
      <header
        class="bg-blue-600 px-3 py-2 border-b border-b-blue-500 shadow-sm"
      >
        <input
          type="search"
          class="bg-blue-700 border border-blue-800 text-blue-50 w-full rounded-md placeholder:text-blue-300"
          placeholder="Search species"
          @input=${this.search}
        />
      </header>
      <div class="relative overflow-y-scroll h-[calc(100%-3rem)]">
        <div class="absolute inset-0">
          ${results}
          <footer class="w-full pb-4">
            ${this._loading ? html`<p class="text-2xl text-center text-white mt-5">Loading...</p>` : nothing}
            ${!this._loading && this._isMore ? html`
            <button class="mx-auto btn variant-slate" @click=${this.fetchMore}>
              Load more
            </button>` : nothing}
          </footer>
        </div>
      </div>
    `
  }

  static styles = TailwindElement.styles
}

declare global {
  interface HTMLElementTagNameMap {
    'species-selection': SpeciesSelection
  }
}
