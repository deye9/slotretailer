<template>
    <section>
        <div class="row">
            <div class="col-8">
                <h3>Create Report</h3>
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
            <div class="form-row">
                <div class="form-group col">
                    <b-card title="Title" header-tag="header" footer-tag="footer">
                    <template v-slot:header>
                        <h6 class="mb-0">Report Builder</h6>
                    </template>
                    <b-card-text>
                        Header and footers using slots.
                    </b-card-text>
                    <template v-slot:footer>
                        <b-button href="#" variant="primary" class="mr-1">
                            <b-icon icon="server" aria-hidden="true"></b-icon>
                            Execute Query
                        </b-button>
                    </template>
                    </b-card>
                </div>
                <div class="form-group col">
                    <h3> Result </h3>
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
            title: '',
            creator: '',
            createDate: '',
            created_by: '',
        };
    },
    mounted() {
        this.createDate = moment().format();
        this.created_by = this.$store.state.user.id;
        this.creator = this.$store.state.user.firstname + " " + this.$store.state.user.lastname;
    }
}
</script>