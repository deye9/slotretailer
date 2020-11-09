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
        <v-select label="name" :options="stores" @input="fetchInventory" placeholder="Kindly select dispatching warehouse"></v-select>
      </div>
      <div class="form-group col">
        <label for="toWHS">To Store</label>
        <v-select label="name" :options="stores" @input="fetchInventory" placeholder="Kindly select receiving warehouse"></v-select>
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
            <th scope="col">Quantity</th>
            <th scope="col"></th>
          </tr>
        </thead>
        <tbody id="LineItems">
          <!-- <tr v-for="(item, i) in items" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>
              {{ item.itemcode }}
            </td>
            <td>
              <v-select label="itemname" @input="(val) => itemSelected(val, i)" @search="(search, loading) => fetchProduct(search, loading, i)" v-model="item.itemname" :options="inventory" :clearable="false" placeholder="Kindly select Product"></v-select>
            </td>
            <td>
              <input type="number" min="1" step="1" class="form-control form-control-sm" :value="item.quantity" @blur="setQuantity(i)" />
            </td>
            <td>
              <button :id="'del' + i" class="btn btn-danger btn-sm mr-2 float-right" @click="deleteRow('row' + i)">Remove Line</button>
              <button :id="'add' + i" class="btn btn-primary btn-sm mr-2 float-right" @click="addRow('add' + i)">New Line</button>
            </td>
          </tr> -->
        </tbody>
      </table>
    </div>

    <button type="submit" id="register" class="btn btn-primary btn-sm float-right" @click="inventoryRequest">
      Create Request
    </button>

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
export default {
  data() {
    return {
      towhs: '',
      stores: [],
      options: [],
      fromwhs: '',
      comment: '',
      inventory: [],
      synced: false,
      isAdmin: false,
      canceled: false,
      created_by: this.$store.state.user.id,
    };
  },
  async created() {
    // Determine the state of the Discount element
    if (this.$store.state.isAdmin)
    {
      this.isAdmin = true;
    }

    // Get all stores
    window.backend.GetStores().then((stores) => {
      this.stores = stores;
      this.addRow();
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });
  },
  methods: {
    addRow(rowID) {
      console.log(rowID);
    },
    deleteRow(rowID) {
      console.log(rowID);
    },
    inventoryRequest() {

    },
    fetchInventory(value) {
      // Get inventory belonging to this store
      window.backend.GetStoreProducts(value.id).then((inventory) => {
        console.log(inventory);
        this.inventory = inventory;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>