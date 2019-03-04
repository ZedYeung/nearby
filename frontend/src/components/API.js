import axios from 'axios';
import Cookies from 'universal-cookie';

const host = process.env.REACT_APP_BACKEND;

axios.interceptors.request.use(
    (config) => {
        const cookies = new Cookies();
        if (!! cookies.get(process.env.REACT_APP_TOKEN_KEY)) {
            config.headers.Authorization = `${process.env.REACT_APP_AUTH_PREFIX} ${cookies.get(process.env.REACT_APP_TOKEN_KEY)}`;
        }
        return config;
    }, (err) => {
        return Promise.reject(err);
    }
)

export const getSearch = (params) => {
    return axios.get(`${host}/search/`, { params: params });
}

export const createPost = (params) => {
    return axios.post(`${host}/post/`, params)
}

export const login = (params) => {
    return axios.post(`${host}/login/`, params)
}

export const register = (params) => {
    return axios.post(`${host}/signup/`, params)
}

export const loginAmazon = `${host}/login/amazon/`
export const loginGoogle = `${host}/login/google-oauth2/`
export const loginTwitter = `${host}/login/twitter/`