import Vue from 'vue'
import Router from 'vue-router'
// import LandingPage from '@/components/LandingPage'
import ExamList from '@/components/ExamList'
import Application from '@/components/Application'
import Impress from '@/components/Impress'
import Privacy from '@/components/Privacy'

Vue.use(Router)

export default new Router({
  routes: [{
      path: '/',
      redirect: '/app/exams'
    },
    {
      path: '/app/',
      name: 'router.app',
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
          component: Impress
        }
        ,
        {
          path: 'privacy',
          name: 'router.privacy',
          component: Privacy
        }
      ]
    },
    
  ],
  mode: 'history'
})