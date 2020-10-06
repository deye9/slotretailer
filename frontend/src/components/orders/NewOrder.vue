<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>New Sales Order</h3>
      </div>
      <div class="col-4">
        <router-link to="/orders/" class="btn btn-info float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row" small>
      <div class="form-group col">
        <label for="cardname">Customer</label>
        <v-select label="cardname" @search="fetchCustomer" :options="customers" v-model="customer"></v-select>
      </div>
      <div class="form-group col">
        <label>Customer Code</label>
        <input v-if="this.customer !== null" type="text" class="form-control" placeholder="City" disabled v-model="this.customer.cardcode" />
      </div>
      <div class="form-group col">
        <label>Create Date</label>
        <input type="text" class="form-control" placeholder="City" disabled :value="this.currentDate" />
      </div>
      <div class="form-group col">
        <label>Created By</label>
        <input type="text" class="form-control" placeholder="City" disabled :value="this.$store.state.user.email" />
      </div>
    </div>

    <h3 style="margin-bottom: 50px">
      <span class="float-left">
        Order Details <small>Double click to remove item.</small>
      </span>
      <b-button v-b-modal.modal-item variant="info" class="float-right">Add Product</b-button>
    </h3>
    <hr />

    <div style="max-height: 300px; overflow: scroll">
      <div class="table-responsive-sm">
        <table class="table table-bordered table-hover table-striped table-sm">
          <thead class="thead-dark">
            <tr>
              <th scope="col">#</th>
              <th scope="col">Item Code</th>
              <th scope="col">Item Name</th>
              <th scope="col">Qty</th>
              <th scope="col">Price</th>
              <th scope="col">Discount (%)</th>
              <th scope="col">Total</th>
            </tr>
          </thead>
          <tbody id="orderedItems"></tbody>
          <tfoot>
            <tr>
              <td colspan="6" class="text-right font-weight-bold">
                Invoice Subtotal:
              </td>
              <td id="subTotal" class="font-weight-bold bg-primary text-white">
                ₦0.00
              </td>
            </tr>
            <tr>
              <td colspan="6" class="text-right font-weight-bold">7.5% VAT:</td>
              <td id="vatAmount" class="font-weight-bold bg-primary text-white">
                ₦0.00
              </td>
            </tr>
            <tr>
              <td colspan="6" class="text-right font-weight-bold">
                Grand Total:
              </td>
              <td id="grandTotal" class="font-weight-bold bg-primary text-white">
                ₦0.00
              </td>
            </tr>
          </tfoot>
        </table>
      </div>

      <br />
      <h3>Payment Details</h3>
      <hr />
      <div class="table-responsive-sm">
        <b-table id="paymentList" name="paymentList" :disabled="true" small striped hover :items="payments" :fields="fields" primary-key="_id">
          <template v-slot:cell(payment_method)="row">
            <select id="paymentMethod" name="paymentMethod" class="form-control" @change="paymentMethod($event, row)">
              <option selected>Cash</option>
              <option>POS</option>
              <option>Matrix</option>
              <option>EasyBuy</option>
              <option>Credpal</option>
              <option>Bank Transfer</option>
            </select>
          </template>
          <template v-slot:cell(banks)>
            <b-form-select id="bankslist" name="bankslist" :options="banks" :disabled="true" class="mb-3" value-field="name" text-field="code"></b-form-select>
          </template>
          <template v-slot:cell(amount_paid)="row">
            <input id="amt" name="amt" type="number" class="form-control" @change="updatePayment(row)" placeholder="Amount Paid" value="0" min="1" />
          </template>
          <template v-slot:cell(actions)="row">
            <b-button size="sm" variant="primary" @click="addRow(row)" class="mr-1">
              <b-icon icon="plus-square-fill" aria-hidden="true">
                Add New Payment
              </b-icon>
            </b-button>
            <b-button size="sm" variant="danger" @click="removePayment(row.index)" class="mr-1">
              <b-icon icon="trash" aria-hidden="true"></b-icon>
            </b-button>
          </template>
          <template v-slot:custom-foot>
            <b-tr>
              <b-td colspan="2">
                Amount Paid: <span id="amtPaid" style="font-weight: bold;">₦0.00</span>
              </b-td>
              <b-td colspan="2" style="text-align:center;">
                Balance Due: <span id="balanceDue" style="font-weight: bold;">₦0.00</span>
              </b-td>
              <b-td>
                Grand Total: <span id="expectedPayment" style="font-weight: bold;">₦0.00</span>
              </b-td>
            </b-tr>
          </template>
        </b-table>
      </div>

      <button type="button" class="btn btn-primary float-right" @click="SaveOrder">
        Save
      </button>

      <button type="button" class="btn btn-dark float-right" style="margin-right: 2px" @click="DraftOrder">
        Save as Draft
      </button>
    </div>

    <!-- Modals -->
    <b-modal id="modal-item" centered title="Add New Item Row" :header-bg-variant="headerBgVariant"
      @hidden="resetModal" @show="resetModal" :header-text-variant="headerTextVariant" @ok="AddItem"
      :body-bg-variant="bodyBgVariant" :body-text-variant="bodyTextVariant"
      :footer-bg-variant="footerBgVariant" :footer-text-variant="footerTextVariant">
      <b-container fluid>
        <b-row class="mb-1">
          <b-col cols="3">
            <label for="product">Product:</label>
          </b-col>
          <b-col>
            <v-select id="product" name="product" label="itemname" @search="fetchItem" @input="getItem" :options="inventory" v-model="item"></v-select>
          </b-col>
        </b-row>

        <hr />
        <b-row class="mb-1">
          <b-col cols="3">
            <label for="quantity">Quantity:</label>
          </b-col>
          <b-col>
            <input id="quantity" name="quantity" type="number" class="form-control" placeholder="Quantity" min="1" value="1" step="1" />
          </b-col>
        </b-row>

        <hr />
        <b-row class="mb-1">
          <b-col cols="3">
            <label for="price">Price:</label>
          </b-col>
          <b-col>
            <input id="price" name="price" type="text" class="form-control" placeholder="₦0.00" disabled value="₦0.00" />
          </b-col>
        </b-row>

        <hr />
        <b-row class="mb-1">
          <b-col cols="3">
            <label for="discount">Discount:</label>
          </b-col>
          <b-col>
            <input id="discount" name="discount" type="number" class="form-control" placeholder="Discount" min="1" step="0.1" value="0" />
          </b-col>
        </b-row>

        <hr />

        <h4><strong>Total</strong> will be shown on the main page.</h4>
      </b-container>
    </b-modal>

  </section>
</template>

<script>
import moment from "moment";

export default {
  data() {
    return {
      rows: [],
      docnum: 0,
      fired: [],
      order: {},
      banks: [],
      item: null,
      payments: [],
      customers: [],
      inventory: [],
      amountPaid: 0,
      grandTotal: 0,
      canPay: false,
      customer: null,
      isDisabled: true,
      currentPayment: [],
      currentDate: moment().format("dddd, MMMM Do YYYY, h:mm:ss a"),
      
      // Modal
      bodyBgVariant: "light",
      bodyTextVariant: "dark",
      headerBgVariant: "dark",
      footerTextVariant: "dark",
      headerTextVariant: "light",
      paymentVariant: "secondary",
      footerBgVariant: "secondary",

      // Payments
      fields: ['_id', 'payment_method', 'banks', 'amount_paid', 'actions'],
    };
  },
  async created() {
    // Get all customers
    window.backend.GetCustomers().then(
      (customers) => {
        this.customers = customers;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );

    // Get all inventory
    window.backend.GetProducts().then(
      (inventory) => {
        this.inventory = inventory;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );

    // Get all banks
    window.backend.GetBanks().then(
      (banks) => {
        this.banks = banks;
        this.payments.push({_id: this.payments.length + 1 });
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );
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
      if (search.length >= 11) {
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
    fetchItem(search, loading) {
      loading(true);
      if (search.length >= 3) {
        this.item = this.inventory.filter((item) => {
          return (
            item.itemcode.toLowerCase() === search.toLowerCase() ||
            item.itemname.toLowerCase() === search.toLowerCase() ||
            item.codebars.toLowerCase() === search.toLowerCase() ||
            item.serialnumber.toLowerCase() === search.toLowerCase()
          );
        })[0];
        document.getElementById("price").value = `₦${this.item.price}`;
      }
      loading(false);
      return;
    },
    getItem(value) {
      this.item = value;
      document.getElementById("price").value = `₦${this.item.price}`;
    },

    // Items Modal
    async removeItemRow(event) {
      let id = event.target.parentNode.id, 
        row = document.getElementById(`${id}`);

      // Remove the table row from the table
      row.parentNode.removeChild(row);

      // Calculate the Grand Total
      await this.TableTotal(document.getElementById("orderedItems"));
    },
    async AddItem(bvModalEvt) {
      if (this.item === null) {
        this.$toast.error("Error! No Product selected.");

        // Prevent modal from closing
        bvModalEvt.preventDefault();

        return;
      }

      if (document.getElementById("quantity").value <= 0) {
        this.$toast.error("Error! Invalid Product quantity.");

        // Prevent modal from closing
        bvModalEvt.preventDefault();

        return;
      }

      // Get a reference to the Table body
      let total = 0,
        discount = document.getElementById("discount").value,
        quantity = document.getElementById("quantity").value,
        tableRef = document.getElementById("orderedItems"),
        rowCnt = tableRef.rows.length + 1,
        // Insert a row in the table at the last row
        row = tableRef.insertRow(),
        price = document.getElementById("price").value.replace("₦", "");

      row.id = `row${rowCnt}`;
      row.addEventListener("dblclick", this.removeItemRow);

      // Insert the needed cells
      for (let index = 0; index <= 6; index++) {
        row.insertCell(index);
      }

      // Update the innerHTML of the cells to the new values.
      row.cells[0].innerHTML = rowCnt;
      row.cells[1].innerHTML = this.item.itemcode;
      row.cells[2].innerHTML = this.item.itemname;
      row.cells[3].innerHTML = quantity;
      row.cells[4].innerHTML = `₦${price}`;
      row.cells[5].innerHTML = discount;

      if (parseInt(discount) === 0) {
        total = `₦${Number(quantity) * parseFloat(price)}`;
      } else {
        let numVal = Number(discount / 100),
          discountVal = price - price * numVal;
        total = `₦${(quantity * price - discountVal).toFixed(2)}`;
      }
      row.cells[6].innerHTML = total;
      
      // Calculate the Grand Total
      await this.TableTotal(tableRef);

      // Reset the modal controls /values
      this.resetModal();

      // Allow for payment as we have a valid item now
      this.isDisabled = false;
      this.paymentVariant = "info";
    },
    async TableTotal(tableRef) {
      let val = 0.0;
      
      for (var i = 0; i < tableRef.rows.length; i++) {
        let data = tableRef.rows[i].cells[6].innerHTML.replace("₦", "");
        val += parseFloat(data);
      }

      this.grandTotal = (7.5 / 100) * val + val;
      document.getElementById("subTotal").innerHTML = `₦${val}`;
      document.getElementById("grandTotal").innerHTML = `₦${parseFloat(this.grandTotal).toFixed(2)}`;
      document.getElementById("vatAmount").innerHTML = `₦${parseFloat((7.5 / 100) * val).toFixed(2)}`;

      // Payment Details
      document.getElementById("balanceDue").innerHTML = document.getElementById("grandTotal").innerHTML;
      document.getElementById("expectedPayment").innerHTML = document.getElementById("grandTotal").innerHTML;
    },
    async resetModal() {
      this.item = null;
      if (document.getElementById("quantity") !== null) {
        document.getElementById("quantity").value = 1;
        document.getElementById("discount").value = 0;
        document.getElementById("price").value = "₦0.00";
      }
    },

    // Payment Section
    async updatePayment() { 
      let amountPaid = 0.0, 
        totalAmt = document.getElementById("expectedPayment").innerHTML,
        tableRef = document.getElementById("paymentList").getElementsByTagName('tbody')[0],
        rowCnt = tableRef.rows.length;

      this.currentPayment = [];

      // Get the amount paid so far.
      for (var i = 0; i < rowCnt; i++) {
        if (tableRef.rows[i].cells[2].childNodes[0].value === "") {
          this.currentPayment.push({
            id: this.id,
            docnum: 0,
            docentry: 0,
            orderid: null,
            canceled: false,
            amount: tableRef.rows[i].cells[3].childNodes[0].value,
            paymenttype: tableRef.rows[i].cells[1].childNodes[0].value,
          });
        } else {
          this.currentPayment.push({
            id: this.id,
            docnum: 0,
            docentry: 0,
            orderid: null,
            canceled: false,
            amount: tableRef.rows[i].cells[3].childNodes[0].value,
            paymenttype: tableRef.rows[i].cells[1].childNodes[0].value,
            paymentdetails: tableRef.rows[i].cells[2].childNodes[0].value,
          });
        }
        amountPaid += parseFloat(tableRef.rows[i].cells[3].childNodes[0].value);
      }
      document.getElementById('amtPaid').innerHTML = `₦${parseFloat(amountPaid).toFixed(2)}`;

      // Calculate the difference and display the balance due
      document.getElementById("balanceDue").innerHTML = `₦${parseFloat(totalAmt.replace("₦", "") - amountPaid).toFixed(2)}`;
    },
    async addRow(row) {
      // Avoid add multiple rows when payment methods change.
      if (!this.fired.includes(row.index)) {
        // Add a new row automatically.
        this.payments.push({_id: this.payments.length + 1 });

        // Log it to the array to avoid it firing twice
        this.fired.push(row.index);

        let rowID = `paymentList__row_${row.item._id}`,
        selectedRow = document.getElementById(rowID);

        selectedRow.cells[4].childNodes[0].disabled = true;
      }
    },
    async removePayment(index) {
      if (this.payments.length > 1) {
        // Remove the row from the table
        document.getElementById("paymentList").getElementsByTagName('tbody')[0].deleteRow(index);

        await this.updatePayment();
        this.$toast.success("Success! Payment has been successfully deleted.");
      }
    },    
    async paymentMethod(event, row) {
      let rowID = `paymentList__row_${row.item._id}`,
        selectedRow = document.getElementById(rowID);

      selectedRow.cells[2].childNodes[0].disabled = true;

      if (event.target.value.toLowerCase() === "pos" || event.target.value.toLowerCase() === "bank transfer") {
        // Display list of banks for selection and allow for value entry.
        selectedRow.cells[2].childNodes[0].disabled = false;
      }
    },
    async SaveOrder() {
      let val = 0,
        items = [],
        tableRef = document.getElementById("orderedItems"),
        rowCnt = tableRef.rows.length,
        balance = document.getElementById("balanceDue").innerHTML;

      if (this.customer === null) {
        this.$toast.error("Error! You are yet to associate a customer to this order.");
        return;
      }

      if (rowCnt === 0) {
        this.$toast.error("Error! You are not permitted to create an empty Sales Order.");
        return;
      }

      if (balance !== '₦0.00' || this.currentPayment.length === 0) {
        this.$toast.error("Error! You are yet to finalize payment on this order.");
        return;
      }

      var r = confirm("Are you sure you want to create this Sales Order!");
      if (r == false) {
        return;
      }

      // Build out the header.
      this.order = {
        id: this.id,
        docentry: 0,
        vatsum: 7.5,
        synced: false,
        canceled: false,
        docnum: this.docnum,
        cardname: this.customer.cardname,
        cardcode: this.customer.cardcode,
        created_by: this.$store.state.user.id,
      };

      // Build out the Line items.
      for (var i = 0; i < rowCnt; i++) {
        let orderItem = {
          orderid: null,
          quantity: tableRef.rows[i].cells[3].innerHTML,
          discount: tableRef.rows[i].cells[5].innerHTML,
          itemname: tableRef.rows[i].cells[2].innerHTML,
          itemcode: tableRef.rows[i].cells[1].innerHTML,
          price: tableRef.rows[i].cells[4].innerHTML.replace("₦", ""),
        };

        // Get the total value of the order
        val += parseFloat(tableRef.rows[i].cells[6].innerHTML.replace("₦", ""));

        // Add the ordered items to the items array
        items.push(orderItem);
      }

      this.order.items = items;
      this.order.payments = this.currentPayment;
      this.order.doctotal = (7.5 / 100) * val + val;
      document.getElementById("subTotal").innerHTML = `₦${val}`;
      document.getElementById("vatAmount").innerHTML = `₦${parseFloat((7.5 / 100) * val).toFixed(2)}`;

      window.backend.NewOrder(this.order).then(() => {
          this.$toast.success("Success! Order has been successfully saved.");
          this.$router.push("/orders/");
        },
        (err) => {
          this.$toast.error("Error! " + err);
        }
      );
    },
    DraftOrder() {},
  },
};
</script>