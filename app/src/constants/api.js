const API = {
  routes: {
    status: '/status',
    clusters: '/clusters',
    service: '/service',
    services: '/services',
    versions: '/versions',
    commit: '/commit',
    deploy: '/deploy',
  },
};

const URL = {
  url: 'localhost:8080/ufo',
  scheme: 'http://',
  ws: 'ws://',
};

export default { ...API, ...URL };
