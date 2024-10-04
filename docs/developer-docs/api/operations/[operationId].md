---
aside: false
outline: false
title: API documentation
---

<script setup lang="ts">
import { useRoute, useData } from 'vitepress'

const route = useRoute()

const { isDark } = useData()

const operationId = route.data.params.operationId

try {
    customElements.define("role-tag", HTMLSpanElement);
} catch (e) {}
</script>

<OAOperation :operationId="operationId" :isDark="isDark" hideDefaultFooter/>

<style scope="global">
role-tag {
    text-transform: uppercase;
    background-color: var(--vp-custom-block-important-bg);
    color: var(--vp-custom-block-important-text);
    font-size: 0.8rem;
    border-radius: 4px;
    padding: 0.25em 0.5em;
}
</style>