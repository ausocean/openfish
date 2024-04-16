import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement } from 'lit/decorators.js'
import resetcss from '../styles/reset.css?raw'

@customElement('site-nav')
export class SiteNav extends LitElement {
  render() {
    return html`
    <h1><a href="/">OpenFish WebApp</a></h1>
    <menu>
        <li><a href="/streams.html">View streams</a></li>
        <li><a href="/capturesources.html">Manage capture sources</a></li>
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
    align-self: end;
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
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'site-nav': SiteNav
  }
}
