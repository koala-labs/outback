<template>
  <div class='main'>
    <div class="header">version</div>
    <select class="select" @change="updateVersion">
      <option value="" disabled selected>Select your version</option>
      <option v-for='version in versions' :key="version">{{ version }}</option>
    </select>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {
  methods: {
    ...mapActions([
      'setVersion',
    ]),
    updateVersion(e) {
      this.setVersion(e.target.value);
    },
  },
  computed: {
    ...mapState({
      versions: state => state.versions,
    }),
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang='scss' scoped>
.main {
 .header {
    text-align: left;
    font-weight: bold;
    padding: 20px 0 10px 20px;
    letter-spacing: 0.04em;
 }

 .select {
    width: 245px;
    height: 25px;
    font-weight: bold;
 }
}
</style>
