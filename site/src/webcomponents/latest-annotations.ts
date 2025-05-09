import { TailwindElement } from './tailwind-element'
import { html } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { client } from '../api'
import type { AnnotationWithJoins } from '@openfish/client'
import './annotation-card'

@customElement('latest-annotations')
export class LatestAnnotations extends TailwindElement {
  @state()
  protected _items: AnnotationWithJoins[] = []

  connectedCallback() {
    super.connectedCallback()
    this.fetchData()
  }

  async fetchData() {
    const { data, error } = await client.GET('/api/v1/annotations', {
      params: {
        query: {
          limit: 8,
          offset: 0,
        },
      },
    })

    if (error !== undefined) {
      console.error(error)
    }

    if (data !== undefined) {
      this._items = data.results
    }
  }

  render() {
    return html`
      <div class="h-full grid grid-cols-3 lg:grid-cols-4 gap-4 justify-start items-center">
        ${repeat(
          this._items,
          (item) => html`
            <annotation-card
                .annotation=${item}
                simple
            >
            </annotation-card>`
        )}
      </div>
    `
  }

  static styles = TailwindElement.styles
}

declare global {
  interface HTMLElementTagNameMap {
    'latest-annotations': LatestAnnotations
  }
}
