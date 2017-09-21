<template>
  <div class='main'>
    <div class="header">service</div>
    <select class="select" @change="updateServiceAndFetchVersions">
      <option value="" disabled selected>Select your service</option>
      <option v-for='service in services' :key="service"> {{ service }} </option>
    </select>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {
  methods: {
    ...mapActions([
      'fetchVersions',
      'setService',
    ]),
    updateServiceAndFetchVersions(e) {
      this.setService(e.target.value);
      this.fetchVersions(this.service);
    },
  },
  computed: {
    ...mapState({
      services: state => state.services,
      service: state => state.deployUnit.service,
    }),
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang='scss' scoped>
.main {
 .header {
    text-align: left;
    padding: 20px 0 10px 20px;
    letter-spacing: 0.04em;
    font-weight: bold;
 }

 .select {
    width: 245px;
    height: 25px;
    font-weight: bold;
 }
}
</style>
