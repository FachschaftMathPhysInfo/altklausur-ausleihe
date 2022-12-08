import Vue from 'vue'
import Router from 'vue-router'
// import LandingPage from '@/components/LandingPage'
import ExamListComponent from '@/components/ExamList'
import ApplicationComponent from '@/components/Application'
import ImpressComponent from '@/components/Impress'
import PrivacyComponent from '@/components/Privacy'

Vue.use(Router)

export default new Router({
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
