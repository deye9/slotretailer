<template>
  <section>
    <h3>Search Results for *{{ this.searchTerm }}*</h3>
    <hr />
    <div class="table-responsive-sm">
      <b-table
        :items="data"
        :busy="isBusy"
        :fields="fields"
        :filter="filter"
        :per-page="perPage"
        @filtered="onFiltered"
        :sort-by.sync="sortBy"
        :sort-desc.sync="sortDesc"
        :current-page="currentPage"
        :sort-direction="sortDirection"
        :filter-included-fields="filterOn"
        hover
        small
        striped
        bordered
        show-empty
        caption-top
        responsive
        sticky-header>
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
        <template v-slot:cell(actions)="row">
            <b-button size="sm" variant="primary" @click="LoadDetails(row.item)" title="Order Details" style="margin-right: 2px">
                <b-icon icon="eye" aria-hidden="true"></b-icon>
            </b-button>
        </template>        
        <template v-slot:custom-foot>
            <b-tr>
            <b-td colspan="2">
                <b-form-group label="Per page" label-cols-sm="6" label-cols-md="4" label-cols-lg="3" label-align-sm="right" label-size="sm" label-for="perPageSelect" class="mb-0">
                <b-form-select v-model="perPage" id="perPageSelect" size="sm" :options="pageOptions"></b-form-select>
                </b-form-group>
            </b-td>
            <b-td colspan="2">
                <b-pagination v-model="currentPage" :total-rows="totalRows" :per-page="perPage"
                align="fill" size="sm" class="my-0"></b-pagination>
            </b-td>
            </b-tr>
        </template>        
      </b-table>
      <!-- <table class="table table-bordered table-hover table-striped table-sm">
          <thead class="thead-dark">
            <tr>
              <th scope="col">#</th>
              <th scope="col">Column</th>
              <th scope="col">Occurrences</th>
              <th scope="col"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(result, index) in this.searchResult" :key="index">
              <th scope="row">
                {{ index + 1 }}
              </th>
              <td>
                {{ result.occurrences }}
              </td>
              <td>
                <b-button
                  size="sm"
                  variant="primary"
                  @click="LoadDetails(result)"
                  title="More details..."
                  style="margin-right: 2px"
                >
                  <b-icon icon="eye" aria-hidden="true"></b-icon>
                </b-button>
              </td>
            </tr>
          </tbody>
        </table> -->
    </div>
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
      searchTerm: "",
      fields: [
        {
          key: "occurrences",
          sortable: true,
        },
        {
          key: "query",
          sortable: true,
        },
        {
          key: "Actions",
          sortable: true,
        },
      ],
      data: [],
      perPage: 10,
      filter: null,
      sortBy: "id",
      filterOn: [],
      totalRows: 1,
      isBusy: false,
      currentPage: 1,
      sortDesc: true,
      sortDirection: "desc",
      pageOptions: [5, 10, 15, 20, 25, 50, 100],
    };
  },
  created() {
    this.isBusy = true;
    var pageURL = location.pathname;
    this.searchTerm = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.Search(this.searchTerm).then((result) => {
        this.isBusy = false;
        this.data = result;

        // Set the initial number of items
        this.totalRows = result.length;
      },
      (err) => {
        this.isBusy = false;
        this.$toast.error("Error! " + err);
      }
    );
  },
  methods: {
    onFiltered(filteredItems) {
      // Trigger pagination to update the number of buttons/pages due to filtering
      this.totalRows = filteredItems.length;
      this.currentPage = 1;
    },
    LoadDetails(data) {

        switch (data.query.toLowerCase()) {
            case "orders":
            case "payments":
            case "ordereditems":
                this.$router.push("/orders/details/" + data.column);    
                break;

            case "products":
                this.$router.push("/products/details/" + data.column);
                break;

            case "customers":
                this.$router.push("/customers/edit/" + data.column);
                break;

            case "users":
                this.$router.push("/users/edit/" + data.column);
                break;            

            case "store":
                this.$router.push("/store/");
                break;

            default:
                break;
        }
    },
  },
};
</script>