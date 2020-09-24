<template>
  <nav class="navbar navbar-dark fixed-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 mr-0 px-3" href="#">
      <router-link :to="{ name: 'dashboard' }" class="nav-link active">
        <img src="../assets/img/slot.png" class="slotLogo" />
      </router-link>
    </a>
    <button
      class="navbar-toggler position-absolute d-md-none collapsed"
      type="button"
      data-toggle="collapse"
      data-target="#sidebarMenu"
      aria-controls="sidebarMenu"
      aria-expanded="false"
      aria-label="Toggle navigation"
    >
      <span class="navbar-toggler-icon"></span>
    </button>
    <input
      class="form-control form-control-dark w-100"
      type="text"
      placeholder="Search"
      aria-label="Search"
    />
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
  </nav>
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