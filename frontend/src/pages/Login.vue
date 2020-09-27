<template>
    <form @submit.prevent="login">
      <div class="form-signin">
        <div class="text-center mb-4">
          <img class="mb-4" src="../assets/img/slot.png" />
          <h1 class="h3 mb-3 font-weight-normal">Sign in</h1>
          <p>
            Kindly sign in here to be able to use the application.
            <br />
            <code>
              For security purposes, this application will log you out after 5
              minutes of inactivity.
            </code>
          </p>
        </div>

        <div class="form-label-group">
          <input
            type="email"
            v-model="email"
            class="form-control"
            placeholder="Email address"
            required
            autofocus
          />
          <label for="inputEmail">Email address</label>
        </div>

        <div class="form-label-group">
          <input
            type="password"
            v-model="password"
            class="form-control"
            placeholder="Password"
            required
          />
          <label for="inputPassword">Password</label>
        </div>

        <button class="btn btn-lg btn-primary btn-block" @click="login">Sign in</button>
      </div>
    </form>
</template>

<script>
export default {
  data() {
    return {
      email: "",
      password: "",
    };
  },
  methods: {
    login() {
      const email = this.$data.email,
        password = this.$data.password;

      if (password === "" || email === "") {
        this.$toast.error("Error! Invalid Credentials Supplied.");

        // // Set the user state to an empty object
        this.$store.state.user = {};
        this.$store.state.isLoggedIn = false;
        return false;
      }

      window.backend.Login(email, password).then((result) => {
        if (result.id !== undefined) {
          this.$store.state.user = result;
          this.$store.state.isLoggedIn = true;
          this.$store.state.isAdmin = result.isadmin;
          this.$router.push("/dashboard");
        } else {
          this.$toast.error("Error! Invalid login Credentials.");
        }
      },(err) => {
          this.$toast.success("Error! " + err);
        });
    },
  },
};
</script>