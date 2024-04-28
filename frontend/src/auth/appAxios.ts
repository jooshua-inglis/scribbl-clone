import { API_HOST } from '@/constants';
import { store } from '@/store/store';
import axios from 'axios';

const appAxios = axios.create({
    baseURL: `${API_HOST}`
 });

appAxios.interceptors.request.use(function (config) {
    const { authToken } = store.getState().game
    if (authToken) {
        config.headers.Authorization = `Bearer ${authToken}`;
    }
    return config
}, function (error) {
    return Promise.reject(error);
});

export default appAxios