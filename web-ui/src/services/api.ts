import axios, { type AxiosRequestConfig } from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8090';

const axiosInstance = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

axiosInstance.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

axiosInstance.interceptors.response.use(
    (response) => response.data,
    (error) => {
        const message = error.response?.data?.message || error.message || 'Request failed';
        return Promise.reject(new Error(message));
    }
);

// Wrapper to fix return types since interceptor unwraps response
export const api = {
    get: <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> =>
        axiosInstance.get(url, config) as unknown as Promise<T>,

    post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> =>
        axiosInstance.post(url, data, config) as unknown as Promise<T>,

    put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> =>
        axiosInstance.put(url, data, config) as unknown as Promise<T>,

    delete: <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> =>
        axiosInstance.delete(url, config) as unknown as Promise<T>,
};

export const endpoints = {
    auth: {
        login: '/auth/login',
        register: '/auth/register',
        me: '/auth/me',
    },
    tweets: {
        feed: '/timeline/global', // Or user timeline
        create: '/tweets/',
        user: (userId: string) => `/tweets/user/${userId}`,
    },
    users: {
        profile: (userId: string) => `/users/${userId}`,
        follow: (userId: string) => `/users/${userId}/follow`,
        unfollow: (userId: string) => `/users/${userId}/unfollow`,
    }
};
