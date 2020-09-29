<template>
  <b-container id="app" fluid>
    <Header v-show="this.$store.state.isLoggedIn" />
    <main role="main" style="margin-top:50px;">
      <router-view></router-view>
      <Footer />
    </main>
  </b-container>
</template>

<script>
import Header from "./components/Header.vue";
import Footer from "./components/Footer.vue";

import "./assets/css/main.css";
import "./assets/css/bootstrap.min.css";

export default {
  name: "app",
  components: {
    Header,
    Footer,
  },
};
</script>

<style lang="css">
@import "./assets/css/bootstrap.min.css";
.bd-placeholder-img {
  font-size: 1.125rem;
  text-anchor: middle;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

@media (min-width: 768px) {
  .bd-placeholder-img-lg {
    font-size: 3.5rem;
  }
}
@import "./assets/css/main.css";
</style>