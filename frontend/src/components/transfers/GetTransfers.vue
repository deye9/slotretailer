<template>
  <section>
    <div class="row mb-2">
      <div class="col-8">
        <h3>Inventory Transfer Request(s)</h3>
      </div>
      <div class="col-4">
        <router-link :to="{ name: 'newtransfer' }" class="btn btn-info btn-sm float-right">New Transfer Request</router-link>
      </div>
    </div>
    
    <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options">
      <div slot="id" slot-scope="{row}" style="text-transform: capitalize;">
        {{ 'trans-' + row.id }}
      </div>
      <div slot="fromwhs" slot-scope="{row}" style="text-transform: capitalize;">
        {{ storeName(row.fromwhs) }}
      </div>
      <div slot="towhs" slot-scope="{row}" style="text-transform: capitalize;">
        {{ storeName(row.towhs) }}
      </div>
      <div slot="created_at" slot-scope="{row}" style="text-transform: capitalize;">
        {{ formatDate(row.created_at) }}
      </div>
      <div slot="synced" slot-scope="{row}" style="text-transform: capitalize;">
        {{ row.synced }}
      </div>
      <div id="actions" slot="actions" slot-scope="{row}">
        <a class="btn btn-primary btn-sm mr-2" title="Edit Record" @click="displayInfo(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Accept
        </a>
        <a class="btn btn-primary btn-sm mr-2" title="Edit Record" @click="displayInfo(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Edit
        </a>
        <a class="btn btn-primary btn-sm mr-2" title="Order Details" @click="detailsInfo(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Details
        </a>
        <a class="btn btn-danger btn-sm" title="Delete Record" @click="removeRow(row, event);" :v-show="allowDelete">
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
      stores: [],
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
    window.backend.GetTransfers().then((transfers) => {
      if (transfers === null) {
        this.$toast.info("Error! No Inventorty Transfer has been raised.");
        this.$refs.myTable.setLoadingState(false);
        return;
      }

      const exempt = [
          "items",
          "deleted_at",
          "updated_at",
          "created_by",
        ],
        keys = Object.keys(transfers[0]);

      keys.forEach((key) => {
        if (!exempt.includes(key)) {
          this.columns.push(key);
        }
      });
      this.columns.push('actions');

      // Set the dataSource
      this.data = transfers;
      this.$refs.myTable.setLoadingState(false);
    }, (err) => {
      this.$toast.error("Error! " + err);
      this.$refs.myTable.setLoadingState(false);
    });

    // Get all stores
    window.backend.GetStores().then((stores) => {
      this.stores = stores;
      // let defaultStoreID = 0,
        // localStore = this.stores.filter((store) => {
        //   return (
        //     store.name.toLowerCase() === this.localStore.sapkey.toLowerCase()
        //   );
        // })[0];

      // if (localStore.id !== this.transfer.fromwhs) {
      //   defaultStoreID = this.transfer.fromwhs;
      //   document.getElementById("toWHS").disabled = true;
      // } else {
      //   defaultStoreID = this.transfer.towhs;
      //   document.getElementById("fromWHS").disabled = true;
      // }
    }, (err) => {
      this.$toast.error("Error! " + err);
    });
    
  },
  methods: {
    formatDate(date) {
      return moment(date.Time).format("Do of MMMM YYYY");
    },
    storeName(storeID) {
      return this.stores.filter((store) => {
        return (
          store.id === storeID
        );
      })[0].name;
    },
    detailsInfo(row) {
      const id = row.id;
      this.$router.push({ name: "transferdetail", params: {id} });
    },
    displayInfo(row) {
      const id = row.id;
      this.$router.push({ name: "edittransfer", params: {id} });
    },
    removeRow(row, index) {
      index = event.srcElement.parentElement.parentElement.parentNode.rowIndex - 1;
      window.backend.RemoveTransfer(parseInt(row.id)).then(() => {
        // Remove the row from the table
        document.getElementById("myTable").getElementsByTagName('tbody')[0].deleteRow(index);
        
        this.$toast.success("Success! Inventory Transfer has been successfully deleted.");
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  }
};
</script>