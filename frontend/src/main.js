import 'core-js/stable';
import 'regenerator-runtime/runtime';
import Vue from 'vue';
import App from './App.vue';
import store from './store';
import router from "./router";
import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";

const options = {
	draggable: false
};
Vue.use(Toast, options);

Vue.config.productionTip = false;
Vue.config.devtools = true;
Vue.mixin({
	methods: {
		userPermission(menuname, action) {
			if (action.toLowerCase() === 'candelete') {
				return false;
			}

			if (this.$store.state.user["acl"] === undefined || action === "") {
				return false;
			}

			return this.$store.state.user["acl"]
				.filter((f) => f.menuname === menuname)
				.map((elem) => ({
					menuname: elem.menuname,
					action: elem[action],
				}))[0].action;
		}
	}
})

Vue.prototype.$appName = "SLOT Systems Retail Managers.";

import * as Wails from '@wailsapp/runtime';

Wails.Init(() => {
	new Vue({
		store,
		router,
		render: h => h(App),
	}).$mount('#app');
});