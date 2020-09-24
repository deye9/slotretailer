<template>
  <section>
    <h1>New Order</h1>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label for="cardcode">Customer</label>
        <select name="customers" id="customers" class="form-control">
          <option selected disabled>Select Customer</option>
          <option v-for="(customer, index) in customers" :key="index" :value="customer.cardcode">{{ customer.cardname }}</option>
        </select>
      </div>
      <div class="form-group col">
        <label for="docnum">SAP Document Number</label>
        <input
          type="text"
          class="form-control"
          placeholder="SAP Document Number"
          disabled
          v-model="docnum"
          required
        />
      </div>
      <div class="form-group col">
        <label for="docnum">Created By</label>
        <br />
        <span
          class="btn btn-info"
          disabled
        >{{this.$store.state.user.firstname + " " + this.$store.state.user.lastname}}</span>
      </div>
    </div>

    <h3>Items</h3>
    <hr />
    <div class="table-responsive">
      <table class="table table-bordered table-hover table-striped">
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item</th>
            <th scope="col">Quantity (Min of 1)</th>
            <th scope="col">Price</th>
            <th scope="col">Discount (%)</th>
            <th scope="col">Total</th>
            <th scope="col"></th>
          </tr>
        </thead>
        <tbody id="orderedItems">
          <tr v-for="(row, index) in rows" :key="index">
            <th scope="row">
              {{ index + 1 }}
            </th>
            <td>
              <select class="form-control" @change="AddRow">
                <option selected disabled>Select Product</option>
                <option v-for="(inventory, index) in inventory" :key="index" :value="inventory.itemcode">
                  {{ inventory.itemname }}
                </option>
              </select>
            </td>
            <td>
              <input type="number" class="form-control" placeholder="Quantity" @blur="UpdateRow" @change="UpdateRow" min="1" value="1" step="1" />
            </td>
            <td>
              <input type="text" class="form-control" placeholder="₦0.00" disabled value="₦0.00"/>
            </td>
            <td>
              <input type="number" class="form-control" placeholder="Discount" min="1" step="0.1" value="0" @blur="UpdateRow" @change="UpdateRow" />
            </td>
            <td>
              <input type="text" class="form-control" placeholder="₦0.00" disabled value="₦0.00" />
            </td>
            <td>
              <button class="btn btn-danger" title="Delete Row" @click="RemoveRow">D</button>
            </td>
          </tr>
        </tbody>
        <tfoot>
          <tr>
            <td colspan="5" class="text-right font-weight-bold">Grand Total:</td>
            <td colspan="2" id="grandTotal" class="font-weight-bold bg-primary text-white">₦0.00</td>
          </tr>
        </tfoot>
      </table>
    </div>

    <button type="button" class="btn btn-primary float-right" @click="SaveOrder">Save Order</button>
  </section>
</template>

<script>
export default {
  data() {
    return {
      rows: [],
      docnum: 0,
      order: {},
      customers: [],
      inventory: [],
    };
  },
  created() {
    // Get all customers
    window.backend.GetCustomers().then((customers) => {
        this.customers = customers;
      },
      (err) => {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! " + err;
      }
    );

    // Get all inventory
    window.backend.GetProducts().then((inventory) => {
        this.inventory = inventory;
      },
      (err) => {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! " + err;
      }
    );

    this.AddRow(0);
  },
  methods: {
    async AddRow(index) {
      // Add new row to the table.
      if (index === 0) {
        this.rows.splice(index + 1, 0, {});
        return;
      }

      this.rows.splice(index + 1, 0, {});
      await this.UpdateRow(index);
    },
    async UpdateRow(index) {
      // Get and write out the Price alongside the discount if any
      let TblRow = index.srcElement.parentElement.parentElement,
        item = TblRow.cells[1].childNodes[0],
        qty = TblRow.cells[2].childNodes[0],
        price = TblRow.cells[3].childNodes[0],
        discount = TblRow.cells[4].childNodes[0],
        total = TblRow.cells[5].childNodes[0];

      // Set the selected Itemcode and Itemname
      let // selText = item.options[item.selectedIndex].text,
        selValue = item.options[item.selectedIndex].value,
        product = this.inventory.filter((products) => {
            return products.itemcode.toLowerCase() === selValue.toLowerCase();
        })[0];

      // Set the needed data
      price.value = `₦${product.price}`;
      if (parseInt(discount.value) === 0) {
        total.value = `₦${Number(qty.value) * parseFloat(product.price)}`;
      } else {
        let numVal = Number(discount.value / 100),
          discountVal = product.price - (product.price * numVal);
        total.value = `₦${((qty.value * product.price) - discountVal).toFixed(2)}`;
      }

      await this.TableTotal();
    },
    async TableTotal() {
      let val = 0.00,
        trs = document.getElementById("orderedItems").getElementsByTagName("tr");

      for (var i = 0; i < trs.length; i++) {
        let data = trs[i].cells[5].childNodes[0].value.replace("₦", "");
        val += parseFloat(data);
      }
      document.getElementById("grandTotal").innerHTML = `₦${val.toFixed(2)}`;
    },
    RemoveRow(index){
      if (this.rows.length > 1) {
        this.rows.splice(index, 1);
      }
    },
    Registration() {
      this.order = {
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
      for (var attribute in this.order) {
        if (this.order[attribute] === "" || this.order[attribute] === null) {
          this.isValid = false;
          this.$store.state.notify.category = "error";
          this.$store.state.notify.message =
            "Error! " + attribute + " cannot be " + this.order[attribute];
          return;
        }
      }

      window.backend.NewOrder(this.order).then(() => {
          this.$router.push("/orders/");
        },
        (err) => {
          this.$store.state.notify.category = "error";
          this.$store.state.notify.message = "Error! " + err;
        }
      );
    },
    SaveOrder() {
      let val = 0,
        items = [], 
        customer = document.getElementById("customers"),
        trs = document.getElementById("orderedItems").getElementsByTagName("tr");
        
      if (customer.options[customer.selectedIndex].text.toLowerCase() === "select customer") {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! You are yet to associate a customer to this order.";
        return
      }

      // Build out the header.
      this.order = {
        id: this.id,
        docentry: 0,
        vatsum: 0.75,
        synced: false,
        canceled: false,
        docnum: this.docnum,
        created_by: this.$store.state.user.id,
        cardname: customer.options[customer.selectedIndex].text,
        cardcode: customer.options[customer.selectedIndex].value,
      }

      // Build out the Line items.
      for (var i = 0; i < trs.length; i++) {
        let item = trs[i].cells[1].childNodes[0],
          qty = trs[i].cells[2].childNodes[0],
          discount = trs[i].cells[4].childNodes[0],
          data = trs[i].cells[5].childNodes[0].value.replace("₦", ""),
          product = this.inventory.filter((products) => {
            return products.itemcode.toLowerCase() === item.options[item.selectedIndex].value.toLowerCase();
          })[0];
        
        if (product != undefined) {
          let orderItem = {
            orderid: null,
            quantity: qty.value,
            price: product.price,
            discount: discount.value,
            itemname: item.options[item.selectedIndex].text,
            itemcode: item.options[item.selectedIndex].value,
          };
          
          // Get the total value of the order
          val += parseFloat(data);

          // Add the ordered items to the items array
          items.push(orderItem);
        }
      }

      if (items.length === 0) {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! You are not permitted to create an empty Sales Order.";
        return
      }

      this.order.items = items;
      this.order.doctotal = `${val.toFixed(2)}`;

      window.backend.NewOrder(this.order).then(() => {
        this.$router.push("/orders/");
      },
      (err) => {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! " + err;
      }
    );
    },
  },
};

</script>