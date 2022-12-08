import { createApp } from 'vue'
import App from './App.vue'

import Router from 'vue-router'
// import LandingPage from '@/components/LandingPage'
import ExamListComponent from '@/components/ExamList'
import ApplicationComponent from '@/components/Application'
import ImpressComponent from '@/components/Impress'
import PrivacyComponent from '@/components/Privacy'


const router = Router.createRouter({
  routes: [{
      path: '/',
      redirect: '/app/exams'
    },
    {
      path: '/app/',
      name: 'router.app',
      component: ApplicationComponent,
      children: [
        {
          path: 'exams',
          name: 'router.altklausurausleihe',
          component: ExamListComponent
        },
        {
          path: 'impress',
          name: 'router.impress',
          component: ImpressComponent
        }
        ,
        {
          path: 'privacy',
          name: 'router.privacy',
          component: PrivacyComponent
        }
      ]
    },
  ],
  mode: 'history'
})

createApp(App).use(router)
