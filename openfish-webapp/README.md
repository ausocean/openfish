# Openfish Webapp
Openfish Webapp is a web interface for interacting with the Openfish API.

## Deploying openfish
1) Build and run the web app using docker:
   ```bash
   docker build ./openfish-webapp -t openfish
   docker run -p 80:80 openfish 
   ```
2) Start go server:
   ```bash
   go run ./api/ 
   ```
3) Open the browser and visit http://localhost.

## Development guide
#### Getting started

Have the following prerequisites installed on your system:
- node v19.2.0 or later
- pnpm v8.5.1 or later
- go 1.20 or later

1) Start go server:
   ```bash
   go run ./api/ 
   ```

2) Install all dependencies in package.json:
   ```bash
   cd openfish-webapp
   pnpm i
   ```

3) Start the live development server:
   ```bash
   pnpm dev
   ```
   Visit http://localhost:5173/. Changes are updated live.

#### Linting / formatting
- `pnpm fmt` to format code.
- `pnpm check` to check for common issues.


## Tools and libraries

Name | Description | Purpose
---|---|---
[Lit](https://lit.dev/) | A library used to create [custom elements / web components](https://developer.mozilla.org/en-US/docs/Web/API/Web_Components/Using_custom_elements) | Webcomponents let us break a site down into smaller reusable sections, while can still be used with vanilla html, css & js, without needing a framework.
[Vite](https://vitejs.dev/) | Development server / bundler | During development changes made can be reflected in the browser instantly because of its hot-module reloading feature.
[Typescript](https://www.typescriptlang.org/) | Static type checking for Javascript | Typescript provides type information in code editors / IDEs and catches mistakes.
[Rome](https://rome.tools/) | Code formatter and linter | Helps us keep code style consistent and catches mistakes.
