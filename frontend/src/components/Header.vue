<template>
  <header>
    <b-navbar type="dark" variant="dark" class="fixed-top w-100" fill>
      <b-navbar-brand href="#">
        <router-link :to="{ name: 'dashboard' }" class="nav-link active">
          <img src="../assets/img/slot.png" class="slotLogo" alt="Logo" />
        </router-link>
      </b-navbar-brand>
      <b-navbar-nav class="w-100" fill>
        <b-nav-item href="#" class="ml-5">
          <router-link :to="{ name: 'dashboard' }" class="nav-link active">Dashboard</router-link>
        </b-nav-item>
        <b-nav-item href="#" class="ml-5">
          <router-link to="/customers/" class="nav-link active">Customers</router-link>
        </b-nav-item>
        <b-nav-item href="#" class="ml-5">
          <router-link to="/orders/" class="nav-link active">Sales Order</router-link>
        </b-nav-item>
        <!-- Navbar dropdowns -->
        <b-dropdown id="dropdown-1" text="Administration" class="m-md-2">
          <b-dropdown-item active>
            <router-link to="/users/" class="nav-link text-black">User Flow</router-link>
          </b-dropdown-item>
          <b-dropdown-item active>
            <router-link to="/products/" class="nav-link">Inventory</router-link>
          </b-dropdown-item>
          <b-dropdown-divider></b-dropdown-divider>
          <b-dropdown-item active>
            <router-link :to="{ name: 'retailstore' }" class="nav-link">Store Details</router-link>
          </b-dropdown-item>
          <b-dropdown-item active>
            <router-link :to="{ name: 'sync' }" class="nav-link">Sync Details</router-link>
          </b-dropdown-item>
        </b-dropdown>
        <b-nav-form class="float-right ml-5">
          <b-form-input class="mr-sm-2" placeholder="Search"></b-form-input>
          <b-button variant="outline-success" class="my-2 my-sm-0" type="submit">Search</b-button>
        </b-nav-form>
        <b-nav-item href="#" class="float-right ml-5" v-show="this.$store.state.isLoggedIn">
          <router-link :to="{ name: 'login' }" @click.native="logout" class="nav-link">
            Sign out
          </router-link>
        </b-nav-item>
      </b-navbar-nav>
    </b-navbar>
    <!-- 
      <ul class="navbar-nav px-3">
        <li class="nav-item text-nowrap">
          <router-link
            :to="{ name: 'login' }"
            @click.native="logout"
            class="nav-link"
            v-show="this.$store.state.isLoggedIn"
          >Sign out</router-link>
        </li>
      </ul> 
    </nav> -->
  </header>
</template>

<script>
export default {
  methods: {
    logout() {
      // Set the user state to an empty object
      this.$store.state.user = {};
      this.$store.state.isLoggedIn = false;
      this.$router.push("/");
    },
  },
};
</script>