const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8090';

interface RequestOptions extends RequestInit {
    headers?: Record<string, string>;
}

export const api = {
    async request(endpoint: string, options: RequestOptions = {}) {
        const token = localStorage.getItem('token');

        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        } as Record<string, string>;

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const config = {
            ...options,
            headers,
        };

        const response = await fetch(`${API_BASE_URL}${endpoint}`, config);

        if (!response.ok) {
            const errorData = await response.text().catch(() => 'Unknown error');
            throw new Error(errorData || `Request failed with status ${response.status}`);
        }

        if (response.status === 204) return null;
        return response.json();
    },

    get(endpoint: string, headers?: Record<string, string>) {
        return this.request(endpoint, { method: 'GET', headers });
    },

    post(endpoint: string, body: any, headers?: Record<string, string>) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(body),
            headers,
        });
    },
};
