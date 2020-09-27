<template>
  <section>
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
      <div class="form-group col">
        <label>Payment Type</label>
        <select class="form-control" @change="PaymentMethod($event)">
          <option>Cash</option>
          <option>POS</option>
          <option>EasyBuy</option>
          <option>Voucher</option>
          <option>Bank Transfer</option>
       </select>
      </div>
      <div class="form-group col">
        <label>Payment Details</label>
        <input type="text" class="form-control" :disabled="isValid" placeholder="Payment Details" />
      </div>
    </div>

    <h3 style="margin-bottom: 50px;">
      <span class="float-left">Items</span>
      <b-button v-b-modal.modal-center variant="info" class="float-right">Add Row</b-button>
    </h3>

    <hr />
    <div class="table-responsive" style="max-height:350px; overflow: scroll;">
      <table class="table table-bordered table-hover table-striped">
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item Code</th>
            <th scope="col">Item Name</th>
            <th scope="col">Quantity (Min of 1)</th>
            <th scope="col">Price</th>
            <th scope="col">Discount (%)</th>
            <th scope="col">Total</th>
          </tr>
        </thead>
        <tbody id="orderedItems"></tbody>
        <tfoot>
          <tr>
            <td colspan="6" class="text-right font-weight-bold">Grand Total:</td>
            <td id="grandTotal" class="font-weight-bold bg-primary text-white">₦0.00</td>
          </tr>
        </tfoot>
      </table>
    </div>

    <button type="button" class="btn btn-info float-right" style="margin-left: 2px;" @click="SaveOrder">Save Order</button>
    <button type="button" class="btn btn-dark float-right" style="margin-right: 2px;" @click="DraftOrder">Save as Draft</button>

    <!-- Modal -->
    <b-modal id="modal-center" centered title="Add New Item Row" :header-bg-variant="headerBgVariant" @hidden="resetModal"
      :header-text-variant="headerTextVariant" :body-bg-variant="bodyBgVariant" @ok="handleOk" @show="resetModal"
      :body-text-variant="bodyTextVariant" :footer-bg-variant="footerBgVariant" :footer-text-variant="footerTextVariant">
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
            <input id="price" name="price" type="text" class="form-control" placeholder="₦0.00" disabled value="₦0.00"/>
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
      item: null,
      isValid: true,
      customers: [],
      inventory: [],
      customer: null,

      // Modal
      headerBgVariant: 'dark',
      headerTextVariant: 'light',
      bodyBgVariant: 'light',
      bodyTextVariant: 'dark',
      footerBgVariant: 'secondary',
      footerTextVariant: 'dark',
    };
  },
  async created() {
    // Get all customers
    window.backend.GetCustomers().then((customers) => {
        this.customers = customers;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );

    // Get all inventory
    window.backend.GetProducts().then((inventory) => {
        this.inventory = inventory;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );

    // await this.AddRow(0);
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
            return customer.phone.toLowerCase() === search.toLowerCase() 
            || customer.phone1.toLowerCase() === search.toLowerCase()
            || customer.email.toLowerCase() === search.toLowerCase()
            || customer.cardname.toLowerCase() === search.toLowerCase()
            || customer.cardcode.toLowerCase() === search.toLowerCase();
          })[0];
        }
        loading(false);
        return;
    },
    fetchItem(search, loading) {
        loading(true);
        if (search.length >= 3) {
          this.item = this.inventory.filter((item) => {
            return item.itemcode.toLowerCase() === search.toLowerCase() 
            || item.itemname.toLowerCase() === search.toLowerCase()
            || item.codebars.toLowerCase() === search.toLowerCase()
            || item.serialnumber.toLowerCase() === search.toLowerCase();
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
    async handleOk(bvModalEvt) {
      if (this.item === null) {
         this.$toast.error("Error! No Product selected.");

       // Prevent modal from closing
        bvModalEvt.preventDefault();

        return;
      }

      // Get a reference to the Table body
      let total = 0, 
        discount = document.getElementById('discount').value,
        quantity = document.getElementById('quantity').value,
        tableRef = document.getElementById("orderedItems"),
        rowCnt = tableRef.rows.length + 1,
        // Insert a row in the table at the last row
        row = tableRef.insertRow(),
        price = document.getElementById('price').value.replace("₦", "");

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
          discountVal = price - (price * numVal);
        total = `₦${((quantity * price) - discountVal).toFixed(2)}`;
      }
      row.cells[6].innerHTML = total;

      // Calculate the Grand Total
      await this.TableTotal(tableRef);

      // Reset the modal controls /values
      this.resetModal();
    },
    async TableTotal(tableRef) {
      let val = 0.00;

      for (var i = 0; i < tableRef.rows.length; i++) {
        let data = tableRef.rows[i].cells[6].innerHTML.replace("₦", "");
        val += parseFloat(data);
      }
      document.getElementById("grandTotal").innerHTML = `₦${val}`;
    },
    PaymentMethod(event) {
      switch (event.target.value.toLowerCase()) {      
        case "pos":
            this.isValid = false;
            break;
        
        case "easybuy":
            this.isValid = false;
            break;
        
        case "voucher":
            this.isValid = false;
            break;
        
        case "bank transfer":
            this.isValid = false;
            break;
      
        default:
          this.isValid = true;
          break;
      }
    },
    DraftOrder() {

    },
    SaveOrder() {
        var r = confirm("Are you sure you want to create this Sales Order!");
        if (r == false) {
          return;
        }

        alert('Saved');
    //   let val = 0,
    //     items = [], 
    //     customer = document.getElementById("customers"),
    //     trs = document.getElementById("orderedItems").getElementsByTagName("tr");
        
    //   if (customer.options[customer.selectedIndex].text.toLowerCase() === "select customer") {
    //     this.$toast.error("Error! You are yet to associate a customer to this order.");
    //     return
    //   }

    //   // Build out the header.
    //   this.order = {
    //     id: this.id,
    //     docentry: 0,
    //     vatsum: 0.75,
    //     synced: false,
    //     canceled: false,
    //     docnum: this.docnum,
    //     created_by: this.$store.state.user.id,
    //     cardname: customer.options[customer.selectedIndex].text,
    //     cardcode: customer.options[customer.selectedIndex].value,
    //   }

    //   // Build out the Line items.
    //   for (var i = 0; i < trs.length; i++) {
    //     let item = trs[i].cells[1].childNodes[0],
    //       qty = trs[i].cells[2].childNodes[0],
    //       discount = trs[i].cells[4].childNodes[0],
    //       data = trs[i].cells[5].childNodes[0].value.replace("₦", ""),
    //       product = this.inventory.filter((products) => {
    //         return products.itemcode.toLowerCase() === item.options[item.selectedIndex].value.toLowerCase();
    //       })[0];
        
    //     if (product != undefined) {
    //       let orderItem = {
    //         orderid: null,
    //         quantity: qty.value,
    //         price: product.price,
    //         discount: discount.value,
    //         itemname: item.options[item.selectedIndex].text,
    //         itemcode: item.options[item.selectedIndex].value,
    //       };
          
    //       // Get the total value of the order
    //       val += parseFloat(data);

    //       // Add the ordered items to the items array
    //       items.push(orderItem);
    //     }
    //   }

    //   if (items.length === 0) {
    //     this.$store.state.notify.category = "error";
    //     this.$store.state.notify.message = "Error! You are not permitted to create an empty Sales Order.";
    //     return
    //   }

    //   this.order.items = items;
    //   this.order.doctotal = `${val.toFixed(2)}`;

    //   window.backend.NewOrder(this.order).then(() => {
    //     this.$router.push("/orders/");
    //   },
    //   (err) => {
    //     this.$toast.error("Error! " + err);
    //   }
    // );
    },
  },
};

</script>

