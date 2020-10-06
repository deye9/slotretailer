<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Editing customer: {{ this.cardname }}</h3>
      </div>
      <div class="col-4">
        <router-link to="/customers/" class="btn btn-info float-right"
          >Back</router-link
        >
      </div>
    </div>
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

    <div class="card">
      <div class="card-header">Contact Information</div>
      <div class="card-body">
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
              @blur="handleBlur"
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
              @blur="handleBlur"
              required
            />
          </div>
        </div>
      </div>
    </div>

    <br />

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
      id="update"
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
      email: null,
      cardcode: 0,
      customer: {},
      address: null,
      phone: "+234",
      phone1: "+234",
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
    handleBlur() {
      if (this.phone.charAt(0) === "0") {
        this.phone = this.phone.replace("0", "+234");
      }

      if (this.phone1.charAt(0) === "0") {
        this.phone1 = this.phone1.replace("0", "+234");
      }

      document.getElementById("update").disabled = false;
      let details = `phone = '${this.phone}' or phone = '${this.phone1}' or phone1 = '${this.phone}' or phone1 = '${this.phone1}'`;

      window.backend.GetCustomerbyPhone(details).then((customer) => {
        this.id = customer.id;
        this.city = customer.city;
        this.phone = customer.phone;
        this.email = customer.email;
        this.phone1 = customer.phone1;
        this.address = customer.address;
        this.cardcode = customer.cardcode;
        this.cardname = customer.cardname;
        this.created_by = this.$store.state.user.id;
        document.getElementById("update").disabled = true;
      });
    },
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

      window.backend.UpdateCustomer(this.customer).then(() => {
          this.$toast.success(
            `Success! Customer ${this.cardname} has been successful updated.`
          );
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