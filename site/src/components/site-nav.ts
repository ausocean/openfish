import { TailwindElement } from '@openfish/ui/components/tailwind-element'
import { html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { consume } from '@lit/context'
import { userContext } from '@openfish/ui/utils/context'
import type { User } from '@openfish/client'

// We have a side-effect dependency on <user-provider> so
// import it here to ensure it gets loaded first in the
// created JavaScript bundle.
import '@openfish/ui/components/user-provider.ts'

@customElement('site-nav')
export class SiteNav extends TailwindElement {
  @consume({ context: userContext, subscribe: true })
  user: User | null = null

  render() {
    let user = html``
    if (this.user !== null) {
      user = html`
      <li>
        <span class="uppercase bg-blue-300 text-blue-900 text-xs rounded px-2 py-1">
          ${this.user?.role}
        </span>
        ${this.user?.display_name}
      </li>
      `
    }

    return html`
    <div class="contain px-8 py-4 h-16 z-1000 flex">
      <h1 class="text-xl self-center flex-1">
        <a href="/" class="link font-bold">
          OpenFish
        </a>
      </h1>
      <menu class="flex justify-end gap-4 self-center">
          <li>
            <a href="/streams" class="link whitespace-nowrap">
              View streams
            </a>
          </li>
          ${
            this.user?.role === 'admin'
              ? html`
            <li>
              <a href="/admin/capturesources" class="link whitespace-nowrap">
                Admin Settings
              </a>
            </li>
          `
              : html``
          }
          |
          ${user}
      </menu>
    </div>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'site-nav': SiteNav
  }
}
