<template>
  <section>
    <div class="row mb-2">
      <div class="col-8">
        <h3>Registered Customers</h3>
      </div>
      <div class="col-4">
        <router-link :to="{ name: 'newcustomer' }" class="btn btn-info btn-sm float-right" v-if="userPermission('customers', 'cancreate')">New Customer</router-link>
      </div>
    </div>
    
    <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options">
      <div id="actions" slot="actions" slot-scope="{row}">
        <a class="btn btn-primary btn-sm mr-2" title="Edit Record" @click="displayInfo(row)" v-if="userPermission('customers', 'canupdate')">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Edit
        </a>
        <a class="btn btn-danger btn-sm" title="Delete Record" @click="removeRow(row, event);" v-if="userPermission('customers', 'candelete')">
          <i class="bi bi-trash-fill">&nbsp;</i>
          Delete
        </a>
      </div>
    </v-client-table>
  </section>
</template>

<script>
import moment from "moment";

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
    window.backend.GetCustomers().then((customers) => {
      if (JSON.stringify(customers) === "{}") {
        this.$refs.myTable.setLoadingState(false);
        return;
      }

      const exempt = [
          "city",
          "address",
          "cardcode",
          "deleted_at",
          "created_at",
          "updated_at",
          "created_by",
        ],
        keys = Object.keys(customers[0]);

      keys.forEach((key) => {
        if (!exempt.includes(key)) {
          this.columns.push(key);
        }
      });
      this.columns.push('actions');

      // Set the dataSource
      this.data = customers;
      this.$refs.myTable.setLoadingState(false);
    }, (err) => {
      this.$toast.error("Error! " + err);
      this.$refs.myTable.setLoadingState(false);
    });
  },
  methods: {
    formatDate(date) {
      return moment(date).format("DD-MM-YYYY HH:mm:ss");
    },
    displayInfo(row) {
      const id = row.id;
      this.$router.push({ name: "editcustomer", params: {id} });
    },
    removeRow(row, index) {
      index = event.srcElement.parentElement.parentElement.parentNode.rowIndex - 1;
      window.backend.RemoveCustomer(parseInt(row.id)).then(() => {
        // Remove the row from the table
        document.getElementById("myTable").getElementsByTagName('tbody')[0].deleteRow(index);
        
        this.$toast.success("Success! Customer record has been successfully deleted.");
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  }
};
</script>