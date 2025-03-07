# Openfish Site
Site is a web interface for interacting with the Openfish API.


## Development guide

#### Linting / formatting
- `pnpm site fmt` to format code.
- `pnpm site check` to check for common issues.


## Tools and libraries

Name | Description | Purpose
---|---|---
[Lit](https://lit.dev/) | A library used to create [custom elements / web components](https://developer.mozilla.org/en-US/docs/Web/API/Web_Components/Using_custom_elements) | Webcomponents let us break a site down into smaller reusable sections, while can still be used with vanilla html, css & js, without needing a framework.
[Vite](https://vitejs.dev/) | Development server / bundler | During development changes made can be reflected in the browser instantly because of its hot-module reloading feature.
[Typescript](https://www.typescriptlang.org/) | Static type checking for Javascript | Typescript provides type information in code editors / IDEs and catches mistakes.
[Biome](https://biomejs.dev/) | Code formatter and linter | Helps us keep code style consistent and catches mistakes.
