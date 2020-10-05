<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Sales Order</h3>
      </div>
      <div class="col-4">
        <router-link to="/orders/new" class="btn btn-info float-right">New Order</router-link>
      </div>
    </div>
    <hr />

    <b-table id="orderList" :items="data" :busy="isBusy" :fields="fields" :sort-by.sync="sortBy" :sort-desc.sync="sortDesc" 
      :per-page="perPage" :current-page="currentPage" :filter="filter" :filter-included-fields="filterOn"
      @filtered="onFiltered" :sort-direction="sortDirection" show-empty striped hover bordered small 
      responsive sticky-header caption-top>
      <template v-slot:table-caption>
        <b-row>
          <b-col>
            <b-form-group label="Filter" label-cols-sm="6" label-align-sm="right" label-size="sm" label-for="filterInput" class="mb-0">
            <b-input-group size="sm">
              <b-form-input v-model="filter" type="search" id="filterInput" placeholder="Type to Search"></b-form-input>
              <b-input-group-append>
                <b-button :disabled="!filter" @click="filter = ''">Clear</b-button>
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
        </b-col>
      </b-row>
      </template>
      <template v-slot:table-busy>
        <div class="text-center text-danger my-2">
          <b-spinner class="align-middle"></b-spinner>
          <strong>Loading...</strong>
        </div>
      </template>
      <template v-slot:cell(id)="row">
        {{ row.value }}
      </template>
      <template v-slot:cell(vatsum)="row">
        {{ `${row.value}%` }}
      </template>
      <template v-slot:cell(doctotal)="row">
        {{ `₦${parseFloat(row.value).toFixed(2)}` }}
      </template>

      <template v-slot:cell(actions)="row" v-if="this.$store.state.isAdmin">
        <b-button size="sm" variant="primary" @click="displayInfo(row.item)" title="Order Details" style="margin-right: 2px">
          <b-icon icon="eye" aria-hidden="true"></b-icon>
        </b-button>
        <b-button size="sm" variant="primary" @click="edit(row.item)" title="Order Details" style="margin-right: 2px">
          <b-icon icon="minecart-loaded" aria-hidden="true"></b-icon>
        </b-button>
        <b-button size="sm" variant="danger" @click="removeRow(row.item, row.index)" title="Sales Order Return" class="mr-1">
          <b-icon icon="trash" aria-hidden="true"></b-icon>
        </b-button>
      </template>      
      <template v-slot:custom-foot>
        <b-tr>
          <b-td colspan="3">
            <b-form-group label="Per page" label-cols-sm="6" label-cols-md="4" label-cols-lg="3" label-align-sm="right" label-size="sm" label-for="perPageSelect" class="mb-0">
              <b-form-select v-model="perPage" id="perPageSelect" size="sm" :options="pageOptions"></b-form-select>
            </b-form-group>
          </b-td>
          <b-td colspan="4">
            <b-pagination v-model="currentPage" :total-rows="totalRows" :per-page="perPage"
              align="fill" size="sm" class="my-0"></b-pagination>
          </b-td>
        </b-tr>
      </template>
    </b-table>
  </section>
</template>

<style lang="css" scoped>
  .b-table-sticky-header {
      overflow-y: auto;
      max-height: 500px;
  }
</style>

<script>
export default {
  data() {
    return {
      data: [],
      fields: [],
      perPage: 10,
      filter: null,
      sortBy: 'id',
      filterOn: [],
      totalRows: 1,
      isBusy: false,
      currentPage: 1,
      sortDesc: true,
      sortDirection: 'desc',
      pageOptions: [5, 10, 15, 20, 25, 50, 100],
    };
  },
  mounted() {
    this.isBusy = true;
    window.backend.GetOrders().then((orders) => {
        if (JSON.stringify(orders) !== "{}") {
          const exempt = [
              "items",
              "docnum",
              "comment",
              "payments",
              "docentry",
              "cardcode",
              "returned",
              "deleted_at",
              "updated_at",
              "created_by",
              "created_at",
            ],
            keys = Object.keys(orders[0]);

          keys.forEach((key) => {
            if (!exempt.includes(key)) {
              this.fields.push({ key: key, sortable: true, });
            }
          });
          this.fields.push({ key: 'actions', label: 'Actions' });

          orders.forEach((order) => {
            this.data.push(order);
          });

           // Set the initial number of items
          this.totalRows = orders.length;
        }
        this.isBusy = false;
      },
      (err) => {
        this.$toast.error("Error! " + err);
        this.isBusy = false;
      }
    );
  },
  methods: {
    edit(row) {
      this.$router.push("/orders/edit/" + row.id);
    },
    displayInfo(row) {
      this.$router.push("/orders/details/" + row.id);
    },
    removeRow(item, index) {
      window.backend.RemoveOrder(parseInt(item.id)).then(() => {
          this.$toast.success("Success! Order has been successfully deleted.");
          // Remove the row from the table
          document.getElementById("orderList").deleteRow(index);
        },
        (err) => {
          this.$toast.error("Error! " + err);
        }
      );
    },
    onFiltered(filteredItems) {
      // Trigger pagination to update the number of buttons/pages due to filtering
      this.totalRows = filteredItems.length
      this.currentPage = 1
    },
    cellValueRenderer(row, column, cellValue) {
      let value = cellValue;

      switch (column.property.toLowerCase()) {
        case "doctotal":
          value = `₦${parseFloat(cellValue).toFixed(2)}`;
          break;

        case "vatsum":
          value = `${cellValue}%`;
          break;

        default:
          break;
      }

      if (typeof row[column.property] === "boolean") {
        value = String(cellValue);
      }
      return value;
    },
  },
};
</script>