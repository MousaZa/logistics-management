export const ORDERS_API_URL = import.meta.env.VITE_ORDERS_API_URL || 'http://localhost:3031/api';
export const INVENTORY_API_URL = import.meta.env.VITE_INVENTORY_API_URL || 'http://localhost:3032/api';

export class ApiError extends Error {
    slug: string;
    constructor(message: string, slug: string = 'unknown_error') {
        super(message);
        this.name = 'ApiError';
        this.slug = slug;
    }
}

export async function fetchApi<T>(service: 'orders' | 'inventory', endpoint: string, options?: RequestInit): Promise<T> {
    const baseUrl = service === 'orders' ? ORDERS_API_URL : INVENTORY_API_URL;
    const url = `${baseUrl}${endpoint}`;

    const headers = {
        'Content-Type': 'application/json',
        ...options?.headers,
    };

    const response = await fetch(url, { ...options, headers });

    if (!response.ok) {
        let errSlug = 'request_failed';
        let errMsg = 'Failed to fetch data';
        try {
            const errorData = await response.json();
            if (errorData.slug && errorData.message) {
                errSlug = errorData.slug;
                errMsg = errorData.message;
            }
        } catch (e) {
            // Ignore JSON parse error on error response
        }
        throw new ApiError(errMsg, errSlug);
    }

    // Handle 204 No Content
    if (response.status === 204) {
        return {} as T;
    }

    return await response.json();
}
