<template>
  <b-container id="app" fluid="md">
    <Header v-show="this.$store.state.isLoggedIn" />
    <b-aspect aspect='16:9'>
      <main role="main">
        <router-view></router-view>
        <Footer />
      </main>
    </b-aspect>
  </b-container>
</template>

<script>
import "./assets/css/main.css";
import Header from "./components/Header.vue";
import Footer from "./components/Footer.vue";

export default {
  name: "app",
  components: {
    Header,
    Footer,
  },
};
</script>