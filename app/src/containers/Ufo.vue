<template>
  <div class="main">
    <div class="header">
      ufo deploy
    </div>
    <div class="container">
      <div class="item">
        <Cluster></Cluster>
      </div>
      <div class="item">
        <Service></Service>
      </div>
      <div class="item">
        <Version></Version>
      </div>
    </div>
    <div class="button">
      <DeployButton></DeployButton>
    </div>
  </div>
</template>

<script>
import DeployButton from '@/components/DeployButton';
import Cluster from '@/components/Cluster';
import Service from '@/components/Service';
import Version from '@/components/Version';

export default {
  components: {
    DeployButton,
    Cluster,
    Service,
    Version,
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang='scss' scoped>
body {
  background-color: #EEEEEE;
}

.main {
  margin: 0 auto;
  max-width: 950px;
}

.container {
  /* flexbox setup */
  display: flex;
  margin: 15px 0 15px 0;

  /* Small screens */
  @media all and (max-width: 880px) {
      /* On small screens, we are no longer using row direction but column */
      flex-direction: column;
  }

  .item {
    width: 100%;
    background-color: white;
    border: 1px solid #EEEEEE;
    height: 150px;
  }
}

.button {
  display: flex;
  justify-content: flex-end;
}

.header {
  display: flex;
  height: 60px;
  color: white;
  background-color: #FF3D00;
  letter-spacing: 2px;
  text-align: left;
  line-height: 60px;
  font-weight: bold;
  font-size: 20px;
  padding: 0 0 0 20px;
  width: auto;
}
</style>
