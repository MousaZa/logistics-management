import { fetchApi } from './api';
import type { Order } from './types';

export const getOrders = async (dateFrom?: string, dateTo?: string) => {
    const params = new URLSearchParams();
    if (dateFrom) params.append('dateFrom', dateFrom);
    if (dateTo) params.append('dateTo', dateTo);

    const query = params.toString() ? `?${params.toString()}` : '';
    const orders = await fetchApi<Order[]>('orders', `/orders${query}`);
    return { orders: orders || [] };
};

export const createOrder = (order: Order) =>
    fetchApi<void>('orders', '/orders', {
        method: 'POST',
        body: JSON.stringify(order)
    });

export const getOrderById = (uuid: string) =>
    fetchApi<Order>('orders', `/orders/${uuid}`);

export const cancelOrder = (uuid: string) =>
    fetchApi<void>('orders', `/orders/${uuid}/cancel`, {
        method: 'POST'
    });
