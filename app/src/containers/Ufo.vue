<template>
  <div class="main">
    <div class="header">
      ufo deploy
    </div>
    <div class="container">
      <div class="item">
        <Cluster :clusters="clusters" :onChange="setCluster"></Cluster>
      </div>
      <div class="item">
        <Service :services="services" :onChange="setService" :serviceDetail="serviceDetail"></Service>
      </div>
      <div class="item">
        <Version :versions="versions" :onChange="setVersion"></Version>
      </div>
    </div>
    <ClipLoader :loading="loading"></ClipLoader>
    <div class="button">
      <DeployButton :onClick="deploy" :disabled="areAllSelected"></DeployButton>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions, mapGetters } from 'vuex';
import Cluster from '@/components/Cluster';
import Service from '@/components/Service';
import Version from '@/components/Version';
import DeployButton from '@/components/DeployButton';
import Socket from '@/utils/socket';
import ClipLoader from 'vue-spinner/src/ClipLoader';

export default {
  components: {
    Cluster,
    Service,
    Version,
    DeployButton,
    ClipLoader,
  },
  data() {
    return {
      isDeployed: false,
    };
  },
  created() {
    this.fetchClusters();
  },
  methods: {
    ...mapActions([
      'fetchClusters', 'fetchServices', 'fetchVersions',
      'fetchServiceDetail', 'createDeployment',
    ]),
    setCluster(e) {
      this.$store.dispatch('setCluster', e.target.value);
      this.fetchServices(this.allSelected);
    },
    setService(e) {
      this.$store.dispatch('setService', e.target.value);
      this.fetchServiceDetail(this.allSelected);
      this.fetchVersions(this.allSelected);
    },
    setVersion(e) {
      this.$store.dispatch('setVersion', e.target.value);
    },
    show(text, type = '') {
      const group = 'ufo';
      const title = 'UFO';
      this.$notify({ group, title, text, type });
    },
    deploy() {
      this.createDeployment(this.allSelected).then(() => {
        this.show('Your deployment has been scheduled', '');
      });
      Socket(this.allSelected).addEventListener('message', (e) => {
        this.isDeployed = JSON.parse(e.data);
      });
    },
  },
  computed: {
    ...mapState({
      loading: state => state.loading,
      clusters: state => state.clusters.list,
      services: state => state.services.list,
      versions: state => state.versions.list,
      serviceDeployedAt: state => state.services.deployedAt,
      serviceCommit: state => state.services.commit,
    }),
    ...mapGetters(['allSelected', 'serviceDetail']),
    areAllSelected() {
      return !(this.allSelected.cluster && this.allSelected.service && this.allSelected.version);
    },
  },
  watch: {
    isDeployed: function checkIfDeployed() {
      if (this.isDeployed) {
        this.show('Deploy Successful', 'success');
      }
    },
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
