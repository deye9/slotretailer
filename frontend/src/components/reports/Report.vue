<template>
  <section>
    <div class="container flex-xl-nowrap2" style="margin-top: 10em;">
      <div class="row">
        <div class="col-12">
          <div class="card">
            <div class="card-body p-0">
              <div class="row p-2">
                <div class="col-12 text-center">
                    <img src="../../assets/img/slot.png" />
                    <br />
                    <h3>{{ this.$store.state.reportTitle }} Report.</h3>
              <hr />
                </div>
              </div>

              <div class="row p-2">
                <div class="col-md-12">
                  <table class="table">
                    <thead>
                      <tr>
                        <th class="border-0 text-uppercase small font-weight-bold">
                          ID
                        </th>
                        <th class="border-0 text-uppercase small font-weight-bold">
                          Item
                        </th>
                        <th class="border-0 text-uppercase small font-weight-bold">
                          Description
                        </th>
                        <th class="border-0 text-uppercase small font-weight-bold">
                          Quantity
                        </th>
                        <th class="border-0 text-uppercase small font-weight-bold">
                          Unit Cost
                        </th>
                        <th class="border-0 text-uppercase small font-weight-bold">
                          Total
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <td>1</td>
                        <td>Software</td>
                        <td>LTS Versions</td>
                        <td>21</td>
                        <td>$321</td>
                        <td>$3452</td>
                      </tr>
                      <tr>
                        <td>1</td>
                        <td>Software</td>
                        <td>Support</td>
                        <td>234</td>
                        <td>$6356</td>
                        <td>$23423</td>
                      </tr>
                      <tr>
                        <td>1</td>
                        <td>Software</td>
                        <td>Sofware Collection</td>
                        <td>4534</td>
                        <td>$354</td>
                        <td>$23434</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <div class="d-flex flex-row-reverse bg-dark text-white p-4">
                <div class="py-3 px-5 text-right">
                  <div class="mb-2">Grand Total</div>
                  <div class="h2 font-weight-light">$234,234</div>
                </div>

                <div class="py-3 px-5 text-right">
                  <div class="mb-2">Discount</div>
                  <div class="h2 font-weight-light">10%</div>
                </div>

                <div class="py-3 px-5 text-right">
                  <div class="mb-2">Sub - Total amount</div>
                  <div class="h2 font-weight-light">$32,432</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="text-light mt-5 mb-5 text-center small">
        by :
        <a class="text-light" target="_blank" href="http://totoprayogo.com"
          >totoprayogo.com</a
        >
      </div>
    </div>
  </section>
</template>

<style lang="css" scoped>
body {
  background: grey;
  margin-top: 120px;
  margin-bottom: 120px;
}
</style>

<script>
export default {
  data() {
    return {
      report: null,
      created_by: null,
    };
  },
  created() {
    var pageURL = location.pathname;
    this.id = pageURL.substr(pageURL.lastIndexOf("/") + 1);

    window.backend.GetReport(parseInt(this.id)).then(
      (report) => {
        this.report = report;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      }
    );
  },
  methods: {},
};
</script>