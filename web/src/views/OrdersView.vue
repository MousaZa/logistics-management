<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Button from '../components/common/Button.vue'
import DataTable from '../components/common/DataTable.vue'
import Modal from '../components/common/Modal.vue'
import { getOrders, cancelOrder, createOrder } from '../services/orders'
import { getProducts } from '../services/inventory'
import type { Order, LineItem, Product } from '../services/types'

const orders = ref<Order[]>([])
const loading = ref(true)
const dateFrom = ref('')
const dateTo = ref('')

const columns = [
  { key: 'orderUUID', label: 'Order ID', format: (v: string) => v ? v.substring(0, 8) + '...' : '-' },
  { key: 'placedBy', label: 'Customer' },
  { key: 'destination', label: 'Destination' },
  { key: 'orderTotal', label: 'Total ($)', format: (v: number) => `$${v.toFixed(2)}` },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions', align: 'right' as const }
]

const loadData = async () => {
  loading.value = true
  try {
    const from = dateFrom.value ? new Date(dateFrom.value).toISOString() : undefined
    const to = dateTo.value ? new Date(dateTo.value).toISOString() : undefined
    const res = await getOrders(from, to)
    orders.value = res.orders || []
  } catch (err) {
    console.error('Failed to load orders', err)
  } finally {
    loading.value = false
  }
}

// Create Order Modal State
const showCreateModal = ref(false)
const products = ref<Product[]>([])
const saving = ref(false)

const newOrderForm = ref({
  placedBy: '',
  destination: '',
  lineItems: [] as { productUUID: string; quantity: number }[]
})

const loadProducts = async () => {
  try {
    const res = await getProducts()
    products.value = res.items || []
  } catch (err) {
    console.error('Failed to load products', err)
  }
}

const addLineItem = () => {
  newOrderForm.value.lineItems.push({ productUUID: '', quantity: 1 })
}

const removeLineItem = (index: number) => {
  newOrderForm.value.lineItems.splice(index, 1)
}

const handleCreateOrder = async () => {
  if (!newOrderForm.value.placedBy || !newOrderForm.value.destination || newOrderForm.value.lineItems.length === 0) {
    alert('Please fill out all required fields and add at least one item.')
    return
  }
  
  if (newOrderForm.value.lineItems.some(item => !item.productUUID || item.quantity <= 0)) {
    alert('Please select a product and valid quantity for all line items.')
    return
  }

  saving.value = true
  try {
    let orderTotal = 0
    let totalWeight = 0
    
    const lineItems: LineItem[] = newOrderForm.value.lineItems.map(item => {
      const product = products.value.find(p => p.productUUID === item.productUUID)!
      const lineTotal = product.price * item.quantity
      const lineWeight = product.weight * item.quantity
      
      orderTotal += lineTotal
      totalWeight += lineWeight
      
      return {
        productUUID: product.productUUID,
        productName: product.name,
        quantity: item.quantity,
        unitPrice: product.price,
        lineTotal: lineTotal,
        unitWeight: product.weight,
        lineWeight: lineWeight
      }
    })

    const orderPayload: Order = {
      placedBy: newOrderForm.value.placedBy,
      destination: newOrderForm.value.destination,
      status: 'pending',
      lineItems,
      orderTotal,
      weight: totalWeight
    }

    await createOrder(orderPayload)
    showCreateModal.value = false
    newOrderForm.value = { placedBy: '', destination: '', lineItems: [] }
    loadData()
  } catch (err) {
    console.error('Failed to create order', err)
    alert('Failed to create order')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadData()
  loadProducts()
})

// Details Modal State
const showDetailsModal = ref(false)
const selectedOrder = ref<Order | null>(null)

const handleViewDetails = (order: Order) => {
  selectedOrder.value = order
  showDetailsModal.value = true
}

const handleCancel = async (uuid: string) => {
  if (!confirm('Are you sure you want to cancel this order?')) return;
  try {
    await cancelOrder(uuid);
    loadData(); // refresh list
  } catch (err) {
    console.error('Failed to cancel order', err);
    alert('Failed to cancel order');
  }
}
</script>

<template>
  <div class="orders-view">
    <div class="page-header flex justify-between items-center mb-6">
      <div>
        <h2 class="mb-2">Order Management</h2>
        <p class="text-secondary">Track and process logistic shipments.</p>
      </div>
      <div class="flex gap-4 items-end">
        <div class="filter-group">
          <label class="text-xs text-secondary block mb-1">From Date</label>
          <input type="date" v-model="dateFrom" @change="loadData" class="input-field" />
        </div>
        <div class="filter-group">
          <label class="text-xs text-secondary block mb-1">To Date</label>
          <input type="date" v-model="dateTo" @change="loadData" class="input-field" />
        </div>
        <Button @click="showCreateModal = true">Create Order</Button> 
      </div>
    </div>

    <DataTable 
      :columns="columns" 
      :data="orders" 
      :loading="loading"
    >
      <template #cell-status="{ value }">
        <span class="status-badge" :class="`status-${value}`">
          {{ value }}
        </span>
      </template>

      <template #cell-actions="{ row }">
        <div class="flex gap-2 justify-end">
          <Button variant="ghost" size="sm" @click="handleViewDetails(row)">Details</Button>
          <Button 
            v-if="['pending', 'confirmed'].includes(row.status)" 
            variant="danger" 
            size="sm" 
            @click="handleCancel(row.orderUUID)"
          >
            Cancel
          </Button>
        </div>
      </template>
    </DataTable>

    <!-- Create Order Modal -->
    <Modal v-if="showCreateModal" @close="showCreateModal = false">
      <template #title>Create New Order</template>
      
      <div class="form-group mb-4">
        <label class="block text-sm mb-1">Placed By</label>
        <input v-model="newOrderForm.placedBy" type="text" class="input-field w-full" placeholder="Customer Name" />
      </div>
      
      <div class="form-group mb-4">
        <label class="block text-sm mb-1">Destination</label>
        <input v-model="newOrderForm.destination" type="text" class="input-field w-full" placeholder="Shipping Address" />
      </div>

      <div class="mb-4">
        <div class="flex justify-between items-center mb-2">
          <label class="block text-sm font-semibold">Line Items</label>
          <Button variant="ghost" size="sm" @click="addLineItem">+ Add Item</Button>
        </div>
        
        <div v-if="newOrderForm.lineItems.length === 0" class="text-sm text-secondary italic mb-2">
          No items added yet.
        </div>
        
        <div v-for="(item, index) in newOrderForm.lineItems" :key="index" class="line-item-row flex gap-2 mb-2 items-center">
          <select v-model="item.productUUID" class="input-field flex-grow">
            <option value="" disabled>Select Product</option>
            <option v-for="p in products" :key="p.productUUID" :value="p.productUUID">
              {{ p.name }} (${{ p.price.toFixed(2) }})
            </option>
          </select>
          <input v-model.number="item.quantity" type="number" min="1" class="input-field w-24" placeholder="Qty" />
          <Button variant="ghost" size="sm" @click="removeLineItem(index)" class="text-danger font-bold text-lg">&times;</Button>
        </div>
      </div>

      <template #footer>
        <Button variant="ghost" @click="showCreateModal = false">Cancel</Button>
        <Button @click="handleCreateOrder" :disabled="saving">
          {{ saving ? 'Saving...' : 'Create Order' }}
        </Button>
      </template>
    </Modal>

    <!-- Order Details Modal -->
    <Modal v-if="showDetailsModal && selectedOrder" @close="showDetailsModal = false">
      <template #title>Order Details</template>
      
      <div class="mb-4">
        <div class="flex justify-between mb-2">
          <span class="text-secondary">Order ID:</span>
          <span class="font-semibold">{{ selectedOrder.orderUUID }}</span>
        </div>
        <div class="flex justify-between mb-2">
          <span class="text-secondary">Status:</span>
          <span class="status-badge" :class="`status-${selectedOrder.status}`">{{ selectedOrder.status }}</span>
        </div>
        <div class="flex justify-between mb-2">
          <span class="text-secondary">Placed By:</span>
          <span class="font-semibold">{{ selectedOrder.placedBy }}</span>
        </div>
        <div class="flex justify-between mb-2">
          <span class="text-secondary">Destination:</span>
          <span class="font-semibold">{{ selectedOrder.destination }}</span>
        </div>
        <div class="flex justify-between mb-2" v-if="selectedOrder.orderedDate">
          <span class="text-secondary">Ordered Date:</span>
          <span class="font-semibold">{{ new Date(selectedOrder.orderedDate).toLocaleString() }}</span>
        </div>
      </div>

      <div class="mb-4">
        <h4 class="mb-2 mt-4">Line Items</h4>
        <div v-if="!selectedOrder.lineItems || selectedOrder.lineItems.length === 0" class="text-sm text-secondary italic">
          No items in this order.
        </div>
        <table v-else class="w-full text-sm">
          <thead>
            <tr class="border-b">
              <th class="text-left pb-2">Product</th>
              <th class="text-right pb-2">Qty</th>
              <th class="text-right pb-2">Price</th>
              <th class="text-right pb-2">Total</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, idx) in selectedOrder.lineItems" :key="idx" class="border-b">
              <td class="py-2">{{ item.productName }}</td>
              <td class="py-2 text-right">{{ item.quantity }}</td>
              <td class="py-2 text-right">${{ item.unitPrice.toFixed(2) }}</td>
              <td class="py-2 text-right font-semibold">${{ item.lineTotal.toFixed(2) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <div class="flex justify-between items-center border-t pt-4 mt-4">
        <div>
          <span class="text-secondary text-sm">Total Weight:</span>
          <span class="font-semibold ml-2">{{ selectedOrder.weight }} kg</span>
        </div>
        <div class="text-lg font-bold">
          <span>Total:</span>
          <span class="ml-2">${{ selectedOrder.orderTotal.toFixed(2) }}</span>
        </div>
      </div>

      <template #footer>
        <Button @click="showDetailsModal = false">Close</Button>
      </template>
    </Modal>
  </div>
</template>

<style scoped>
.mb-2 { margin-bottom: var(--spacing-2); }
.mb-6 { margin-bottom: var(--spacing-6); }
.text-secondary { color: var(--color-text-secondary); margin: 0; }
.justify-end { justify-content: flex-end; }

.status-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-pending { background: rgba(245, 158, 11, 0.2); color: var(--color-accent-warning); }
.status-confirmed { background: rgba(59, 130, 246, 0.2); color: var(--color-brand-primary); }
.status-shipped { background: rgba(16, 185, 129, 0.2); color: var(--color-accent-success); }
.status-delivered { background: rgba(16, 185, 129, 0.2); color: var(--color-accent-success); }
.status-completed { background: rgba(16, 185, 129, 0.2); color: var(--color-accent-success); }
.status-cancelled { background: rgba(239, 68, 68, 0.2); color: var(--color-accent-danger); }

.flex { display: flex; }
.gap-4 { gap: var(--spacing-4); }
.gap-2 { gap: var(--spacing-2); }
.items-end { align-items: flex-end; }
.items-center { align-items: center; }
.text-xs { font-size: 0.75rem; }
.text-sm { font-size: 0.875rem; }
.text-lg { font-size: 1.125rem; }
.text-danger { color: var(--color-accent-danger); }
.block { display: block; }
.mb-1 { margin-bottom: var(--spacing-1); }
.mb-2 { margin-bottom: var(--spacing-2); }
.mb-4 { margin-bottom: var(--spacing-4); }

.mt-4 { margin-top: var(--spacing-4); }
.pt-4 { padding-top: var(--spacing-4); }
.pb-2 { padding-bottom: var(--spacing-2); }
.py-2 { padding-top: var(--spacing-2); padding-bottom: var(--spacing-2); }
.ml-2 { margin-left: var(--spacing-2); }

.text-left { text-align: left; }
.text-right { text-align: right; }

.border-b { border-bottom: 1px solid var(--color-border); }
.border-t { border-top: 1px solid var(--color-border); }

.w-full { width: 100%; }
.w-24 { width: 6rem; }
.flex-grow { flex-grow: 1; }
.font-semibold { font-weight: 600; }
.font-bold { font-weight: 700; }
.italic { font-style: italic; }

.input-field {
  padding: 8px 12px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
  font-family: inherit;
  font-size: 0.875rem;
  outline: none;
  transition: border-color 0.2s;
}

.input-field:focus {
  border-color: var(--color-brand-primary);
}</style>
