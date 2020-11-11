<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Viewing Inventory Transfer Request {{this.transfer.id}}</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'orderlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label>From Store</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.transfer.fromwhs" />
      </div>

      <div class="form-group col">
        <label>To Store</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.transfer.towhs" />
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
      <div class="form-group col-8">
        <label>Comments</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.transfer.comment" v-if="this.transfer.comment !== undefined" />
      </div>

      <div class="form-group col">
        <label for="docnum">Created On</label>
        <br />
        <span class="btn btn-info btn-sm" disabled v-if="this.user !== null" >
          {{ this.transfer.created_at }}
        </span>
      </div>
    </div>

    <div class="table-responsive">
      <table class="table table-fixed table-bordered table-hover table-sm">
        <caption>
          <h5 style="display:inline-flex;" class="float-left">Transfer Items</h5>
        </caption>
        <thead class="thead-dark">
          <tr>
            <th scope="col">#</th>
            <th scope="col">Item No.</th>
            <th scope="col">Item Description</th>
            <th scope="col">Quantity Available when requested.</th>
            <th scope="col">Quantity Requested</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, i) in transfer.items" :key="'row' + i" :id="'row' + i">
            <th scope="row">{{ i + 1 }}</th>
            <td>
              {{ item.itemcode }}
            </td>
            <td>
              {{ item.itemname }}
            </td>
            <td>
              {{ item.onhand }}
            </td>
            <td>
              {{ item.quantity }}
            </td>
          </tr>
        </tbody>
      </table>
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
import moment from "moment";

export default {
  data() {
    return {
      user: null,
      stores: [],
      transfer: {},
    };
  },
  created() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    // Get all stores
    window.backend.GetStores().then((stores) => {
      this.stores = stores;
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });
  
    window.backend.GetTransfer(parseInt(this.id)).then((transfer) => {
        this.transfer = transfer;
        let toWhs = this.transfer.towhs,
          fromWhs = this.transfer.fromwhs,
          create_time = this.transfer.created_at.Time;

        this.transfer.fromwhs = this.stores.filter((store) => {
          return (
            store.id === fromWhs
          );
        })[0].name;

        this.transfer.towhs = this.stores.filter((store) => {
          return (
            store.id === toWhs
          );
        })[0].name;

        this.transfer.created_at = moment(create_time).format("Do of MMMM YYYY");

        // Get the user who created the order
        window.backend.GetUser(parseInt(transfer.created_by)).then((user) => {
          if (JSON.stringify(user) !== "{}") {
            this.user = user;
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