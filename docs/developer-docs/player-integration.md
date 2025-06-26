# Openfish Player Integration

Openfish's video player and annotation features can be integrated into your website or application using the `@openfish/ui` and `@openfish/client` packages. `@openfish/ui` provides a set of custom elements for embedding the Openfish player and displaying video streams thumbnails, species thumbnails and annotations in your site. This guide will walk you through the process of integrating the Openfish player into your website or application.

::: warning
Openfish's components are distributed as uncompiled TypeScript. You will need to set up your project to use TypeScript and Vite. Other bundlers and methods of installation may work but they are untested. If you have success with another method, please make a pull-request to our documentation.
:::

## 1: Create a vite project

```bash
pnpm create vite
```

## 2: Installation
To install the Openfish UI library and API client library, run the following commands:

::: code-group
```bash [pnpm]
pnpm add github:ausocean/openfish#path:/client
pnpm add github:ausocean/openfish#path:/ui
pnpm add tailwind
```
```bash [npm]
TODO
```
:::

## 3. Configure tailwind

In `src/app.css`, add the following:
```css
@import '@openfish/ui/theme.css';
@source "../node_modules/@openfish/ui/components/**/*.ts";
```

If you want to theme the user interface with different colours, here is where you can override the theme.

## 4. Configure Vite
Create a file: `vite.config.ts` with the following:
```ts
import config from "@openfish/ui/vite.config"
import { mergeConfig } from 'vite'

export default mergeConfig(config, {
  root: '.'
})
```

The `mergeConfig` function is used here to merge your config with the one exported from `@openfish/ui/vite.config`.

## 5. Update typescript options
In your `tsconfig.json`, set the following.
```json
"experimentalDecorators": false,
"useDefineForClassFields": true,
```

This is so typescript code with decorators can be compiled correctly.

## 6. Configure the API client
Create a file `src/api-provider.ts` with the following:

```ts
import { useApiProvider } from '@openfish/ui/components/api-provider.ts'
import { createClient } from '@openfish/client'

const client = createClient({baseUrl: 'http://localhost:8080/'})

useApiProvider(client)
```

`createClient` returns a new instance of the OpenFish API client, which is a simple wrapper around the OpenFish API. You should configure `baseUrl` if you are hosting the OpenFish API on a different origin to where you are serving the HTML from.


## 7. Use the custom elements in your site.

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>OpenFish player example</title>
    <link rel="stylesheet" href="./src/app.css" />
    <script type="module">
    import "./src/api-provider.ts"

    import "@openfish/ui/components/user-provider.ts"
    import "@openfish/ui/components/watch-stream.ts"

    const el = document.querySelector('watch-stream')
    if (el) {
      el.streamID = 5723003237171200
    }
    </script>


  </head>
  <body class="bg-blue-800 p-12">
    <api-provider>
        <user-provider>
            <div class="max-w-320 w-full mx-auto">
                <watch-stream></watch-stream>
            </div>
        </user-provider>
    </api-provider>
  </body>
</html>
```

## Example
See https://github.com/scott97/openfish-example
