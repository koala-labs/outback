import API from '@/constants/api';

const socket = ({ cluster, service }) => {
  return new WebSocket(`${API.ws}${API.url}${API.routes.status}?cluster=${cluster}&service=${service}`);
};

export default socket;
