<template>
    <section class="div2PDF">
        <div class="flex-xl-nowrap2 mb-2">
            <div class="row">
                <div class="col-12">
                    <div class="card">
                        <div class="card-body p-0">
                            <div class="row p-2">
                                <div class="col-12 text-center">
                                    <img id="img" src="../../assets/images/slot.png" />
                                    <br />
                                    <h3>{{ this.title }}</h3>
                                    <!-- <button id="downloadBtn" class="btn btn-primary btn-sm float-right" @click="generateReport">Download Report</button> -->
                                    <download-excel id="downloadBtn" class="btn btn-primary btn-sm float-right" :fetch="exportTableData" :type="dataExportType" :name="title+dataExportType">
                                        Download Report
                                    </download-excel>
                                    <br />
                                    <br />

                                    <div class="radio form-check-inline float-right">
                                        <input id="exportCheckAll" v-model="dataExportAll" checked type="checkbox" value="true">
                                        <label for="exportCheckAll"> All Data </label>
                                    </div>

                                    <div class="radio form-check-inline float-right">
                                        <input id="exportRadioPDF" v-model="dataExportType" type="radio" value="pdf">
                                        <label for="exportRadioPDF"> PDF </label>
                                    </div>

                                    <div class="radio form-check-inline float-right">
                                        <input id="exportRadioXls" v-model="dataExportType" type="radio" value="xls">
                                        <label for="exportRadioXls"> XLS </label>
                                    </div>

                                    <div class="radio form-check-inline float-right">
                                        <input id="exportRadioCsv" v-model="dataExportType" type="radio" value="csv" checked>
                                        <label for="exportRadioCsv"> CSV </label>
                                    </div>
                                </div>
                            </div>   
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options"></v-client-table>
        <div id="loader"></div>
    </section>
</template>

<style scoped>
    td {
        text-transform: capitalize;
    }

    #loader {
        display: none;
        border-top: 16px solid blue;
        border-right: 16px solid green;
        border-bottom: 16px solid red;
        border-left: 16px solid pink;
        border-radius: 50%;
        width: 120px;
        height: 120px;
        animation: spin 2s linear infinite;
        /* Center */
        display: flex;
        justify-content: center;
        align-items: center;
    }

    @keyframes spin {
        0% { transform: rotate(0deg); }
        100% { transform: rotate(360deg); }
    }
</style>

<script>
import $ from "jquery";
import jsPDF from "jspdf";
import html2canvas from "html2canvas";
import JsonExcel from 'vue-json-excel';

export default {
components: {
    downloadExcel: JsonExcel,
    },
    data() {
        return {
            data: [],
            columns: [],
            options: {},
            title: null,
            dataExportAll: true,
            dataExportType: 'csv',
            dateColumns:['created_at','updated_at', 'deleted_at']
        };
    },
    created() {
        var pageURL = location.pathname;
        this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);
        window.backend.GetReport(parseInt(this.id)).then((report) => {  
            if (report === null) {
                document.getElementById("loader").style.display = "none";
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
            this.title = this.$store.state.reportTitle + " Report.";
            document.getElementById("loader").style.display = "none";
        }, (err) => {
            this.$toast.error("Error! " + err);
            document.getElementById("loader").style.display = "none";
        });
    },
    methods: {
        generateReport() {
            let title = this.title,
                top_left_margin = 15,
                HTML_Width = $(".div2PDF").width(),
                canvas_image_width = HTML_Width,
                HTML_Height = $(".div2PDF").height(),
                canvas_image_height = HTML_Height,
                // imgSrc = document.getElementById("img").src,
                PDF_Width = HTML_Width + (top_left_margin * 2),
                PDF_Height = (PDF_Width * 1.5) + (top_left_margin * 2),
                totalPDFPages = Math.ceil(HTML_Height / PDF_Height) - 1;

                $('.radio').hide(); // hides
                document.getElementById("downloadBtn").style.visibility = 'hidden';
                document.getElementsByClassName('VueTables__search')[0].style.visibility = 'hidden';

            html2canvas($('.div2PDF')[0]).then(function(canvas) {
                canvas.getContext('2d');
                
                var imgData = canvas.toDataURL("image/jpeg", 1.0);
                var pdf = new jsPDF('p', 'pt',  [PDF_Width, PDF_Height]);
                pdf.addImage(imgData, 'JPG', top_left_margin, top_left_margin, canvas_image_width, canvas_image_height);
                
                for (var i = 1; i <= totalPDFPages; i++) {
                    pdf.addPage(PDF_Width, PDF_Height);
                    pdf.addImage(imgData, 'JPG', top_left_margin, - (PDF_Height * i) + (top_left_margin * 4), canvas_image_width, canvas_image_height);
                }
                
                pdf.save(`${title}pdf`);

                $('.radio').show(); // shows
                document.getElementById("downloadBtn").style.visibility = 'visible';
                document.getElementsByClassName('VueTables__search')[0].style.visibility = 'visible';
            });
        },
        exportTableData() {
            if (this.dataExportType.toLowerCase() === "pdf") {
                return this.generateReport();
            }

            if (this.dataExportAll === true) {
                return this.$refs.myTable.data;
            } else {
                return this.$refs.myTable.filteredData;
            }
        },
    }
};
</script>