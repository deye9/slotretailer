<template>
  <section style="margin-top: 3em ;">
    <div class="row">
      <div class="col-12">
        <h3>Inventory</h3>
      </div>
    </div>
    <hr />

    <div class="dataList">
      <data-tables
        :data="data"
        :action-col="actionCol"
        :page-size="pageSize"
        :pagination-props="{ pageSizes: [5, 10, 15, 20] }"
        :table-props="tableProps"
        style="min-width: 90%; width: 100%"
      >
        <div slot="empty" style="color: red">
          There is currently no data to show
        </div>
        <el-table-column
          fixed
          v-for="title in titles"
          :prop="title.prop"
          :label="title.label"
          :key="title.label"
          sortable="custom"
          :formatter="cellValueRenderer"
        ></el-table-column>
      </data-tables>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      data: [],
      titles: [],
      pageSize: 5,
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
          },
        ],
      },
    };
  },
  mounted() {
    window.backend.GetProducts().then(
      (products) => {
        if (JSON.stringify(products) !== "{}") {
          const exempt = [
              "id",
              "vat",
              "minlevel",
              "codebars",
              "manufacturer",
              "serialnumber",
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
        }
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );
  },
  methods: {
    cellValueRenderer(row, column, cellValue) {
      let value = cellValue;

      switch (column.property.toLowerCase()) {
        case "price":
          value = `₦${parseFloat(cellValue).toFixed(2)}`;
          break;

        default:
          break;
      }
      return value;
    },
  },
};
</script>