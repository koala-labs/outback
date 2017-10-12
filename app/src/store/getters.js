export const allSelected = (state) => {
  return {
    cluster: state.clusters.selected,
    service: state.services.selected,
    version: state.versions.selected,
  };
};

export const serviceDetail = (state) => {
  return {
    deployedAt: state.service.detail.Deployments[0].UpdatedAt,
    commit: state.service.commit,
  };
};
