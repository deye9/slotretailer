<template>
  <section>
    <h3>Editing Store [{{name}} : {{email}}]</h3>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="name">Name</label>
        <input type="text" class="form-control" placeholder="Store name" v-model="name" />
      </div>
      <div class="form-group col">
        <label for="email">Email</label>
        <input type="email" class="form-control" placeholder="Email Address" v-model="email" />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group col">
        <label for="phone">Phone</label>
        <input type="text" class="form-control" placeholder="Phone Number" v-model="phone" />
      </div>
      <div class="form-group col">
        <label for="city">City</label>
        <input type="text" class="form-control" placeholder="City" v-model="city" />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group col">
        <label for="address">Address</label>
        <input type="text" class="form-control" placeholder="Address" v-model="address" />
      </div>
      <div class="form-group col">
        <label for="address">SAP Store Identifier</label>
        <input type="text" class="form-control" placeholder="SAP Store Identifier" v-model="sapkey" />
      </div>
    </div>
    <h3>API Endpoints:</h3>
    <hr />
    <div class="form-row">
      <div class="form-group col">
        <label for="orders">Orders API</label>
        <input type="text" class="form-control" placeholder="Orders API" v-model="orders" />
      </div>
      <div class="form-group col">
        <label for="products">Products API</label>
        <input type="text" class="form-control" placeholder="Products API" v-model="products" />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group col">
        <label for="customers">Customers API</label>
        <input type="text" class="form-control" placeholder="Customers API" v-model="customers" />
      </div>
      <div class="form-group col">
        <label for="sync_interval">Sync Interval in Minutes</label>
        <input type="text" class="form-control" placeholder="Sync Interval" v-model="sync_interval" />
      </div>
    </div>
    <button type="submit" class="btn btn-primary float-right" @click="StoreDetails">{{buttontext}}</button>
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
      phone: null,
      email: null,
      orders: null,
      sapkey: null,
      address: null,
      products: null,
      customers: null,
      created_by: null,
      sync_interval: 5,

      buttontext: "Register Store",
    };
  },
  mounted() {
    window.backend.GetStore().then((store) => {
        this.id = store.id;
        this.city = store.city;
        this.name = store.name;
        this.phone = store.phone;
        this.email = store.email;
        this.sapkey = store.sapkey;
        this.orders = store.orders;
        this.address = store.address;
        this.products = store.products;
        this.customers = store.customers;
        this.sync_interval = store.sync_interval;
        this.created_by = this.$store.state.user.id;

        if (this.id === 0) {
          this.buttontext = "Register Store";
        } else {
          this.buttontext = "Update Store";
        }
      },
      (err) => {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! " + err;
      }
    );
  },
  methods: {
    StoreDetails() {
      this.store = {
        id: this.id,
        city: this.city,
        name: this.name,
        phone: this.phone,
        email: this.email,
        sapkey: this.sapkey,
        orders: this.orders,
        address: this.address,
        products: this.products,
        customers: this.customers,
        updated_at: moment().format(),
        sync_interval: this.sync_interval,
        created_by: this.$store.state.user.id,
      };

      // Validate the payload.
      for (var attribute in this.store) {
        if (this.store[attribute] === "" || this.store[attribute] === null) {
          this.$store.state.notify.category = "error";
          this.$store.state.notify.message =
            "Error! " + attribute + " cannot be " + this.store[attribute];
          return;
        }
      }

      window.backend.SaveStore(this.store).then(() => {
          this.$store.state.notify.category = "success";
          this.$store.state.notify.message =
            "Success! Store details successfully modified.";
        },
        (err) => {
          this.$store.state.notify.category = "error";
          this.$store.state.notify.message = "Error! " + err;
        }
      );
    },
  },
};
</script>