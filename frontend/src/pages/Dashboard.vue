<template>
  <section>
    <h3>Dashboard</h3>
    <hr />

    <div class="card">
      <div class="card-header">
        <strong>Analysis</strong>
      </div>
      <div class="card-body">
        <div class="card-deck">
          <div class="card text-white bg-primary" style="max-width: 18rem">
            <div class="card-header">Daily Sales</div>
            <div class="card-body">
              <h3 class="card-title">₦{{ this.salesTotal.toFixed(2) }}</h3>
            </div>
          </div>
          <div class="card text-white bg-info" style="max-width: 18rem">
            <div class="card-header">Last Sync</div>
            <div class="card-body">
              <h5 class="card-title">Daily Sales</h5>
              <p class="card-text">None for now</p>
            </div>
          </div>
          <div class="card text-white bg-secondary" style="max-width: 18rem">
            <div class="card-header">Sync Errors</div>
            <div class="card-body">
              <h5 class="card-title">Special title treatment</h5>
              <p class="card-text">Count of Errors.</p>
            </div>
          </div>
        </div>
      </div>
      <div class="card-footer">
        <small class="text-muted"
          >Last updated <strong>{{ this.lastSyncTime }}</strong></small
        >
      </div>
    </div>

    <hr />
    <div class="card">
      <div class="card-header">
        <strong>Today's Orders</strong>
      </div>
      <div class="card-body">
        <v-client-table
          ref="ordersTable"
          :columns="columns"
          v-model="orders"
          :options="options"
        >
          <div
            slot="canceled"
            slot-scope="{ row }"
            style="text-transform: capitalize"
          >
            {{ row.canceled }}
          </div>
          <div slot="doctotal" slot-scope="{ row }">₦{{ row.doctotal }}</div>
          <div
            slot="synced"
            slot-scope="{ row }"
            style="text-transform: capitalize"
          >
            {{ row.synced }}
          </div>
          <div
            slot="returned"
            slot-scope="{ row }"
            style="text-transform: capitalize"
          >
            {{ row.returned }}
          </div>
        </v-client-table>
      </div>
      <div class="card-footer">
        <small class="text-muted"
          >Last updated <strong>{{ this.lastSyncTime }}</strong></small
        >
      </div>
    </div>

    <hr />
    <div class="card">
      <div class="card-header">
        <strong>Week's Top Sellers</strong>
      </div>
      <div class="card-body">
        <v-client-table
          ref="itemsTable"
          :columns="itemscolumn"
          v-model="items"
          :options="options"
        >
          <div
            slot="discount"
            slot-scope="{ row }"
            style="text-transform: capitalize"
          >
            {{ row.discount }}%
          </div>
          <div slot="price" slot-scope="{ row }">₦{{ row.price }}</div>
        </v-client-table>
      </div>
      <div class="card-footer">
        <small class="text-muted"
          >Last updated <strong>{{ this.lastSyncTime }}</strong></small
        >
      </div>
    </div>

    <!-- Modal -->
    <div class="modal" id="syncModal" data-backdrop="static" data-keyboard="false" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered modal-lg modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header text-center bg-dark text-white">
            <h5 class="modal-title" id="storeModalLabel">{{this.headerText}}</h5>
            <button type="button" class="close text-white" @click="dismiss" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body" id="staticBackdropLabel">
            <button class="btn btn-primary btn-lg btn-block" type="button" disabled v-if="displaySync">
              <span class="spinner-grow spinner-grow-sm" role="status" aria-hidden="true"></span>
              Loading data from SAP. 
              <br />
              Please wait while we load the latest data.
            </button>
            <retailstore v-else></retailstore> 
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import moment from "moment";
import $ from "jquery";
import RetailStore from "@/pages/RetailStore";

export default {
  data() {
    return {
      items: [],
      orders: [],
      columns: [],
      options: {},
      itemscolumn: [],
      salesTotal: 0.0,
      displaySync: true,
      headerText: "Please wait...",
      lastSyncTime: new Date().toLocaleString(),
      dateColumns: ["created_at", "updated_at", "deleted_at"],
    };
  },
  components: {
    'retailstore': RetailStore
  },
  mounted: async function() {
    // Display the sync Modal
    $('#syncModal').modal('show'); 

    try {
      await window.backend.SAPSync();
      this.endSync();
    } catch (err) {
      if (err.toLowerCase() === "missing endpoint from application") {
       this.handleModals();
      } else {
        this.$toast.error("Error! " + err);
      }
    }
  },
  methods: {
    handleModals() {
      // Hide the sync section and display the store setup component
      this.displaySync = false;
      this.headerText= "Enter Store Details";
    },
    async dismiss() {
      // Hide the store setup component and display the sync section
      this.displaySync = true;
      this.headerText= "Please wait...";

      try {
        await window.backend.SAPSync();
      } catch (err) {
        if (err.toLowerCase() === "missing endpoint from application") {
          this.handleModals();
        } else {
          this.$toast.error("Error! " + err);
        }
      }
    },
    endSync() {
      // Dispose of the sync Modal
      this.handleModals();

      $('#syncModal').modal('hide');

      this.getData();
      // Run every 5 minutes [5 * 60 * 1000 = 300000]
      setTimeout(() => {
        this.$refs.myTable.setLoadingState(true);
        this.getData();
        this.$refs.myTable.setLoadingState(false);
      }, 300000);
    },
    getData() {
      window.backend.Dashboard().then((response) => {
          let total = 0.0;
          if (response.orders !== null) {
            // Today's Orders
            this.orders = [];
            const exempt = [
                "id",
                "items",
                "vatsum",
                "docnum",
                "docentry",
                "comment",
                "payments",
                "cardcode",
                "deleted_at",
                "updated_at",
                "created_by",
                "created_at",
              ],
              keys = Object.keys(response.orders[0]);

            keys.forEach((key) => {
              if (!exempt.includes(key)) {
                this.columns.push(key);
              }
            });

            response.orders.forEach((order) => {
              total += parseFloat(order.doctotal);
              this.orders.push(order);
            });
          }

          if (response.items !== null) {
            // Top Sellers
            this.items = [];
            const exempt = [
                "id",
                "orderid",
                "itemcode",
                "deleted_at",
                "updated_at",
                "created_by",
                "created_at",
              ],
              itemkeys = Object.keys(response.items[0]);

            itemkeys.forEach((key) => {
              if (!exempt.includes(key)) {
                this.itemscolumn.push(key);
              }
            });

            response.items.forEach((order) => {
              this.items.push(order);
            });
          }
          this.salesTotal = total;
          this.lastSyncTime = new Date().toLocaleString();
        },
        (err) => {
          this.$toast.error("Error! " + err);
        }
      );
    },
    formatDate(date) {
      return moment(date).format("DD-MM-YYYY HH:mm:ss");
    },
  },
};
</script>