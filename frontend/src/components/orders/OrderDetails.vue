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
          {{ this.order.created_at.Time }}
        </span>
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Comments</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.order.comment.String" v-if="this.order.comment !== undefined" />
      </div>
      <div class="form-group col">
        <label>Discount Approved By</label>
        <br />
        <span class="btn btn-info btn-sm" disabled v-if="this.discountby !== null" >
          {{ this.discountby.firstname + " " + this.discountby.lastname }}
        </span>
      </div>
    </div>

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
            <th scope="col">Discount ₦</th>
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
              {{ item.serialnumber }}
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
            <td>
              {{ item.total.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
            </td>
          </tr>
        </tbody>
        <tfoot id="ItemsFooter">
          <tr>
            <td colspan="7" class="text-right font-weight-bold">
              Subtotal:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              {{ subTotal.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
            </td>
          </tr>
          <tr v-show="canVat">
            <td colspan="7" class="text-right font-weight-bold">7.5% VAT:</td>
            <td class="font-weight-bold bg-primary text-white">
              {{ vatAmount.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
            </td>
          </tr>
          <tr>
            <td colspan="7" class="text-right font-weight-bold">
              Grand Total:
            </td>
            <td class="font-weight-bold bg-primary text-white">
              {{ grandTotal.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
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
            <th scope="col">Payment Details</th>
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
              {{ payment.amount.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
            </td>
          </tr>
        </tbody>
        <tfoot>
          <tr>
            <td colspan="4" class="text-right font-weight-bold">Amount Paid</td>
            <td class="font-weight-bold bg-primary text-white">{{ grandTotal.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}</td>
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
      canVat: false,

      order: [],
      user: null,
      payment: null,
      discountby: null,
      subTotal: '0.00',
      vatAmount: '0.00',
      grandTotal: '0.00',
      balanceDue: '0.00',
    };
  },
  created() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    // Check if the store is allowed to calculate VAT
    if (this.$store.state.userStore.vat)
    {
      this.canVat = true;
    }

    window.backend.GetOrder(parseInt(this.id)).then((order) => {
        order.doctotal = order.doctotal.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' });
        this.order = order;
        let runningTotal = 0.0;

        this.order.items.forEach(item => {
          // Calculate the discount. (quantity * price) - discount value
          let quantity = item.quantity,
            price = parseFloat(item.price),
            discount = parseFloat(item.discount),
            currentTotal = (quantity * price) - discount;

          item.total = currentTotal;
          runningTotal += parseFloat(currentTotal);

          // Calculate footer details
          this.subTotal = parseFloat(runningTotal);
          if (this.canVat === true) {
            this.vatAmount = parseFloat((7.5 / 100) * runningTotal);
            this.grandTotal = parseFloat((7.5 / 100) * runningTotal + runningTotal);
          } else {
            this.grandTotal = parseFloat(runningTotal);
          }
          this.balanceDue = parseFloat(this.grandTotal - this.amtPaid);
        });
        
        window.backend.GetUser(parseInt(order.created_by)).then((user) => {
          if (JSON.stringify(user) !== "{}") {
            this.user = user;
          }
        }, (err) => {
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

        if (parseInt(order.discountapprovedby) === 0) {
          this.discountby = {
              firstname: "Super",
              lastname: "Admin",
            };
        } else {
          window.backend.GetUser(parseInt(order.discountapprovedby)).then((user) => {
            if (JSON.stringify(user) !== "{}") {
              this.discountby = user;
            }
          }, (err) => {
            this.$toast.error("Error! " + err);
          });
        }
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
  }
};
</script>