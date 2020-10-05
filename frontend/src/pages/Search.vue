<template>
    <section>
        <h3>Search Results for *{{this.searchTerm}}*</h3>
        <hr />

        <div style="max-height: 350px; overflow: scroll">
            <div class="table-responsive-sm">
                <table class="table table-bordered table-hover table-striped table-sm">
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
                                {{ result.column }}
                            </td>
                            <td>
                                {{ result.occurrences }}
                            </td>
                            <td>
                                <b-button size="sm" variant="primary" @click="LoadDetails(result)" title="More details..." style="margin-right: 2px">
                                    <b-icon icon="eye" aria-hidden="true"></b-icon>
                                </b-button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </section>
</template>

<script>
export default {
    data() {
        return {
            searchTerm: '',
            searchResult: []
        }
    },
    created() {
        var pageURL = location.pathname;
        this.searchTerm = pageURL.substr(pageURL.lastIndexOf("/") + 1);

        window.backend.Search(this.searchTerm).then((result) => {
            this.searchResult = result;
        },
        (err) => {
            this.$toast.error("Error! " + err);
        });
    },
    methods: {
        LoadDetails(data) {
            console.log(data);
        },
  },
};
</script>