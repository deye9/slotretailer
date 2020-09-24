<template>
  <section>
    <h3>Editing User [{{firstname}} {{lastname}} : {{email}}]</h3>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="firstname">First Name</label>
        <input
          type="text"
          class="form-control"
          placeholder="First name"
          v-model="firstname"
          required
        />
      </div>
      <div class="form-group col">
        <label for="lastname">Last Name</label>
        <input type="text" class="form-control" placeholder="Last name" v-model="lastname" required />
      </div>
    </div>
    <div class="form-group">
      <label for="email">Email</label>
      <input type="email" class="form-control" placeholder="Email Address" v-model="email" required />
    </div>
    <div class="form-row">
      <div class="form-group col">
        <label for="password">Password</label>
        <input
          type="password"
          class="form-control"
          placeholder="Password"
          v-model="password"
          required
        />
      </div>
      <div class="form-group col">
        <label for="confirmpassword">Confirm Password</label>
        <input
          type="password"
          class="form-control"
          placeholder="Confirm Password"
          v-model="confirmpassword"
          required
        />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group col">
        <label for="isadmin">Make user a System Administrator</label>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
        <input type="checkbox" v-model="isadmin" required />
      </div>
      <div class="form-group col">
        <button type="submit" class="btn btn-primary float-right" @click="UpdateUser">Update User</button>
      </div>
    </div>
  </section>
</template>

<script>
import moment from 'moment';

export default {
  data() {
    return {
      user: {},
      id: 0,
      email: null,
      isadmin: false,
      password: null,
      lastname: null,
      firstname: null,
      created_by: null,
      confirmpassword: null,
    };
  },
  mounted() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.GetUser(parseInt(this.id)).then((user) => {
        this.email = user.email;
        this.isadmin = user.isadmin;
        this.password = user.password;
        this.lastname = user.lastname;
        this.firstname = user.firstname;
        this.created_by = user.created_by;
        this.confirmpassword = user.password;
      },
      (err) => {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! " + err;
      }
    );
  },
  methods: {
    UpdateUser() {
      if (this.password !== this.confirmpassword) {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! Passwords do not match.";
        return;
      }

      this.user = {
        id: this.id,
        email: this.email,
        isadmin: this.isadmin,
        password: this.password,
        lastname: this.lastname,
        firstname: this.firstname,
        updated_at: moment().format(),
        created_by: this.$store.state.user.id,
      };

      // Validate the payload.
      for (var attribute in this.user) {
        if (this.user[attribute] === "" || this.user[attribute] === null) {
          this.$store.state.notify.category = "error";
          this.$store.state.notify.message =
            "Error! " + attribute + " cannot be " + this.user[attribute];
          return;
        }
      }
      
      window.backend.UpdateUser(this.user).then(
        () => {
          this.$router.push("/users/");
        },
        (err) => {
          this.$store.state.notify.category = "error";
          this.$store.state.notify.message = "Error! " + err;
        }
      );
    },
  },
};
</script>