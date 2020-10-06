<template>
    <section>
        <div class="row flex-xl-nowrap2">
            <div class="col-8">
                <h3>Registered Reports</h3>
            </div>
            <div class="col-4">
                <router-link to="/reports/new" class="btn btn-info float-right">Create Report</router-link>
            </div>
        </div>
        <hr />


    </section>
</template>

<script>
import moment from "moment";

export default {
    data() {
        return {
            report: null,
            created_by: null,
        };
    },
    mounted() {
        var pageURL = location.pathname;
        this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

        window.backend.GetReport(parseInt(this.id)).then((report) => {
            this.report = report;
        },(err) => {
            this.$toast.error("Error! " + err);
        });
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