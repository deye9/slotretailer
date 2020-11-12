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

    <div class="card">
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
        </div>

        <div class="form-row">
          <div class="form-group col">
            <label for="banks">Banks API</label>
            <input type="text" class="form-control form-control-sm" placeholder="Banks API" v-model="banks" />
          </div>
          <div class="form-group col">
            <label for="sync_interval">Sync Interval in Minutes</label>
            <input min="30" step="1" type="number" class="form-control" placeholder="Sync Interval" v-model="sync_interval" />
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
        </div>
      </div>
    </div>

    <br />
    <button type="submit" class="btn btn-primary btn-sm float-right" @click="StoreDetails">
      {{ buttontext }}
    </button>    
  </section>
</template>

<script>
import moment from "moment";

export default {
  data() {
    return {
      id: 0,
      store: {},
      city: null,
      name: null,
      banks: null,
      phone: null,
      email: null,
      orders: null,
      sapkey: null,
      address: null,
      products: null,
      customers: null,
      created_by: null,
      sync_interval: 5,
      logrotation: null,
      buttontext: "Register Store",
    };
  },
  mounted() {
    window.backend.GetStore().then((store) => {
      this.id = store.id;
      this.city = store.city;
      this.name = store.name;
      this.banks = store.banks;
      this.phone = store.phone;
      this.email = store.email;
      this.sapkey = store.sapkey;
      this.orders = store.orders;
      this.address = store.address;
      this.products = store.products;
      this.customers = store.customers;
      this.logrotation = store.logrotation;
      this.sync_interval = store.sync_interval;
      this.created_by = this.$store.state.user.id;

      if (this.id === 0) {
        this.buttontext = "Register Store";
      } else {
        this.buttontext = "Update Store";
      }
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });
  },
  methods: {
    StoreDetails() {
      if (this.logrotation === null) {
        this.$toast.error("Error! Invalid value set for log Rotation.");
        return;
      }

      this.store = {
        id: this.id,
        city: this.city,
        name: this.name,
        phone: this.phone,
        email: this.email,
        sapkey: this.sapkey,
        orders: this.orders,
        address: this.address,
        banks: this.banks,
        products: this.products,
        customers: this.customers,
        updated_at: moment().format(),
        logrotation: this.logrotation,
        sync_interval: this.sync_interval,
        created_by: this.$store.state.user.id,
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
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>