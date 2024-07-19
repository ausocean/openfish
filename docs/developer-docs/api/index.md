---
layout: page
sidebar: false
---

<script setup lang="ts">
import { ApiReference } from '@scalar/api-reference'
import '@scalar/api-reference/style.css'
import swagger from './swagger.json?url'
import Container from '../../components/Container.vue'
</script>


<Container>
<ApiReference
:configuration="{
    spec: {
    url: swagger,
    },
    defaultHttpClient: { targetKey:'http', clientKey:'http1.1'},
    theme: 'purple',
    hideModels: true,
    servers: [
        {
            url: 'http://localhost:8080',
            description: 'Localhost development server'
        },
        {
            url: 'https://openfish.appspot.com',
            description: 'Live server'
        },
    ],
}" />
</Container>