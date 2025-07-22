import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/layout',
      name: 'layout',
      component: () => import('@/views/Layout.vue'),
      children: [
        {
          path: "/profile",
          name: "profile",
          component: () => import("@/views/ProfileView.vue")
        }]
    }],
})

router.beforeEach((to, _, next) => {
  const token = localStorage.getItem('jwt')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.meta.userTypes && Array.isArray(to.meta.userTypes)) {
    // Check user type if required
    const userType = parseInt(localStorage.getItem('userType') || '0')
    if (!to.meta.userTypes.includes(userType)) {
      next('/dashboard')
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
