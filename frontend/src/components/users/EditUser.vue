<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Editing user: {{ firstname }} {{ lastname }}</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'userlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="firstname">First Name</label>
        <input type="text" class="form-control" placeholder="First name" v-model="firstname" required />
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

    <div class="card">
      <div class="card-header">Contact Information</div>
      <div class="card-body">
        <div class="form-row">
          <div class="form-group col">
            <label for="password">Password</label>
            <input type="password" class="form-control" placeholder="Password" v-model="password" required />
          </div>
          <div class="form-group col">
            <label for="confirmpassword">Confirm Password</label>
            <input type="password" class="form-control" placeholder="Confirm Password" v-model="confirmpassword" required />
          </div>
        </div>
      </div>
    </div>

    <br />
    <div class="form-row">
      <div class="form-group col">
        <label for="isadmin">Make user a System Administrator</label>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
        <input type="checkbox" v-model="isadmin" required />
      </div>
      <div class="form-group col">
        <button type="submit" class="btn btn-primary btn-sm float-right" @click="UpdateUser">
          Update User
        </button>
      </div>
    </div>
  </section>
</template>

<script>
import moment from "moment";

export default {
  data() {
    return {
      id: 0,
      user: {},
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
      this.$toast.error("Error! " + err);
    });
  },
  methods: {
    UpdateUser() {
      if (this.password !== this.confirmpassword) {
        this.$toast.error("Error! Passwords do not match.");
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
          this.$toast.error("Error! " + attribute + " cannot be " + this.user[attribute]);
          return;
        }
      }

      window.backend.UpdateUser(this.user).then(() => {
        this.$router.push("{name: 'userlist'}");
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>