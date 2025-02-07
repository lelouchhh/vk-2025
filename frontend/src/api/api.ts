import axios, { AxiosInstance } from 'axios';

const BASE_URL = 'http://localhost:8080'; // URL вашего бэкенда

// Создаем экземпляр Axios
const apiClient: AxiosInstance = axios.create({
    baseURL: BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Функция для установки токена авторизации
export const setAuthToken = (token: string | null) => {
    if (token) {
        apiClient.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        localStorage.setItem('token', token);
    } else {
        delete apiClient.defaults.headers.common['Authorization'];
        localStorage.removeItem('token');
    }
};

// Проверяем, есть ли сохраненный токен при загрузке приложения
const savedToken = localStorage.getItem('token');
if (savedToken) {
    setAuthToken(savedToken);
}

// API методы
export const register = async (login: string, password: string) => {
    const response = await apiClient.post('/register', { login, password });
    return response.data;
};

export const login = async (login: string, password: string) => {
    const response = await apiClient.post('/login', { login, password });
    const token = response.data.token;
    setAuthToken(token); // Сохраняем токен
    return token;
};

export const getContainers = async () => {
    const response = await apiClient.get('/protected/containers');
    return response.data;
};