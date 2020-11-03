<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>New Customer</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'customerlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="cardcode">Customer Code</label>
        <input type="text" class="form-control" placeholder="Customer Code" v-model="cardcode" disabled/>
      </div>
      <div class="form-group col">
        <label for="cardname">Customer Name</label>
        <input type="text" class="form-control" placeholder="Customer name" v-model="cardname" required />
      </div>
    </div>

    <div class="form-group">
      <label for="email">Email Address</label>
      <input type="email" class="form-control" placeholder="Email Address" v-model="email" required />
    </div>

    <div class="card mb-3">
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
            <label for="phone1">
              Phone Number 1
              <span style="color: red">*</span>
            </label>
            <input
              id="phone"
              :key="phone"
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
              id="phone1"
              :key="phone1"
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

    <button type="submit" id="register" class="btn btn-primary btn-sm float-right" @click="Registration">
      Register Customer
    </button>

    <br />
  </section>
</template>

<script>
export default {
  data() {
    return {
      customer: {},
      city: null,
      email: null,
      cardcode: 0,
      address: null,
      phone: "+234",
      phone1: "+234",
      cardname: null,
      created_by: null,
    };
  },
  methods: {
    handleBlur() {
      if (this.phone.charAt(0) === "0") {
        this.phone = this.phone.replace("0", "+234");
      }

      if (this.phone1.charAt(0) === "0") {
        this.phone1 = this.phone1.replace("0", "+234");
      }
      document.getElementById("register").disabled = false;
      let details = `phone = '${this.phone}' or phone = '${this.phone1}' or phone1 = '${this.phone}' or phone1 = '${this.phone1}'`;

      window.backend.GetCustomerbyPhone(details).then((customer) => {
        if (customer.id !== undefined) {
          this.id = customer.id;
          this.city = customer.city;
          this.phone = customer.phone;
          this.email = customer.email;
          this.phone1 = customer.phone1;
          this.address = customer.address;
          this.cardcode = customer.cardcode;
          this.cardname = customer.cardname;
          this.created_by = this.$store.state.user.id;
          document.getElementById("register").disabled = true;
        }
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    },
    Registration() {
      this.customer = {
        synced: false,
        city: this.city,
        phone: this.phone,
        email: this.email,
        phone1: this.phone1,
        address: this.address,
        cardname: this.cardname,
        cardcode: this.cardcode,
        created_by: this.$store.state.user.id,
      };

      // Validate the payload.
      for (var attribute in this.customer) {
        if (this.customer[attribute] === "" || this.customer[attribute] === null) {
          this.isValid = false;
          this.$toast.error( "Error! " + attribute + " cannot be " + this.customer[attribute] );
          return;
        }
      }

      window.backend.NewCustomer(this.customer).then(() => {
        this.$toast.success( `Success! Customer ${this.cardname} has been successful registered.` );
        this.$router.push({name: "customerlist"});
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>