<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>New User</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'userlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="firstname">First Name</label>
        <input type="text" class="form-control form-control-sm" placeholder="First name" v-model="firstname" required />
      </div>
      <div class="form-group col">
        <label for="lastname">Last Name</label>
        <input type="text" class="form-control form-control-sm" placeholder="Last name" v-model="lastname" required />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label for="email">Email</label>
        <input type="email" class="form-control form-control-sm" placeholder="Email Address" v-model="email" required />
      </div>

      <div class="form-group col">
        <label for="role">Role</label>
        <select class="form-control form-control-sm" v-model="role" >
          <option :key="r.id" :value="r.id" v-for="r in roles">{{ r.rolename }}</option>
        </select>
      </div>
    </div>

    <div class="card">
      <div class="card-header">Contact Information</div>
      <div class="card-body">
        <div class="form-row">
          <div class="form-group col">
            <label for="password">Password</label>
            <input type="password" class="form-control form-control-sm" placeholder="Password" v-model="password" required />
          </div>
          <div class="form-group col">
            <label for="confirmpassword">Confirm Password</label>
            <input type="password" class="form-control form-control-sm" placeholder="Confirm Password" v-model="confirmpassword" required />
          </div>
        </div>
      </div>
    </div>

    <br />

    <div class="form-row">
      <div class="form-group col">
        <button type="submit" class="btn btn-primary btn-sm float-right" @click="RegisterUser">
          Register User
        </button>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      user: {},
      roles: [],
      role: null,
      email: null,
      isValid: true,
      password: null,
      lastname: null,
      firstname: null,
      created_by: null,
      confirmpassword: null,
    };
  },
  async mounted() {
    // Get all roles alongside their ID's
    window.backend.GetRoleswithID().then((roles) => {
        this.roles = roles;
    }, (err) => {
        this.$toast.error("Error! " + err);
    });
  },  
  methods: {
    RegisterUser() {
      if (this.password !== this.confirmpassword) {
        this.$toast.error("Error! Passwords do not match.");
        return;
      }

      this.user = {
        role: this.role,
        email: this.email,
        password: this.password,
        lastname: this.lastname,
        firstname: this.firstname,
        created_by: this.$store.state.user.id,
      };

      // Validate the payload.
      for (var attribute in this.user) {
        if (this.user[attribute] === "" || this.user[attribute] === null) {
          this.isValid = false;
          this.$toast.error("Error! " + attribute + " cannot be " + this.user[attribute]);
          return;
        }
      }

      window.backend.NewUser(this.user).then(() => {
        this.$router.push({name: 'userlist'});
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>