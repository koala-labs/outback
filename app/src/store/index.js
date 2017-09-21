import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    clusters: [],
    services: [],
    versions: [],
    deployUnit: {
      cluster: '',
      service: '',
      version: '',
      result: {},
    },
  },
  actions: {
    async fetchClusters({ commit }) {
      const res = await axios.get('http://localhost:8080/ufo/clusters');
      commit('FETCH_CLUSTERS', res.data);
    },
    async fetchServices({ commit }, cluster) {
      const res = await axios.get(`http://localhost:8080/ufo/services?cluster=${cluster}`);
      commit('FETCH_SERVICES', res.data);
    },
    async fetchVersions({ commit }, service) {
      const res = await axios.get(`http://localhost:8080/ufo/versions?service=${service}`);
      commit('FETCH_VERSIONS', res.data);
    },
    async runDeploy({ commit }, { cluster, service, version }) {
      const res = await axios.post('http://localhost:8080/ufo/deploy', { cluster, service, version });
      commit('SET_RESULT', res.data);
    },
    setCluster({ commit }, payload) {
      commit('SET_CLUSTER', payload);
    },
    setService({ commit }, payload) {
      commit('SET_SERVICE', payload);
    },
    setVersion({ commit }, payload) {
      commit('SET_VERSION', payload);
    },
    clearVersions({ commit }) {
      commit('CLEAR_VERSIONS');
    },
  },
  mutations: {
    FETCH_CLUSTERS(state, payload) {
      state.clusters = payload;
    },
    SET_CLUSTER(state, payload) {
      state.deployUnit.cluster = payload;
    },
    FETCH_SERVICES(state, payload) {
      state.services = payload;
    },
    SET_SERVICE(state, payload) {
      state.deployUnit.service = payload;
    },
    FETCH_VERSIONS(state, payload) {
      state.versions = payload;
    },
    SET_VERSION(state, payload) {
      state.deployUnit.version = payload;
    },
    CLEAR_VERSIONS(state) {
      state.versions = [];
    },
    SET_RESULT(state, payload) {
      state.deployUnit.result = payload;
    },
  },
  getters: {

  },
});

export default store;
