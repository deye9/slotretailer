<template>
  <section>
    <br />
    <h1>Synchronization Log</h1>
    <hr />

    <span>Please select desired Synchronization log date to view: </span>
    <b-form-select v-model="selected" :options="options" :autofocus="true" @change="getlog" class="mb-3">
      <b-form-select-option :value="null">Please select desired log date</b-form-select-option>
    </b-form-select>

    <br />
    <div id="displayContent" style="white-space: pre-line;max-height: 350px; overflow: scroll;"></div>
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

      if (this.selected === null) {
        div.innerHTML = '';
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