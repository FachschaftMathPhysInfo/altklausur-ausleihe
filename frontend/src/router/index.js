import { createApp } from "vue";
import App from "./App.vue";

import { createRouter, createWebHistory } from 'vue-router'
// import LandingPage from '@/components/LandingPage'
import ExamListComponent from "@/components/ExamList";
import ApplicationComponent from "@/components/Application";
import ImpressComponent from "@/components/Impress";
import PrivacyComponent from "@/components/Privacy";

const routes = [
  {
    path: "/",
    redirect: "/app/exams",
  },
  {
    path: "/app/",
    name: "router.app",
    component: ApplicationComponent,
    children: [
      {
        path: "exams",
        name: "router.altklausurausleihe",
        component: ExamListComponent,
      },
      {
        path: "impress",
        name: "router.impress",
        component: ImpressComponent,
      },
      {
        path: "privacy",
        name: "router.privacy",
        component: PrivacyComponent,
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes: routes,
});

const app = createApp(App);
app.use(router);
app.mount("#app");
