// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';
import VueProgressBar from 'vue-progressbar';
import Notifications from 'vue-notification';
import App from './App';
import router from './router';
import store from './store';

const progressBarOptions = {
  color: '#2196F3',
  failedColor: '#F44336',
};

Vue.use(VueProgressBar, progressBarOptions);
Vue.use(Notifications);
Vue.config.productionTip = false;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  template: '<App/>',
  components: { App },
});
