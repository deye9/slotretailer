import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const store = new Vuex.Store({
    state: {
        user: {},
        userStore: {},
        isAdmin: false,
        reportTitle: '',
        isLoggedIn: false,
    },
    getters: {},
    mutations: {},
    actions: {}
});

export default store;