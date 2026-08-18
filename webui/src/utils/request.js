import axios from 'axios';

const service = axios.create({
    baseURL: '/api_v1', // 设置基础URL
    timeout: 1000000, // 设置请求超时时间
    headers: {
        'Content-Type': 'application/json', // 设置请求头
    },
});

service.interceptors.request.use(
    config => {
        // 在发送请求之前做些什么
        return config;
    },
    error => {
        // 对请求错误做些什么
        return Promise.reject(error);
    }
);

service.interceptors.response.use(
    response => {
        // 对响应数据做点什么
        return response.data;
    },
    error => {
        // 对响应错误做点什么
        return Promise.reject(error);
    }
);

export default service;