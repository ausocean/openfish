<template>
    <div class="page">
        <main class="content">
            <slot />
        </main>
    </div>
</template>

<style scoped>
.page {
    width: 100%;
    min-height: 100%;
    display: flex;
    justify-content: center
}

.content {
    width: 100%;
    max-width: 1440px
}
</style>
