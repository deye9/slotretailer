import Vue from 'vue';
import store from "../store";
import Router from 'vue-router';

// Module Routes
import Sync from "@/pages/Sync";
import Login from "@/pages/Login";
import Search from "@/pages/Search";

import Dashboard from "@/pages/Dashboard";
import AuditLogs from "@/pages/AuditLogs";
import RetailStore from "@/pages/RetailStore";

import Users from "@/pages/Users";
import NewUser from "@/components/users/NewUser";
import EditUser from "@/components/users/EditUser";
import GetUsers from "@/components/users/GetUsers";

import Orders from "@/pages/Orders";
import NewOrder from "@/components/orders/NewOrder";
import GetOrders from "@/components/orders/GetOrders";
import ReturnOrder from "@/components/orders/ReturnOrder";
import OrderDetails from "@/components/orders/OrderDetails";

import Reports from "@/pages/Reports";
import Report from "@/components/reports/Report";
import NewReport from "@/components/reports/NewReport";
import EditReport from "@/components/reports/EditReport";
import GetReports from "@/components/reports/GetReports";

import Products from "@/pages/Products";
import GetProducts from "@/components/products/GetProducts";
import ProductDetails from "@/components/products/ProductDetails";

import Customers from "@/pages/Customers";
import NewCustomer from "@/components/customers/NewCustomer";
import GetCustomers from "@/components/customers/GetCustomers";
import EditCustomer from "@/components/customers/EditCustomer";

import StockTransfers from "@/pages/StockTransfers";
import NewTransfer from "@/components/transfers/NewTransfer";
import GetTransfers from "@/components/transfers/GetTransfers";
import EditTransfer from "@/components/transfers/EditTransfer";
import TransferDetails from "@/components/transfers/TransferDetails";
import FinalizeTransfer from "@/components/transfers/FinalizeTransfer";

import vSelect from "vue-select";
import {ClientTable} from 'vue-tables-2';

Vue.use(Router);
Vue.component("v-select", vSelect);
Vue.use(ClientTable, {}, false, 'bootstrap4');

const originalPush = Router.prototype.push;
Router.prototype.push = function push(location) {
    return originalPush.call(this, location).catch(err => err);
}

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
            name: "userlist",
            component: GetUsers
        }, {
            path: "new",
            name: "newuser",
            component: NewUser
        },
        {
            path: "edit/:id",
            name: "edituser",
            component: EditUser
        }
        ]
    },
    {
        path: "/orders",
        name: "orders",
        component: Orders,
        children: [{
                path: "/",
                name: "orderlist",
                component: GetOrders
            },
            {
                path: "new",
                name: "neworder",
                component: NewOrder
            },
            {
                path: "return",
                name: "returnorder",
                component: ReturnOrder
            },
            {
                path: "details/:id",
                name: "orderdetail",
                component: OrderDetails
            }
        ]
    },
    {
        path: "/reports",
        name: "reports",
        component: Reports,
        children: [{
                path: "/",
                name: "reportlist",
                component: GetReports
            },
            {
                path: "display/:id",
                name: "reportdetails",
                component: Report
            },
            {
                path: "new",
                name: "newreport",
                component: NewReport
            },
            {
                path: "edit/:id",
                name: "editreport",
                component: EditReport
            }
        ]
    },
    {
        path: "/products",
        name: "products",
        component: Products,
        children: [{
            path: "/",
            name: "productlist",
            component: GetProducts
        },
        {
            path: "details/:id",
            name: "productdetails",
            component: ProductDetails
        }]
    },
    {
        path: "/customers",
        name: "customers",
        component: Customers,
        children: [{
            path: "/",
            name: "customerlist",
            component: GetCustomers
        }, {
            path: "new",
            name: "newcustomer",
            component: NewCustomer
        },
        {
            path: "edit/:id",
            name: "editcustomer",
            component: EditCustomer
        }]
    },
    {
        path: "/transfers",
        name: "transfers",
        component: StockTransfers,
        children: [{
            path: "/",
            name: "transferlist",
            component: GetTransfers
        },
        {
            path: "new",
            name: "newtransfer",
            component: NewTransfer
        },
        {
            path: "edit/:id",
            name: "edittransfer",
            component: EditTransfer
        },
        {
            path: "finalize/:id",
            name: "finalizeTransfer",
            component: FinalizeTransfer
        },
        {
            path: "details/:id",
            name: "transferdetail",
            component: TransferDetails
        }]
    },
    {
        path: "/dashboard",
        name: "dashboard",
        component: Dashboard
    },
    {
        path: "/store/",
        name: "store",
        component: RetailStore
    }, 
    {
        path: "/auditlogs",
        name: "auditlogs",
        component: AuditLogs
    }
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