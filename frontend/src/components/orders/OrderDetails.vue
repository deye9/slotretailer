<template>
  <section>
    <h1>Order Details</h1>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label>Order Number</label>
        <input type="text" class="form-control" disabled v-model="this.order.id"/>
      </div>
      <div class="form-group col">
        <label>Customer</label>
        <input type="text" class="form-control" disabled v-model="this.order.cardname"/>
      </div>
      <div class="form-group col">
        <label for="docnum">SAP Document Number</label>
        <input type="text" class="form-control" disabled v-model="this.order.docnum" />
      </div>
      <div class="form-group col">
        <label for="docnum">Created By</label>
        <br />
        <span class="btn btn-info" disabled>
          {{this.user.firstname + " " + this.user.lastname}}
          </span>
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Synced</label>
        <input type="text" class="form-control" disabled v-model="this.order.synced"/>
      </div>
      <div class="form-group col">
        <label>Canceled</label>
        <input type="text" class="form-control" disabled v-model="this.order.canceled"/>
      </div>
      <div class="form-group col">
        <label for="docnum">VAT Amount</label>
        <input type="text" class="form-control" disabled v-model="this.order.vatsum" />
      </div>
      <div class="form-group col">
        <label for="docnum">Created On</label>
        <br />
        <span class="btn btn-info" disabled>
          {{this.user.created_at.String}}
          </span>
      </div>
    </div>
    <h3>Ordered Items</h3>
    <hr />

    <div class="table-responsive">
      <table class="table table-bordered table-hover table-striped">
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item</th>
            <th scope="col">Quantity</th>
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
              {{item.itemname}}
            </td>
            <td>
              {{item.quantity}}
            </td>
            <td>
              ₦{{item.price}}
            </td>
            <td>
              {{item.discount}}
            </td>
            <td v-if="item.discount === 0">
              ₦{{parseFloat(item.quantity * item.price).toFixed(2)}}
            </td>
            <td v-else>
              ₦{{
                parseFloat((item.quantity * item.price) - (item.price - (item.price * (item.discount / 100)))).toFixed(2)
              }}
            </td>
          </tr>
        </tbody>
        <tfoot>
          <tr>
            <td colspan="5" class="text-right font-weight-bold">Grand Total:</td>
            <td colspan="2" id="grandTotal" class="font-weight-bold bg-primary text-white">₦{{parseFloat(this.order.doctotal).toFixed(2)}}</td>
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
    };
  },
  mounted() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.GetOrder(parseInt(this.id)).then((order) => {
        this.order = order;
        window.backend.GetUser(parseInt(order.created_by)).then((user) => {
            this.user = user;
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
  },
};
</script>