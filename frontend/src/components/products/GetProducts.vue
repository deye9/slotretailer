<template>
  <section>
    <div class="row">
      <div class="col-12">
        <h3>Inventory</h3>
      </div>
    </div>
    <hr />

    <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options">
      <div slot="actions" slot-scope="{row}">
        <a class="btn btn-primary btn-sm mr-2" title="Edit Record" @click="displayInfo(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Details
        </a>
      </div>
    </v-client-table>
  </section>
</template>

<script>
export default {
  data() {
    return {
      data: [],
      columns: [],
      options: {},
      dateColumns:['created_at','updated_at', 'deleted_at']
    };
  },
  mounted() {
    this.$refs.myTable.setLoadingState(true);
    window.backend.GetProducts().then((products) => {
      if (products === null) {
        this.$toast.info("Error! No Product was returned.");
        this.$refs.myTable.setLoadingState(false);
        return;
      }

      const exempt = [
          "vat",
          "minlevel",
          "codebars",
          "manufacturer",
          "serialnumber",
        ],
        keys = Object.keys(products[0]);

      keys.forEach((key) => {
        if (!exempt.includes(key)) {
          this.columns.push(key);
        }
      });
      this.columns.push('actions');

      // Set the dataSource
      this.data = products;
    },
    (err) => {
      this.$toast.error("Error! " + err);
      this.$refs.myTable.setLoadingState(false);
    });
  },
  methods: { 
    displayInfo(row) {
      const id = row.id;
      this.$router.push({ name: "productdetails", params: {id} });
    },
  },
};
</script>