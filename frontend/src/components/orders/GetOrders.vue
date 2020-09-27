<template>
  <div>
    <div class="row">
      <div class="col-8">
        <h3>Orders</h3>
      </div>
      <div class="col-4">
        <router-link to="/orders/new" class="btn btn-info float-right">New Order</router-link>
      </div>
    </div>
    <hr />

    <div class="dataList">
      <data-tables ref="ordersTable" :data="data" :action-col="actionCol"
        :page-size="pageSize" :pagination-props="{ pageSizes: [5, 10, 15, 20] }"
        :table-props="tableProps" style="min-width:90%; width:100%;">
        <div slot="empty" style="color: red">There is currently no data to show</div>
        <el-table-column fixed :formatter="cellValueRenderer" v-for="title in titles" 
          :prop="title.prop" :label="title.label" :key="title.label" sortable="custom"
        ></el-table-column>
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
              this.$router.push("/orders/details/" + row.id);
            },
            label: "View",
          },
          // {
          //   props: {
          //     size: "mini",
          //     type: "primary",
          //     icon: "el-icon-edit",
          //   },
          //   handler: (row) => {
          //     this.$router.push("/orders/edit/" + row.id);
          //   },
          //   label: "Edit",
          // },
          {
            props: {
              size: 'mini',
              type: 'danger',
              icon: 'el-icon-delete'
            },
            handler: row => {
              window.backend.RemoveOrder(parseInt(row.id)).then(() => {
                  // Remove the row from the table
                  this.data.splice(this.data.indexOf(row), 1);
                },
                (err) => {
                  this.$store.state.notify.category = "error";
                  this.$store.state.notify.message = "Error! " + err;
                }
              );
            },
            label: 'delete',
          }
        ],
      },
    };
  },
   mounted() {
    window.backend.GetOrders().then((orders) => {
      const exempt = [
          "id",
          "docnum",
          "docentry",
          "cardcode",
          "deleted_at",
          "updated_at",
          "created_by",
          "created_at",
        ],
        keys = Object.keys(orders[0]).sort();
        keys.forEach((key) => {
          if (!exempt.includes(key)) {
            this.titles.push({
              prop: key,
              label: key.toUpperCase(),
            });
          }
        });

        orders.forEach((order) => {
          this.data.push(order);
        });
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );
  },
  methods: {
    cellValueRenderer(row, column, cellValue) {
        let value = cellValue;
        
        switch (column.property.toLowerCase() ) {
          case 'doctotal':
            value = `₦${parseFloat(cellValue).toFixed(2)}`;
            break;

          case 'vatsum':
            value = `${cellValue}%`;
            break;
        
          default:
            break;
        }

        if (typeof row[column.property] === 'boolean') {
            value = String(cellValue);
        }
        return value;
    },
  }
};
</script>