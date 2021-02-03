<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Access Control List</h3>
      </div>
      <div class="col-4">
        <router-link :to="{ name: 'dashboard' }" class="btn btn-info btn-sm float-right">Home</router-link>
      </div>
    </div>
    <hr />
  </section>
</template>

<script>
export default {
  data() {
    return {
      user: {},
      email: null,
      isValid: true,
      isadmin: false,
      password: null,
      lastname: null,
      firstname: null,
      created_by: null,
      confirmpassword: null,
    };
  },
  async mounted() {
    // Determine the state of the Discount element
    if (this.$store.state.isAdmin) {
      this.isAdmin = true;
    }

    // Check if the store is allowed to calculate VAT
    if (this.$store.state.userStore.vat) {
      this.canVat = true;
    }

    // Get all stores
    window.backend.GetStores().then(
      (stores) => {
        this.stores = stores;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );

    // Get all Payment Details
    window.backend.GetPaymentDetails().then(
      (PaymentDetails) => {
        this.paymentDetails = {};
        this.paymentDetails["pos"] = PaymentDetails[0]["pos"];
        this.paymentDetails["cheque"] = PaymentDetails[0]["cheque"];
        this.paymentDetails["banktransfer"] = PaymentDetails[0]["banktransfer"];
        this.addRow();
        this.addItemRow();
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );
  },
  methods: {
    RegisterUser() {
      if (this.password !== this.confirmpassword) {
        this.$toast.error("Error! Passwords do not match.");
        return;
      }

      this.user = {
        email: this.email,
        password: this.password,
        lastname: this.lastname,
        firstname: this.firstname,
        created_by: this.$store.state.user.id,
      };

      // Validate the payload.
      for (var attribute in this.user) {
        if (this.user[attribute] === "" || this.user[attribute] === null) {
          this.isValid = false;
          this.$toast.error(
            "Error! " + attribute + " cannot be " + this.user[attribute]
          );
          return;
        }
      }

      window.backend.NewUser(this.user).then(
        () => {
          this.$router.push("{name: 'userlist'");
        },
        (err) => {
          this.$toast.error("Error! " + err);
        }
      );
    },
  },
};
</script>