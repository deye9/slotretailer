<template>
  <section>
    <h3>Search Results for *{{ this.searchTerm }}*</h3>
    <hr />

    <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options">
      <div id="actions" slot="actions" slot-scope="{row}">
          <a class="btn btn-primary btn-sm mr-2" title="Load Details" @click="LoadDetails(row)">
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
      options: {},
      searchTerm: "",
      columns: ['occurrences', 'query', 'actions'],
    };
  },
  mounted() {
    this.$refs.myTable.setLoadingState(true);
    var pageURL = location.pathname;
    this.searchTerm = pageURL.substr(pageURL.lastIndexOf("/") + 1);
    window.backend.Search(this.searchTerm).then((result) => {
        this.$refs.myTable.setLoadingState(false);
        this.data = result;
      },
      (err) => {
        this.$refs.myTable.setLoadingState(false);
        this.$toast.error("Error! " + err);
      });
  },
  methods: {
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

        case "reports":
          var id = data.column;
          this.$store.state.reportTitle = data.occurrences.split("title = ")[1];
          this.$router.push({ name: "reportdetails", params: {id} });
          break;

        default:
          break;
      }
    },
  },
};
</script>