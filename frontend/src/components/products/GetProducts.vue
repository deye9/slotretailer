<template>
  <section>
    <div class="row">
      <div class="col-12">
        <h3>Inventory</h3>
      </div>
    </div>
    <hr />

    <b-table id="productList" :items="data" :busy="isBusy" :fields="fields" :sort-by.sync="sortBy" :sort-desc.sync="sortDesc" 
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
        <template v-slot:cell(price)="row">
          {{ `₦${parseFloat(row.value).toFixed(2)}` }}
        </template>
        <template v-slot:cell(actions)="row" v-if="this.$store.state.isAdmin">
          <b-button size="sm" variant="primary" @click="displayInfo(row.item)" style="margin-right: 2px">
            <b-icon icon="eye" aria-hidden="true"></b-icon>
          </b-button>
        </template>      
        <template v-slot:custom-foot>
          <b-tr>
            <b-td colspan="3">
              <b-form-group label="Per page" label-cols-sm="6" label-cols-md="4" label-cols-lg="3" label-align-sm="right" label-size="sm" label-for="perPageSelect" class="mb-0">
                <b-form-select v-model="perPage" id="perPageSelect" size="sm" :options="pageOptions"></b-form-select>
              </b-form-group>
            </b-td>
            <b-td colspan="3">
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
    window.backend.GetProducts().then((products) => {
        if (JSON.stringify(products) !== "{}") {
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
              this.fields.push({ key: key, sortable: true, });
            }
          });
          this.fields.push({ key: 'actions', label: 'Actions' });

          // Set the dataSource
          this.data = products;

          // Set the initial number of products
          this.totalRows = products.length;
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
    displayInfo(row) {
      this.$router.push("/products/details/" + row.id);
    },
    onFiltered(filteredItems) {
      // Trigger pagination to update the number of buttons/pages due to filtering
      this.totalRows = filteredItems.length
      this.currentPage = 1
    }
  },
};
</script>