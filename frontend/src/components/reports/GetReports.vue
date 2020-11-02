<template>
    <section class="container">
        <div class="row flex-xl-nowrap2">
            <div class="col-8">
                <h3>Registered Reports</h3>
            </div>
            <div class="col-4">
                <router-link :to="{name: 'newreport'}" class="btn btn-info float-right">Create Report</router-link>
            </div>
        </div>
        <hr />

        <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options">
            <div id="actions" slot="actions" slot-scope="{row}">
                <a class="btn btn-primary btn-sm mr-2" title="Edit Record" @click="displayInfo(row)">
                    <i class="bi bi-pencil-fill">&nbsp;</i>
                    Edit
                </a>
                <a class="btn btn-danger btn-sm mr-2" title="Delete Record" @click="removeRow(row, event);" :v-show="allowDelete">
                    <i class="bi bi-trash-fill">&nbsp;</i>
                    Delete
                </a>
                <a class="btn btn-primary btn-sm mr-2" title="Execute Query" @click="loadReport(row)">
                    <i class="bi bi-pencil-fill">&nbsp;</i>
                    Execute Query
                </a>
            </div>
        </v-client-table>          
    </section>
</template>

<script>
    import moment from "moment";

    export default {
        data() {
            return {
            data: [],
            columns: [],
            options: {},
            allowDelete: false,
            dateColumns:['created_at','updated_at', 'deleted_at']
            };
        },
        created() {
            this.allowDelete = this.$store.state.isLoggedIn;
        },
        mounted() {
            this.$refs.myTable.setLoadingState(true);
            window.backend.GetReports().then((reports) => {
                if (reports === null) {
                    this.$toast.info("Error! No Report was returned.");
                    this.$refs.myTable.setLoadingState(false);
                    return;
                }
                    
                const exempt = [
                    "query",
                    "deleted_at",
                    "updated_at",
                    ],
                keys = Object.keys(reports[0]);

                keys.forEach((key) => {
                    if (!exempt.includes(key)) {
                        this.columns.push(key);
                    }
                });
                this.columns.push('actions');

                // Set the dataSource
                reports.forEach((report) => {
                    report.created_at = moment(report.created_at.Time).utc().format('Do of MMMM YYYY, h:mm:ss a');
                    this.data.push(report);
                });
            },
            (err) => {
                this.$toast.error("Error! " + err);
                this.$refs.myTable.setLoadingState(false);
            });
        },
        methods: {
            displayInfo(row) {
                const id = row.id;
                this.$store.state.reportTitle = row.title;
                this.$router.push({ name: "editreport", params: {id} });
            },
            loadReport(row) {
                const id = row.id;
                this.$store.state.reportTitle = row.title;
                this.$router.push({ name: "reportdetails", params: {id} });
            },
            removeRow(item, index) {
                if (item.created_by === 1 && this.$store.state.user.id !== item.created_by) {
                    this.$toast.info("ALERT! You are not permitted to delete a default report.");
                    return
                }

                index = event.srcElement.parentElement.parentElement.parentNode.rowIndex - 1;
                window.backend.RemoveReport(parseInt(item.id)).then(() => {
                    // Remove the row from the table
                    document.getElementById("myTable").getElementsByTagName('tbody')[0].deleteRow(index);

                    this.$toast.success("Success! Report has been successfully deleted.");
                }, (err) => {
                    this.$toast.error("Error! " + err);
                });
            }
        }
    };
</script>