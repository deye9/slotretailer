<template>
  <div>
    <h1>Inventory</h1>
    <hr />

    <div class="dataList">
      <data-tables :data="data" :action-col="actionCol" :page-size="pageSize" :pagination-props="{ pageSizes: [5, 10, 15, 20] }" :table-props="tableProps" style="min-width:90%; width:100%;">
        <div slot="empty" style="color: red">There is currently no data to show</div>
        <el-table-column fixed v-for="title in titles" :prop="title.prop" :label="title.label" :key="title.label" sortable="custom"></el-table-column>
      </data-tables>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      data: [],
      titles: [],
      pageSize: 1,
      tableProps: {
        defaultSort: {
          prop: "id",
          order: "ascending",
        },
      },
      actionCol: {
        props: {
          label: "Actions",
          align: "center",
        },
        buttons: [
          {
            props: {
              size: "mini",
              type: "primary",
              icon: "el-icon-view",
            },
            handler: (row) => {
              this.$router.push("/products/details/" + row.id);
            },
            label: "Details",
          }
        ],
      },
    };
  },
  mounted() {
    window.backend.GetProducts().then((products) => {
        const exempt = [
            "id",
            "vat",
            "minlevel",
            "codebars",
            "manufacturer",
          ],
          keys = Object.keys(products[0]);
        keys.forEach((key) => {
          if (!exempt.includes(key)) {
            this.titles.push({
              prop: key,
              label: key.toUpperCase(),
            });
          }
        });

        products.forEach((product) => {
          this.data.push(product);
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