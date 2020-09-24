import 'core-js/stable';
import 'regenerator-runtime/runtime';
import Vue from 'vue';
import App from './App.vue';
import store from './store';
import router from './router';
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

Vue.config.productionTip = false;
Vue.config.devtools = true;
Vue.use(ElementUI);

import * as Wails from '@wailsapp/runtime';

Wails.Init(() => {
	new Vue({
		store,
		router,
		render: h => h(App)
	}).$mount('#app');
});
