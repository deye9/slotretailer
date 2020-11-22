<template>
    <section>
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
                                    <button class="btn btn-primary btn-sm float-right" @click="generateReport">Download Report</button>
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

<style scoped>
    td {
        text-transform: capitalize;
    }
</style>

<script>
import $ from "jquery";
import jsPDF from "jspdf";
import html2canvas from "html2canvas";

export default {
    data() {
        return {
            data: [],
            columns: [],
            options: {},
            title: null,
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
            this.title = this.$store.state.reportTitle + " Report.";
        }, (err) => {
            this.$refs.myTable.setLoadingState(false);
            this.$toast.error("Error! " + err);
        });
    },
    methods: {
        generateReport() {
            let title = this.title,
                top_left_margin = 15,
                HTML_Width = $("#myTable").width(),
                canvas_image_width = HTML_Width,
                HTML_Height = $("#myTable").height(),
                canvas_image_height = HTML_Height,
                imgSrc = document.getElementById("img").src,
                PDF_Width = HTML_Width + (top_left_margin * 2),
                PDF_Height = (PDF_Width * 1.5) + (top_left_margin * 2),
                totalPDFPages = Math.ceil(HTML_Height / PDF_Height) - 1;

            html2canvas($('.VueTables__table')[0]).then(function(canvas) {
                canvas.getContext('2d');
                
                var imgData = canvas.toDataURL("image/jpeg", 1.0);
                var pdf = new jsPDF('p', 'pt',  [PDF_Width, PDF_Height]);

                pdf.setFontSize(20);
                pdf.setTextColor(40);
                pdf.addImage(imgSrc, 'PNG', PDF_Width / 2, 40, 10, 10);
                pdf.text(title, PDF_Width / 2, 50, null, null, 'center');
                pdf.addImage(imgData, 'JPG', top_left_margin, top_left_margin + 100, canvas_image_width, canvas_image_height);
                
                for (var i = 1; i <= totalPDFPages; i++) { 
                    pdf.addPage(PDF_Width, PDF_Height);
                    pdf.addImage(imgSrc, 'PNG', PDF_Width / 2, 40, 10, 10);
                    pdf.text(title, PDF_Width / 2, 50, null, null, 'center');
                    pdf.addImage(imgData, 'JPG', top_left_margin, - (PDF_Height * i) + (top_left_margin * 4), canvas_image_width, canvas_image_height);
                }
                
                pdf.save(`${title}pdf`);
            });
        }
    }
};
</script>