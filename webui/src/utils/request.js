import axios from 'axios';
import { ElMessageBox } from 'element-plus'


const service = axios.create({
    baseURL: '/api_v1', // 设置基础URL
    timeout: 1000000, // 设置请求超时时间
    headers: {
        'Content-Type': 'application/json', // 设置请求头
    },
});

service.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    error => {
        // 对请求错误做些什么
        return Promise.reject(error);
    }
);

service.interceptors.response.use(
    response => {
        response.status === 200 && console.log('请求成功');
        return response.data;
    },
    error => {
        // 对响应错误做点什么
        if (error.response && error.response.status === 401) {
            ElMessageBox.alert('未登录或登录已过期，请重新登录', '未登录', {
                // if you want to disable its autofocus
                // autofocus: false,
                confirmButtonText: '重新登陆',
                callback: (action) => {
                    if (action === 'confirm') {
                        localStorage.removeItem('token');
                        window.location.href = '/login';
                    }
                },
            })
        } else if (error.response && error.response.status === 403) {
            console.error('禁止访问', error.response);
            ElNotification({
                title: '错误',
                message: '禁止访问',
                type: 'error',
            })
        } else if (error.response && error.response.status === 404) {
            console.error('请求的资源不存在', error.response);
            ElNotification({
                title: '错误',
                message: '请求的资源不存在',
                type: 'error',
            })
        } else if (error.response && error.response.status === 500) {
            console.error('服务器内部错误', error.response);
            ElNotification({
                title: '错误',
                message: '服务器内部错误',
                type: 'error',
            })
        } else {
            console.error('请求失败', error.response);
            ElNotification({
                title: '错误',
                message: '请求失败',
                type: 'error',
            })
        }
        return Promise.reject(error);
    }
);

export default service;