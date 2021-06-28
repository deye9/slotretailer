<template>
  <section>
    <h3>Store Setup</h3>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="name">Name</label>
        <input type="text" class="form-control form-control-sm" placeholder="Store name" v-model="name" />
      </div>
      <div class="form-group col">
        <label for="email">Email</label>
        <input type="email" class="form-control form-control-sm" placeholder="Email Address" v-model="email" />
      </div>
      <div class="form-group col">
        <label for="phone">Phone</label>
        <input type="text" class="form-control form-control-sm" placeholder="Phone Number" v-model="phone" />
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group col">
        <label for="city">City</label>
        <input type="text" class="form-control form-control-sm" placeholder="City" v-model="city" />
      </div>
      <div class="form-group col">
        <label for="address">Address</label>
        <input type="text" class="form-control form-control-sm" placeholder="Address" v-model="address" />
      </div>
      <div class="form-group col">
        <label for="address">SAP Store Identifier</label>
        <input type="text" class="form-control form-control-sm" placeholder="SAP Store Identifier" v-model="sapkey" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label for="sync_interval">Sync Interval in Minutes</label>
        <input min="30" step="1" type="number" class="form-control form-control-sm" placeholder="Sync Interval" v-model="sync_interval" />
      </div>

      <div class="form-group col">
        <label>Log Rotation Details</label>
        <select v-model="logrotation" class="form-control form-control-sm">
          <option value="null">Please select Log Rotation Frequency</option>
          <option value="Daily">Daily</option>
          <option value="Weekly">Weekly</option>
          <option value="Bi-weekly">Bi-weekly</option>
          <option value="Monthly">Monthly</option>
          <option value="Quarterly">Quarterly</option>
          <option value="Yearly">Yearly</option>
        </select>
      </div>

      <div class="form-group col">
        <label for="vatCheck">
          Allow VAT
        </label>
        <div class="input-group">
          <div id="radioBtn" class="btn-group">
            <a class="btn btn-primary btn-sm notActive" data-toggle="canVat" data-title="Y" @click="allowVAT('Y')">YES</a>
            <a class="btn btn-primary btn-sm notActive" data-toggle="canVat" data-title="N" @click="allowVAT('N')">NO</a>
          </div>
            <input type="hidden" name="canVat" id="canVat">
        </div>
      </div>

    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Set Store Price List</label>
        <select class="form-control form-control-sm" v-model="productpricelist" >
          <option :key="price.id" :value="price.name" v-for="price in pricelistArray">{{ price.code }}</option>
        </select>
      </div>

      <div class="form-group col">
        <label>Set Store Cash Account</label>
        <select class="form-control form-control-sm" v-model="storecashaccount" >
          <option :key="acct.id" :value="acct.name" v-for="acct in cashaccounts">{{ acct.code }}</option>
        </select>
      </div>
    </div>

    <div class="card mb-2">
      <div class="card-header">API Endpoints:</div>

      <div class="card-body">
      
        <div class="form-row">
          <div class="form-group col">
            <label for="orders">Orders API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Orders API" v-model="orders" />
          </div>

          <div class="form-group col">
            <label for="products">Products API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Products API" v-model="products" />
          </div>

          <div class="form-group col">
            <label for="customers">Customers API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Customers API" v-model="customers" />
          </div>

          <div class="form-group col">
            <label>Price List API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Price List API" v-model="pricelist" />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group col">
            <label>Cash Account API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Cash Account API" v-model="cashaccount" />
          </div>

          <div class="form-group col">
            <label>Credit Cards API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Banks API" v-model="creditcard" />
          </div>

          <div class="form-group col">
            <label>Bank Transfer API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Bank Transfer API" v-model="banktransfer" />
          </div>

          <div class="form-group col">
            <label>Cheques API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Cheques API" v-model="cheques" />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group col">
            <label for="banks">Transfers API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Transfers API" v-model="transfers" />
          </div>

          <div class="form-group col">
            <label for="banks">Warehouses API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Warehouses API" v-model="warehouses" />
          </div>
        </div>
      </div>
    </div>

    <button type="submit" class="btn btn-primary btn-sm float-right" @click="storeDetails" v-if="userPermission('store', 'cancreate')">
      {{ buttontext }}
    </button>    
  </section>
</template>

<style scoped>
  #radioBtn .notActive{
      color: #3276b1;
      background-color: #fff;
  }
</style>

<script>
import $ from "jquery";
import moment from "moment";

export default {
  data() {
    return {
      id: 0,
      store: {},
      vat: false,
      city: null,
      name: null,
      phone: null,
      email: null,
      orders: null,
      sapkey: null,
      cheques: null,
      address: null,
      products: null,
      customers: null,
      transfers: null,
      pricelist: null,
      cashaccount: null,
      warehouses: null,
      creditcard: null,
      created_by: null,
      sync_interval: 5,
      logrotation: null,
      banktransfer: null,
      cashaccounts: [],
      pricelistArray: [],
      productpricelist: null,
      storecashaccount: null,
      buttontext: "Register Store",
    };
  },
  mounted() {
    // Get all Product Price List
    window.backend.GetPriceList().then((pricelistArray) => {
      this.pricelistArray = pricelistArray;
    }, (err) => {
      this.$toast.error("Error! " + err);
    });

    // Get all Cash Accounts
    window.backend.GetCashAccounts().then((cashaccounts) => {
      this.cashaccounts = cashaccounts;
    }, (err) => {
      this.$toast.error("Error! " + err);
    });

    // Get store details
    window.backend.GetStore().then((store) => {
      this.id = store.id;
      this.city = store.city;
      this.name = store.name;
      this.phone = store.phone;
      this.email = store.email;
      this.sapkey = store.sapkey;
      this.orders = store.orders;
      this.address = store.address;
      this.cheques = store.cheques;
      this.products = store.products;
      this.pricelist = store.pricelist;
      this.customers = store.customers;
      this.transfers = store.transfers;
      this.warehouses = store.warehouses;
      this.creditcard = store.creditcard;
      this.logrotation = store.logrotation;
      this.cashaccount = store.cashaccount;
      this.banktransfer = store.banktransfer;
      this.sync_interval = store.sync_interval;
      this.created_by = this.$store.state.user.id;
      this.productpricelist = store.productpricelist;
      this.storecashaccount = store.storecashaccount;

      if (this.id === 0) {
        this.buttontext = "Register Store";
      } else {
        this.buttontext = "Update Store";
      }

      let vatCheck = store.vat === true ? 'Y' : 'N';
      this.allowVAT(vatCheck);
    }, (err) => {
      this.$toast.error("Error! " + err);
    });
  },
  methods: {
    allowVAT(sel) {
      var tog = 'canVat';
      $('#'+tog).prop('value', sel);
      $('a[data-toggle="'+tog+'"]').not('[data-title="'+sel+'"]').removeClass('active').addClass('notActive');
      $('a[data-toggle="'+tog+'"][data-title="'+sel+'"]').removeClass('notActive').addClass('active');

      if (sel === 'Y') {
        this.vat = true;
      } else { 
        this.vat = false;
      }
    },
    storeDetails() {
      if (this.logrotation === null) {
        this.$toast.error("Error! Invalid value set for log Rotation.");
        return;
      }

      this.store = {
        id: this.id,
        vat: this.vat,
        city: this.city,
        name: this.name,
        phone: this.phone,
        email: this.email,
        sapkey: this.sapkey,
        orders: this.orders,
        address: this.address,
        cheques: this.cheques,
        products: this.products,
        pricelist: this.pricelist,
        customers: this.customers,
        transfers: this.transfers,
        warehouses: this.warehouses,
        creditcard: this.creditcard,
        updated_at: moment().format(),
        logrotation: this.logrotation,
        cashaccount: this.cashaccount,
        banktransfer: this.banktransfer,
        sync_interval: this.sync_interval,
        created_by: this.$store.state.user.id,
        productpricelist: this.productpricelist,
        storecashaccount: this.storecashaccount,
      };

      // Validate the payload.
      for (var attribute in this.store) {
        if (this.store[attribute] === "" || this.store[attribute] === null) {
          this.$toast.error("Error! " + attribute + " cannot be " + this.store[attribute]);
          return;
        }
      }

      // Keep the store details in vuex
      this.$store.state.userStore = this.store;

      window.backend.SaveStore(this.store).then(() => {
        this.$toast.success("Success! Store details successfully modified.");
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>