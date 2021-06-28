<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Registered Users</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'newuser'}" class="btn btn-info btn-sm float-right" v-if="userPermission('users', 'cancreate')">New User</router-link>
      </div>
    </div>
    <hr />

    <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options">
      <div slot="actions" slot-scope="{row}">
        <a class="btn btn-primary btn-sm mr-2" title="Edit Record" @click="displayInfo(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Edit
        </a>
        <a class="btn btn-danger btn-sm" title="Delete Record" @click="removeRow(row, event);" v-if="userPermission('users', 'candelete')">
          <i class="bi bi-trash-fill">&nbsp;</i>
          Delete
        </a>
      </div>
    </v-client-table>
  </section>
</template>

<script>
  export default {
    data() {
      return {
        data: [],
        columns: [],
        options: {},
        dateColumns:['created_at','updated_at', 'deleted_at']
      };
    },
    mounted() {
      this.$refs.myTable.setLoadingState(true);
      window.backend.GetUsers().then((users) => {
        if (users !== null) {
          this.columns = Object.keys(users[0]);
          this.columns.push('actions');

          // Set the dataSource
          this.data = users;
        } else {
          this.$refs.myTable.setLoadingState(false);
        }
        this.isBusy = false;
        }, (err) => {
          this.$toast.error("Error! " + err);
          this.isBusy = false;
      });
    },
    methods: {
      displayInfo(row) {
        const id = row.id;
        this.$router.push({ name: "edituser", params: {id} });
      },
      removeRow(row, index) {
        index = event.srcElement.parentElement.parentElement.parentNode.rowIndex - 1;
        window.backend.RemoveUser(parseInt(row.id)).then(() => {
          // Remove the row from the table
          document.getElementById("myTable").getElementsByTagName('tbody')[0].deleteRow(index);
          
          this.$toast.success("Success! User record has been successfully deleted.");
        }, (err) => {
          this.$toast.error("Error! " + err);
        });
      },
    }
  };
</script>