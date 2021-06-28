<template>
    <section class="container">
        <br />
        <div class="row">
            <div class="col-8">
                <h3>Create Report</h3>
            </div>
            <div class="col-4">
                <router-link :to="{ name: 'reportlist' }" class="btn btn-info btn-sm float-right">Back</router-link>
            </div>
        </div>
        <hr />

        <div class="form-row">
            <div class="form-group col">
                <label for="title">Report Title</label>
                <input id="title" type="text" class="form-control form-control-sm" placeholder="Report Title" required v-model="title"/>
            </div>

            <div class="form-group col">
                <label for="creator">Created By</label>
                <input id="creator" type="text" class="form-control form-control-sm" disabled required v-model="creator"/>
            </div>

             <div class="form-group col">
                <label for="createDate">Created On</label>
                <input id="createDate" type="text" class="form-control form-control-sm" disabled required v-model="createDate"/>
            </div>
        </div>

        <hr />
        <div>
            <div class="card text-center">
                <div class="card-header">
                    <h6 class="mb-0">Report Builder</h6>
                </div>
                <div class="card-body">
                    <div class="form-row">
                        <div class="col-3 text-left" style="max-height: 250px; overflow: scroll">
                            <span v-for="(item, index) in this.schema" :key="index" style="text-transform: capitalize;">
                                <strong class="text-primary">
                                    Table: {{ item.TableName }}
                                </strong>
                                <br />
                                <div class="form-check" v-for="option in item.columns" :key="option.column_name" :title="option.column_type">
                                    <input class="form-check-input" type="checkbox" v-model="selected" @change="toggleAll" :value="item.TableName + '.' + option.column_name" :id="item.TableName + '_' + option.column_name">
                                    <label class="form-check-label" :for="item.TableName + '_' + option.column_name">
                                        {{option.column_name}}
                                    </label>
                                </div>
                                <hr />
                            </span>
                        </div>
                        <div class="col-9 text-center">
                            <h3> Generated Query </h3>
                            <br />
                            <textarea v-model="qry" placeholder="SQL Command" class="form-control form-control-sm" rows="5" max-rows="6"></textarea>
                        </div>
                    </div>
                </div>
                <div class="card-footer text-muted">
                    <button class="btn btn-primary mr-1 float-right" @click="saveQuery">
                        Save Query
                    </button>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
import moment from "moment";

export default {
    data() {
        return {
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
        this.created_by = this.$store.state.user.id;
        this.createDate = moment(new Date()).format('yyyy-MM-DD');
        this.creator = this.$store.state.user.firstname + " " + this.$store.state.user.lastname;

        window.backend.GetTableSchema().then((schema) => {
            // Set the dataSource
            schema.forEach((key) => {
                this.schema.push({ TableName: key.TABLE_NAME, columns: JSON.parse(key.columns) });
            });
            this.dynamicQry = Array.from(this.schema.length);
        },
        (err) => {
            this.$toast.error("Error! " + err);
        });
    },
    methods: {
        toggleAll() {
            this.qry = '';
            this.dynamicQry = [];
            let checked = document.querySelectorAll('input[type=checkbox]:checked');

            // Get the Details
            checked.forEach((check) => {
                let datasplit = check.value.split(".");
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
        saveQuery() {
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

            window.backend.NewReport(this.title, this.qry, this.created_by).then(() => {
                this.$toast.success(`Success! Report ${this.title} has been successfully registered.`);
                this.$router.push({ name: 'reportlist' });
            },
            (err) => {
                this.$toast.error("Error! " + err);
            });
        }
    }
}
</script>