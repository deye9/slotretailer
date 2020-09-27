<template>
  <div>
    <div class="row">
      <div class="col-8">
        <h3>Registered Users</h3>
      </div>
      <div class="col-4">
        <router-link to="/users/new" class="btn btn-info float-right">New User</router-link>
      </div>
    </div>
    <hr />

    <data-tables
      :data="data"
      :action-col="actionCol"
      :page-size="pageSize"
      :pagination-props="{ pageSizes: [5, 10, 15, 20] }"
      :table-props="tableProps"
      style="min-width:90%; width:100%;">
      <div slot="empty" style="color: red">There is currently no data to show</div>
      <el-table-column
        fixed
        :formatter="cellValueRenderer"
        v-for="title in titles"
        :prop="title.prop"
        :label="title.label"
        :key="title.label"
        sortable="custom"
      ></el-table-column>
      <el-table-column width="55" property="isadmin"></el-table-column>
    </data-tables>
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
              icon: "el-icon-edit",
            },
            handler: (row) => {
              this.$router.push("/users/edit/" + row.id);
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
              window.backend.RemoveUser(parseInt(row.id)).then(
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
    window.backend.GetUsers().then((users) => {
        const exempt = [
            "id",
            "deleted_at",
            "password",
            "created_at",
            "updated_at",
            "created_by",
          ],
          keys = Object.keys(users[0]);
        keys.forEach((key) => {
          if (!exempt.includes(key)) {
            this.titles.push({
              prop: key,
              label: key.toUpperCase(),
            });
          }
        });

        users.forEach((user) => {
          this.data.push(user);
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
      if (typeof row[column.property] === "boolean") {
        value = String(cellValue);
      }
      return value;
    },
  },
};
</script>