<template>
    <section>
        <img src="../assets/images/slot.png" class="rounded mx-auto d-block" />
        <br />
        <h4 class="text-center" v-html="this.title"></h4>
        
        <div class="row">
            <div class="col-5">
                <label for="startDate">Choose a start date</label>
                <datepicker v-model="startdate" id="startDate" name="startDate" :minimumView="'day'" :maximumView="'year'" :initialView="'month'" :format="dateFormatter" :bootstrap-styling="true"></datepicker>
            </div>
            <div class="col-5">
                <label for="endDate">Choose a end date</label>
                <datepicker v-model="enddate" id="endDate" name="endDate" :minimumView="'day'" :maximumView="'year'" :initialView="'month'" :format="dateFormatter" :bootstrap-styling="true"></datepicker>
            </div>
            <div class="col-2">
                <br />
                <button type="button" class="btn btn-outline-primary mt-2" @click="loadLog">Load Audit Log</button>
            </div>
        </div>
        <hr />

        <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options"></v-client-table>
    </section>
</template>

<script>
import moment from "moment";
import Datepicker from 'vuejs-datepicker';

export default {
    components: {
        Datepicker
    },
    data() {
        const now = new Date(),
            today = this.dateFormatter(new Date(now.getFullYear(), now.getMonth(), now.getDate())),
            yesterday = this.dateFormatter(new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1));
        
        return {
            data: [],
            options: {},
            title: 'Audit Log Report',
            dateColumns:['created_at','updated_at', 'deleted_at', 'timestamp'],
            columns: ['id', 'old_row_data', 'new_row_data', 'dml_type', 'created_by', 'timestamp'],
            
            // Date Section
            enddate: today,
            startdate: yesterday,

        };
    },
    created() {
        this.$refs.myTable.setLoadingState(true);
        window.backend.GetAuditLogs().then((auditLog) => {
            if (auditLog === null) {
                this.$refs.myTable.setLoadingState(false);
                this.$toast.info("Info! Report returned no data.");
                return;
            }

            // Set the dataSource
            auditLog.forEach((log) => {
                log.timestamp = moment(log.timestamp.Time).utc().format('Do of MMMM YYYY, h:mm:ss a');
                log.old_row_data = log.old_row_data.String;
                this.data.push(log);
            });

            this.title = `Audit Log Report for <i class="text-primary">${moment(this.startdate).format('dddd, MMMM Do YYYY')}</i> to <i class="text-primary">${moment(this.enddate).format('dddd, MMMM Do YYYY')}</i>`;
        },
        (err) => {
            this.isBusy = false;
            this.$toast.error("Error! " + err);
        });
        },
    methods: {
        loadLog() {
            this.$refs.myTable.setLoadingState(true);
            window.backend.GetAuditLog(this.startdate, this.enddate).then((auditLog) => {
                if (auditLog === null) {
                    this.$refs.myTable.setLoadingState(false);
                    this.$toast.info("Info! There is no Audit Log for selected Period.");
                    return;
                }

                // Clear the datasource
                this.data = [];

                // Set the dataSource
                auditLog.forEach((log) => {
                    log.timestamp = moment(log.timestamp.Time).utc().format('Do of MMMM YYYY, h:mm:ss a');
                    log.old_row_data = log.old_row_data.String;
                    this.data.push(log);
                });

                this.title = `Audit Log Report for <i class="text-primary">${moment(this.startdate).format('dddd, MMMM Do YYYY')}</i> to <i class="text-primary">${moment(this.enddate).format('dddd, MMMM Do YYYY')}</i>`;
            },
            (err) => {
                this.$toast.error("Error! " + err);
                this.$refs.myTable.setLoadingState(false);
            });
        },
        dateFormatter(date) {
            return moment(date).format('yyyy-MM-D');
        },
    },
};
</script>