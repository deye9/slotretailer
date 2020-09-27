<template>
  <div style="margin-top: 23em;">
    <div class="row">
      <div class="col-12">
        <h1>Dashboard</h1>
      </div>
    </div>
    <hr />
    
      <div class="card">
        <div class="card-header">
          <strong>Analysis</strong>
        </div>
        <div class="card-body">
          <div class="card-deck">
            <div class="card text-white bg-primary" style="max-width: 18rem">
              <div class="card-header">Daily Sales</div>
              <div class="card-body">
                <h3 class="card-title">₦{{this.salesTotal.toFixed(2)}}</h3>
              </div>
            </div>
            <div class="card text-white bg-info" style="max-width: 18rem">
              <div class="card-header">Last Sync</div>
              <div class="card-body">
                <h5 class="card-title">Daily Sales</h5>
                <p class="card-text">
                  None for now
                </p>
              </div>
            </div>
            <div class="card text-white bg-secondary" style="max-width: 18rem">
              <div class="card-header">Sync Errors</div>
              <div class="card-body">
                <h5 class="card-title">Special title treatment</h5>
                <p class="card-text">
                  Count of Errors.
                </p>
              </div>
            </div>
          </div>
        </div>
        <div class="card-footer">
          <small class="text-muted">Last updated  <strong>{{this.lastSyncTime}}</strong></small>
        </div>
      </div>

      <hr />

      <div class="card">
        <div class="card-header">
          <strong>Today's Orders</strong>
        </div>
        <div class="card-body">
          <data-tables ref="ordersTable" :data="orders"
            :page-size="pageSize" :pagination-props="{ pageSizes: [5, 10, 15, 20] }"
            :table-props="tableProps" style="min-width:90%; width:100%;">
            <div slot="empty" style="color: red">There is currently no data to show</div>
            <el-table-column fixed :formatter="cellValueRenderer" v-for="title in titles" 
              :prop="title.prop" :label="title.label" :key="title.label" sortable="custom"
            ></el-table-column>
          </data-tables>
        </div>
        <div class="card-footer">
          <small class="text-muted">Last updated  <strong>{{this.lastSyncTime}}</strong></small>
        </div>
      </div>

      <hr />
      
      <div class="card">
        <div class="card-header">
          <strong>Week's Top Sellers</strong>
        </div>
        <div class="card-body">
          <data-tables ref="itemsTable" :data="items"
            :page-size="pageSize" :pagination-props="{ pageSizes: [5, 10, 15, 20] }"
            :table-props="tableProps" style="min-width:90%; width:100%;">
            <div slot="empty" style="color: red">There is currently no data to show</div>
            <el-table-column fixed :formatter="cellValueRenderer" v-for="title in itemstitle" 
              :prop="title.prop" :label="title.label" :key="title.label" sortable="custom"
            ></el-table-column>
          </data-tables>
        </div>
        <div class="card-footer">
          <small class="text-muted">Last updated <strong>{{this.lastSyncTime}}</strong></small>
        </div>
      </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      items: [],
      orders: [],
      titles: [],
      pageSize: 1,
      itemstitle: [],
      salesTotal: 0.00,
      lastSyncTime: new Date().toLocaleString(),
    };
  },
  mounted() {
    this.getData();

    // Run every 5 minutes [5 * 60 * 1000 = 300000]
    setTimeout(() => {
      this.getData();
    }, 300000);
  },
  methods: {
    getData() {
      window.backend.Dashboard().then((response) => {
        // Today's Orders
        this.titles = [];
        const exempt = [
          "id",
          "items",
          "docnum",
          "docentry",
          "canceled",
          "cardcode",
          "deleted_at",
          "updated_at",
          "created_by",
          "created_at",
        ],
        keys = Object.keys(response.orders[0]).sort();
        keys.forEach((key) => {
          if (!exempt.includes(key)) {
            // key = key.replace("vatsum", "Vat");
            // key = key.replace("doctotal", "Total");
            // key = key.replace("cardname", "Customer");
            this.titles.push({
              prop: key,
              label: key.toUpperCase(),
            });
          }
        });

        let total = 0.00;
        response.orders.forEach((order) => {
          total += parseFloat(order.doctotal);
          this.orders.push(order);
        });
        
        // Top Sellers
        this.itemsexempt = [];
        const itemsexempt = [
          "id",
          "orderid",
          "itemcode",
          "deleted_at",
          "updated_at",
          "created_by",
          "created_at",
        ],
        itemkeys = Object.keys(response.items[0]).sort();
        itemkeys.forEach((key) => {
          if (!itemsexempt.includes(key)) {
            // key = key.replace("vatsum", "Vat");
            // key = key.replace("doctotal", "Total");
            // key = key.replace("cardname", "Customer");
            this.itemstitle.push({
              prop: key,
              label: key.toUpperCase(),
            });
          }
        });

        response.items.forEach((order) => {
          this.items.push(order);
        });
        
        this.salesTotal = total;
        this.lastSyncTime = new Date().toLocaleString();
        },
        (err) => {
          this.$toast.error("Error! " + err);
      });
    },
    cellValueRenderer(row, column, cellValue) {
        let value = cellValue;
        
        switch (column.property.toLowerCase() ) {
          case 'doctotal':
            value = `₦${parseFloat(cellValue).toFixed(2)}`;
            break;

          case 'price':
            value = `₦${parseFloat(cellValue).toFixed(2)}`;
            break;

          case 'vatsum':
            value = `${cellValue}%`;
            break;
        
          default:
            break;
        }

        if (typeof row[column.property] === 'boolean') {
            value = String(cellValue);
        }
        return value;
    },
  },
};
</script>