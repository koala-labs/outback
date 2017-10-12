export default {
  loading: false,
  clusters: {
    list: [],
    selected: '',
    loading: false,
  },
  services: {
    list: [],
    selected: '',
    loading: false,
  },
  service: {
    detail: {
      Deployments: [{
        UpdatedAt: '',
      }],
    },
    loading: false,
    commit: '',
  },
  versions: {
    list: [],
    selected: '',
    loading: false,
  },
  deployment: {
    loading: false,
    result: {},
  },
};
