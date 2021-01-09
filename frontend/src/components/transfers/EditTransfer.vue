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
          <h5 style="display:inline-flex;" class="float-left">Ordered Items</h5>
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
            <td><input type="number" min="1" step="1" class="form-control form-control-sm" :value="item.quantity" @blur="setQuantity(i)" /></td>
            <td>{{ item.serialnumbers }}</td>
            <td><button :id="'del' + i" class="btn btn-danger btn-sm mr-2 float-right" @click="deleteRow(i)">Remove Line</button></td>
          </tr>
        </tbody>
        <tfoot id="ItemsFooter">
          <tr>
            <td colspan="7" class="text-right font-weight-bold">&nbsp;</td>
          </tr>
        </tfoot>
      </table>
      <div id="loader"></div>
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
import moment from "moment";

export default {
  data() {
    return {
      stores: [],
      comment: "",
      options: [],
      transfer: [],
      inventory: [],
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
      console.log(transfer);

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
  }
};
</script>