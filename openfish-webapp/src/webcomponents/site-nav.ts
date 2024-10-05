import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement } from 'lit/decorators.js'
import resetcss from '../styles/reset.css?raw'
import type { User } from '../api/user'
import { consume } from '@lit/context'
import { userContext } from '../utils/context'

@customElement('site-nav')
export class SiteNav extends LitElement {
  @consume({ context: userContext, subscribe: true })
  user: User | null = null

  render() {
    let user = html``
    if (this.user !== null) {
      user = html`
      <li>
        <span class="tag">${this.user?.role}</span>
        ${this.user?.email} 
      </li>
      `
    }

    return html`
    <h1><a href="/">OpenFish</a></h1>
    <menu>
        <li><a href="/streams.html">View streams</a></li>
        ${
          this.user?.role === 'admin'
            ? html`
          <li><a href="/capturesources.html">Manage capture sources</a></li>   
        `
            : html``
        }
        |
        ${user}
    </menu>
    `
  }

  static styles = css`
  ${unsafeCSS(resetcss)}
  
  :host {
    grid-column: fullwidth;
    grid-row-start: 1;
    grid-row-end: 2;
    display: grid;
    grid-template-columns: subgrid;
    border-bottom: 1px solid;
    padding: 1rem 0;
  }
  h1 {
    grid-column: left-aside;
    font-size: 1.25rem;
    font-weight: 600;
    align-self: end;
  }
  menu {
    grid-column: right-content;
    display: flex;
    justify-content: flex-end;
    align-self: center;
    gap: 1rem;
    font-weight: 500;

    & li {
      list-style-type: none;
      padding: 0;
    }
  }
  a {
    text-decoration: none;
    color: currentColor;
  }
  a:hover {
    color: var(--bright-blue-500);
    text-decoration: underline;
  }
  .tag {
      text-transform: uppercase;
      background-color: var(--blue-300);
      color: var(--blue-900);
      font-size: 0.7rem;
      border-radius: 4px;
      padding: 0.25em 0.5em;
  }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'site-nav': SiteNav
  }
}
