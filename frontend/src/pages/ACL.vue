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
        <v-select label="rolename" @input="(val) => itemSelected(val)" :options="roles" v-model="rolename" :clearable="false" placeholder="Kindly select desired Role"></v-select>
      </div>
      <div class="col-4">
        <label>Enter New Role</label>
        <input type="text" class="form-control form-control-sm" placeholder="Enter New Role Name" v-model="newrole" />
      </div>
      <div class="col-auto">
        <br />
        <button type="button" class="btn btn-primary btn-sm mt-2 float-right" @click="NewRole">Create Role</button>
      </div>
    </div>

    <div class="table-responsive">
      <table class="table table-fixed table-bordered table-hover table-sm">
        <caption>
          <h5 style="display:inline-flex;" class="float-left">Define {{this.currentrole}} Permissions.</h5>
        </caption>
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Menu Name</th>
            <th scope="col">Access Menu</th>
            <th scope="col">Create Records</th>
            <th scope="col">Update Records</th>
            <th scope="col">Delete Records</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(detail, i) in roledetail" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>
              {{ titleCase(detail.menuname) }}
            </td>
            <td>
              <div id="radioBtn" class="btn-group">
                <a :class="(detail.canview === true) ? 'btn btn-primary btn-sm active' : 'btn btn-primary btn-sm notActive'" :id="detail.menuname + '_canviewY'"  @click="toggleLink(i, true, detail.menuname + '_canview')">YES</a>
                <a :class="(detail.canview === false) ? 'btn btn-primary btn-sm active' : 'btn btn-primary btn-sm notActive'" :id="detail.menuname + '_canviewN'"  @click="toggleLink(i, false, detail.menuname + '_canview')">NO</a>
              </div>
            </td>
            <td>
              <div id="radioBtn" class="btn-group">
                <a :class="(detail.cancreate === true) ? 'btn btn-primary btn-sm active' : 'btn btn-primary btn-sm notActive'" :id="detail.menuname + '_cancreateY'"  @click="toggleLink(i, true, detail.menuname + '_cancreate')">YES</a>
                <a :class="(detail.cancreate === false) ? 'btn btn-primary btn-sm active' : 'btn btn-primary btn-sm notActive'" :id="detail.menuname + '_cancreateN'"  @click="toggleLink(i, false, detail.menuname + '_cancreate')">NO</a>
              </div>
            </td>
            <td>
              <div id="radioBtn" class="btn-group">
                <a :class="(detail.canupdate === true) ? 'btn btn-primary btn-sm active' : 'btn btn-primary btn-sm notActive'" :id="detail.menuname + '_canupdateY'"  @click="toggleLink(i, true, detail.menuname + '_canupdate')">YES</a>
                <a :class="(detail.canupdate === false) ? 'btn btn-primary btn-sm active' : 'btn btn-primary btn-sm notActive'" :id="detail.menuname + '_canupdateN'"  @click="toggleLink(i, false, detail.menuname + '_canupdate')">NO</a>
              </div>
            </td>
            <td>
              <div id="radioBtn" class="btn-group">
                <a :class="(detail.candelete === true) ? 'btn btn-primary btn-sm active' : 'btn btn-primary btn-sm notActive'" :id="detail.menuname + '_candeleteY'"  @click="toggleLink(i, true, detail.menuname + '_candelete')">YES</a>
                <a :class="(detail.candelete === false) ? 'btn btn-primary btn-sm active' : 'btn btn-primary btn-sm notActive'" :id="detail.menuname + '_candeleteN'"  @click="toggleLink(i, false, detail.menuname + '_candelete')">NO</a>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <button type="button" class="btn btn-primary float-right" @click="saveRole">Save Role</button>
    </div>
  </section>
</template>

<style scoped>
  caption {
    padding-top: .75rem;
    padding-bottom: .75rem;
    color: Black;
    text-align: left;
    caption-side: top;
  }

  #radioBtn .notActive{
    color: #3276b1;
    background-color: #fff;
  }  
</style>

<script>
import $ from "jquery";

export default {
  data() {
    return {
        roles: [],
        newrole: null,
        rolename: null,
        roledetail: [],
        currentrole: null,
    };
  },
  async mounted() {
    // Get all roles
    window.backend.GetRoles().then((roles) => {
        this.roles = roles;
    }, (err) => {
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
      this.rolename = this.newrole;

      // Reset the property back to it's inital value
      this.newrole = null;

      // Fire the event
      this.itemSelected(this.rolename);
    },
    titleCase(str) {
      return str.toLowerCase().split(' ').map(function(word) {
        return word.replace(word[0], word[0].toUpperCase());
      }).join(' ');
    },
    async itemSelected(val) {
      this.currentrole = val;

      // Get role detail
      await window.backend.GetRoleByName(val).then(async (roledetail) => {
        this.roledetail = [];
        if (roledetail !== null) {
          this.roledetail = roledetail;
        }
        var routeNames = this.$router.app._router.options.routes.map(({ name }) => name);

        if (roledetail === null || roledetail.length < routeNames.length) {
          await this.fillUp(routeNames);
        }
      }, (err) => {
        this.roledetail = [];
        this.$toast.error("Error! " + err);
      });
    },
    async toggleLink(rowIndex, title, toggle) {
      switch (title) {
        case true:
          $("#" + toggle + "N").removeClass('active').addClass('notActive');
          $("#" + toggle + "Y").removeClass('notActive').addClass('active');
          break;
          
        case false:
          $("#" + toggle + "Y").removeClass('active').addClass('notActive');
          $("#" + toggle + "N").removeClass('notActive').addClass('active');
          break;

        default:
          break;
      }

      // Set the value based on the key calculated
      this.roledetail[rowIndex][toggle.split("_")[1]] = title;
    },
    async fillUp(routesGotten) {
      let routeNames = this.roledetail.map(({ menuname }) => menuname),
       elmts = routesGotten.filter(f => !routeNames.includes(f));

      elmts.forEach(elem => {
        this.roledetail.push({
          id: 0,
          menuname: elem,
          canview: false,
          canupdate: false,
          cancreate: false,
          candelete: false,
          rolename: this.rolename,
        });
      });
    },
    async saveRole() {
      if (this.roledetail === []) {
        this.$toast.error("Error! You are not permitted to save an EMPTY Access Level");
        return;
      }

      // Send the data to the backend for processing.
      window.backend.SaveRole(this.roledetail).then(() => {
        this.$toast.success("Success! Permissions has been successfully set.");
        this.roledetail = [];
      }, (err) => {
          this.$toast.error("Error! " + err);
      });
    }
  }
};
</script>