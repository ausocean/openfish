import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement } from 'lit/decorators.js'
import resetcss from '../styles/reset.css?lit'

@customElement('admin-menu')
export class AdminMenu extends LitElement {
  render() {
    return html`
    <aside>
      <menu>
          <li><a href="/admin/capturesources.html">Manage Capture Sources</a></li>
          <li><a href="/admin/users.html">Manage Users</a></li>
          <li><a href="/admin/capturesources.html">Manage Species</a></li>
      </menu>
    </aside>
    `
  }

  static styles = css`
    ${unsafeCSS(resetcss)}
    aside {
      border-right: 1px solid var(--gray-200);
      padding-right: 1rem;
      padding-top: 1rem;
      height: 100%
    }   
    menu {
      font-weight: 500;
      display: flex;
      flex-direction: column;
      gap: .25rem;

      & li {
        list-style-type: none;
        padding: 0;
      }
    }
    a {
      text-decoration: none;
      color: currentColor;
      padding: .25rem .5rem;
      border-radius: .25rem;
      display: block;
    }
    a:hover {
      color: var(--blue-500);
      background-color: var(--gray-100);
    }
    `
}

declare global {
  interface HTMLElementTagNameMap {
    'admin-menu': AdminMenu
  }
}
