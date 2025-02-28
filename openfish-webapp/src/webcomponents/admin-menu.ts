import { TailwindElement } from './tailwind-element'
import { css, html } from 'lit'
import { customElement } from 'lit/decorators.js'

@customElement('admin-menu')
export class AdminMenu extends TailwindElement {
  render() {
    return html`
    <aside class="pr-4 h-full">
      <menu class="flex flex-col gap-1 text-slate-800 *:transition-colors *:rounded-md *:overflow-clip *:hover:bg-slate-200 *:hover:text-blue-700">
          <li><a href="/admin/capturesources.html">Manage Capture Sources</a></li>
          <li><a href="/admin/users.html">Manage Users</a></li>
          <li><a href="/admin/species.html">Manage Species</a></li>
          <li><a href="/admin/videostreams.html">Manage Video Streams</a></li>
      </menu>
    </aside>
    `
  }

  static styles = [
    TailwindElement.styles!,
    css`
      a {
        padding: .25rem .5rem;
        display: block;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'admin-menu': AdminMenu
  }
}
