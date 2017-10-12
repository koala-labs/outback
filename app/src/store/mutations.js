import TYPES from '@/store/mutation-types';

export default {
  // Clusters Mutations
  [TYPES.FETCH_CLUSTERS](state) {
    state.clusters.loading = true;
  },

  [TYPES.FETCH_CLUSTERS_SUCCESS](state, payload) {
    state.clusters.loading = false;
    state.clusters.list = payload;
  },

  [TYPES.SET_CLUSTER](state, cluster) {
    state.clusters.selected = cluster;
  },

  // Services Mutations
  [TYPES.FETCH_SERVICES](state) {
    state.services.loading = true;
  },

  [TYPES.FETCH_SERVICES_SUCCESS](state, payload) {
    state.services.loading = false;
    state.services.list = payload;
  },

  [TYPES.FETCH_SERVICE](state) {
    state.service.loading = true;
  },

  [TYPES.FETCH_SERVICE_SUCCESS](state, payload) {
    state.service.loading = false;
    state.service.detail = payload;
  },

  [TYPES.FETCH_SERVICE_DETAIL](state) {
    state.service.loading = true;
  },

  [TYPES.FETCH_SERVICE_DETAIL_SUCCESS](state, payload) {
    state.service.loading = false;
    state.service.commit = payload;
  },

  [TYPES.SET_SERVICE](state, service) {
    state.services.selected = service;
  },

  // Versions Mutations
  [TYPES.FETCH_VERSIONS](state) {
    state.services.loading = true;
  },

  [TYPES.FETCH_VERSIONS_SUCCESS](state, payload) {
    state.services.loading = false;
    state.versions.list = payload;
  },

  [TYPES.SET_VERSION](state, version) {
    state.versions.selected = version;
  },

  // Deploy Mutations
  [TYPES.CREATE_DEPLOYMENT](state) {
    state.deployment.loading = true;
  },

  [TYPES.CREATE_DEPLOYMENT_SUCCESS](state, payload) {
    state.deployment.loading = false;
    state.deployment.result = payload;
  },
};
