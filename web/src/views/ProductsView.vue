<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Card from '../components/common/Card.vue'
import Button from '../components/common/Button.vue'
import DataTable from '../components/common/DataTable.vue'
import { getProducts, createProduct, getProductLocations } from '../services/inventory'
import type { Product, ProductLocationInventory } from '../services/types'

const products = ref<Product[]>([])
const loading = ref(true)
const showModal = ref(false)
const showLocationsModal = ref(false)

const locationsData = ref<ProductLocationInventory[]>([])
const locationsLoading = ref(false)
const locationsColumns = [
  { key: 'name', label: 'Location Name' },
  { key: 'city', label: 'City' },
  { key: 'address', label: 'Address' },
  { key: 'availableQuantity', label: 'Available' },
  { key: 'reservedQuantity', label: 'Reserved' },
  { key: 'damagedQuantity', label: 'Damaged' }
]

const newProduct = ref<Product>({
  name: '',
  price: 0,
  weight: 0
})

const columns = [
  { key: 'name', label: 'Product Name' },
  { key: 'price', label: 'Unit Price ($)', format: (val: number) => `$ ${val.toFixed(2)}` },
  { key: 'weight', label: 'Weight (kg)', format: (val: number) => `${val} kg` },
  { key: 'actions', label: 'Actions', align: 'right' as const }
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await getProducts()
    products.value = res.items || []
  } catch (err) {
    console.error('Failed to load products', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})

const submitProduct = async () => {
  if (!newProduct.value.name || newProduct.value.price <= 0) return
  
  try {
    await createProduct(newProduct.value)
    showModal.value = false
    newProduct.value = { name: '', price: 0, weight: 0 }
    loadData()
  } catch (err) {
    console.error('Failed to create product', err)
    alert('Failed to create product')
  }
}

const openLocationsModal = async (uuid?: string) => {
  if (!uuid) return
  showLocationsModal.value = true
  locationsLoading.value = true
  try {
    const res = await getProductLocations(uuid)
    locationsData.value = res.locations || []
  } catch (err) {
    console.error('Failed to load locations', err)
    locationsData.value = []
  } finally {
    locationsLoading.value = false
  }
}
</script>

<template>
  <div class="products-view">
    <div class="page-header flex justify-between items-center mb-6">
      <div>
        <h2 class="mb-2">Product Catalog</h2>
        <p class="text-secondary">Manage items available in the logistics network.</p>
      </div>
      <Button @click="showModal = true">Add Product</Button>
    </div>

    <DataTable 
      :columns="columns" 
      :data="products" 
      :loading="loading"
    >
      <template #cell-actions="{ row }">
        <div class="flex gap-2 justify-end">
          <Button variant="ghost" size="sm" @click="openLocationsModal(row.productUUID)">View Locations</Button>
          <Button variant="ghost" size="sm">Edit</Button>
        </div>
      </template>
    </DataTable>

    <div v-if="showModal" class="modal-overlay">
      <Card class="modal-content">
        <h3 class="mb-4">Add New Product</h3>
        
        <div class="form-group">
          <label>Product Name</label>
          <input v-model="newProduct.name" type="text" class="form-input" placeholder="e.g. Industrial Widget" />
        </div>
        
        <div class="form-group">
          <label>Price ($)</label>
          <input v-model.number="newProduct.price" type="number" step="0.01" class="form-input" />
        </div>
        
        <div class="form-group">
          <label>Weight (kg)</label>
          <input v-model.number="newProduct.weight" type="number" step="0.1" class="form-input" />
        </div>

        <div class="modal-actions flex justify-between gap-4 mt-6">
          <Button variant="secondary" @click="showModal = false">Cancel</Button>
          <Button @click="submitProduct">Save Product</Button>
        </div>
      </Card>
    </div>

    <!-- Locations Modal Overlay -->
    <div v-if="showLocationsModal" class="modal-overlay">
      <Card class="modal-content locations-modal">
        <h3 class="mb-4">Product Locations</h3>
        
        <DataTable 
          :columns="locationsColumns" 
          :data="locationsData" 
          :loading="locationsLoading"
        />

        <div class="modal-actions flex justify-end mt-6">
          <Button variant="secondary" @click="showLocationsModal = false">Close</Button>
        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.mb-2 { margin-bottom: var(--spacing-2); }
.mb-4 { margin-bottom: var(--spacing-4); }
.mb-6 { margin-bottom: var(--spacing-6); }
.mt-6 { margin-top: var(--spacing-6); }
.text-secondary { color: var(--color-text-secondary); margin: 0; }
.justify-end { justify-content: flex-end; }

.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(15, 23, 42, 0.8);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}

.modal-content {
  width: 100%;
  max-width: 500px;
  border: 1px solid var(--color-border);
}

.locations-modal {
  max-width: 600px;
}

.form-group {
  margin-bottom: var(--spacing-4);
}

.form-group label {
  display: block;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-2);
}

.form-input {
  width: 100%;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
  padding: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-md);
  font-family: var(--font-family-sans);
  transition: border-color var(--transition-fast);
  outline: none;
}

.form-input:focus {
  border-color: var(--color-brand-primary);
}

.modal-actions {
  justify-content: flex-end;
}
</style>
