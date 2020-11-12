<template>
  <section class="login text-center">
    <img src="../assets/images/slot.png" alt="Slot Logo" />
    <div class="login-form">
      <div class="main-div">
        <div class="panel">
          <h2>Welcome</h2>
          <p>Please enter your email and password</p>
        </div>
        <form id="Login" @submit.prevent="login">
          <div class="form-group">
            <input
              type="email"
              v-model="email"
              class="form-control form-control-sm"
              placeholder="Email Address"
              required
              autofocus
            />
          </div>
          <div class="form-group">
            <input
              type="password"
              v-model="password"
              class="form-control form-control-sm"
              placeholder="Password"
              required
            />
          </div>
          <button @click="login" class="btn btn-lg btn-sm btn-primary btn-block mt-3">
            Login
          </button>
        </form>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      email: "",
      password: "",
    };
  },
    mounted() {
    window.backend.GetStore().then((store) => {
      // Keep the store details in vuex
      this.$store.state.userStore = store;
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });
  },
  methods: {
    login() {
      const email = this.$data.email,
        password = this.$data.password;

      if (password === "" || email === "") {
        this.$toast.error("Error! Invalid Credentials Supplied.");

        // Set the user state to an empty object
        this.$store.state.user = {};
        this.$store.state.isLoggedIn = false;
        return false;
      }

      window.backend.Login(email, password).then((result) => {
        if (result.id !== undefined) {
          this.$store.state.user = result;
          this.$store.state.isLoggedIn = true;
          this.$store.state.isAdmin = result.isadmin;

          this.$router.push({ name: "dashboard" });
        } else {
          this.$toast.error("Error! Invalid login Credentials.");
        }
      },
      (err) => {
        this.$toast.success("Error! " + err);
      });
    },
  },
};
</script>

<style scoped>
.panel h2 {
  color: #444444;
  font-size: 18px;
  margin: 0 0 8px 0;
}

.panel p {
  color: #777777;
  font-size: 14px;
  margin-bottom: 30px;
  line-height: 24px;
}

.login-form .form-control {
  background: #f7f7f7 none repeat scroll 0 0;
  border: 1px solid #d4d4d4;
  border-radius: 4px;
  height: 50px;
  line-height: 50px;
}

.main-div {
  background: #ffffff none repeat scroll 0 0;
  border-radius: 2px;
  margin: 10px auto 30px;
  max-width: 38%;
  padding: 50px 70px 70px 71px;
}

.login-form .form-group {
  margin-bottom: 10px;
}

.login-form {
  text-align: center;
}

.login-form .btn.btn-primary {
  color: #ffffff;
  height: 50px;
  line-height: 50px;
  padding: 0;
}
</style>