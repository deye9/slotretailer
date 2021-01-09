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

    <ul class="nav nav-tabs">
      <li class="nav-item" role="presentation">
        <a class="nav-link active" id="pending-tab" data-toggle="tab" href="#pending" role="tab" aria-controls="pending" aria-selected="true">Pending Request(s)</a>
      </li>
      <li class="nav-item" role="presentation">
        <a class="nav-link" id="incoming-tab" data-toggle="tab" href="#incoming" role="tab" aria-controls="incoming" aria-selected="false">Incoming Request(s)</a>
      </li>
      <li class="nav-item" role="presentation">
        <a class="nav-link" id="outgoing-tab" data-toggle="tab" href="#outgoing" role="tab" aria-controls="outgoing" aria-selected="false">Outgoing Transfers</a>
      </li>
    </ul>

    <div class="tab-content" id="transfersTab">
      <div class="tab-pane fade active mt-4" id="pending" role="tabpanel" aria-labelledby="pending-tab">
        <v-client-table ref="pendingTransfers" :columns="columns" v-model="pending" :options="options">
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
        </v-client-table>
      </div>
      <div class="tab-pane fade show mt-4" id="incoming" role="tabpanel" aria-labelledby="incoming-tab">
        <v-client-table ref="myTable" id="myTable" :columns="incomingcolumns" v-model="incoming" :options="options">
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
            <a class="btn btn-primary btn-sm mr-2" title="Edit Record" @click="detailsInfo(row)">
              <i class="bi bi-pencil-fill">&nbsp;</i>
              Process
            </a>
            <a class="btn btn-primary btn-sm mr-2" title="Order Details" @click="displayInfo(row)">
              <i class="bi bi-pencil-fill">&nbsp;</i>
              Details
            </a>
          </div>
        </v-client-table>
      </div>
      <div class="tab-pane fade mt-4" id="outgoing" role="tabpanel" aria-labelledby="outgoing-tab">
        <v-client-table ref="outgoingTransfers" :columns="columns" v-model="outgoing" :options="options">
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
        </v-client-table>
      </div>
    </div>
  </section>
</template>

<script>
import moment from "moment";

export default {
  data() {
    return {
      stores: [],
      columns: [],
      pending: [],
      incoming: [],
      outgoing: [],
      options: {},
      allowDelete: false,
      incomingcolumns: [],
      dateColumns:['created_at','updated_at', 'deleted_at']
    };
  },
  created() {
    this.allowDelete = this.$store.state.isLoggedIn;
  },
  mounted() {
    this.$refs.myTable.setLoadingState(true);
    this.$refs.pendingTransfers.setLoadingState(true);
    this.$refs.outgoingTransfers.setLoadingState(true);

    window.backend.GetTransfers().then(async (transfers) => {
      if (transfers === null) {
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
          this.incomingcolumns.push(key);
        }
      });
      this.incomingcolumns.push('actions');

      // Set the dataSource
      this.pending = await transfers.filter((transfer) => {
        return (
          transfer.status.toLowerCase() === "pending"
        )
      });
      this.incoming = await transfers.filter((transfer) => {
        return (
          transfer.status.toLowerCase() === "incoming"
        )
      });
      this.outgoing = await transfers.filter((transfer) => {
        return (
          transfer.status.toLowerCase() === "accepted" || transfer.status.toLowerCase() === "approved"
        )
      });
      this.$refs.myTable.setLoadingState(false);
      this.$refs.pendingTransfers.setLoadingState(false);
      this.$refs.outgoingTransfers.setLoadingState(false);
    }, (err) => {
      this.$toast.error("Error! " + err);
      this.$refs.myTable.setLoadingState(false);
      this.$refs.pendingTransfers.setLoadingState(false);
      this.$refs.outgoingTransfers.setLoadingState(false);
    });

    // Get all stores
    window.backend.GetStores().then((stores) => {
      this.stores = stores;
    }, (err) => {
      this.$toast.error("Error! " + err);
    });
  },
  methods: {
    formatDate(date) {
      return moment(date.Time).format("Do of MMMM YYYY");
    },
    async storeName(storeID) {
      return await this.stores.filter((store) => {
        return (
          store.code.toLowerCase() === storeID.toLowerCase()
        );
      })[0].name;
    },
    detailsInfo(row) {
      const id = row.id;
      this.$router.push({ name: "edittransfer", params: {id} });
    },
    displayInfo(row) {
      const id = row.id;
      this.$router.push({ name: "transferdetail", params: {id} });
    },
  }
};
</script>