<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Editing Inventory Transfer [Trans-{{this.transfer.id}}]</h3>
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
        <input type="text" class="form-control form-control-sm" :value="transfer.comment" title="Request Comment" />
      </div>
    </div>

    <div class="table-responsive">
      <table class="table table-fixed table-bordered table-hover table-sm">
        <caption>
          <h5 style="display:inline-flex;" class="float-left">Requested Items</h5>
        </caption>
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item Code</th>
            <th scope="col">Item Description</th>
            <th scope="col">Current Inventory</th>
            <th scope="col">Quantity Requested</th>
            <th scope="col">Serial Number</th>
            <th scope="col"></th>
          </tr>
        </thead>
        <tbody id="LineItems">
          <tr v-for="(item, i) in transfer.items" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>{{ item.itemcode }}</td>
            <td>{{item.itemname}}</td>
            <td>{{ item.onhand }}</td>
            <td><input type="number" min="1" step="1" class="form-control form-control-sm" :value="item.quantity" @blur="setQuantity(i)" style="width: 60px;" /></td>
            <td>{{ item.serial }}</td>
            <td>
              <button :id="'del' + i" class="btn btn-danger btn-sm mr-2 float-right" @click="deleteRow(i)">Remove Line</button>
              <button class="btn btn-info btn-sm mr-2 float-right" v-if="item.serialnumber !== '[]'" @click="serialsDialog(item.serialnumber, i)">SetSerials</button>
              </td>
          </tr>
        </tbody>
      </table>
    </div>

    <button type="button" class="btn btn-danger btn-sm float-right ml-2" @click="RejectTransfer">
      Reject
    </button>
    
    <button type="button" class="btn btn-primary btn-sm float-right" @click="AcceptTransfer">
      Approve
    </button>

    <!-- Modals -->
    <div class="modal fade" id="serialsModal" data-backdrop="static" tabindex="-1" aria-labelledby="serialsModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header text-center bg-dark text-white">
            <h5 class="modal-title" id="serialsModalLabel">Select Serial Number.</h5>
            <button type="button" class="close text-white" @click="dismiss" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <div class="row">
              <div v-for="(val, key) in serialnumbers" :key="key" class="col-3">
                <b>{{key + 1}}</b>.&nbsp;&nbsp;
                <label class="form-check-label" :for="'btncheck_' + key">
                <input class="form-check-input" type="checkbox" :id="'btncheck_' + key" :name="'btncheck_' + key" v-model="selectedSerials" :value="val" autocomplete="off">
                {{val}}
              </label>
              </div>
            </div>
            <hr />
            <span>Selected Serial Number: <strong>{{selectedSerials}}</strong></span>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary btn-sm mr-2" @click="setSerial">Set Serial Number</button>
            <button type="button" class="btn btn-secondary btn-sm" @click="dismiss">Close</button>
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
      currentIndex: 0,
      serialnumbers: [],
      selectedSerials: [],
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
    async dismiss() {
      if (this.transfer.items[this.currentIndex].serial === undefined || this.transfer.items[this.currentIndex].serial === null) {
        this.$toast.error("Error! Product Serial Number is required in order to proceed.");
        return;
      }

      $('#serialsModal').modal('hide');
    },
    async setSerial() {
      if (this.transfer.items[this.currentIndex].serial === null || this.transfer.items[this.currentIndex].serial === " ") {
        this.$toast.error("Error! Product Serial Number is required in order to proceed.");
        return;
      }

      if (this.selectedSerials.length < this.transfer.items[this.currentIndex].quantity) {
        this.transfer.items[this.currentIndex].quantity = this.selectedSerials.length;
        this.$toast.info("Info! Transferred Serial Numbers is less than quantity requested.");
      }

      if (this.selectedSerials.length > this.transfer.items[this.currentIndex].quantity) {
        this.transfer.items[this.currentIndex].quantity = this.selectedSerials.length;
        this.$toast.info("Info! Transferred Serial Numbers is more than quantity requested.");
      }

      // Set the serial number of the item.
      this.transfer.items[this.currentIndex].serial = this.selectedSerials;
      
      // Hide the modal and reset the serial variable to an empty string
      $('#serialsModal').modal('hide');
      
      this.selectedSerials = [];
    },
    async setQuantity(rowIndex) {
      this.transfer.items[rowIndex].quantity = event.target.value;
    },
    async serialsDialog(serialnumbers, index) {
      this.currentIndex = index;

      // Prompt for serial number of item.
      this.serialnumbers = serialnumbers.substring(1, serialnumbers.length-1).split(" ");

      // Display Modal
      $('#serialsModal').modal('show');
    },
    async deleteRow(index) {
      this.$delete(this.transfer.items, index);
    },
    async AcceptTransfer() {
      // Confirm we have at least one item to transfer
      if (this.transfer.items.length === 0) {
        this.$toast.error("Error! You are not permitted to approve an empty Inventory Transfer.");
        return;
      }

      let transfer = {};
      transfer.items = [];
      transfer.id = this.transfer.id;
      transfer.comment = this.transfer.comment;

      // Loop through the items to validate Inventory
      this.transfer.items.forEach(item => {
        let serial = "";
        if (item.serialnumber === "") {
          serial = "null";
        } else {
          serial = item.serialnumber;
        }

        // Adjust the transferred quantity based on inventory available
        if (item.quantity > item.onhand || item.quantity <= 0) {
          item.quantity = item.onhand;
        }

        // Add the items to the items array
        transfer.items.push({
          id: item.id,
          serialnumber: serial,
          itemcode: item.itemcode,
          quantity: item.quantity,
        });
      });
      
      var r = confirm("Are you sure you want to approve this Inventory Transfer Request!");
      if (r == false) {
        return;
      }

      window.backend.UpdateTransfer(transfer).then(() => {
        this.$toast.success("Success! Tranfer has been successfully approved and accepted.");
        this.$router.push({name: 'transferlist'});
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
    },
    async RejectTransfer() {
      window.backend.RejectTransfer(parseInt(this.id)).then((response) => {
        if (response === null || response === undefined) {
          this.$toast.success("Success! Transfer Request has been successfully rejected.");
          this.$router.push({name: 'transferlist'});
        }
      }, (err) => {
          this.$toast.error("Error! " + err);
        }
      );
    },
  }
};
</script>