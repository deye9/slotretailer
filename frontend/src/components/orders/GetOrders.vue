<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Sales Order</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'neworder'}" class="btn btn-info btn-sm float-right">New Order</router-link>
      </div>
    </div>
    <hr />

    <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options">
      <div slot="canceled" slot-scope="{row}" style="text-transform: capitalize;">
        {{row.canceled}}
      </div>
      <div slot="doctotal" slot-scope="{row}">
        ₦{{row.doctotal}}
      </div>
      <div slot="synced" slot-scope="{row}" style="text-transform: capitalize;">
        {{row.synced}}
      </div>
      <div slot="returned" slot-scope="{row}" style="text-transform: capitalize;">
        {{row.returned}}
      </div>
      <div id="actions" slot="actions" slot-scope="{row}" :v-show="allowDelete">
        <a class="btn btn-primary btn-sm mr-2" title="Print Sales Order" @click="printOrder(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Print
        </a>
        <a class="btn btn-primary btn-sm mr-2" title="Order Details" @click="displayInfo(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Details
        </a>        
        <a class="btn btn-primary btn-sm mr-2" title="Return Order" @click="returns(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Return
        </a>
        <a class="btn btn-danger btn-sm" title="Delete Sales Order" @click="removeRow(row, event);">
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
      allowDelete: false,
      dateColumns:['created_at','updated_at', 'deleted_at']
    };
  },
  created() {
    this.allowDelete = this.$store.state.isLoggedIn;
  },
  mounted() {
    this.$refs.myTable.setLoadingState(true);
    window.backend.GetOrders().then((orders) => {
      if (JSON.stringify(orders) === "{}" || orders === null) {
        this.$refs.myTable.setLoadingState(false);
        return;
      }

      const exempt = [
          "vatsum",
          "items",
          "docnum",
          "comment",
          "payments",
          "docentry",
          "deleted_at",
          "updated_at",
          "created_by",
          "created_at",
        ],
        keys = Object.keys(orders[0]);

      keys.forEach((key) => {
        if (!exempt.includes(key)) {
          this.columns.push(key);
        }
      });
      this.columns.push('actions');

      // Set the dataSource
      this.data = orders;
      this.$refs.myTable.setLoadingState(false);
    },
    (err) => {
      this.$toast.error("Error! " + err);
      this.$refs.myTable.setLoadingState(false);
    });
  },
  methods: {
    returns(row) {
      const id = row.id;
      this.$router.push({ name: "returnorder", params: {id} });
    },
    printOrder(row) {
      const id = row.id;
      this.$router.push({ name: "orderdetail", params: {id} });
    },
    displayInfo(row) {
      const id = row.id;
      this.$router.push({ name: "orderdetail", params: {id} });
    },
    removeRow(item, index) {
      window.backend.RemoveOrder(parseInt(item.id)).then(() => {
        // Remove the row from the table
        document.getElementById("myTable").getElementsByTagName('tbody')[0].deleteRow(index);
        
        this.$toast.success("Success! Order has been successfully deleted.");
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
  //   var doc = new jsPDF();

  // // Set the document to automatically print via JS
  // doc.autoPrint();
</script>