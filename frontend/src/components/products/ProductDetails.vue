<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h1>{{this.product['itemname']}}</h1>
      </div>
      <div class="col-4">
        <button type="button" class="btn btn-primary float-right">
          <router-link to="/products/" class="nav-link text-white">Inventory</router-link>
        </button>
      </div>
    </div>
    <hr />
    <ul class="list-group list-group-flush" v-for="title in this.titles" :label="title.label" :key="title.label">
      <li class="list-group-item">{{title.label}} : {{title.value}}</li>
    </ul>
  </section>
</template>

<script>
export default {
  data() {
    return {
      product: {},
      titles: [],
    };
  },
  mounted() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.ProductDetails(parseInt(this.id)).then((product) => {
        this.product = product;
        let keys = Object.keys(product);

        keys.forEach((key) => {
          this.titles.push({
            prop: key,
            value: product[key],
            label: key.toUpperCase(),
          });
        });
      },
      (err) => {
        this.$store.state.notify.category = "error";
        this.$store.state.notify.message = "Error! " + err;
      }
    );
  },
};
</script>