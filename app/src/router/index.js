import Vue from 'vue';
import Router from 'vue-router';
import UFO from '@/pages/Ufo';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'UFO',
      component: UFO,
    },
  ],
});
