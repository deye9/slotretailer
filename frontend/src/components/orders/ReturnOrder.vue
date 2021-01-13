<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Create Returns for Order</h3>
      </div>
      <div class="col-4">
        <router-link :to="{ name: 'orderlist' }" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label>Order Number</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.id" />
      </div>
      <div class="form-group col">
        <label for="docnum">SAP Document Number</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.docNum" />
      </div>      
      <div class="form-group col">
        <label for="cardname">Customer</label>
        <input type="text" class="form-control form-control-sm" placeholder="Customer Name" disabled :value="this.order.cardName" />
      </div>
      <div class="form-group col">
        <label>Customer Code</label>
        <input type="text" class="form-control form-control-sm" placeholder="Customer Code" disabled v-model="this.order.cardCode" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Document Total</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.docTotal" />
      </div>
      <div class="form-group col">
        <label>Create Date</label>
        <input type="text" class="form-control form-control-sm" placeholder="Create Date" disabled :value="this.currentDate" />
      </div>
      <div class="form-group col">
        <label>Created By</label>
        <input type="text" class="form-control form-control-sm" placeholder="Created By" disabled :value="this.$store.state.user.email" />
      </div>
    </div>

    <div class="row">
      <div class="col-8">
        <h3>Order Details</h3>
        <hr />
        <div class="table-responsive-sm">
          <table class="table table-bordered table-hover table-striped table-sm">
            <thead class="thead-dark">
              <tr>
                <th scope="col">#</th>
                <th scope="col">Item No.</th>
                <th scope="col">Item Description</th>
                <th scope="col">Serial Number</th>
                <th scope="col">Quantity</th>
                <th scope="col">Price</th>
                <th scope="col">Discount</th>
              </tr>
            </thead>
            <tbody id="orderedItems">
              <tr v-for="(item, index) in this.order.items" :key="index">
                <th scope="row">
                  {{ index + 1 }}
                </th>
                <td>
                  {{ item.itemCode }}
                </td>
                <td>
                  {{ item.itemName }}
                </td>
                <td>
                  {{ item.serialNumber }}
                </td>
                <td>
                  {{ item.quantity }}
                </td>
                <td>
                  {{ item.price.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
                </td>
                <td>
                  {{ item.discount.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div class="col-4">
        <h3>Payment Details</h3>
        <hr />
        <div class="table-responsive-sm">
          <table class="table table-bordered table-hover table-striped table-sm">
            <thead class="thead-dark">
              <tr>
                <th scope="col">#</th>
                <th scope="col">Payment Method</th>
                <th scope="col">Payment Details</th>
                <th scope="col">Amount Paid</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(payment, index) in this.order.payments" :key="index">
                <th scope="row">
                  {{ index + 1 }}
                </th>
                <td>
                  {{ payment.paymentType }}
                </td>
                <td>
                  {{ transformPayment(payment) }}
                </td>
                <td>
                  {{ payment.amount.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div class="table-responsive">
      <table class="table table-fixed table-bordered table-hover table-sm">
        <caption>
          <h3 style="display:inline-flex;" class="float-left">Replacement Details:</h3>
          <button class="btn btn-primary btn-sm mr-2 float-right" data-toggle="modal" data-target="#paymentModal">Payment</button>
        </caption>
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item No.</th>
            <th scope="col">Item Description</th>
            <th scope="col">Serial Number</th>
            <th scope="col">Quantity</th>
            <th scope="col">Price</th>
            <th scope="col">Discount ₦</th>
            <th scope="col">Total</th>
            <th scope="col"></th>
          </tr>
        </thead>
        <tbody id="LineItems">
          <tr v-for="(item, i) in items" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>{{ item.itemcode }}</td>
            <td>
              <!-- <v-select label="itemName" code="itemCode" @input="(val) => itemSelected(val, i)" v-model="item.itemName" :options="inventory" :clearable="false" placeholder="Kindly select Product"></v-select> -->
              <v-select @search="fetchOptions" :filterable="false" label="itemname" :options="inventory" :clearable="false" @input="(val) => itemSelected(val, i)" >
                <template slot="no-options">
                  type to search for Product
                </template>
                <template slot="option" slot-scope="option">
                  <div class="d-center">
                    {{ option.itemname }}
                  </div>
                </template>
                <template slot="selected-option" slot-scope="option">
                  <div class="selected d-center">
                    {{ option.itemname }}
                  </div>
                </template>                
              </v-select>
            </td>
            <td>{{ item.serialnumber }}</td>
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
              <button :id="'del' + i" class="btn btn-danger btn-sm mr-2 float-right" @click="deleteItemRow(i)">Remove Line</button>
            </td>
          </tr>
        </tbody>
        <tfoot id="ItemsFooter">
          <tr>
            <td colspan="8" class="text-right font-weight-bold">
              Subtotal:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{subTotal}}
            </td>
          </tr>
          <tr v-show="canVat">
            <td colspan="8" class="text-right font-weight-bold">7.5% VAT:</td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{vatAmount}}
            </td>
          </tr>
          <tr>
            <td colspan="8" class="text-right font-weight-bold">
              Amount Paid:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{amtPaid}}
            </td>
          </tr>
          <tr>
            <td colspan="8" class="text-right font-weight-bold">
              Grand Total:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{grandTotal}}
            </td>
          </tr>
        </tfoot>
      </table>
    </div>
    
    <button type="button" class="btn btn-primary btn-sm float-right" @click="SaveOrder" :disabled="isDisabled">
      Save
    </button>

    <!-- Modal -->
    <div class="modal fade" id="returnsModal" data-backdrop="static" tabindex="-1" aria-labelledby="returnsModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header text-center bg-dark text-white">
            <h5 class="modal-title" id="returnsModalLabel">
              Set Returns Criteria
            </h5>
            <button type="button" class="close text-white" @click="dismiss" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <div class="container-fluid">
              <div class="form-row mb-3">
                <div class="form-group col">
                  <label for="fromWHS">Purchased From</label>
                  <select id="fromWHS" class="form-control form-control-sm" placeholder="Kindly select store item was purchased from" v-model="StoreID">
                    <option value="" selected>Select store item was purchased from</option>
                    <option :key="store.name" :value="store.code" v-for="store in stores">
                      {{ store.name }}
                    </option>
                  </select>
                </div>
                <div class="form-group col">
                  <label class="form-check-label" for="orderid">Order ID</label>
                  <input id="orderid" class="form-control form-control-sm mt-2" type="text" v-model="OrderID" />
                </div>
              </div>
              <hr />
              <div class="form-row mb-3">
                <div class="form-group col">
                  <label class="form-check-label" for="serialnumber">
                    Serial ID
                  </label>
                  <input id="serialnumber" class="form-control form-control-sm" type="text" v-model="SerialNumber" />
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary btn-sm mr-2" @click="GetOrder">
              Get Order
            </button>
            <button type="button" class="btn btn-secondary btn-sm" @click="dismissModal">
              Close
            </button>
          </div>
        </div>
      </div>
    </div>

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
                  <input type="email" v-model="email" class="form-control form-control-sm" placeholder="Email address" required autofocus />
                </div>

                <div class="form-group">
                  <label for="inputPassword">Password</label>
                  <input type="password" v-model="password" class="form-control form-control-sm" placeholder="Password" required />
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

    <div class="modal fade" id="serialsModal" data-backdrop="static" tabindex="-1" aria-labelledby="serialsModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header text-center bg-dark text-white">
            <h5 class="modal-title" id="serialsModalLabel">Select Serial Number.</h5>
            <button type="button" class="close text-white" @click="dismiss" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <div class="row">
              <div v-for="(val, key) in serialnumbers" :key="key" class="col-3">
                <b>{{key + 1}}</b>. &nbsp;&nbsp;
                <label class="form-check-label" :for="'btnradio_' + val">
                <input type="radio" class="btn-check" name="serialnumber" v-model="serial" :id="'btnradio_' + val" :value="val" autocomplete="off">
                {{val}}
              </label>
              </div>
            </div>
            <hr />
            <span>Selected Serial Number: <strong>{{serial}}</strong></span>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary btn-sm mr-2" @click="setSerial">Set Serial Number</button>
            <button type="button" class="btn btn-secondary btn-sm" @click="dismiss">Close</button>
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
                      <th scope="col">Payment Details</th>
                      <th scope="col">Amount Paid</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody id="paymentLines">
                    <tr v-for="(payment, i) in payments" :key="'prow' + i" :id="'prow' + i">
                      <th scope="row">{{ i + 1 }}</th>
                      <td>
                        <select class="form-control form-control-sm" @change="paymentMethod($event, i)" v-model="payment.method">
                          <option value="null" selected>Select Payment Method</option>
                          <option>Cash</option>
                          <option>POS</option>
                          <option>Cheque</option>
                          <option>Bank Transfer</option>
                        </select>
                      </td>
                      <td>
                        <select disabled="true" :id="'pdet' + i" class="form-control form-control-sm" v-model="payment.bank" @change="setDetail($event, i)">
                          <option value="null" selected>Select Inflow Details</option>
                          <optgroup :key="'k_' + name" :id="'k_' + name" v-for="(group, name) in paymentDetails" :label="name">
                            <option :key="option.id" v-for="option in group" :value="option.code">
                              {{ option.name }}
                            </option>
                          </optgroup>
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
import _ from 'lodash';
import $ from "jquery";
import moment from "moment";
import "vue-select/dist/vue-select.css";

export default {
  data() {
    return {
      stores: [],
      StoreID: null,
      OrderID: null,
      SerialNumber: null,

      cardname: null,
      cardcode: null,

      // Login
      email: "",
      password: "",
      discountby: 0,
      
      canVat: false,
      isAdmin: false,
      
      items: [],
      order: {},
      item: null,
      serial: null,
      payments: [],
      inventory: [],
      currentIndex: 0,
      amtPaid: 0.00,
      isDisabled: true,
      selecteditem: '',
      ordercolumns: [],
      subTotal: 0.00,
      serialnumbers: [],
      vatAmount: 0.00,
      grandTotal: 0.00,
      balanceDue: 0.00,
      paymentcolumns: [],
      paymentDetails: {},
      currentDate: moment().format("Do of MMMM YYYY"),
      dateColumns:['created_at','updated_at', 'deleted_at'],
    };
  },
  async mounted() {
    // Determine the state of the Discount element
    if (this.$store.state.isAdmin)
    {
      this.isAdmin = true;
    }

    // Check if the store is allowed to calculate VAT
    if (this.$store.state.userStore.vat)
    {
      this.canVat = true;
    }

    // // Get all inventory
    // window.backend.GetProducts().then((inventory) => {
    //   this.inventory = inventory;
    //   this.addItemRow();
    // }, (err) => {
    //   this.$toast.error("Error! " + err);
    // });

    // Get all stores
    window.backend.GetStores().then((stores) => {
      this.stores = stores;
    }, (err) => {
      this.$toast.error("Error! " + err);
    });

    // Get all Payment Details
    window.backend.GetPaymentDetails().then((PaymentDetails) => {
      this.paymentDetails = {};
      this.paymentDetails["pos"] = PaymentDetails[0]["pos"];
      this.paymentDetails["cheque"] = PaymentDetails[0]["cheque"];
      this.paymentDetails["banktransfer"] = PaymentDetails[0]["banktransfer"];
    }, (err) => {
      this.$toast.error("Error! " + err);
    });

    // Show the modal
    $("#returnsModal").modal("show");
  },
  methods: {
    dismissModal() {
      // Hide the modal
      $("#returnsModal").modal("hide");

      // Hide the information and redirect to the Transfer Request list.
      this.$toast.info("Info! You can't continue without setting your role.");
      this.$router.push({ name: "orderlist" });
    },
    GetOrder() {
      let runningTotal = 0.0, 
        path = this.$store.state.userStore.orders;

      if (this.StoreID !== null && this.OrderID !== null) {
        path += `/${this.OrderID}?storeId=${this.StoreID}`;
      } else if (this.SerialNumber !== null) {
        path += `/Orders/SerialNumber?serialNumber=${this.SerialNumber}`;
      }

      var requestOptions = {
        method: 'GET',
        redirect: 'follow'
      };

      fetch(path, requestOptions)
        .then(response => response.json())
        .then(result => {
          if (result !== null && result !== undefined) {
            result.payments.forEach((payment) => {
              // Add the payment to the payments array
              this.payments.push({
                amount: payment.amount,
                method: payment.paymentType,
                name: payment.paymentDetails,
                id: this.payments.length + 1,
              });

              // Loop through the array and perform all needed calculations
              runningTotal += parseFloat(payment.amount.toString().replace("₦", ""));
        
              // Calculate footer details
              this.amtPaid = parseFloat(runningTotal).toFixed(2);
              this.grandTotal = parseFloat(runningTotal).toFixed(2);
            });
            this.addRow();
            this.order = result;
          }
        })
        .catch(error => { 
          this.$toast.error("Error! Unable to retrieve data for store.  \n" + error);
          this.$router.push({name: 'orderlist'});
        });
      
      // Hide the modal
      $("#returnsModal").modal("hide");
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
        price: '₦0.00',
        total: '₦0.00',
        serialnumber: '',
        discount: '₦0.00',
      });

      const el = document.getElementById("ItemsFooter");
      if (el) {
        el.scrollIntoView();
      }
    },
    async deleteItemRow(index) {
      this.$delete(this.items, index);
      
      if (this.items.length === 0 || this.items.length === index) {
        await this.addItemRow(index);
      }
    },
    async setSerial() {
      if (this.serial === null || this.serial === " ") {
        this.$toast.error("Error! Product Serial Number is required in order to proceed.");
        return;
      }

      this.items[this.currentIndex].serialnumber = this.serial; 
      $('#serialsModal').modal('hide');
    },
    dismiss() {
      if (this.serial === null || this.serial === " ") {
        this.$toast.error("Error! Product Serial Number is required in order to proceed.");
        return;
      }
    },
    async populateRow(rowIndex) {
      if (this.item === undefined) {
        return;
      }
      this.currentIndex = rowIndex;

      // Set the table data
      if (this.items[rowIndex].itemcode === "") {
        this.items[rowIndex].quantity = 1;
        this.items[rowIndex].id = rowIndex;
        this.items[rowIndex].discount = '₦0.00';
        this.items[rowIndex].itemcode = this.item.itemcode;
        this.items[rowIndex].itemname = this.item.itemname;
        this.items[rowIndex].price = `₦${parseFloat(this.item.price).toFixed(2)}`;
      } else {
        this.items[rowIndex].itemcode = this.item.itemcode;
        this.items[rowIndex].itemname = this.item.itemname;
        this.items[rowIndex].price = `₦${parseFloat(this.item.price).toFixed(2)}`;
      }

      if (this.item.serialnumbers !== "[]") {
        // Prompt for serial number of item.
        this.serialnumbers = this.item.serialnumbers.substring(1, this.item.serialnumbers.length-1).split(" ");
        $('#serialsModal').modal('show');
      } else {
        this.items[rowIndex].serialnumber = "";
      }

      // Calculate the totals [invoice subtotal, grand total, vat]
      await this.totals();
    },
    /**
     * Triggered when the search text changes.
     *
     * @param search  {String}    Current search text
     * @param loading {Function}	Toggle loading class
     */
    async fetchOptions (search, loading) {
      if (search.length >= 3) {
        loading(true);
        this.search(loading, search, this);
      }
    },
    search: _.debounce((loading, search, vm) => {
      // Get all inventory
      window.backend.GetProduct(search).then((inventory) => {
        if (inventory !== null) {
          vm.inventory = inventory;
        }
        loading(false);
      }, (err) => {
        this.$toast.error("Error! " + err);
        loading(false);
      });      
    }, 350),    
    async itemSelected(val, rowIndex) {
      this.item = val;
      await this.populateRow(rowIndex);
      await this.addItemRow(rowIndex);
    },
    async myFilter(option, label, search, rowIndex) {
      this.item = await this.inventory.filter((item) => {
          return (
            item.itemcode.toLowerCase().indexOf(search.toLowerCase()) !== -1 ||
            item.itemname.toLowerCase().indexOf(search.toLowerCase()) !== -1 ||
            item.serialnumbers.toLowerCase().indexOf(search.toLowerCase()) !== -1
          );
        })[0];

        await this.populateRow(rowIndex);
        await this.addItemRow(rowIndex);
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
          this.discountby = result.id;

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
      if (event.target.innerText === '₦0.00') {
        event.target.innerText = '';
      }
      
      // Allow for only one decimal point "." and interger characters
      if (event.which == '110' || event.which == '46') {
        if (event.target.innerText.includes(".")) {
          event.preventDefault();
        }
      } else if (isNaN(String.fromCharCode(event.which)) || event.which == '13' || event.which == '32') {
        // Reject all keypress characters that are not numeric
        event.preventDefault();
      }
    },
    async applyDiscount(rowIndex) {
      if (this.items.serialnumbers !== "[]" && this.items[rowIndex].serialnumber === "") {
        this.$toast.error("Error! Product Serial Number is required in order to proceed.");
        $('#serialsModal').modal('show');
        return;
      }

      // Perform a quick clean up
      if (this.items[rowIndex].price === '₦0.00') {
        event.target.innerText = '₦0.00';
        this.$toast.error("Error! Discount cannot be applied on ₦0.00. \nKindly select a product with a valid Price.");
        return;
      }

      if (event.target.innerText === '' || event.target.innerText < 0 || event.target.innerText > parseFloat(this.items[rowIndex].total.replace("₦", ""))) {
        event.target.innerText = '₦0.00';
        this.$toast.error(`Error! Discount must be between ₦1.00 and ${this.items[rowIndex].price}`);
        return;
      }

      // Reset the Discount element to its inital state.
      if (!this.$store.state.isAdmin)
      {
        this.isAdmin = false;
      }

      // Store the data into the items array
      this.items[rowIndex].discount = '₦' + parseFloat(event.target.innerText.replace("₦", "")).toFixed(2);

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
        // Calculate the discount. (quantity * price) - discount value
        let quantity = element.quantity,
          price = parseFloat(element.price.replace("₦", "")),
          discount = parseFloat(element.discount.replace("₦", "")),
          currentTotal = (quantity * price) - discount;

        element.total = "₦" + currentTotal;
        runningTotal += parseFloat(currentTotal);
      });
      
      // Calculate footer details
      this.subTotal = parseFloat(runningTotal).toFixed(2);
      if (this.canVat === true) {
        this.vatAmount = parseFloat((7.5 / 100) * runningTotal).toFixed(2);
        this.grandTotal = parseFloat((7.5 / 100) * runningTotal + runningTotal).toFixed(2);
      } else {
        this.grandTotal = parseFloat(runningTotal) + parseFloat(this.amtPaid);
        this.grandTotal = parseFloat(this.grandTotal).toFixed(2);
      }
      this.balanceDue = parseFloat(this.grandTotal - this.amtPaid).toFixed(2);
    },
    // Payment Section
    transformPayment(payment) {
      let pay = this.paymentDetails[payment.paymentType.toLowerCase()];

      return pay.filter((data) => {
        return (data.code === payment.paymentDetails);
      })[0].name;
    },
    async setDetail(event, index) {
      this.payments[index].name = event.target.value;
    },
    async paymentMethod(event, index) {
      let ctrl = event.target.value,
        amt = event.target.parentElement.parentElement.cells[3].childNodes[0],
        bank = event.target.parentElement.parentElement.cells[2].childNodes[0];

      bank.disabled = false;
      amt.disabled = false;
      
      switch (ctrl.toLowerCase()) {
        case "pos":
          $(`#${bank.id} optgroup#k_pos`).removeAttr('disabled');
          $(`#${bank.id} optgroup#k_cheque`).attr('disabled','disabled');
          $(`#${bank.id} optgroup#k_banktransfer`).attr('disabled','disabled');
          break;

        case "bank transfer":
          $(`#${bank.id} optgroup#k_pos`).attr('disabled','disabled');
          $(`#${bank.id} optgroup#k_cheque`).attr('disabled','disabled');
          $(`#${bank.id} optgroup#k_banktransfer`).removeAttr('disabled');
          break;

        case "cheque":
          $(`#${bank.id} optgroup#k_cheque`).removeAttr('disabled');
          $(`#${bank.id} optgroup#k_pos`).attr('disabled','disabled');
          $(`#${bank.id} optgroup#k_banktransfer`).attr('disabled','disabled');
          break;
      
        default:
          bank.disabled = true;
          $(`#${bank.id} optgroup#k_pos`).attr('disabled','disabled');
          $(`#${bank.id} optgroup#k_cheque`).attr('disabled','disabled');
          $(`#${bank.id} optgroup#k_banktransfer`).attr('disabled','disabled');
          break;
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
      
      if (this.payments.length === 0 || this.items.length === index) {
        this.addRow(index);
      }
    },
    PaymentTotal(index, event) {
      let runningTotal = 0.0,
        amt = event.target.value,
        method = event.target.parentElement.parentElement.cells[1].childNodes[0].value,
        details = event.target.parentElement.parentElement.cells[2].childNodes[0].value;

      if (details === null) {
        this.$toast.error("Error! Payment Details is required in order to proceed.");
        return;
      }

      if (amt === '0.00' || amt === 0) {
        this.$toast.error("Error! Invalid Amount inputted.");
        return;
      }

      this.payments[index].amount = amt;
      this.payments[index].name = details;
      this.payments[index].method = method;

       // Loop through the array and perform all needed calculations
      this.payments.forEach(element => {
        runningTotal += parseFloat(element.amount.toString().replace("₦", ""));
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
        if (element.name === "" && (element.method.toLowerCase() === "pos" || element.method.toLowerCase() === "bank transfer")) {
          isvalid = false;
        }
      });

      if (isvalid === false) {
        this.$toast.error("Error! Bank Details is required in order to proceed.");
        return;
      }

      this.isDisabled = false;
      $('#paymentModal').modal('hide');
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
        let serial = "";
        if (element.serialnumber === "") {
          serial = "null";
        } else {
          serial = element.serialnumber;
        }

        let orderItem = {
          orderid: null,
          serialnumber: serial,
          itemname: element.itemname,
          itemcode: element.itemcode,
          quantity: element.quantity,
          price: element.price.replace("₦", ""),
          discount: element.discount.replace("₦", ""),
        };

        // Add the items to the items array
        items.push(orderItem);
      });
      this.order.items.forEach(element => {
        let serial = "";
        if (element.serialNumber === "") {
          serial = "null";
        } else {
          serial = element.serialNumber;
        }

        let orderItem = {
          orderid: null,
          serialnumber: serial,
          itemname: element.itemName,
          itemcode: element.itemCode,
          quantity: "-" + element.quantity,
          price: element.price.toString().replace("₦", ""),
          discount: element.discount.toString().replace("₦", ""),
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
      let order = {
        id: 0,
        docnum: this.order.docNum,
        docentry: this.order.docEntry,
        vatsum: 7.5,
        synced: false,
        canceled: false,
        cardname: this.order.cardName,
        cardcode: this.order.cardCode,
        created_by: this.$store.state.user.id,
      };
      order.items = items;
      order.returned = true;
      order.payments = payments;
      order.doctotal = this.grandTotal;
      order.discountapprovedby = parseInt(this.discountby);
      
      window.backend.NewOrder(order).then(() => {
        this.$toast.success("Success! Order has been successfully saved.");
        this.$router.push({name: 'orderlist'});
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
</script>