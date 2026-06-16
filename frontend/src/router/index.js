import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/seafarers' },
  { path: '/seafarers', name: 'Seafarers', component: () => import('../views/Seafarers.vue') },
  { path: '/certificates', name: 'Certificates', component: () => import('../views/Certificates.vue') },
  { path: '/ships', name: 'Ships', component: () => import('../views/Ships.vue') },
  { path: '/assignments', name: 'Assignments', component: () => import('../views/Assignments.vue') },
  { path: '/transfers', name: 'Transfers', component: () => import('../views/Transfers.vue') },
  { path: '/alerts', name: 'Alerts', component: () => import('../views/Alerts.vue') },
  { path: '/records', name: 'Records', component: () => import('../views/Records.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
