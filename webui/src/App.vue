<template>
    <div id="app">
        <router-view/>
    </div>
</template>

<style lang="scss">
#app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
}

.flash__message {
    margin-right: 1rem;
}

.container {
    position: relative;
}

.flash-pool {
    position: absolute;
    z-index: 1;
    top: 0;
    right: 0;
}

</style>
