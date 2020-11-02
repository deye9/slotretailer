<template>
    <section>
        <div class="flex-xl-nowrap2 mb-2">
            <div class="row">
                <div class="col-12">
                    <div class="card">
                        <div class="card-body p-0">
                            <div class="row p-2">
                                <div class="col-12 text-center">
                                    <img src="../../assets/images/slot.png" />
                                    <br />
                                    <h3>{{ this.$store.state.reportTitle }} Report.</h3>
                                </div>
                            </div>   
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options"></v-client-table>
    </section>
</template>

<script>
export default {
    data() {
        return {
            data: [],
            columns: [],
            options: {},
            dateColumns:['created_at','updated_at', 'deleted_at']
        };
    },
    created() {
        var pageURL = location.pathname;
        this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);
        
        window.backend.GetReport(parseInt(this.id)).then((report) => {    
            this.$refs.myTable.setLoadingState(true);
            if (report === null) {
                this.$refs.myTable.setLoadingState(false);
                this.$toast.info("Info! Report returned no data.");
                return;
            }
            
            const exempt = [
                "deleted_at",
                "created_at",
                "updated_at",
                "created_by",
            ],
            keys = Object.keys(report[0]);
            keys.forEach((key) => {
                if (!exempt.includes(key)) {
                    this.columns.push(key);
                }
            });

            // Set the dataSource
            this.data = report;
        },
        (err) => {
            this.$refs.myTable.setLoadingState(false);
            this.$toast.error("Error! " + err);
        });
    }
};
</script>