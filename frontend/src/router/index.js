import Vue from 'vue'
import Router from 'vue-router'
import LandingPage from '@/components/LandingPage'
import ExamList from '@/components/ExamList'


Vue.use(Router)

export default new Router({
  routes: [{
      path: '/',
      name: 'LandingPage',
      component: LandingPage
    },
    {
      path: '/exams',
      name: 'ExamList',
      component: ExamList
    }
  ],
  mode: 'history'
})