import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { title: 'home' }
    },
    {
      path: '/articles/:title',
      name: 'article-detail',
      component: () => import('../views/ArticleDetailView.vue'),
      meta: { title: 'article_detail' },
      props: (route) => ({ title: route.params.title })
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { title: 'login' }
    },
    {
      path: '/admin',
      redirect: '/admin/articles',
      component: () => import('../views/admin/AdminLayout.vue'),
      meta: { title: 'admin_backend', requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: 'articles',
          name: 'admin-articles',
          component: () => import('../views/admin/ArticleManagement.vue'),
          meta: { title: 'article_management' }
        },
        {
          path: 'articles/new',
          name: 'admin-article-new',
          component: () => import('../views/admin/ArticleEditor.vue'),
          meta: { title: 'new_article' }
        },
        {
          path: 'articles/:id/edit',
          name: 'admin-article-edit',
          component: () => import('../views/admin/ArticleEditor.vue'),
          meta: { title: 'edit_article' },
          props: (route) => ({ id: Number(route.params.id) })
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/NotFoundView.vue'),
      meta: { title: 'page_not_found' }
    }
  ],
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  await authStore.initAuth()
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }
  
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'home' })
    return
  }
  
  if (to.name === 'login' && authStore.isLoggedIn) {
    next({ name: 'home' })
    return
  }
  
  if (to.meta.title) {
    document.title = `${to.meta.title} - Pea Blog`
  }
  
  next()
})

export default router
