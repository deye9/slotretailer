<template>
  <div id="app">
    <Header v-show="this.$store.state.isLoggedIn" />
    <main role="main" class="container mt-5">
      <router-view />
      <br /><br />
    </main>
    <footer class="footer fixed-bottom bg-secondary text-light text-center">
      <div class="container">&copy; SLOT Systems Limited. 2017-{{ new Date().getFullYear() }}</div>
    </footer>
  </div>
</template>

<script>
import Header from "./components/Header.vue";
import "./assets/css/main.css";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.min.js";

export default {
  name: "app",
  components: {
    Header,
  },
  mounted() {
    document.title = "SLOT Systems Retail Manager.";
  },
};
</script>
  