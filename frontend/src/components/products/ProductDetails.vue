<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>{{ this.product["itemname"] }}</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'productlist'}" class="btn btn-info btn-sm float-right">Back</router-link>
      </div>
    </div>
    <hr />

    <div class="form-row">
      <div class="form-group col">
        <label>Item Code</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.itemcode" />
      </div>
      <div class="form-group col">
        <label>Item Name</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.itemname" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>SAP Code Bars</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.codebars" />
      </div>
      <div class="form-group col">
        <label>Serial Number</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.serialnumber" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Warehouse</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.warehouse" />
      </div>
      <div class="form-group col">
        <label>Manufacturer</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.manufacturer" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Current Inventory</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.onhand" />
      </div>
      <div class="form-group col">
        <label>Price</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.price" />
      </div>
    </div>

    <div class="form-row">
      <div class="form-group col">
        <label>Vat</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.vat" />
      </div>
      <div class="form-group col">
        <label>Min SAP Level</label>
        <input type="text" class="form-control form-control-sm" disabled v-model="this.product.minlevel" />
      </div>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      product: {},
    };
  },
  mounted() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.ProductDetails(parseInt(this.id)).then((product) => {
      this.product = product;
      this.product.price = `₦${parseFloat(product.price).toFixed(2)}`;
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });
  },
};
</script>