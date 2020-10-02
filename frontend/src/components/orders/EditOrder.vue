<template>
  <section>
    <h3>Editing Customer [{{ firstname }} {{ lastname }} : {{ email }}]</h3>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="cardcode">Customer Code</label>
        <input
          type="text"
          class="form-control"
          placeholder="Customer Code"
          v-model="cardcode"
          disabled
        />
      </div>
      <div class="form-group col">
        <label for="cardname">Customer Name</label>
        <input
          type="text"
          class="form-control"
          placeholder="Customer name"
          v-model="cardname"
          required
        />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group col">
        <label for="city">City</label>
        <input
          type="text"
          class="form-control"
          placeholder="City"
          v-model="city"
          required
        />
      </div>
      <div class="form-group col">
        <label for="address">Contact Address</label>
        <input
          type="text"
          class="form-control"
          placeholder="Contact Address"
          v-model="address"
          required
        />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group col">
        <label for="phone1"
          >Phone Number 1 <span style="color: red">*</span></label
        >
        <input
          type="text"
          class="form-control"
          placeholder="Phone Number"
          v-model="phone"
          required
        />
      </div>
      <div class="form-group col">
        <label for="phone2">Phone Number 2</label>
        <input
          type="text"
          class="form-control"
          placeholder="Phone Number"
          v-model="phone1"
          required
        />
      </div>
    </div>
    <div class="form-group">
      <label for="email">Email Address</label>
      <input
        type="email"
        class="form-control"
        placeholder="Email Address"
        v-model="email"
        required
      />
    </div>
    <button
      type="submit"
      class="btn btn-primary float-right"
      @click="Modification"
    >
      Update Customer
    </button>
  </section>
</template>

<script>
import moment from "moment";

export default {
  data() {
    return {
      id: null,
      city: null,
      phone: null,
      phone1: null,
      email: null,
      cardcode: 0,
      customer: {},
      address: null,
      cardname: null,
      created_by: null,
    };
  },
  mounted() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.GetCustomer(parseInt(this.id)).then(
      (customer) => {
        this.city = customer.city;
        this.phone = customer.phone;
        this.email = customer.email;
        this.phone1 = customer.phone1;
        this.address = customer.address;
        this.cardcode = customer.cardcode;
        this.cardname = customer.cardname;
        this.created_by = this.$store.state.user.id;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );
  },
  methods: {
    Modification() {
      this.customer = {
        id: this.id,
        city: this.city,
        phone: this.phone,
        email: this.email,
        phone1: this.phone1,
        address: this.address,
        cardcode: this.cardcode,
        cardname: this.cardname,
        updated_at: moment().format(),
        created_by: this.$store.state.user.id,
      };

      // Validate the payload.
      for (var attribute in this.customer) {
        if (
          this.customer[attribute] === "" ||
          this.customer[attribute] === null
        ) {
          this.$toast.error(
            "Error! " + attribute + " cannot be " + this.customer[attribute]
          );
          return;
        }
      }

      window.backend.UpdateCustomer(this.customer).then(
        () => {
          this.$router.push("/customers/");
        },
        (err) => {
          this.$toast.error("Error! " + err);
        }
      );
    },
  },
};
</script>