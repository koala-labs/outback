<template>
  <div class='main'>
    <div class="header">cluster</div>
    <select class="select"  @change="updateClusterFetchServices">
      <option value="" disabled selected>Select your cluster</option>
      <option v-for='cluster in clusters' :key="cluster">{{ cluster }}</option>
    </select>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {
  beforeMount() {
    this.$Progress.start();
    this.fetchClusters();
    this.$Progress.finish();
  },
  methods: {
    ...mapActions([
      'fetchClusters',
      'fetchServices',
      'setCluster',
      'setService',
      'setVersion',
      'clearVersions',
    ]),
    updateClusterFetchServices(e) {
      this.clearServiceAndVersions();
      this.setCluster(e.target.value);
      this.$Progress.start();
      this.fetchServices(e.target.value);
      this.$Progress.finish();
    },
    clearServiceAndVersions() {
      this.setService('');
      this.setVersion('');
      this.clearVersions();
    },
  },
  computed: {
    ...mapState({
      clusters: state => state.clusters,
      cluster: state => state.deployUnit.cluster,
    }),
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang='scss' scoped>
.main {
  margin-left: 20px;

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
