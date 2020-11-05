<template>
  <section>
    <br />
    <h1>Synchronization Log</h1>
    <hr />

    <span>Please select desired Synchronization log date to view: </span>
    <select v-model="selected" @change="getlog" class="form-control">
      <option value="null">Please select desired log date</option>
      <option :key="logdate" v-for="logdate in options">{{ logdate }}</option>
    </select>

    <br />
    <div id="displayContent" style="white-space: pre-line;"></div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      options: [],
      selected: null,
    };
  },
  created() {
    window.backend.GetLogs().then((logs) => {
      this.options = logs;
    },
    (err) => {
      this.$toast.error("Error! " + err);
    });
  },
  methods: {
    getlog() {
      let div = document.getElementById('displayContent');
      if (this.selected === 'null') {
        return;
      }
      
      window.backend.GetLog(this.selected + ".log").then((logContent) => {
        logContent = logContent.replace(/banks/gi, "<strong>Banks</strong>");
        logContent = logContent.replace(/products/gi, "<strong>Products</strong>");
        logContent = logContent.replace(/customers/gi, "<strong>Customers</strong>");
        div.innerHTML = logContent;
      },
      (err) => {
        this.$toast.error("Error! " + err);
      });
    }
  }
}
</script>