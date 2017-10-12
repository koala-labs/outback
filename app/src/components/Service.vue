<template>
  <div class='main'>
    <div class="header">service</div>
    <select class="select" @change="onChange">
      <option value="" disabled selected>Select your service</option>
      <option v-for='service in services' :key="service"> {{ service }} </option>
    </select>
    <div v-if="serviceDetail.commit" class="message">
      commit <b>{{ serviceDetail.commit }}</b> deployed <b>{{ timeAgo }}</b>
    </div>
  </div>
</template>

<script>
import distanceInWordsToNow from 'date-fns/distance_in_words_to_now';

export default {
  props: {
    services: {
      type: Array,
      required: true,
    },
    serviceDetail: {
      type: Object,
    },
    onChange: {
      type: Function,
      required: true,
    },
  },
  computed: {
    timeAgo() {
      return `${distanceInWordsToNow(this.serviceDetail.deployedAt)} ago`;
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

 .message {
   font-size: 14px;
   padding: 10px 0 0 0;
   text-align: left;
 }
}
</style>
