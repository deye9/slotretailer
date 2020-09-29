<template>
  <div style="margin-top: 3em ;">
    <div class="row">
      <div class="col-8">
        <h3>Registered Customers</h3>
      </div>
      <div class="col-4">
        <router-link to="/customers/new" class="btn btn-info float-right"
          >New Customer</router-link
        >
      </div>
    </div>
    <hr />

    <data-tables
      :data="data"
      :action-col="actionCol"
      :page-size="pageSize"
      :pagination-props="{ pageSizes: [5, 10, 15, 20] }"
      :table-props="tableProps"
    >
      <div slot="empty" style="color: red">
        There is currently no data to show
      </div>
      <el-table-column
        v-for="title in titles"
        :prop="title.prop"
        :label="title.label"
        :key="title.label"
        sortable="custom"
      ></el-table-column>
    </data-tables>
  </div>
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
              icon: "el-icon-edit",
            },
            handler: (row) => {
              this.$router.push("/customers/edit/" + row.id);
            },
            label: "Edit",
          },
          {
            props: {
              size: "mini",
              type: "danger",
              icon: "el-icon-delete",
            },
            handler: (row) => {
              window.backend.RemoveCustomer(parseInt(row.id)).then(
                () => {
                  // Remove the row from the table
                  this.data.splice(this.data.indexOf(row), 1);
                },
                (err) => {
                  this.$store.state.notify.category = "error";
                  this.$store.state.notify.message = "Error! " + err;
                }
              );
            },
            label: "delete",
          },
        ],
      },
    };
  },
  mounted() {
    window.backend.GetCustomers().then(
      (customers) => {
        if (JSON.stringify(customers) !== "{}") {
          const exempt = [
              "cardcode",
              "id",
              "deleted_at",
              "created_at",
              "updated_at",
              "created_by",
            ],
            keys = Object.keys(customers[0]);

          keys.forEach((key) => {
            if (!exempt.includes(key)) {
              this.titles.push({
                prop: key,
                label: key.toUpperCase(),
              });
            }
          });

          customers.forEach((customer) => {
            this.data.push(customer);
          });
        }
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );
  },
};
</script>