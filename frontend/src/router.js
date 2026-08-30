import { createRouter, createWebHashHistory } from 'vue-router'
import { getToken } from './api'

const routes = [
  { path: '/enter', name: 'enter', component: () => import('./views/EnterView.vue'), meta: { public: true } },
  { path: '/', name: 'home', component: () => import('./views/HomeView.vue') },
  { path: '/create', name: 'create', component: () => import('./views/CreateView.vue') },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

router.beforeEach((to) => {
  if (!to.meta.public && !getToken()) return { name: 'enter' }
  if (to.name === 'enter' && getToken()) return { name: 'home' }
})

export default router
