<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>New Sales Order</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'orderlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="cardname">Customer</label>
        <v-select label="cardname" @search="fetchCustomer" :options="customers" v-model="customer" :clearable="false" placeholder="Kindly select Customer"></v-select>
      </div>
      <div class="form-group col">
        <label>Customer Code</label>
        <input v-if="this.customer !== null" type="text" class="form-control form-control-sm" placeholder="City" disabled v-model="this.customer.cardcode" />
      </div>
      <div class="form-group col">
        <label>Create Date</label>
        <input type="text" class="form-control form-control-sm" placeholder="City" disabled :value="this.currentDate" />
      </div>
      <div class="form-group col">
        <label>Created By</label>
        <input type="text" class="form-control form-control-sm" placeholder="City" disabled :value="this.$store.state.user.email" />
      </div>
    </div>

    <div class="table-responsive">
      <table class="table table-fixed table-bordered table-hover table-sm">
        <caption>
          <h5 style="display:inline-flex;" class="float-left">Ordered Items</h5>
          <button class="btn btn-primary btn-sm mr-2 float-right" data-toggle="modal" data-target="#paymentModal">Payment</button>
        </caption>
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item No.</th>
            <th scope="col">Item Description</th>
            <th scope="col">Quantity</th>
            <th scope="col">Price</th>
            <th scope="col">Discount %</th>
            <th scope="col">Total</th>
            <th scope="col"></th>
          </tr>
        </thead>
        <tbody id="LineItems">
          <tr v-for="(item, i) in items" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>{{ item.itemcode }}</td>
            <td>
              <v-select label="itemname" @input="(val) => itemSelected(val, i)" @search="(search, loading) => fetchProduct(search, loading, i)" v-model="item.itemname" :options="inventory" :clearable="false" placeholder="Kindly select Product"></v-select>
            </td>
            <td>
              <input type="number" min="1" step="1" class="form-control form-control-sm" :value="item.quantity" @blur="setQuantity(i)" />
            </td>
            <td>
              {{ item.price }}
            </td>
            <td>
              <div :contenteditable="isAdmin" @click="getPermission" @keypress="validateInput" @blur="applyDiscount(i)" v-text="item.discount"></div>
            </td>
            <td>
              {{ item.total }}
            </td>
            <td>
              <button :id="'del' + i" class="btn btn-danger btn-sm mr-2 float-right" @click="deleteItemRow( i)">Remove Line</button>
            </td>
          </tr>
        </tbody>
        <tfoot id="ItemsFooter">
          <tr>
            <td colspan="7" class="text-right font-weight-bold">
              Subtotal:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{subTotal}}
            </td>
          </tr>
          <tr>
            <td colspan="7" class="text-right font-weight-bold">7.5% VAT:</td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{vatAmount}}
            </td>
          </tr>
          <tr>
            <td colspan="7" class="text-right font-weight-bold">
              Amount Paid:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{amtPaid}}
            </td>
          </tr>
          <tr>
            <td colspan="7" class="text-right font-weight-bold">
              Grand Total:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{grandTotal}}
            </td>
          </tr>
        </tfoot>
      </table>
    </div>

    <button type="button" class="btn btn-dark btn-sm float-right ml-2" @click="DraftOrder">
      Save as Draft
    </button>
    
    <button type="button" class="btn btn-primary btn-sm float-right" @click="SaveOrder" :disabled="isDisabled">
      Save
    </button>

    <!-- Modals -->
    <div class="modal fade" id="adminModal" data-backdrop="static" tabindex="-1" aria-labelledby="adminModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header text-center bg-dark text-white">
            <h5 class="modal-title" id="adminModalLabel">Authorize Discount</h5>
            <button type="button" class="close text-white" data-dismiss="modal" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="adminLogin">
              <div class="form-signin">
                <div class="text-center mb-4">
                  <img class="mb-4" src="../../assets/images/slot.png" />
                  <h1 class="h3 mb-3 font-weight-normal">Sign in</h1>
                  <p>
                    Kindly sign in here to be able to authorize the Discount.
                    <br />
                    <code>
                      Only a Valid System Administrator can authorize this action.
                    </code>
                  </p>
                </div>

                <div class="form-group">
                  <label for="inputEmail">Email address</label>
                  <input type="email" v-model="email" class="form-control" placeholder="Email address" required autofocus />
                </div>

                <div class="form-group">
                  <label for="inputPassword">Password</label>
                  <input type="password" v-model="password" class="form-control" placeholder="Password" required />
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary btn-sm mr-2" @click="adminLogin">Authorize</button>
            <button type="button" class="btn btn-secondary btn-sm" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" id="paymentModal" data-backdrop="static" tabindex="-1" aria-labelledby="paymentModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header text-center bg-dark text-white">
            <h5 class="modal-title" id="paymentModalLabel">Payment Details</h5>
            <button type="button" class="close text-white" data-dismiss="modal" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <div class="container-fluid">
              <div class="table-responsive">
                <table class="table table-fixed table-bordered table-hover table-sm">
                  <caption>Registered Payments for this Order.</caption>
                  <thead class="thead-dark">
                    <tr>
                      <th scope="col">#</th>
                      <th scope="col">Payment Method</th>
                      <th scope="col">Bank Used</th>
                      <th scope="col">Amount Paid</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody id="paymentLines">
                    <tr v-for="(payment, i) in payments" :key="'prow' + i" :id="'prow' + i">
                      <th scope="row">{{ i + 1 }}</th>
                      <td>
                        <select class="form-control" @change="paymentMethod($event, i)" :v-model="payment.method">
                          <option value="null" selected>Select Payment Method</option>
                          <option>Cash</option>
                          <option>POS</option>
                          <option>Matrix</option>
                          <option>EasyBuy</option>
                          <option>Credpal</option>
                          <option>Bank Transfer</option>
                        </select>
                      </td>
                      <td>
                        <select disabled="true" class="form-control form-control-sm" :v-model="payment.bank">
                          <option value="null" selected>Select Inflow Bank</option>
                          <option :key="bank.id" :value="bank.name" v-for="bank in banks">{{ bank.code }}</option>
                        </select>
                      </td>
                      <td>
                        <input type="number" class="form-control form-control-sm" placeholder="Amount Paid" :value="payment.amount" min="0.00" @blur="PaymentTotal(i, $event)" disabled="true" />
                      </td>
                      <td style="white-space: nowrap;">
                        <button :id="'pdel' + i" class="btn btn-danger btn-sm mr-2 float-right" @click="deleteRow(i)">Remove Line</button>
                      </td>
                    </tr>
                  </tbody>
                  <tfoot id="PaymentFooter">
                    <tr>
                      <td colspan="2" class="text-right font-weight-bold">Balance Due</td>
                      <td class="font-weight-bold bg-primary text-white">₦{{balanceDue}}</td>
                      <td class="text-right font-weight-bold">Amount Paid</td>
                      <td class="font-weight-bold bg-primary text-white">₦{{amtPaid}}</td>
                    </tr>
                  </tfoot>
                </table>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary btn-sm mr-2" @click="isValid">Add Payment</button>
            <button type="button" class="btn btn-secondary btn-sm" data-dismiss="modal">Close</button>
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
import moment from "moment";
import "vue-select/dist/vue-select.css";

export default {
  data() {
    return {
      // Login
      email: "",
      password: "",
      
      banks: [],
      items: [],
      order: {},
      item: null,
      payments: [],
      customers: [],
      inventory: [],
      isAdmin: false,
      customer: null,
      amtPaid: '0.00',
      isDisabled: true,
      selecteditem: '',
      ordercolumns: [],
      subTotal: '0.00',
      vatAmount: '0.00',
      grandTotal: '0.00',
      balanceDue: '0.00',
      paymentcolumns: [],
      currentDate: moment().format("Do of MMMM YYYY"),
      dateColumns:['created_at','updated_at', 'deleted_at'],
    };
  },
  async created() {
    // Determine the state of the Discount element
    if (this.$store.state.isAdmin)
    {
      this.isAdmin = true;
    }

    // Get all customers
    window.backend.GetCustomers().then((customers) => {
      this.customers = customers;
      this.customers.forEach((customer) => {
        if (customer.cardcode === "0") {
          customer.cardcode = `R-${customer.id}`;
        }
      });
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });

    // Get all inventory
    window.backend.GetProducts().then((inventory) => {
      this.inventory = inventory;
      this.addItemRow();
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });

    // Get all banks
    window.backend.GetBanks().then((banks) => {
      this.banks = banks;
      this.addRow();
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });
  },
  methods: {
    /**
     * Triggered when the search text changes.
     *
     * @param search  {String}    Current search text
     * @param loading {Function}	Toggle loading class
     */
    fetchCustomer(search, loading) {
      loading(true);
      if (search.length >= 3) {
        this.customer = this.customers.filter((customer) => {
          return (
            customer.phone.toLowerCase() === search.toLowerCase() ||
            customer.phone1.toLowerCase() === search.toLowerCase() ||
            customer.email.toLowerCase() === search.toLowerCase() ||
            customer.cardname.toLowerCase() === search.toLowerCase() ||
            customer.cardcode.toLowerCase() === search.toLowerCase()
          );
        })[0];
      }
      loading(false);
      return;
    },
    async addItemRow(index) {
      if ((index + 1) < this.items.length) {
        return;
      }

      this.selecteditem = '';
      this.items.push({
        id: this.items.length + 1,
        quantity: 1,
        itemcode: '',
        itemname: '',
        discount: '0%',
        price: '₦0.00',
        total: '₦0.00',
      });

      const el = document.getElementById("ItemsFooter");
      if (el) {
        el.scrollIntoView();
      }
    },
    async deleteItemRow(index) {
      this.$delete(this.items, index);
      
      if (this.items.length === 0) {
        await this.addItemRow(index);
      }
    },
    async populateRow(rowIndex) {
      if (this.item === undefined) {
        return;
      }

      // Set the table data
      if (this.items[rowIndex].itemcode === "") {
        this.items[rowIndex].quantity = 1;
        this.items[rowIndex].id = rowIndex;
        this.items[rowIndex].discount = '0%';
        this.items[rowIndex].price = `₦${this.item.price}`;
        this.items[rowIndex].itemcode = this.item.itemcode;
        this.items[rowIndex].itemname = this.item.itemname;
      } else {
        this.items[rowIndex].price = `₦${this.item.price}`;
        this.items[rowIndex].itemcode = this.item.itemcode;
        this.items[rowIndex].itemname = this.item.itemname;
      }

      // Calculate the totals [invoice subtotal, grand total, vat]
      await this.totals();
    },
    async itemSelected(val, rowIndex) {
      this.item = val;
      await this.populateRow(rowIndex);
      await this.addItemRow(rowIndex);
    },
    async fetchProduct(search, loading, rowIndex) {
      loading(true);
      if (search.length >= 3) {  
        this.item = await this.inventory.filter((item) => {
          return (
            item.itemcode.toLowerCase() === search.toLowerCase() ||
            item.itemname.toLowerCase() === search.toLowerCase() ||
            item.codebars.toLowerCase() === search.toLowerCase() ||
            item.serialnumber.toLowerCase() === search.toLowerCase()
          );
        })[0];
        await this.populateRow(rowIndex);
      }
      loading(false);
      return;
    },
    async adminLogin() {
      const email = this.email,
        password = this.password;

      if (password === "" || email === "") {
        this.$toast.error("Error! Invalid Credentials Supplied.");
        return false;
      }

      window.backend.Login(email, password).then((result) => {
        if (result.id !== undefined && result.isadmin === true) {
          this.email = "";
          this.password = "";
          this.isAdmin = true;
          // Hide the modal
          $('#adminModal').modal('hide');
        } else {
          this.$toast.error("Error! Invalid login Credentials.");
        }
      },(err) => {
          this.$toast.success("Error! " + err);
      });
    },
    async getPermission() {
      // Check for permissions
      if (this.isAdmin === false) {
        // Display admin login Modal
        $('#adminModal').modal('show');
      }
    },
    async validateInput() {
      // Perform a quick clean up
      if (event.target.innerText === '0%') {
        event.target.innerText = '';
      }

      // Reject all keypress characters that are not numeric
      if (isNaN(String.fromCharCode(event.which)) || event.which == '13') {
        event.preventDefault();
      }
    },
    async applyDiscount(rowIndex) {
      // Perform a quick clean up
      if (this.items[rowIndex].price === '₦0.00') {
        event.target.innerText = '0%';
        this.$toast.error("Error! Discount cannot be applied on ₦0.00. \nKindly select a product with a valid Price.");
        return;
      }

      if (event.target.innerText === '' || event.target.innerText < 0 || event.target.innerText > 100) {
        event.target.innerText = '0%';
        this.$toast.error("Error! Discount must be between 0% and 100%;");
        return;
      }

      // Reset the Discount element to its inital state.
      if (this.$store.state.isAdmin)
      {
        this.isAdmin = false;
      }

      if (event.target.innerText !== '0%') {
        // Add the percentage Symbol
        event.target.innerText += "%";
      }

      // Store the data into the items array
      this.items[rowIndex].discount = event.target.innerText;

      // Calculate the totals [invoice subtotal, grand total, vat]
      await this.totals();
    },
    async setQuantity(rowIndex) {
      this.items[rowIndex].quantity = event.target.value;
      await this.totals();
    },
    async totals() {
      this.subTotal = 0;
      this.vatAmount = 0;
      this.grandTotal = 0;
      let runningTotal = 0.0;

      // Loop through the array and perform all needed calculations
      this.items.forEach(element => {
        // Calculate the % discount.
        let numVal1 = parseFloat(element.price.replace("₦", "")),
          numVal2 = parseFloat(element.discount.replace("%", "")) / 100,
          currentTotal = parseFloat(Number(element.quantity) * (numVal1 - (numVal1 * numVal2))).toFixed(2);
        
        element.total = "₦" + currentTotal;
        runningTotal += parseFloat(currentTotal);
      });
      
      // Calculate footer details
      this.subTotal = parseFloat(runningTotal);
      this.vatAmount = parseFloat((7.5 / 100) * runningTotal).toFixed(2);
      this.grandTotal = parseFloat((7.5 / 100) * runningTotal + runningTotal).toFixed(2);
      this.balanceDue = parseFloat(this.grandTotal - this.amtPaid).toFixed(2);
    },
    // Payment Section
    async paymentMethod(event, index) {
      let ctrl = event.target.value,
        amt = event.target.parentElement.parentElement.cells[3].childNodes[0],
        bank = event.target.parentElement.parentElement.cells[2].childNodes[0];

      bank.disabled = true;
      amt.disabled = false;

      if (ctrl.toLowerCase() === "pos" || ctrl.toLowerCase() === "bank transfer") {
        // Display list of banks for selection and allow for value entry.
        bank.disabled = false;
      }

      await this.addRow(index);
    },
    async addRow(index) {
      if ((index + 1) < this.payments.length) {
        return;
      }
      
      this.payments.push({
        name: null,
        method: null,
        amount: '0.00',
        id: this.payments.length + 1,
      });

      const el = document.getElementById("PaymentFooter");
      if (el) {
        el.scrollIntoView();
      }
    },
    async deleteRow(index) {
      this.$delete(this.payments, index);
      
      if (this.payments.length === 0) {
        this.addRow(index);
      }
    },
    PaymentTotal(index, event) {
      let runningTotal = 0.0,
        amt = event.target.value,
        bank = event.target.parentElement.parentElement.cells[2].childNodes[0].value,
        method = event.target.parentElement.parentElement.cells[1].childNodes[0].value;

      if (bank === null) {
        this.$toast.error("Error! Bank Details is required in order to proceed.");
        return;
      }

      if (amt === '0.00' || amt === 0) {
        this.$toast.error("Error! Invalid Amount inputted.");
        return;
      }

      this.payments[index].name = bank;
      this.payments[index].amount = amt;
      this.payments[index].method = method;

       // Loop through the array and perform all needed calculations
      this.payments.forEach(element => {
        runningTotal += parseFloat(element.amount.replace("₦", ""));
      });

      // Calculate footer details
      this.amtPaid = parseFloat(runningTotal).toFixed(2);
      this.balanceDue = parseFloat(this.grandTotal - this.amtPaid).toFixed(2);
    },
    isValid() {
      let isvalid = true;
      if (this.amtPaid !== this.grandTotal) {
        this.isDisabled = true;
        this.$toast.info("Error! Amount paid is less than Grand Total.\nKindly correct.");
        return;
      }

      this.payments.forEach(element => {
        if (element.name === "null" && (element.method.toLowerCase() === "pos" || element.method.toLowerCase() === "bank transfer")) {
          isvalid = false;
        }
      });

      if (isvalid === false) {
        this.$toast.error("Error! Bank Details is required in order to proceed.");
        return;
      }

      this.isDisabled = false;
      $('#paymentModal').modal('hide')
    },
    SaveOrder() {
      if (this.customer === null) {
        this.$toast.error("Error! You are yet to associate a customer to this order.");
        return;
      }

      if (this.items.length === 0) {
        this.$toast.error("Error! You are not permitted to create an empty Sales Order.");
        return;
      }

      if (this.amtPaid !== this.grandTotal) {
        this.$toast.error("Error! You are yet to finalize payment on this order.");
        return;
      }

      var r = confirm("Are you sure you want to create this Sales Order!");
      if (r == false) {
        return;
      }

      let items = [],
        payments = [];

      this.$delete(this.items, this.items.length - 1);
      this.$delete(this.payments, this.payments.length - 1);

      // Loop through the array and perform all needed calculations
      this.items.forEach(element => {
        let orderItem = {
          orderid: null,
          itemname: element.itemname,
          itemcode: element.itemcode,
          quantity: element.quantity,
          price: element.price.replace("₦", ""),
          discount: element.discount.replace("%", ""),
        };

        // Add the items to the items array
        items.push(orderItem);
      });
      
      // Loop through the array and perform all needed calculations
      this.payments.forEach(element => {
        let payment = {
          id: 0,
          docnum: 0,
          docentry: 0,
          orderid: null,
          canceled: false,
          amount: element.amount,
          paymenttype: element.method,
          paymentdetails: element.name,
        };

        // Add the payments to the payments array
        payments.push(payment);
      });

      // Build out the header.
      this.order = {
        id: 0,
        docnum: 0,
        docentry: 0,
        vatsum: 7.5,
        synced: false,
        canceled: false,
        cardname: this.customer.cardname,
        cardcode: this.customer.cardcode,
        created_by: this.$store.state.user.id,
      };
      this.order.items = items;
      this.order.payments = payments;
      this.order.doctotal = this.grandTotal;
            
      window.backend.NewOrder(this.order).then(() => {
        this.$toast.success("Success! Order has been successfully saved.");
        this.$router.push({name: 'orderlist'});
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    },
    DraftOrder() {},
  },
};
</script>