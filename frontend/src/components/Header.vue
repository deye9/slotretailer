<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-dark fixed-top">
    <a class="navbar-brand" href="#">
      <img
        src="../assets/images/slot.png"
        width="100"
        height="30"
        alt=""
        loading="lazy"
      />
    </a>
    <button
      class="navbar-toggler"
      type="button"
      data-toggle="collapse"
      data-target="#navbarSupportedContent"
      aria-controls="navbarSupportedContent"
      aria-expanded="false"
      aria-label="Toggle navigation"
    >
      <span class="navbar-toggler-icon"></span>
    </button>

    <div class="collapse navbar-collapse" id="navbarSupportedContent">
      <ul class="navbar-nav mr-auto">
        <li class="nav-item">
          <router-link
            :to="{ name: 'dashboard' }"
            v-if="canView('dashboard')"
            class="nav-link"
          >
            Dashboard
          </router-link>
        </li>
        <li class="nav-item">
          <router-link
            :to="{ name: 'customerlist' }"
            v-if="canView('customers')"
            class="nav-link"
          >
            Customers
          </router-link>
        </li>
        <li class="nav-item">
          <router-link
            :to="{ name: 'orderlist' }"
            v-if="canView('orders')"
            class="nav-link"
          >
            Sales Order
          </router-link>
        </li>
        <li class="nav-item">
          <router-link
            :to="{ name: 'reportlist' }"
            v-if="canView('reports')"
            class="nav-link"
          >
            Reports
          </router-link>
        </li>
        <li class="nav-item dropdown">
          <a
            class="nav-link dropdown-toggle"
            href="#"
            id="navbarDropdown"
            role="button"
            data-toggle="dropdown"
            aria-haspopup="true"
            aria-expanded="false"
          >
            Administration
          </a>
          <div class="dropdown-menu" aria-labelledby="navbarDropdown">
            <router-link
              :to="{ name: 'userlist' }"
              v-if="canView('users')"
              class="dropdown-item"
            >
              User Flow
            </router-link>
            <router-link
              :to="{ name: 'productlist' }"
              v-if="canView('products')"
              class="dropdown-item"
            >
              Inventory
            </router-link>
            <div class="dropdown-divider"></div>
            <router-link
              :to="{ name: 'store' }"
              v-if="canView('store')"
              class="dropdown-item"
            >
              Store Details
            </router-link>
            <router-link
              :to="{ name: 'transferlist' }"
              v-if="canView('transfers')"
              class="dropdown-item"
            >
              Stock Transfers
            </router-link>
            <div class="dropdown-divider"></div>
            <router-link
              :to="{ name: 'auditlogs' }"
              v-if="canView('auditlogs')"
              class="dropdown-item"
            >
              Audit Logs
            </router-link>
            <router-link
              :to="{ name: 'sync' }"
              v-if="canView('sync')"
              class="dropdown-item"
            >
              Sync Details
            </router-link>
            <div class="dropdown-divider"></div>
            <router-link
              :to="{ name: 'acl' }"
              v-if="canView('acl')"
              class="dropdown-item"
            >
              Access Control
            </router-link>
          </div>
        </li>
      </ul>
      <form
        class="form-inline my-2 my-lg-0 mr-2"
        @submit.prevent="search"
        v-if="canView('search')"
      >
        <input
          class="form-control mr-sm-2"
          type="search"
          placeholder="Search"
          aria-label="Search"
          v-model="searchTerm"
        />
        <button class="btn btn-outline-success my-2 my-sm-0" type="submit">
          Search
        </button>
      </form>
      <router-link
        :to="{ name: 'login' }"
        @click.native="logout"
        v-show="this.$store.state.isLoggedIn"
        style="float: right"
        class="nav-link"
      >
        <i class="fa fa-fw fa-sign-out"></i> Sign out
      </router-link>
    </div>
  </nav>
</template>

<script>
export default {
  data() {
    return {
      searchTerm: "",
    };
  },
  methods: {
    search(evt) {
      evt.preventDefault();

      if (this.searchTerm === "") {
        this.$toast.error("Error! You cannot search for an empty string.");
        return;
      }

      if (this.searchTerm.length <= 2) {
        this.$toast.info("Info! You need 3 characters or more for a search.");
        return;
      }

      const q = this.searchTerm;
      this.$router.push({ name: "search", params: { q } });
    },
    logout() {
      // Set the user state to an empty object
      this.$store.state.user = {};
      this.$store.state.isLoggedIn = false;
      this.$router.push({ name: "login" });
    },
    canView(link) {
        if (this.$store.state.user["acl"] === undefined) {
            return false;
        }

        return this.$store.state.user["acl"]
          .filter((f) => f.menuname === link)
          .map((elem) => ({
            menuname: elem.menuname,
            canview: elem.canview,
          }))[0].canview;
    },
  },
};
</script>