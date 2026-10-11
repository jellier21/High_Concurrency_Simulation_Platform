import {createRouter, createMemoryHistory, createWebHistory} from "vue-router";

const routes = [
    {
        path:'/',
        name:'index',
        component: ()=>import('../App.vue')
    },
    {
        path: '/user',
        name: 'user',
        component: ()=>import('/src/components/user/userMain.vue')
    },
    {
        path: '/manager',
        name: 'manager',
        component: ()=>import('/src/components/manager/managerMain.vue')
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router