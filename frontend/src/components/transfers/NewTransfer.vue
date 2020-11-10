<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>New Inventory Transfer</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'transferlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="fromWHS">From Store</label>
        <select @input="fetchInventory" label="name" ref="fromWHS" class="form-control form-control-sm" placeholder="Kindly select dispatching warehouse">
          <option value="null">Kindly select dispatching warehouse</option>
          <option :key="store.name" :value="store.id" v-for="store in stores">{{ store.name }}</option>
        </select>
      </div>
      <div class="form-group col">
        <label for="toWHS">To Store</label>
        <select @input="fetchInventory" label="name" ref="toWHS" class="form-control form-control-sm" placeholder="Kindly select receiving warehouse">
          <option value="null">Kindly select receiving warehouse</option>
          <option :key="store.name" :value="store.id" v-for="store in stores">{{ store.name }}</option>
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
          <h5 style="display:inline-flex;" class="float-left">Transfer Items</h5>
        </caption>
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item No.</th>
            <th scope="col">Item Description</th>
            <th scope="col">Quantity Available</th>
            <th scope="col">Quantity</th>
            <th scope="col"></th>
          </tr>
        </thead>
        <tbody id="LineItems">
          <tr v-for="(item, i) in items" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>
              {{ item.itemcode }}
            </td>
            <td>
              <select class="form-control form-control-sm">
                <option value="null" selected>Select Desired Product</option>
                <option :key="item.id" :value="item.itemcode" v-for="item in inventory">{{ item.itemname }}</option>
              </select>
            </td>
            <td>
              {{ item.onHand }}
            </td>
            <td>
              <input type="number" min="1" max="1" step="1" class="form-control form-control-sm" :value="1" />
            </td>
            <td>
              <button :id="'del' + i" class="btn btn-danger btn-sm mr-2 float-right" @click="deleteRow('row' + i)">Remove Line</button>
              <button :id="'add' + i" class="btn btn-primary btn-sm mr-2 float-right" @click="addRow('add' + i)">New Line</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <button type="submit" id="register" class="btn btn-primary btn-sm float-right" @click="inventoryRequest">
      Create Request
    </button>

    <!-- Modal -->
    <div class="modal fade" id="roleModal" data-backdrop="static" tabindex="-1" aria-labelledby="roleModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header text-center bg-dark text-white">
            <h5 class="modal-title" id="roleModalLabel">Set Inventory Transfer Role</h5>
            <button type="button" class="close text-white" @click="dismiss" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <div class="container-fluid">
              <div class="form-check mb-3">
                <input class="form-check-input" type="radio" name="RoleRadios" id="reqRadio" value="requester" v-model="picked" />
                <label class="form-check-label" for="reqRadio">
                  Requesting Store
                </label>
              </div>
              <div class="form-check">
                <input class="form-check-input" type="radio" name="RoleRadios" id="recRadio" value="receiver" v-model="picked" />
                <label class="form-check-label" for="recRadio">
                  Receiving Store
                </label>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary btn-sm mr-2" @click="setRole">Set Role</button>
            <button type="button" class="btn btn-secondary btn-sm" @click="dismiss">Close</button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
  caption {
    padding-top: .75rem;
    padding-bottom: .75rem;
    color: Black;
    text-align: left;
    caption-side: top;
  }
</style>

<script>
import $ from "jquery";

export default {
  data() {
    return {
      items: [],
      towhs: '',
      picked: '',
      stores: [],
      options: [],
      fromwhs: '',
      comment: '',
      inventory: [],
      synced: false,
      isAdmin: false,
      canceled: false,
      created_by: this.$store.state.user.id,
      localStore: this.$store.state.userStore,
    };
  },
  async mounted() {
    // Determine the state of the Discount element
    if (this.$store.state.isAdmin)
    {
      this.isAdmin = true;
    }

    // Get all stores
    window.backend.GetStores().then((stores) => {
      this.stores = stores;
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });
  
    // Show the modal
    $('#roleModal').modal('show');
  },
  methods: {
    setRole() {
      let localStore = this.stores.filter((store) => {
          return (
            store.name.toLowerCase() === this.localStore.sapkey.toLowerCase()
          );
        })[0];

      if (this.picked === "requester") {
        this.$refs.fromWHS.disabled = true;
        this.$refs.fromWHS.value = localStore.id;
      } else if (this.picked === "receiver") {
        this.$refs.toWHS.disabled = true;
        this.$refs.toWHS.value = localStore.id;
      }

      // Hide the modal
      $('#roleModal').modal('hide');
    },
    dismiss() {
      // Hide the modal
      $('#roleModal').modal('hide');

      // Hide the information and redirect to the Transfer Request list.
      this.$toast.info("Info! You can't continue without setting your role.");
      this.$router.push({name: 'transferlist'});
    },
    addRow(butID) {
      if (butID !== undefined) {
        // Remove Button from table cell
        var myobj = document.getElementById(butID);
        myobj.remove();
      }

      this.items.push({
        onHand: 1,
        quantity: 1,
        itemcode: '',
        itemname: '',
      });

      const el = document.getElementById("register");
      if (el) {
        el.scrollIntoView();
      }
    },
    deleteRow(rowID) {
      console.log(rowID);
    },
    inventoryRequest() {

    },
    fetchInventory(event) {
      let selectedtext = event.path[0].selectedOptions[0].text,
        selectedValue = parseInt(event.path[0].selectedOptions[0].value);

      if (selectedtext === this.localStore.sapkey) {
        this.$toast.error("Error! You cannot select your store. Kindly select another store.");
        return;
      }

      // Get inventory belonging to this store
      window.backend.GetStoreProducts(selectedValue).then((inventory) => {        
        this.items = [];
        this.inventory = inventory;
        this.addRow();
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>