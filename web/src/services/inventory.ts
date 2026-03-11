import { fetchApi } from './api';
import type {
    Location,
    ProductLocationInventory,
    Product,
    ProductStock,
    TransferProductRequest,
    AddProductToLocationRequest
} from './types';

// LOCATIONS
export const getLocations = async () => {
    const locations = await fetchApi<Location[]>('inventory', '/locations');
    return { locations: locations || [] };
};

export const createLocation = (location: Location) =>
    fetchApi<void>('inventory', '/locations', {
        method: 'POST',
        body: JSON.stringify(location),
    });

export const getLocationByUUID = (uuid: string) =>
    fetchApi<Location>('inventory', `/locations/${uuid}`);

export const updateLocation = (uuid: string, location: Location) =>
    fetchApi<void>('inventory', `/locations/${uuid}`, {
        method: 'PUT',
        body: JSON.stringify(location)
    });

// LOCATION PRODUCTS (Inventory)
export const getLocationContents = async (uuid: string) => {
    const items = await fetchApi<ProductStock[]>('inventory', `/locations/${uuid}/products`);
    return { items: items || [] };
};

export const addProductsToLocation = (uuid: string, products: AddProductToLocationRequest[]) =>
    fetchApi<void>('inventory', `/locations/${uuid}/products`, {
        method: 'PUT',
        body: JSON.stringify(products)
    });

// PRODUCTS
export const getProducts = async () => {
    const items = await fetchApi<Product[]>('inventory', '/products');
    return { items: items || [] };
};

export const createProduct = (product: Product) =>
    fetchApi<void>('inventory', '/products', {
        method: 'POST',
        body: JSON.stringify(product),
    });

export const getProductByUUID = (uuid: string) =>
    fetchApi<Product>('inventory', `/products/${uuid}`);

export const updateProduct = (uuid: string, product: Product) =>
    fetchApi<void>('inventory', `/products/${uuid}`, {
        method: 'PUT',
        body: JSON.stringify(product)
    });

export const getProductLocations = async (uuid: string) => {
    const locations = await fetchApi<ProductLocationInventory[]>('inventory', `/product/${uuid}/locations`);
    return { locations: locations || [] };
};

export const transferProducts = (transferReq: TransferProductRequest) =>
    fetchApi<void>('inventory', '/products/transfer', {
        method: 'PUT',
        body: JSON.stringify(transferReq)
    });
