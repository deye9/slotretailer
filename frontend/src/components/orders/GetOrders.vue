<template>
  <section>
    <div class="row">
      <div class="col-8">
        <h3>Sales Order</h3>
      </div>
      <div class="col-4">
        <router-link :to="{name: 'neworder'}" class="btn btn-info btn-sm mr-2" v-if="userPermission('orders', 'cancreate')">New Order</router-link>
        <router-link :to="{name: 'returnorder'}" class="btn btn-info btn-sm mr-2" v-if="userPermission('orders', 'cancreate')">Returns</router-link>
      </div>
    </div>
    <hr />

    <v-client-table ref="myTable" id="myTable" :columns="columns" v-model="data" :options="options">
      <div slot="canceled" slot-scope="{row}" style="text-transform: capitalize;">
        {{row.canceled}}
      </div>
      <div slot="doctotal" slot-scope="{row}">
        {{ row.doctotal.toLocaleString('en-NG', { style: 'currency', currency: 'NGN' }) }}
      </div>
      <div slot="synced" slot-scope="{row}" style="text-transform: capitalize;">
        {{row.synced}}
      </div>
      <div slot="returned" slot-scope="{row}" style="text-transform: capitalize;">
        {{row.returned}}
      </div>
      <div id="actions" slot="actions" slot-scope="{row}">
        <a class="btn btn-primary btn-sm mr-2" title="Print Sales Order" @click="printOrder(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Print
        </a>
        <a class="btn btn-primary btn-sm mr-2" title="Order Details" @click="displayInfo(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Details
        </a>        
        <!-- <a class="btn btn-primary btn-sm mr-2" title="Return Order" @click="returns(row)">
          <i class="bi bi-pencil-fill">&nbsp;</i>
          Return
          Allow return if the difference btw the created_date and today is less than 7 days
        </a> -->
        <a class="btn btn-danger btn-sm" title="Delete Sales Order" @click="removeRow(row, event);" v-if="userPermission('orders', 'candelete')">
          <i class="bi bi-trash-fill">&nbsp;</i>
          Delete
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
      dateColumns:['created_at','updated_at', 'deleted_at']
    };
  },
  mounted() {
    this.$refs.myTable.setLoadingState(true);

    window.backend.GetOrders().then((orders) => {
      if (JSON.stringify(orders) === "{}" || orders === null) {
        this.$refs.myTable.setLoadingState(false);
        return;
      }

      const exempt = [
          "vatsum",
          "items",
          "docnum",
          "comment",
          "payments",
          "docentry",
          "deleted_at",
          "updated_at",
          "created_by",
          "created_at",
          "discountapprovedby",
        ],
        keys = Object.keys(orders[0]);

      keys.forEach((key) => {
        if (!exempt.includes(key)) {
          this.columns.push(key);
        }
      });
      this.columns.push('actions');

      // Set the dataSource
      this.data = orders;
      this.$refs.myTable.setLoadingState(false);
    }, (err) => {
      this.$toast.error("Error! " + err);
      this.$refs.myTable.setLoadingState(false);
    });
  },
  methods: {
    returns(row) {
      const id = row.id;
      this.$router.push({ name: "returnorder", params: {id} });
    },
    printOrder() {
      // const id = row.id;
      let IDString = Math.random().toString(36).substring(7);

      var salesReceipt = window.open('', 'Sales Receipt', 'height=400,width=600');
      salesReceipt.document.write('<html><head><title>Sales Receipt</title>');
      salesReceipt.document.write('</head><body >');
      salesReceipt.document.writeln("<h3>SLOT SYSTEMS LTD.</h3>");
      salesReceipt.document.writeln(this.$store.state.userStore.address);
      salesReceipt.document.writeln(`<br />TEL: ${this.$store.state.userStore.phone}`);
      salesReceipt.document.writeln('<hr />');
      salesReceipt.document.writeln('<h3>Copy</h3>');
      salesReceipt.document.writeln('<h3>Receipt</h3>');
      salesReceipt.document.writeln(`<br /><b>ID: </b>${IDString}`);
      salesReceipt.document.writeln(`<br /><b>Date: </b>${moment().format("Do.MMMM.YYYY HH:MM:SS")}<br />`);

      this.data.items.forEach(item => {
        salesReceipt.document.writeln("<br />" + item.itemname + " " + item.quantity + " piece * ₦" + item.price);
        salesReceipt.document.writeln("<br />Serial No.:" + item.serialnumber) + "<br />";
      });

      salesReceipt.document.writeln("<br /><b>Total: ₦</b>" + this.grandTotal);
      salesReceipt.document.writeln("<br /><b>Net amount: ₦</b>" + this.grandTotal);
      salesReceipt.document.writeln("<br /><b>Tax amounts</b>: ₦0");

      salesReceipt.document.writeln("<br />You were served by: " + this.$store.state.user.firstname + " " + this.$store.state.user.lastname);
      salesReceipt.document.writeln('<br /><hr />');
      salesReceipt.document.writeln('<br /><hr />');
      salesReceipt.document.writeln('<br />**ITEMS BROUGHT IN GOOD CONDITION**');
      salesReceipt.document.writeln('<br />*****ARE NOT RETURNABLE*****');
      
      salesReceipt.document.write('</body></html>');

      salesReceipt.print();
      salesReceipt.close();
      return true;
      // this.$router.push({ name: "orderdetail", params: {id} });
    },
    displayInfo(row) {
      const id = row.id;
      this.$router.push({ name: "orderdetail", params: {id} });
    },
    removeRow(item, index) {
      window.backend.RemoveOrder(parseInt(item.id)).then(() => {
        // Remove the row from the table
        document.getElementById("myTable").getElementsByTagName('tbody')[0].deleteRow(index);
        
        this.$toast.success("Success! Order has been successfully deleted.");
      }, (err) => {
        this.$toast.error("Error! " + err);
      });
    },
  },
};
  //   var doc = new jsPDF();

  // // Set the document to automatically print via JS
  // doc.autoPrint();
</script>