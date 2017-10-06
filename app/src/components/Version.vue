<template>
  <div class='main'>
    <div class="header">version</div>
    <select class="select" @change="updateVersion">
      <option value="" disabled selected>Select your version</option>
      <option v-for='version in sortedVersions' :key="version.RegistryId">{{ version.ImageTags[0] }}</option>
    </select>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import getTime from 'date-fns/get_time';

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
    sortedVersions() {
      return this.versions.sort((x, y) => {
        return getTime(y.ImagePushedAt) - getTime(x.ImagePushedAt);
      });
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
    font-weight: bold;
    padding: 20px 0 10px 0px;
    letter-spacing: 0.04em;
 }

 .select {
    width: 245px;
    height: 25px;
    font-weight: bold;
    display: table;
 }
}
</style>
