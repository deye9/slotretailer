<template>
  <section style="margin-top: 3em">
    <div class="row">
      <div class="col-8">
        <h3>New Order</h3>
      </div>
      <div class="col-4">
        <router-link to="/orders/" class="btn btn-info float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="cardcode">Customer</label>
        <v-select label="cardname" @search="fetchCustomer" :options="customers" v-model="customer"></v-select>
      </div>
      <!-- <div class="form-group col">
        <label>Payment Type</label>
        <select class="form-control" @change="PaymentMethod($event)">
          <option selected>Cash</option>
          <option>POS</option>
          <option>Matrix</option>
          <option>EasyBuy</option>
          <option>Credpal</option>
          <option>Bank Transfer</option>
        </select>
      </div>
      <div class="form-group col">
        <label>Payment Details</label>
        <input type="text" class="form-control" :disabled="isDisabled" placeholder="Payment Details" />
      </div> -->
    </div>

    <h3 style="margin-bottom: 50px">
      <span class="float-left">
        Items <small class="text-primary">Double click on row to edit.</small>
      </span>
      <b-button v-b-modal.modal-item variant="info" class="float-right">Add Product</b-button>
    </h3>

    <hr />
    <div class="table-responsive" style="max-height: 325px; overflow: scroll">
      <table class="table table-bordered table-hover table-striped">
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

    <!-- <button type="button" class="btn btn-info float-right" style="margin-left: 2px" @click="SaveOrder" :disabled="isDisabled">
      Save Order
    </button> -->

    <b-button v-b-modal.modal-payment :variant="paymentVariant" class="float-right" :disabled="isDisabled">Payment</b-button>

    <button type="button" class="btn btn-dark float-right" style="margin-right: 2px" @click="DraftOrder">
      Save as Draft
    </button>

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

    <b-modal id="modal-payment" centered title="Add New Item Row" :header-bg-variant="headerBgVariant"
      :header-text-variant="headerTextVariant" @ok="AddPayment"  :body-bg-variant="bodyBgVariant"
      :body-text-variant="bodyTextVariant" :footer-bg-variant="footerBgVariant"  :footer-text-variant="footerTextVariant">
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
export default {
  data() {
    return {
      rows: [],
      docnum: 0,
      order: {},
      banks: [],
      item: null,
      customers: [],
      inventory: [],
      customer: null,
      isDisabled: true,

      // Modal
      bodyBgVariant: "light",
      bodyTextVariant: "dark",
      headerBgVariant: "dark",
      footerTextVariant: "dark",
      headerTextVariant: "light",
      paymentVariant: "secondary",
      footerBgVariant: "secondary",
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
    resetModal() {
      this.item = null;
      document.getElementById("price").value = 1;
      document.getElementById("discount").value = 0;
      document.getElementById("price").value = "₦0.00";
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

      document.getElementById("subTotal").innerHTML = `₦${val}`;
      document.getElementById("grandTotal").innerHTML = `₦${
        (7.5 / 100) * val + val
      }`;
      document.getElementById("vatAmount").innerHTML = `₦${parseFloat(
        (7.5 / 100) * val
      ).toFixed(2)}`;
    },
    async AddPayment(bvModalEvt) {

    },
    PaymentMethod(event) {
      // It can be a mixture of payment methods.
      switch (event.target.value.toLowerCase()) {
        case "matrix": // Allow for entry of value
        case "easybuy":
        case "creditpal":
          this.isDisabled = false;
          break;

        case "pos":
        case "bank transfer": // Display list of banks for selection and allow for value entry.
          this.isDisabled = true;
          break;

        default:
          this.isDisabled = true;
          break;
      }
    },
    DraftOrder() {},
    SaveOrder() {
      let val = 0,
        items = [],
        tableRef = document.getElementById("orderedItems"),
        rowCnt = tableRef.rows.length;

      if (this.customer === null) {
        this.$toast.error(
          "Error! You are yet to associate a customer to this order."
        );
        return;
      }

      if (rowCnt === 0) {
        this.$toast.error(
          "Error! You are not permitted to create an empty Sales Order."
        );
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
          price: tableRef.rows[i].cells[4].innerHTML,
          quantity: tableRef.rows[i].cells[3].innerHTML,
          discount: tableRef.rows[i].cells[5].innerHTML,
          itemname: tableRef.rows[i].cells[2].innerHTML,
          itemcode: tableRef.rows[i].cells[1].innerHTML,
        };

        // Get the total value of the order
        val += parseFloat(tableRef.rows[i].cells[6].innerHTML);

        // Add the ordered items to the items array
        items.push(orderItem);
      }

      this.order.payments = null;

      this.order.items = items;
      this.order.doctotal = `₦${(7.5 / 100) * val + val}`;
      document.getElementById("subTotal").innerHTML = `₦${val}`;
      document.getElementById("vatAmount").innerHTML = `₦${parseFloat(
        (7.5 / 100) * val
      ).toFixed(2)}`;

      window.backend.NewOrder(this.order).then(
        () => {
          this.$toast.success("Success! Order has been successfully saved.");
          this.$router.push("/orders/");
        },
        (err) => {
          this.$toast.error("Error! " + err);
        }
      );
    },
  },
};
</script>

