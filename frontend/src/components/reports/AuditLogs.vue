<template>
    <section>
        <div class="container flex-xl-nowrap2" style="margin-top: 30em;">
            <div class="row">
                <div class="col-12">
                <div class="card">
                    <div class="card-body p-0">
                        <div class="row p-2">
                            <div class="col-12 text-center">
                                <img src="../../assets/img/slot.png" />
                                <br />
                                <h3>Audit Log Report.</h3>
                                <hr />
                            </div>
                        </div>   
                    </div>
                </div>
                </div>
            </div>
        </div>

        <b-table id="report" :items="data" :busy="isBusy" :fields="fields" :sort-by.sync="sortBy" :sort-desc.sync="sortDesc" 
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
            <template v-slot:custom-foot>
                <b-tr>
                <b-td :colspan="span">
                    <b-form-group label="Per page" label-cols-sm="6" label-cols-md="4" label-cols-lg="3" label-align-sm="right" label-size="sm" label-for="perPageSelect" class="mb-0">
                        <b-form-select v-model="perPage" id="perPageSelect" size="sm" :options="pageOptions"></b-form-select>
                    </b-form-group>
                </b-td>
                <b-td :colspan="span1">
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
            span: 0,
            span1: 0,
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
            created_by: null,
        };
    },
    created() {
        alert(12345);
        this.isBusy = true;

        window.backend.GetAuditLog().then((report) => {
            this.isBusy = false;
            console.log(report);
            if (report === null) {
                this.$toast.info("Info! Report returned no data.");
                return;
            }
            
            const keys = Object.keys(report[0]);
            keys.forEach((key) => {
                this.fields.push({ key: key, sortable: true, });
            });

            // Set the dataSource
            this.data = report;
            
            // Set the column span
            this.span = Math.floor(this.fields.length / 2);
            this.span1 = this.fields.length - parseInt(this.span);
            
            // Set the initial number of items
            this.totalRows = report.length;
        },
        (err) => {
            this.isBusy = false;
            this.$toast.error("Error! " + err);
        });
    },
    methods: {
        onFiltered(filteredItems) {
            // Trigger pagination to update the number of buttons/pages due to filtering
            this.totalRows = filteredItems.length
            this.currentPage = 1
        }
    },
};
</script>