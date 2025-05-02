import { TailwindElement } from './tailwind-element'
import { css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { parseVideoTime } from '../utils/datetime.ts'
import type { AnnotationWithJoins, Identification, User } from '@openfish/client'
import user from '../icons/user.svg?raw'
import caretUp from '../icons/caret-up.svg?raw'
import { unsafeSVG } from 'lit/directives/unsafe-svg.js'
import './tooltip'

// We have a side-effect dependency on <user-provider> so
// import it here to ensure it gets loaded first in the
// created JavaScript bundle.
import './user-provider.ts'
import { consume } from '@lit/context'
import { userContext } from '../utils/context.ts'

@customElement('annotation-card')
export class AnnotationCard extends TailwindElement {
  @property({ type: Object })
  annotation: AnnotationWithJoins | undefined

  @property({ type: Boolean })
  active = false

  @consume({ context: userContext, subscribe: true })
  user: User | null = null

  dispatchSeekEvent(time: number) {
    this.dispatchEvent(
      new CustomEvent('seek', {
        detail: time,
        bubbles: true,
        composed: true,
      })
    )
  }

  render() {
    if (this.annotation === undefined) {
      return html`<div class="card"></div>`
    }

    const start = this.annotation.start
    const end = this.annotation.end

    const isMine = (iden: Identification) => iden.identified_by.some((u) => u.id === this.user?.id)

    const renderIdentification = (iden: Identification) => html`
      <li class="flex justify-between items-center">
        <div>
          <p class="font-bold">${iden.species.common_name}</p>
          <p class="text-sm text-slate-700">${iden.species.scientific_name}</p>
        </div>
        <div class="flex items-center gap-2">
          <button
            id=${`btn-list-${iden.species.id}`}
            class="flex items-center btn variant-subtle size-sm"
          >
            <span class="text-sm text-slate-700">
              ${iden.identified_by.length}
            </span>
            <div class="*:w-4 *:h-4 *:fill-slate-700">${unsafeSVG(user)}</div>
          </button>
          <tooltip-elem
            for=${`btn-list-${iden.species.id}`}
            type="click"
            placement="bottom"
            class="bg-blue-700 text-slate-700 ring ring-blue-700"
          >
            <header class="px-1 pb-2 text-slate-50 rounded-t-sm">
              Identified by
            </header>
            <ul class="-mx-2 -mb-2 px-3 py-2 bg-slate-200 rounded-b-sm">
              ${repeat(
                iden.identified_by,
                (u) => html`
                  <li class="block flex gap-1 items-center">
                    <span class="*:w-4 *:h-4 *:fill-slate-700">
                      ${unsafeSVG(user)}
                    </span>
                    <span>${u.display_name}</span>
                  </li>
                `
              )}
            </ul>
          </tooltip-elem>

          <button
            id=${`btn-upvote-${iden.species.id}`}
            class="btn variant-blue size-sm"
            ?disabled=${isMine(iden)}
          >
            <div class="*:w-4 *:h-4 *:fill-slate-50 mr-1 -ml-1">
              ${unsafeSVG(caretUp)}
            </div>
            <span>Upvote</span>
          </button>
          ${
            isMine(iden)
              ? html`<tooltip-elem
                for=${`btn-upvote-${iden.species.id}`}
                type="hover"
                placement="bottom"
                class="bg-slate-900 max-w-48"
              >
                Upvotes limited to one per user per identification. You cannot
                upvote your own identifications.
              </tooltip-elem>`
              : html``
          }
        </div>
      </li>
    `

    return html`
      <div
        id="arrow"
        class="w-4 h-4 absolute z-20 top-4 -left-2 ${this.active ? 'bg-blue-50' : 'bg-slate-200'}""
        style="clip-path: polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%);"
      ></div>
      <article
        class="card p-0 overflow-clip border-none ${this.active ? 'bg-blue-50' : 'bg-slate-200'}"
      >
        <header class="flex justify-between items-baseline text-sm px-4 py-3">
          <span class="text-nowrap">
            <button
              class="link cursor-pointer"
              @click=${() =>
                this.dispatchSeekEvent(this.annotation ? parseVideoTime(this.annotation.start) : 0)}
            >
              ${start}
            </button>
            -
            <button
              class="link cursor-pointer"
              @click=${() =>
                this.dispatchSeekEvent(this.annotation ? parseVideoTime(this.annotation.end) : 0)}
            >
              ${end}
            </button>
          </span>
          <span>
            Created by:
            <span class="bg-blue-200 rounded-sm py-0.5 px-2">
              ${this.annotation.created_by.display_name}
            </span>
          </span>
        </header>
        <ul class="space-y-2 px-4 py-2 ${!this.active ? 'pb-4' : ''}">
          ${repeat(this.annotation.identifications, renderIdentification)}
        </ul>
        ${
          this.active
            ? html`<footer
                class="flex justify-between items-baseline px-4 py-3 bg-slate-200"
              >
                <span class="text-sm">Not correct?</span>
                <button class="btn variant-blue size-sm">
                  Add identification
                </button>
              </footer>`
            : html``
        }

      </article>
    `
  }

  static styles = [
    TailwindElement.styles!,
    css`
      :host {
        position: relative;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'annotation-card': AnnotationCard
  }
}
