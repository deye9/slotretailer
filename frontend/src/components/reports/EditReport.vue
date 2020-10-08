<template>
    <section>
        <br />
        <div class="row">
            <div class="col-8">
                <h3>Editing Report <strong>{{this.title}}</strong></h3>
            </div>
            <div class="col-4">
                <router-link to="/reports/" class="btn btn-info float-right">Back</router-link>
            </div>
        </div>
        <hr />

        <div class="form-row">
            <div class="form-group col">
                <label for="title">Report Title</label>
                <input id="title" type="text" class="form-control" placeholder="Report Title" required v-model="title"/>
            </div>

            <div class="form-group col">
                <label for="creator">Created By</label>
                <input id="creator" type="text" class="form-control" disabled required v-model="creator"/>
            </div>

             <div class="form-group col">
                <label for="createDate">Created On</label>
                <input id="createDate" type="text" class="form-control" disabled required v-model="createDate"/>
            </div>
        </div>

        <hr />
        <div>
            <b-card header-tag="header" footer-tag="footer">
            <template v-slot:header>
                <h6 class="mb-0">Report Builder</h6>
            </template>
            <b-card-text>
                <div class="form-row">
                    <div class="form-group col-4">
                        <div style="max-height: 250px; overflow: scroll">
                            <span v-for="(item, index) in this.schema" :key="index" style="text-transform: capitalize;">
                                <strong class="text-primary">
                                    {{ item.TableName }}
                                </strong>
                                <br />
                                <b-form-checkbox stacked style="text-transform: capitalize;" 
                                    v-model="selected" @change="toggleAll"
                                    v-for="option in item.columns"
                                    :key="option.column_name"
                                    :title="option.column_type"
                                    :value="item.TableName + '.' + option.column_name">
                                    {{option.column_name}}
                                </b-form-checkbox>
                                <hr />
                            </span>
                        </div>
                    </div>
                    <div class="form-group col-8">
                        <h3> Generated Query </h3>
                        <b-form-textarea v-model="qry" placeholder="Enter something..." size="lg" rows="5" max-rows="6"></b-form-textarea>
                    </div>
                </div>
            </b-card-text>
            <template v-slot:footer>
                <!-- <b-button href="#" variant="primary" class="mr-1 float-right">
                    <b-icon icon="server" aria-hidden="true"></b-icon>
                    Execute Query
                </b-button> -->
                <b-button href="#" variant="primary" class="mr-1 float-right" @click="updateQuery">
                    <b-icon icon="folder" aria-hidden="true"></b-icon>
                    Update Query
                </b-button>
            </template>
            </b-card>
        </div>
 
    </section>
</template>

<script>
// import moment from "moment";

export default {
    data() {
        return {
            id: 0,
            qry: '',
            title: '',
            schema: [],
            creator: '',
            selected: [],
            createDate: '',
            created_by: '',
            dynamicQry: [],
        };
    },
    created() {
        var pageURL = location.pathname;
        this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

        window.backend.GetTableSchema().then((schema) => {
            // Set the dataSource
            schema.forEach((key) => {
                this.schema.push({ TableName: key.TABLE_NAME, columns: JSON.parse(key.columns) });
            });
            this.dynamicQry = Array.from(this.schema.length);

            window.backend.GetDML(parseInt(this.id)).then((report) => {
                if (JSON.stringify(report) === "{}") {
                    this.$toast.info("Info! Unable to get Report details.");
                    return;
                }
                
                // Set the dataSources
                this.qry = report.query;
                this.title = report.title;
                this.creator = report.created_by;
                this.createDate = report.created_at.String;

                window.backend.GetUser(parseInt(this.creator)).then((user) => {
                    if (JSON.stringify(user) !== "{}") {
                        this.creator = user.firstname + " " + user.lastname;
                    }
                },
                (err) => {
                    this.$toast.error("Error! " + err);
                });
            },
            (err) => {
                this.$toast.error("Error! " + err);
            });
        },
        (err) => {
            this.$toast.error("Error! " + err);
        });
    },
    methods: {
        toggleAll(checked) {
            this.qry = '';
            this.dynamicQry = [];

            // Get the Details
            checked.forEach((check) => {
                let datasplit = check.split(".");
                if (!this.dynamicQry[datasplit[0]]) {
                    this.dynamicQry[datasplit[0]] = `SELECT ${datasplit[1]} FROM ${datasplit[0]};`;
                } else {                    
                    this.dynamicQry[datasplit[0]] += `${datasplit[1]} `;
                    this. dynamicQry[datasplit[0]] = this.dynamicQry[datasplit[0]].replace(` FROM ${datasplit[0]};`, ", ") + `FROM ${datasplit[0]};`;
                } 
            });

            // Display the query to the user
            for (var key in this.dynamicQry) {
                this.qry += this.dynamicQry[key] + "\n";
            } 
        },
        updateQuery() {
            if (this.title === '') {
                this.$toast.error("Error! Report Title cannot be empty.");
                return
            }

            if (this.qry === '') {
                this.$toast.error("Error! Query body cannot be empty.");
                return
            }

            const invalidCommands = [
                "create",
                "drop",
                "alter",
                "truncate",
                "comment",
                "rename",
                "update",
                "delete",
                "replace",
                "savepoint",
                "merge",
            ];

            invalidCommands.forEach((cmd) => {
                if (this.qry.toLowerCase().includes(cmd)) {
                    this.$toast.error(`Error! Your query cannot contain the word ${cmd}.`);
                    return
                }
            });

            window.backend.UpdateReport(this.title, this.qry, parseInt(this.id)).then(() => {
                this.$toast.success(`Success! Report ${this.title} has been successfully updated.`);
                this.$router.push("/reports/");
            },
            (err) => {
                this.$toast.error("Error! " + err);
            });
        }
    }
}
</script>