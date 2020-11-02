<template>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark fixed-top">
        <a class="navbar-brand" href="#">
            <img src="../assets/images/slot.png" width="100" height="30" alt="" loading="lazy">
        </a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav mr-auto">
                <li class="nav-item">
                    <router-link :to="{ name: 'dashboard' }" class="nav-link">
                        Dashboard
                    </router-link>
                </li>
                <li class="nav-item">
                    <router-link :to="{ name: 'customerlist' }" class="nav-link">
                        Customers
                    </router-link>
                </li>
                <li class="nav-item">
                    <router-link :to="{ name: 'orderlist' }" class="nav-link">
                        Sales Order
                    </router-link>
                </li>
                <li class="nav-item">
                    <router-link :to="{ name: 'orders' }" class="nav-link">
                        Stock Transfers
                    </router-link>
                </li>
                <li class="nav-item">
                    <router-link :to="{ name: 'reportlist' }" class="nav-link">
                        Reports
                    </router-link>
                </li>
                <li class="nav-item dropdown" v-show="this.$store.state.isLoggedIn">
                    <a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        Administration
                    </a>
                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <router-link :to="{ name: 'userlist' }" class="dropdown-item">
                            User Flow
                        </router-link>
                        <router-link :to="{ name: 'productlist' }" class="dropdown-item">
                            Inventory
                        </router-link>
                        <router-link :to="{ name: 'store' }" class="dropdown-item">
                            Store Details
                        </router-link>
                        <div class="dropdown-divider"></div>
                        <router-link :to="{ name: 'auditlogs' }" class="dropdown-item">
                            Audit Logs
                        </router-link>
                        <router-link :to="{ name: 'sync' }" class="dropdown-item">
                            Sync Details
                        </router-link>
                    </div>
                </li>
            </ul>
            <form class="form-inline my-2 my-lg-0 mr-2">
                <input class="form-control mr-sm-2" type="search" placeholder="Search" aria-label="Search">
                <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
            </form>
            <router-link :to="{ name: 'login' }" @click.native="logout" v-show="this.$store.state.isLoggedIn" style="float:right;" class="nav-link">
                <i class="fa fa-fw fa-sign-out"></i> Sign out
            </router-link>
        </div>
    </nav>
</template>

<script>
export default {
  data() {
      return {
        searchTerm: '',
      }
    },
  methods: {
    search(evt) {
      evt.preventDefault();

      if (this.searchTerm === '') {
        this.$toast.error("Error! You cannot search for an empty string.");
        return;
      }

      if (this.searchTerm.length <= 2) {
        this.$toast.info("Info! You need 3 characters or more for a search.");
        return;
      }

      this.$router.push("/search/" + this.searchTerm);
    },
    logout() {
      // Set the user state to an empty object
      this.$store.state.user = {};
      this.$store.state.isLoggedIn = false;
      this.$router.push({name: "login"});
    },
  },
};
</script>

<style scoped></style>