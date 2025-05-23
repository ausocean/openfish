<script setup>
import "@openfish/ui/components/stream-thumb"
</script>

# Openfish Player Integration

Openfish's video player and annotation features can be integrated into your website or application using the `@openfish/ui` and `@openfish/client` packages. `@openfish/ui` provides a set of custom elements for embedding the Openfish player and displaying video streams thumbnails, species thumbnails and annotations in your site. This guide will walk you through the process of integrating the Openfish player into your website or application.

::: warning
Openfish's components are distributed as uncompiled TypeScript. You will need to set up your project to use TypeScript and Vite. Other bundlers and methods of installation may work but they are untested. If you have success with another method, please make a pull-request to our documentation.
:::

## 1: Create a vite project (optional)

```bash
pnpm create vite
```

## 2: Installation
To install the Openfish UI library and API client library, run the following commands:

::: code-group
```bash [pnpm]
pnpm add github:ausocean/openfish#path:/client
pnpm add github:ausocean/openfish#path:/ui
```
```bash [npm]
TODO
```
:::

## 3. Use custom elements in your HTML
Custom elements can be used in your HTML, by importing the component from `@openfish/ui`. As a simple example, this is the `<stream-thumb>` element which displays a thumbnail of a video stream.

<demo html="./code-examples/basic.html" title="Basic example"  />

## 4. Configure the API client
In the previous example, we set the `stream` property to a hardcoded value, now we will show how to use the API client to fetch the stream data. `createClient` returns a new instance of the OpenFish API client, which is a simple wrapper around the OpenFish API. You should configure `baseUrl` if you are hosting the OpenFish API on a different origin to where you are serving the HTML from.

<demo html="./code-examples/api-client.html" title="API client example"  />

## 5. Create a custom API provider component
Many components, such as the `<watch-stream>` element interact with the OpenFish API on your behalf, so you do not need to use the API client directly. To configure these components to use the API using our settings, create an `<api-provider>` element from your API client using `useApiProvider`.

<demo html="./code-examples/api-provider.html" title="API provider example"  />
