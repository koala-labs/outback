<template>
  <div class='main'>
    <div class="header">service</div>
    <select class="select" @change="updateServiceAndFetchVersions">
      <option value="" disabled selected>Select your service</option>
      <option v-for='service in services' :key="service"> {{ service }} </option>
    </select>
    <div v-if="currentService" class="message">
      commit <b>{{ currentCommit }}</b> deployed <b>{{ timeAgo }}</b>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import distanceInWordsToNow from 'date-fns/distance_in_words_to_now';

export default {
  data() {
    return {

    };
  },
  methods: {
    ...mapActions([
      'fetchCommit',
      'fetchService',
      'fetchVersions',
      'setService',
      'clearVersions',
    ]),
    updateServiceAndFetchVersions(e) {
      this.$Progress.start();
      this.clearVersions();
      this.setService(e.target.value).then(() => {
        this.fetchVersions(this.deployUnit.service);
      });
      this.fetchService(this.deployUnit).then(() => {
        this.fetchCommit(this.currentService.TaskDefinition);
      });
      this.$Progress.finish();
    },
  },
  computed: {
    ...mapState({
      deployUnit: state => state.deployUnit,
      services: state => state.services,
      currentService: state => state.service,
      currentCommit: state => state.commit,
    }),
    timeAgo() {
      return `${distanceInWordsToNow(this.currentService.Deployments[0].UpdatedAt)} ago`;
    },
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

 .message {
   font-size: 14px;
   padding: 10px 0 0 0;
   text-align: left;
 }
}
</style>
