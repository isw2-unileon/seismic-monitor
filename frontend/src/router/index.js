import { createRouter, createWebHistory } from 'vue-router'
import SeismicMap from '../components/SeismicMap.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      name: 'map',
      component: SeismicMap,
      meta: { requiresAuth: true }
    },
    {
      path: '/account',
      name: 'account',
      component: () => import('../views/AccountView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/history',
      name: 'history',
      component: () => import('../views/HistoryView.vue'),
      meta: { requiresAuth: true }
    },
    {
      // Intercepta URLs no válidas
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
})

// Guardia de Navegación Moderno (Vue Router v4)
// Observa que ya no declaramos el parámetro 'next'
router.beforeEach((to, from) => {
  const isAuthenticated = !!localStorage.getItem('auth_token')

  if (to.meta.requiresAuth && !isAuthenticated) {
    // Si requiere sesión y no hay token, aborta y devuelve a login
    return { name: 'login' }
  } 
  
  if (to.name === 'login' && isAuthenticated) {
    // Si ya tiene sesión e intenta ir al login, devuélvelo al mapa
    return { name: 'map' }
  }

  // Si no se cumple ninguna restricción, retorna true para permitir el paso
  return true
})

export default router