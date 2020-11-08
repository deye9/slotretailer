<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Viewing Order {{this.order.id}}</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'orderlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label>Order Number</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.id" />
      </div>
      <div class="form-group col">
        <label>Customer</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.cardname" />
      </div>
      <div class="form-group col">
        <label for="docnum">SAP Document Number</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.docnum" />
      </div>
      <div class="form-group col">
        <label for="docnum">Created By</label>
        <br />
        <span class="btn btn-info btn-sm" disabled v-if="this.user !== null">
          {{ this.user.firstname + " " + this.user.lastname }}
        </span>
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Synced</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.synced" />
      </div>
      <div class="form-group col">
        <label>Canceled</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.canceled" />
      </div>
      <div class="form-group col">
        <label for="docnum">Document Total</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.doctotal" />
      </div>
      <div class="form-group col">
        <label for="docnum">Created On</label>
        <br />
        <span class="btn btn-info btn-sm" disabled v-if="this.user !== null" >
          {{ this.user.created_at.Time }}
        </span>
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Comments</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.comment.String" v-if="this.order.comment !== undefined" />
      </div>
    </div>

    <h3>Order Details</h3>
    <hr />
    
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
        <tbody id="orderedItems">
          <tr v-for="(item, index) in this.order.items" :key="index">
            <th scope="row">
              {{ index + 1 }}
            </th>
            <td>
              {{ item.itemcode }}
            </td>
            <td>
              {{ item.itemname }}
            </td>
            <td>
              {{ item.quantity }}
            </td>
            <td>₦{{ item.price }}</td>
            <td>
              {{ item.discount }}%
            </td>
            <td>
              ₦{{ item.total }}
            </td>
          </tr>
        </tbody>
        <tfoot id="ItemsFooter">
          <tr>
            <td colspan="6" class="text-right font-weight-bold">
              Subtotal:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{subTotal}}
            </td>
          </tr>
          <tr>
            <td colspan="6" class="text-right font-weight-bold">7.5% VAT:</td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{vatAmount}}
            </td>
          </tr>
          <tr>
            <td colspan="6" class="text-right font-weight-bold">
              Grand Total:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              ₦{{grandTotal}}
            </td>
          </tr>
        </tfoot>
      </table>
    </div>

    <h3>Payment Details</h3>
    <hr />
    <div class="table-responsive-sm">
      <table class="table table-bordered table-hover table-striped table-sm">
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Payment Method</th>
            <th scope="col">Bank Used</th>
            <th scope="col">Canceled</th>
            <th scope="col">Amount Paid</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(payment, index) in this.payment" :key="index">
            <th scope="row">
              {{ index + 1 }}
            </th>
            <td>
              {{ payment.paymenttype }}
            </td>
            <td>
              {{ payment.paymentdetails }}
            </td>
            <td>
              {{ payment.canceled }}
            </td>
            <td>
              ₦{{ payment.amount }}
            </td>
          </tr>
        </tbody>
        <tfoot>
          <tr>
            <td colspan="4" class="text-right font-weight-bold">Amount Paid</td>
            <td class="font-weight-bold bg-primary text-white">₦{{grandTotal}}</td>
          </tr>
        </tfoot>
      </table>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      order: [],
      user: null,
      payment: null,
      subTotal: '0.00',
      vatAmount: '0.00',
      grandTotal: '0.00',
      balanceDue: '0.00',
    };
  },
  created() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.GetOrder(parseInt(this.id)).then((order) => {
        order.doctotal = `₦${order.doctotal}`;
        this.order = order;
        let runningTotal = 0.0;

        this.order.items.forEach(item => {
          // Loop through the array and perform all needed calculations
          let numVal1 = parseFloat(item.price),
            numVal2 = parseFloat(item.discount) / 100,
            currentTotal = parseFloat(Number(item.quantity) * (numVal1 - (numVal1 * numVal2))).toFixed(2);
          
          item.total = currentTotal;
          runningTotal += parseFloat(currentTotal);

          // Calculate footer details
          this.subTotal = parseFloat(runningTotal);
          this.vatAmount = parseFloat((7.5 / 100) * runningTotal).toFixed(2);
          this.grandTotal = parseFloat((7.5 / 100) * runningTotal + runningTotal).toFixed(2);
          this.balanceDue = parseFloat(this.grandTotal - this.amtPaid).toFixed(2);
        });
        
        window.backend.GetUser(parseInt(order.created_by)).then((user) => {
          if (JSON.stringify(user) !== "{}") {
            this.user = user;
          }
        },
        (err) => {
          this.$toast.error("Error! " + err);
        });
        
        window.backend.PaymentOnOrder(parseInt(this.id)).then((payment) => {
          if (JSON.stringify(payment) !== "{}") {
            this.payment = payment;
          }
        },
        (err) => {
          this.$toast.error("Error! " + err);
        });
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
  }
};
</script>