<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Access Control List</h3>
      </div>
      <div class="col-4">
        <router-link :to="{ name: 'dashboard' }" class="btn btn-info btn-sm float-right">Home</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="col-6">
        <label for="rolename">Role Name</label>
        <v-select label="rolename" @search="RoleDetails" :options="roles" v-model="rolename" :clearable="false" placeholder="Kindly select desired Role"></v-select>
      </div>
      <div class="col-4">
        <label>Enter New Role</label>
        <input type="text" class="form-control form-control-sm" placeholder="Enter New Role Name" v-model="newrole" />
      </div>
      <div class="col-auto">
        <br />
        <button type="button" class="btn btn-primary float-right" @click="NewRole">Save</button>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
        roles: [],
        newrole: null,
        rolename: null,
    };
  },
  async mounted() {
    // Get all roles
    window.backend.GetRoles().then((roles) => {
        this.roles = roles;
    },
    (err) => {
        this.$toast.error("Error! " + err);
    });
  },
  methods: {
    NewRole() {
      if (this.newrole === null) {
        this.$toast.error("Error! Role Name cannot be blank.");
        return;
      }

      // Add the new Role to the list of available roles
      this.roles.push(this.newrole);

      // Set the selected Index to the newly added Item

      // Reset the property back to it's inital value
      this.newrole = null;
    },
    /**
     * Triggered when the search text changes.
     *
     * @param search  {String}    Current search text
     * @param loading {Function}	Toggle loading class
     */
    RoleDetails(search, loading) {
      loading(true);
      if (search.length >= 3) {
        alert(search);
      }
      loading(false);
      return;
    }
  }
};
</script>