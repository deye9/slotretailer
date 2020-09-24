import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const store = new Vuex.Store({
    state: {
        user: {},
        isAdmin: false,
        isLoggedIn: false,
        notify: {
            message: "",
            category: "",
        }
    },
    getters: {},
    mutations: {},
    actions: {}
});

export default store;