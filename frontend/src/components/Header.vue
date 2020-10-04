<template>
  <header>
    <b-navbar toggleable="lg" type="dark" variant="dark" fixed="top">
      <b-navbar-brand href="#">
        <router-link :to="{ name: 'dashboard' }" class="nav-link active">
          <img src="../assets/img/slot.png" class="slotLogo" alt="Logo" />
        </router-link>
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item href="#">
            <router-link :to="{ name: 'dashboard' }" class="nav-link active">Dashboard</router-link>
          </b-nav-item>
          <b-nav-item href="#">
            <router-link to="/customers/" class="nav-link active">Customers</router-link>
          </b-nav-item>
          <b-nav-item href="#">
            <router-link to="/orders/" class="nav-link active">Sales Order</router-link>
          </b-nav-item>
          
          <b-nav-item-dropdown text="Administration" class="m-md-2"  v-if="this.$store.state.isAdmin">
            <b-dropdown-item href="#">
              <router-link to="/users/" class="text-dark">User Flow</router-link>
            </b-dropdown-item>
            <b-dropdown-item href="#">
              <router-link to="/products/" class="text-dark">Inventory</router-link>
            </b-dropdown-item>
            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-item href="#">
              <router-link :to="{ name: 'retailstore' }" class="text-dark">Store Details</router-link>
            </b-dropdown-item>
            <b-dropdown-item href="#">
              <router-link :to="{ name: 'sync' }" class="text-dark">Sync Details</router-link>
            </b-dropdown-item>
          </b-nav-item-dropdown>
        </b-navbar-nav>

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-form @submit="search">
          <b-form-input size="sm" class="mr-sm-2" placeholder="Search" v-model="searchTerm"></b-form-input>
          <b-button size="sm" class="my-2 my-sm-0" type="submit">Search</b-button>
        </b-nav-form>

        <b-nav-item href="#" right>
          <router-link :to="{ name: 'login' }" @click.native="logout" v-show="this.$store.state.isLoggedIn" class="nav-link">
            Sign out
          </router-link>
        </b-nav-item>
      </b-navbar-nav>
      </b-collapse>
    </b-navbar>
  </header>
</template>

<script>
export default {
  data() {
    return {
      searchTerm: 'Hi',
    };
  },
  methods: {
    search(evt) {
      evt.preventDefault();
      alert("Search Term is: ", this.searchTerm);
    },
    logout() {
      // Set the user state to an empty object
      this.$store.state.user = {};
      this.$store.state.isLoggedIn = false;
      this.$router.push("/");
    },
  },
};
</script>