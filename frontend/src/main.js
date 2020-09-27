import 'core-js/stable';
import 'regenerator-runtime/runtime';
import Vue from 'vue';
import App from './App.vue';
import store from './store';
import router from './router';
import vSelect from "vue-select";
import ElementUI from 'element-ui';
import Toast from "vue-toastification";
import "vue-select/dist/vue-select.css";
import "vue-toastification/dist/index.css";
import 'element-ui/lib/theme-chalk/index.css';

const options = {
	draggable: false
};

Vue.config.productionTip = false;
Vue.config.devtools = true;
Vue.use(ElementUI);
Vue.use(Toast, options);
Vue.component("v-select", vSelect);

import * as Wails from '@wailsapp/runtime';

Wails.Init(() => {
	new Vue({
		store,
		router,
		render: h => h(App)
	}).$mount('#app');
});