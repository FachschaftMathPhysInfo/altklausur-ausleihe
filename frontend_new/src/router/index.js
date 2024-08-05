import { createRouter, createWebHistory } from 'vue-router'
import Application from '../components/Application.vue'
import ExamList from '../components/ExamList.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/app/exams'
    },
    {
      path: '/app/',
      name: 'router.application',
      component: Application,
      children: [
        {
          path: 'exams',
          name: 'router.altklausurausleihe',
          component: ExamList
        },
        {
          path: 'impress',
          name: 'router.impress',
          component: () => import('../components/Impress.vue')
        }
        ,
        {
          path: 'privacy',
          name: 'router.privacy',
          component: () => import('../components/Privacy.vue')
        }
      ]
    }
  ]
})

export default router
