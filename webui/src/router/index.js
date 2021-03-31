import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import("@/views/Home.vue")
    },
    {
        path: '/s/:sk',
        name: 'Show',
        component: () => import("@/views/Show.vue")
    },
    {
        path: '/e/:sk',
        name: 'Edit',
        component: () => import("@/views/Edit.vue")
    }
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

export default router
