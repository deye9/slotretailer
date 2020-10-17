<template>
    <section class="container">
        <div class="flex-xl-nowrap2" style="margin-top: 10em;">
            <div class="row">
                <div class="col-12">
                <div class="card">
                    <div class="card-body p-0">
                        <div class="row p-2">
                            <div class="col-12 text-center">
                                <img src="../assets/img/slot.png" />
                                <br />
                                <h4 v-html="this.title"></h4>
                            </div>
                        </div>
                        <div class="row">
                            <div class="col-5">
                                <label for="startDate">Choose a start date</label>
                                <b-input-group class="mb-3">
                                    <b-form-input v-model="startdate" type="text" autocomplete="off" placeholder="Choose a start log date to view" readonly></b-form-input>
                                    <b-input-group-append>
                                        <b-form-datepicker id="startDate" v-model="startdate" button-only right locale="en-US" aria-controls="startDate" @context="startContext" :hide-header="true" selected-variant="success" 
                                        nav-button-variant="primary" show-decade-nav today-button close-button no-flip></b-form-datepicker>
                                    </b-input-group-append>
                                </b-input-group>
                            </div>
                            <div class="col-5">
                                <label for="endDate">Choose a end date</label>
                                <b-input-group class="mb-3">
                                    <b-form-input v-model="enddate" type="text" autocomplete="off" placeholder="Choose a end log date to view" readonly></b-form-input>
                                    <b-input-group-append>
                                        <b-form-datepicker id="endDate" v-model="enddate" button-only right locale="en-US" aria-controls="endDate" @context="endContext" :hide-header="true" selected-variant="success" 
                                        nav-button-variant="primary" show-decade-nav today-button close-button no-flip></b-form-datepicker>
                                    </b-input-group-append>
                                </b-input-group>
                            </div>
                            <div class="col-2">
                                <br />
                                <b-button variant="outline-primary" @click="loadLog">Load Audit Log</b-button>
                            </div>
                        </div>
                        <hr />
                    </div>
                </div>
                </div>
            </div>
        </div>

        <b-table id="report" :items="data" :busy="isBusy" :fields="fields" :sort-by.sync="sortBy" :sort-desc.sync="sortDesc" 
        :per-page="perPage" :current-page="currentPage" :filter="filter" :filter-included-fields="filterOn"
        @filtered="onFiltered" :sort-direction="sortDirection" show-empty striped hover bordered small fixed
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
        max-height: 400px;
    }
</style>

<script>
import moment from "moment";

export default {
    data() {
        const now = new Date(),
            today = new Date(now.getFullYear(), now.getMonth(), now.getDate()),
            yesterday = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1);
        
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
            title: 'Audit Log Report',
            pageOptions: [5, 10, 15, 20, 25, 50, 100],

            // Date Section
            enddate: today,
            startdate: yesterday,
        };
    },
    created() {
        this.isBusy = true;

        window.backend.GetAuditLogs().then((auditLog) => {
            this.isBusy = false;
            if (auditLog === null) {
                this.$toast.info("Info! Report returned no data.");
                return;
            }
            
            const keys = Object.keys(auditLog[0]);
            keys.forEach((key) => {
                this.fields.push({ key: key, sortable: true, });
            });

            // Set the dataSource
            auditLog.forEach((log) => {
                log.timestamp = moment(log.timestamp.Time).utc().format('MMMM Do YYYY, h:mm:ss a');
                log.old_row_data = log.old_row_data.String;
                this.data.push(log);
            });
            
            // Set the column span
            this.span = Math.floor(this.fields.length / 2);
            this.span1 = this.fields.length - parseInt(this.span);
            
            // Set the initial number of items
            this.totalRows = auditLog.length;

            this.title = `Audit Log Report for <i>${this.startdate}</i> to <i>${this.enddate}</i>`;
        },
        (err) => {
            this.isBusy = false;
            this.$toast.error("Error! " + err);
        });
    },
    methods: {
        loadLog() {
            this.isBusy = true;
            window.backend.GetAuditLog(this.startdate, this.enddate).then((auditLog) => {
                if (auditLog === null) {
                    this.isBusy = false;
                    this.$toast.info("Info! There is no Audit Log for selected Period.");
                    return;
                }
            
                const keys = Object.keys(auditLog[0]);
                keys.forEach((key) => {
                    this.fields.push({ key: key, sortable: true, });
                });

                // Set the dataSource
                auditLog.forEach((log) => {
                    log.timestamp = moment(log.timestamp.Time).utc().format('MMMM Do YYYY, h:mm:ss a');
                    log.old_row_data = log.old_row_data.String;
                    this.data.push(log);
                });
                
                // Set the column span
                this.span = Math.floor(this.fields.length / 2);
                this.span1 = this.fields.length - parseInt(this.span);

                // Set the initial number of items
                this.totalRows = auditLog.length;
                this.isBusy = false;
            },
            (err) => {
                this.$toast.error("Error! " + err);
                this.isBusy = false;
            });
        },
        startContext(ctx) {
            // The following will be an empty string until a valid date is entered
            this.startdate = ctx.selectedYMD;
            this.title = `Audit Log Report for <i>${this.startdate}</i> to <i>${this.enddate}</i>`;
        },
        endContext(ctx) {
            // The following will be an empty string until a valid date is entered
            this.enddate = ctx.selectedYMD;
            this.title = `Audit Log Report for <i>${this.startdate}</i> to <i>${this.enddate}</i>`;
        },
        onFiltered(filteredItems) {
            // Trigger pagination to update the number of buttons/pages due to filtering
            this.totalRows = filteredItems.length
            this.currentPage = 1
        }
    },
};
</script>