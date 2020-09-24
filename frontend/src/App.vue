<template>
  <div id="app" class="container-fluid">
    <Header />
    <div class="container-fluid">
      <div class="row">
        <Sidebar v-if="this.$store.state.isLoggedIn" />
        <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-md-4">
          <div :class="this.$store.state.notify.category" id="notify" style="margin-botton:2px;">
            <span class="closebtn" @click="toggle">&times;</span>
            {{this.$store.state.notify.message}}
          </div>
          <router-view></router-view>
        </main>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script>
import Header from "./components/Header.vue";
import Sidebar from "./components/SideBar.vue";
import Footer from "./components/Footer.vue";

import "./assets/css/main.css";
import "./assets/css/bootstrap.min.css";

export default {
  name: "app",
  components: {
    Header,
    Sidebar,
    Footer,
  },
  methods: {
    toggle() {
      var x = document.getElementById("notify");
      if (x.style.display === "none") { 
        x.style.display = "block"; 
      } else { 
        x.style.display = "none"; 
      }
    },
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