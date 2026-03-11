// Data Models from inventory.yml & orders.yml

export interface ErrorResponse {
    slug: string;
    message: string;
}

// --- Inventory Models ---

export interface Product {
    productUUID?: string;
    name: string;
    weight: number;
    price: number;
    createdAt?: string; // date-time
    updatedAt?: string; // date-time
}

export interface ProductStock {
    productUUID?: string;
    name: string;
    availableQuantity: number;
    reservedQuantity: number;
    damagedQuantity: number;
    weight: number;
    price: number;
    createdAt?: string;
    updatedAt?: string;
}

export interface Location {
    locationUUID?: string;
    name: string;
    city: string;
    address: string;
    longitude: number;
    latitude: number;
    createdAt?: string;
    updatedAt?: string;
}

export interface ProductLocationInventory extends Location {
    availableQuantity: number;
    reservedQuantity: number;
    damagedQuantity: number;
}

export interface TransferProductRequest {
    sourceLocationUUID: string;
    destLocationUUID: string;
    productUUID: string;
    quantity: number;
}

export interface AddProductToLocationRequest {
    productUUID: string;
    quantity: number;
}


// --- Order Models ---

export interface LineItem {
    productUUID?: string;
    productName: string;
    quantity: number;
    unitPrice: number;
    lineTotal: number;
    unitWeight: number;
    lineWeight: number;
}

export interface Order {
    orderUUID?: string;
    placedBy: string;
    lineItems: LineItem[];
    orderTotal: number;
    weight: number;
    status: 'pending' | 'confirmed' | 'shipped' | 'delivered' | 'completed' | 'cancelled';
    destination: string;
    orderedDate?: string;
    shippedDate?: string | null;
    deliveredDate?: string | null;
    completedDate?: string | null;
}
