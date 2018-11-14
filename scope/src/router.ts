import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import Mirrors from './views/Mirrors.vue'
import Stations from './views/Stations.vue'

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/',
            name: 'home',
            component: Home
        },
        {
            path: '/mirrors',
            name: 'mirrors',
            component: Mirrors
        },
        {
            path: '/stations',
            name: 'stations',
            component: Stations
        }
    ]
})
