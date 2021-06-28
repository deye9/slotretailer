<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Viewing Inventory Transfer [Trans-{{this.transfer.id}}]</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'transferlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label>From Store</label>
        <input type="text" class="form-control form-control-sm" disabled :value="transfer.fromwhs" title="Requesting Store" />
      </div>
      <div class="form-group col">
        <label>To Store</label>
        <input type="text" class="form-control form-control-sm" disabled :value="transfer.towhs" title="Receiving Store" />
      </div>
      <div class="form-group col">
        <label>Status</label>
        <input type="text" class="form-control form-control-sm" disabled :value="transfer.status" title="Request Status" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>SAP Doc Number</label>
        <input type="text" class="form-control form-control-sm" disabled :value="transfer.docnum" title="Request Comment" />
      </div>
      <div class="form-group col">
        <label>SAP Doc Entry Number</label>
        <input type="text" class="form-control form-control-sm" disabled :value="transfer.docentrty" title="Request Comment" />
      </div>
      <div class="form-group col">
        <label>Requested On</label>
        <input type="text" class="form-control form-control-sm" disabled :value="transfer.created_at" title="Requested On" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Remarks</label>
        <input type="text" class="form-control form-control-sm" disabled :value="transfer.comment" title="Request Comment" />
      </div>
    </div>

    <div class="table-responsive">
      <table class="table table-fixed table-bordered table-hover table-sm">
        <caption>
          <h5 style="display:inline-flex;" class="float-left">Transfered Items</h5>
        </caption>
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item Code</th>
            <th scope="col">Item Description</th>
            <th scope="col">Quantity Transferred</th>
            <th scope="col">Serial Number</th>
          </tr>
        </thead>
        <tbody id="LineItems">
          <tr v-for="(item, i) in transfer.items" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>{{ item.itemcode }}</td>
            <td>{{item.itemname}}</td>
            <td>{{ item.quantity }}</td>
            <td>
              <a class="btn btn-primary btn-sm mr-2" title="Order Details" @click="viewSerials(item.id)">
                <i class="bi bi-pencil-fill">&nbsp;</i>
                View Serial Numbers
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <button type="button" class="btn btn-danger btn-sm float-right ml-2" @click="FinalizeTransfer('Reject')">
      Reject
    </button>
    
    <button type="button" class="btn btn-primary btn-sm float-right" @click="FinalizeTransfer('Accepted')">
      Approve
    </button>

    <!-- Modals -->
    <div class="modal fade" id="serialsModal" data-backdrop="static" tabindex="-1" aria-labelledby="serialsModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header text-center bg-dark text-white">
            <h5 class="modal-title" id="serialsModalLabel">Select Serial Number.</h5>
            <button type="button" class="close text-white" data-dismiss="modal" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <div class="row">
              <div v-for="(val, key) in serialnumbers" :key="key" class="col-3">
                <b>{{key + 1}}</b>.&nbsp;&nbsp;
                <label class="form-check-label">
                {{val.toUpperCase()}}
              </label>
              </div>
            </div>
            <hr />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary btn-sm" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </section>  
</template>

<style scoped>
  caption {
    padding-top: 0.75rem;
    padding-bottom: 0.75rem;
    color: Black;
    text-align: left;
    caption-side: top;
  }
</style>

<script>
import $ from "jquery";
import moment from "moment";

export default {
  data() {
    return {
      id: 0,
      stores: [],
      transfer: [],
      serialnumbers: [],
      created_by: this.$store.state.user.id,
      localStore: this.$store.state.userStore,
    };
  },
  created() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.GetTransfer(parseInt(this.id)).then((transfer) => {
      let toWhs = transfer.towhs,
        fromWhs = transfer.fromwhs;

      transfer.created_at = moment(transfer.created_at.Time).format("Do of MMMM YYYY");
      this.transfer = transfer;

      // Get all stores
      window.backend.GetStores().then((stores) => {
        this.stores = stores;

        stores.filter((store) => {
          if (store.code === fromWhs) {
            this.transfer.fromwhs = store.name;
          }
        });
        
        stores.filter((store) => {
          if (store.code === toWhs) {
            this.transfer.towhs = store.name;
          }
        });
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
    }, (err) => {
        this.$toast.error("Error! " + err);
    });
  },
  methods: {
    async viewSerials(id) {
      window.backend.GetSerials(parseInt(id)).then((serialnumbers) => {
        if (serialnumbers !== null || serialnumbers !== undefined) {
          // Split the serial numbers by space
          this.serialnumbers = serialnumbers.substring(1, serialnumbers.length-1).split(" ");

          // Show the modal and populate it
          $('#serialsModal').modal('show');
        }
      }, (err) => {
          this.$toast.error("Error! " + err);
      });
    },
    async FinalizeTransfer(status) {
      window.backend.FinalizeTransfer(parseInt(this.id), status).then((response) => {
        if (response !== null || response !== undefined) {
          this.$toast.success("Success! Transfer Request has been successfully Finalized.");
          this.$router.push({name: 'transferlist'});
        }
      }, (err) => {
          this.$toast.error("Error! " + err);
      });
    }
  }
};
</script>