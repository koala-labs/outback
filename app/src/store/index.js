import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    clusters: [],
    services: [],
    commit: '',
    service: null,
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
    async fetchService({ commit }, { cluster, service }) {
      const res = await axios.get(`http://localhost:8080/ufo/service?cluster=${cluster}&service=${service}`);
      commit('FETCH_SERVICE', res.data);
    },
    async fetchVersions({ commit }, service) {
      const res = await axios.get(`http://localhost:8080/ufo/versions?service=${service}`);
      commit('FETCH_VERSIONS', res.data);
    },
    async fetchCommit({ commit }, definition) {
      const res = await axios.get(`http://localhost:8080/ufo/commit?definition=${definition}`);
      commit('FETCH_COMMIT', res.data);
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
    FETCH_COMMIT(state, payload) {
      state.commit = payload;
    },
    FETCH_CLUSTERS(state, payload) {
      state.clusters = payload;
    },
    SET_CLUSTER(state, payload) {
      state.deployUnit.cluster = payload;
    },
    FETCH_SERVICES(state, payload) {
      state.services = payload;
    },
    FETCH_SERVICE(state, payload) {
      state.service = payload;
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
