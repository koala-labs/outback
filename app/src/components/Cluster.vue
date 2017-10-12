<template>
  <div class='main'>
    <div class="header">cluster</div>
    <select class="select" @change="onChange">
      <option value="" disabled selected>Select your cluster</option>
      <option v-for='cluster in clusters' :key="cluster">{{ cluster }}</option>
    </select>
  </div>
</template>

<script>

export default {
  name: 'cluster',
  props: {
    clusters: {
      type: Array,
      required: true,
    },
    onChange: {
      type: Function,
      required: true,
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang='scss' scoped>
.main {
  margin: 0 20px 0 20px;

 .header {
    text-align: left;
    padding: 20px 0 10px 0px;
    letter-spacing: 0.04em;
    font-weight: bold;
 }

 .select {
    width: 245px;
    height: 25px;
    font-weight: bold;
    display: table;
 }
}
</style>
