import Vue from 'vue';
import store from "../store";
import Router from 'vue-router';

// Module Routes
import Sync from "@/pages/Sync";
import Login from "@/pages/Login";
import Search from "@/pages/Search";

import Users from "@/pages/Users";
import NewUser from "@/components/users/NewUser";
import EditUser from "@/components/users/EditUser";
import GetUsers from "@/components/users/GetUsers";

import Reports from "@/pages/Reports";
import Report from "@/components/reports/Report";
import NewReport from "@/components/reports/NewReport";
import EditReport from "@/components/reports/EditReport";
import GetReports from "@/components/reports/GetReports";

import Customers from "@/pages/Customers";
import NewCustomer from "@/components/customers/NewCustomer";
import EditCustomer from "@/components/customers/EditCustomer";
import GetCustomers from "@/components/customers/GetCustomers";

import Orders from "@/pages/Orders";
import NewOrder from "@/components/orders/NewOrder";
import GetOrders from "@/components/orders/GetOrders";
import ReturnOrder from "@/components/orders/ReturnOrder";
import OrderDetails from "@/components/orders/OrderDetails";

import Products from "@/pages/Products";
import GetProducts from "@/components/products/GetProducts";
import ProductDetails from "@/components/products/ProductDetails";

import Dashboard from "@/pages/Dashboard";
import RetailStore from "@/pages/RetailStore";

// Library Imports
import {
    BootstrapVue,
    ModalPlugin,
    BootstrapVueIcons
} from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'bootstrap/dist/css/bootstrap.css';
import {
    DataTables,
    DataTablesServer
} from "vue-data-tables";
import 'pc-bootstrap4-datetimepicker/build/css/bootstrap-datetimepicker.css';
import moment from 'moment';

// set language to EN
import lang from 'element-ui/lib/locale/lang/en';
import locale from 'element-ui/lib/locale';

locale.use(lang);

Vue.use(Router);
Vue.use(BootstrapVue);
Vue.use(ModalPlugin);
Vue.use(BootstrapVueIcons);
Vue.use(DataTables);
Vue.use(DataTablesServer);

Vue.filter("capitalize", str => str.charAt(0).toUpperCase() + str.slice(1));
Vue.filter("moment", date => moment(date).format('DD MMMM YYYY'));
Vue.filter("moment2", date => moment(date).format('DD-MM-YYYY'));

const router = new Router({
    mode: 'history',
    linkExactActiveClass: 'active',
    routes: [
        {
            path: "/",
            name: "login",
            component: Login
        },
        {
            path: "/search/:q",
            name: "search",
            component: Search
        },
        {
            path: "/sync",
            name: "sync",
            component: Sync
        },
        {
            path: "/users",
            name: "users",
            component: Users,
            children: [{
                    path: "/",
                    component: GetUsers
                },
                {
                    path: "new",
                    component: NewUser
                },
                {
                    path: "edit/:id",
                    component: EditUser
                }
            ]
        },
        {
            path: "/customers",
            name: "customers",
            component: Customers,
            children: [{
                    path: "/",
                    component: GetCustomers
                },
                {
                    path: "new",
                    component: NewCustomer
                },
                {
                    path: "edit/:id",
                    component: EditCustomer
                }
            ]
        },
        {
            path: "/orders",
            name: "orders",
            component: Orders,
            children: [{
                    path: "/",
                    component: GetOrders
                },
                {
                    path: "new",
                    component: NewOrder
                },
                {
                    path: "return/:id",
                    component: ReturnOrder
                },
                {
                    path: "details/:id",
                    component: OrderDetails
                }
            ]
        },
        {
            path: "/products",
            name: "products",
            component: Products,
            children: [{
                    path: "/",
                    component: GetProducts
                },
                {
                    path: "details/:id",
                    component: ProductDetails
                }
            ]
        },
        {
            path: "/store",
            name: "retailstore",
            component: RetailStore
        },
        {
            path: "/dashboard",
            name: "dashboard",
            component: Dashboard
        },
        {
            path: "/reports",
            name: "reports",
            component: Reports,
            children: [
                {
                    path: "/",
                    component: GetReports
                },
                {
                    path: "display/:id",
                    component: Report
                },
                {
                    path: "new",
                    component: NewReport
                },
                {
                    path: "edit/:id",
                    component: EditReport
                }
            ]
        },
    ]
});

router.beforeEach((to, from, next) => {
    // Redirect all unauthenticated users to the login page
    if (to.name !== 'login' && !store.state.isLoggedIn) next({
        name: 'login'
    })
    else next()
});

export default router;