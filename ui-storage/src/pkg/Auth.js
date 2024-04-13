import Cookies from 'js-cookie';
import axios from 'axios';

class Auth {
    constructor() {
        this.axiosInstance = axios.create({
            baseURL: process.env.REACT_APP_STORAGE_API_URL,
            timeout: 5000,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + Cookies.get('token')
            }
        });

        this.axiosInstance.interceptors.request.use((config) => {
                const token = Cookies.get('token');
                if (token) {
                    config.headers.Authorization = 'Bearer ' + token;
                }
                return config;
            },
            (error) => {
                return Promise.reject(error);
            }
        )

        this.axiosInstance.interceptors.response.use(
            (response) => {
                return response;
            },
            (error) => {
                if (error.response.status === 401 && window.location.pathname !== '/login') {
                    window.location = '/login';
                }
                return Promise.reject(error);
            }
        )
    }
}

export default new Auth();
