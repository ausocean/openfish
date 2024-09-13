---
outline: deep
---


# Openfish CLI

The Openfish program starts the OpenFish API. Options can be configured using command line arguments or environment variables.

## Options


### Port
Port for server to listen on.

::: info Usage
CLI flag: **`-port=<port>`**
Environment variable: **`PORT=<port>`**
Type: **`integer`**
Default value: **`8080`**
:::

### Filestore
Uses the local datastore instead of Google Cloud Datastore.

::: info Usage
CLI flag: **`-filestore`**
Environment variable: **`FILESTORE=true`**
Type: **`boolean`**
Default value: **`false`**
:::

### Identity Aware Proxy
Use Google's Identity Aware Proxy for authentication.

::: info Usage
CLI flag: **`-iap`**
Environment variable: **`IAP=true`**
Type: **`boolean`**
Default value: **`false`**
:::

### JWT Audience
Audience to use to validate JWT token.

::: info Usage
CLI flag: **`-jwt-audience`**
Environment variable: **`JWT_AUDIENCE`**
Type: **`string`**
Default value: 
:::

<style>
.info.custom-block>p:not(.custom-block-title) {
    display: grid;
    grid-template-columns: auto 1fr;
    column-gap: 1rem;

    &>code {
        justify-self: left;
    }
}
</style>