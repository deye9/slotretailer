<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>New Inventory Transfer</h3>
      </div>
      <div class="col-4">
        <router-link :to="{ name: 'transferlist' }" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="fromWHS">From Store</label>
        <input id="createdBy" type="text" class="form-control form-control-sm" placeholder="Requested By" disabled v-model="dispatching" />
      </div>
      <div class="form-group col">
        <label for="toWHS">To Store</label>
        <select @input="fetchInventory" v-model="receiving" label="name" id="toWHS" ref="toWHS" class="form-control form-control-sm" placeholder="Kindly select receiving warehouse">
          <option value="" selected>Select Receiving Warehouse</option>
          <option :key="store.name" :value="store.code" v-for="store in stores">{{ store.name }}</option>
        </select>
      </div>
      <div class="form-group col">
        <label for="createdBy">Requested by</label>
        <input id="createdBy" type="text" class="form-control form-control-sm" placeholder="Requested By" disabled :value="this.$store.state.user.email" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label for="comment">Remarks</label>
        <input id="comment" type="text" class="form-control form-control-sm" placeholder="Comment" v-model="comment" required />
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
            <th scope="col">Item</th>
            <th scope="col">Item No.</th>
            <th scope="col">Item Description</th>
            <th scope="col">Availale</th>
            <th scope="col">Quantity</th>
            <th scope="col"></th>
          </tr>
        </thead>
        <tbody id="LineItems">
          <tr v-for="(item, i) in items" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>
              <v-select label="itemName" @input="(val) => itemSelected(val, i)" v-model="item.itemname" :options="inventory" :clearable="false" placeholder="Kindly select Product"></v-select>
            </td>
            <td>{{ item.itemcode }}</td>
            <td>{{item.itemname}}</td>
            <td>{{ item.onhand }}</td>
            <td><input type="number" min="1" step="1" class="form-control form-control-sm" :value="item.quantity" @blur="setQuantity(i)" /></td>
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

    <button type="submit" id="register" class="btn btn-primary btn-sm float-right" :disabled="isDisabled" @click="inventoryRequest">
      Create Request
    </button>

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

  #loader {
    display: none;
    border-top: 16px solid blue;
    border-right: 16px solid green;
    border-bottom: 16px solid red;
    border-left: 16px solid pink;
    border-radius: 50%;
    width: 120px;
    height: 120px;
    animation: spin 2s linear infinite;
    /* Center */
    display: flex;
    justify-content: center;
    align-items: center;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
</style>

<script>
export default {
  data() {
    return {
      items: [],
      picked: "",
      stores: [],
      comment: "",
      options: [],
      inventory: [],
      isAdmin: false,
      receiving: null,
      isDisabled: true,
      dispatching: null,
      created_by: this.$store.state.user.id,
      localStore: this.$store.state.userStore,
    };
  },
  async mounted() {
    // Determine the state of the Discount element
    if (this.$store.state.isAdmin) {
      this.isAdmin = true;
    }

    // Get all stores
    window.backend.GetStores().then((stores) => {
      // Remove the localstore
      this.stores = stores.filter((el) => {
        return el.code.toLowerCase() !== this.localStore.sapkey.toLowerCase();
      });
    }, (err) => {
        this.$toast.error("Error! " + err);
      }
    );

    document.getElementById("loader").style.display = "none";
    this.dispatching = this.localStore.sapkey;
  },
  methods: {
    fetchInventory(event) {
      this.isDisabled = true;

      let products = this.$store.state.userStore.products,
        selectedValue = event.target.selectedOptions[0].value;

      if (selectedValue === this.localStore.sapkey) {
        this.$toast.error("Error! You cannot select your store. Kindly select another store.");
        return;
      }

      document.getElementById("loader").style.display = "block";
      var requestOptions = {
        method: 'GET',
        redirect: 'follow'
      };
      
      // Get inventory belonging to this store
      fetch(`${products}/Quantity?PageSize=10000&StoreId=${selectedValue}`, requestOptions)
        .then(response => response.json())
        .then(result => {
          this.inventory = result;
          this.items.push({
            onhand: 1,
            quantity: 1,
            itemcode: '',
            itemname: '',
            transferid: null,
          });
          this.isDisabled = false;
          document.getElementById("loader").style.display = "none";
        })
        .catch(error => {
          this.isDisabled = true;
          this.$toast.error("Error! " + error);
          document.getElementById("loader").style.display = "none";
        });
    },
    async itemSelected(val, rowIndex) {
      this.item = this.inventory.filter((inventory) => {
        return (
          inventory.itemCode.toLowerCase() === val.itemCode.toLowerCase()
        );
      })[0];

      if (this.item === undefined) {
        return;
      }

      // Set the table data
      this.items[rowIndex].quantity = 1;
      this.items[rowIndex].transferid = null;
      this.items[rowIndex].onhand = this.item.onHand;
      this.items[rowIndex].itemcode = this.item.itemCode;
      this.items[rowIndex].itemname = this.item.itemName;
      
      // Add a new row to the table
      this.items.push({
        onhand: 1,
        quantity: 1,
        itemcode: '',
        itemname: '',
        transferid: null,
      });
    },
    async setQuantity(rowIndex) {      
      if (parseInt(event.target.value) > parseInt(this.items[rowIndex].onhand)) {
        this.isDisabled = true;
        this.$toast.info("Error! Requested Quantity is more than available Inventory at store.");
        return
      }

      this.isDisabled = false;
      this.items[rowIndex].quantity = event.target.value;
    },   
    async deleteRow(index) {
      this.$delete(this.items, index);
      
      if (this.items.length === 0 || this.items.length === index) {
        // Add a new row to the table
        this.items.push({
          onhand: 1,
          quantity: 1,
          itemcode: '',
          itemname: '',
          transferid: null,
        });
      }
    },     
    async inventoryRequest() {
      if (this.items.length === 0) {
        this.$toast.error("Error! You are not permitted to create an empty Inventory Transfer Request.");
        return;
      }

      var r = confirm("Are you sure you want to create this Inventory Transfer Request!");
      if (r == false) {
        return;
      }

      // Remove the last row as this is not needed.
      this.$delete(this.items, this.items.length - 1);

      let transfer = {};
      transfer.id = 0;
      transfer.docnum = 0;
      transfer.docentry = 0;
      transfer.requestid = 0;
      transfer.synced = false;
      transfer.canceled = false;
      transfer.status = "Pending";
      transfer.items = this.items;
      transfer.comment = this.comment;
      transfer.towhs = this.receiving;
      transfer.fromwhs = this.dispatching;
      transfer.created_by = this.created_by;

      window.backend.NewTransfer(transfer).then(() => {
        this.$toast.success("Success! Inventory Transfer request has been successfully registered.");
        this.$router.push({ name: "transferlist" });
      }, (err) => {
          this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>