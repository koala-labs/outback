import TYPES from '@/store/mutation-types';

export default {
  // Clusters Mutations
  [TYPES.FETCH_CLUSTERS](state) {
    state.loading = true;
  },

  [TYPES.FETCH_CLUSTERS_SUCCESS](state, payload) {
    state.loading = false;
    state.clusters.list = payload;
  },

  [TYPES.FETCH_CLUSTERS_FAILURE](state) {
    state.loading = false;
  },

  [TYPES.SET_CLUSTER](state, cluster) {
    state.clusters.selected = cluster;
  },

  // Services Mutations
  [TYPES.FETCH_SERVICES](state) {
    state.loading = true;
  },

  [TYPES.FETCH_SERVICES_SUCCESS](state, payload) {
    state.loading = false;
    state.services.list = payload;
  },

  [TYPES.FETCH_SERVICES_FAILURE](state) {
    state.loading = false;
  },

  [TYPES.FETCH_SERVICE](state) {
    state.loading = true;
  },

  [TYPES.FETCH_SERVICE_SUCCESS](state, payload) {
    state.loading = false;
    state.service.detail = payload;
  },

  [TYPES.FETCH_SERVICE_FAILURE](state) {
    state.loading = false;
  },

  [TYPES.FETCH_SERVICE_DETAIL](state) {
    state.loading = true;
  },

  [TYPES.FETCH_SERVICE_DETAIL_SUCCESS](state, payload) {
    state.loading = false;
    state.service.commit = payload;
  },

  [TYPES.FETCH_SERVICE_DETAIL_FAILURE](state) {
    state.loading = false;
  },

  [TYPES.SET_SERVICE](state, service) {
    state.services.selected = service;
  },

  // Versions Mutations
  [TYPES.FETCH_VERSIONS](state) {
    state.loading = true;
  },

  [TYPES.FETCH_VERSIONS_SUCCESS](state, payload) {
    state.loading = false;
    state.versions.list = payload;
  },

  [TYPES.FETCH_VERSIONS_FAILURE](state) {
    state.loading = false;
  },

  [TYPES.SET_VERSION](state, version) {
    state.versions.selected = version;
  },

  // Deploy Mutations
  [TYPES.CREATE_DEPLOYMENT](state) {
    state.loading = true;
  },

  [TYPES.CREATE_DEPLOYMENT_SUCCESS](state, payload) {
    state.loading = false;
    state.deployment.result = payload;
  },

  [TYPES.CREATE_DEPLOYMENT_FAILURE](state) {
    state.loading = false;
  },
};
