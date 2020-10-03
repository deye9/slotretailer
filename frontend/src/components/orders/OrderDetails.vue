<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Viewing Order {{this.order.id}}</h3>
      </div>
      <div class="col-4">
        <router-link to="/orders/" class="btn btn-info float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label>Order Number</label>
        <input
          type="text"
          class="form-control"
          disabled
          v-model="this.order.id"
        />
      </div>
      <div class="form-group col">
        <label>Customer</label>
        <input
          type="text"
          class="form-control"
          disabled
          v-model="this.order.cardname"
        />
      </div>
      <div class="form-group col">
        <label for="docnum">SAP Document Number</label>
        <input
          type="text"
          class="form-control"
          disabled
          v-model="this.order.docnum"
        />
      </div>
      <div class="form-group col">
        <label for="docnum">Created By</label>
        <br />
        <span class="btn btn-info" disabled>
          {{ this.user.firstname + " " + this.user.lastname }}
        </span>
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Synced</label>
        <input
          type="text"
          class="form-control"
          disabled
          v-model="this.order.synced"
        />
      </div>
      <div class="form-group col">
        <label>Canceled</label>
        <input
          type="text"
          class="form-control"
          disabled
          v-model="this.order.canceled"
        />
      </div>
      <div class="form-group col">
        <label for="docnum">VAT Amount</label>
        <input
          type="text"
          class="form-control"
          disabled
          v-model="this.order.vatsum"
        />
      </div>
      <div class="form-group col">
        <label for="docnum">Created On</label>
        <br />
        <span class="btn btn-info" disabled>
          {{ this.user.created_at.String }}
        </span>
      </div>
    </div>

    <h3>Order Details</h3>
    <hr />
    <div style="max-height: 250px; overflow: scroll">
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
                {{ item.discount }}
              </td>
              <td v-if="item.discount === 0">
                ₦{{ parseFloat(item.quantity * item.price).toFixed(2)}}
              </td>
              <td v-else>
                ₦{{
                  parseFloat(
                    item.quantity * item.price -
                      (item.price - item.price * (item.discount / 100))
                  ).toFixed(2)
                }}
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr>
              <td colspan="6" class="text-right font-weight-bold">
                Invoice Subtotal:
              </td>
              <td class="font-weight-bold bg-primary text-white">
                ₦{{this.subtotal}}
              </td>
            </tr>
            <tr>
              <td colspan="6" class="text-right font-weight-bold">7.5% VAT:</td>
              <td class="font-weight-bold bg-primary text-white">
                ₦{{this.vatAmount}}
              </td>
            </tr>          
            <tr>
              <td colspan="6" class="text-right font-weight-bold">
                Grand Total:
              </td>
              <td id="grandTotal" class="font-weight-bold bg-primary text-white">
                ₦{{ parseFloat(this.order.doctotal).toFixed(2) }}
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
              <th scope="col">Payment Type</th>
              <th scope="col">Amount Paid</th>
              <th scope="col">Canceled</th>
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
                ₦{{ payment.amount }}
              </td>
              <td>
                {{ payment.canceled }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
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
      subtotal: 0.0,
      vatAmount: 0.0,
    };
  },
  created() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.GetOrder(parseInt(this.id)).then((order) => {
        this.order = order;

        this.order.items.forEach(item => {
          if (item.discount === 0) {
            this.subtotal = parseFloat(this.subtotal) + parseFloat(item.quantity) * parseFloat(item.price).toFixed(2);
          } else {
            this.subtotal = parseFloat(parseFloat(item.quantity) * parseFloat(item.price) - (parseFloat(item.price) - parseFloat(item.price) * (parseFloat(item.discount) / 100))).toFixed(2);
          }

          this.vatAmount = parseFloat((7.5 / 100) * this.subtotal).toFixed(2);
        });

        window.backend.PaymentOnOrder(parseInt(this.id)).then((payment) => {
            if (JSON.stringify(payment) !== "{}") {
              this.payment = payment;
            }
          },
          (err) => {
            this.$toast.error("Error! " + err);
          }
        );

        window.backend.GetUser(parseInt(order.created_by)).then((user) => {
            if (JSON.stringify(user) !== "{}") {
              this.user = user;
            }
          },
          (err) => {
            this.$toast.error("Error! " + err);
          }
        );
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );
  }
};
</script>