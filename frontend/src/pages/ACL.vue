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
      <div class="form-group col">
        <label for="rolename">Role Name</label>
        <v-select label="rolename" @search="RoleDetails" :options="roles" v-model="rolename" :clearable="false" placeholder="Kindly select desired Role"></v-select>
      </div>
      <div class="form-group col">
        <label>Enter New Role</label>
        <input type="text" class="form-control form-control-sm" placeholder="Enter New Role Name" v-model="this.newrole" />
        <button type="button" class="btn btn-primary btn-sm float-right" @click="NewRole">Save</button>
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
      
    },
  }
};
</script>