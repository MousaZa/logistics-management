<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Card from '../components/common/Card.vue'
import Button from '../components/common/Button.vue'
import DataTable from '../components/common/DataTable.vue'
import { getLocations, createLocation, getProducts, addProductsToLocation, getLocationContents } from '../services/inventory'
import type { Location, Product, ProductStock } from '../services/types'

const locations = ref<Location[]>([])
const products = ref<Product[]>([])
const loading = ref(true)
const showModal = ref(false)
const showAddProductModal = ref(false)
const showInventoryModal = ref(false)
const selectedLocationUUID = ref('')
const addProductData = ref({
  productUUID: '',
  quantity: 1
})

const inventoryData = ref<ProductStock[]>([])
const inventoryLoading = ref(false)
const inventoryColumns = [
  { key: 'name', label: 'Product Name' },
  { key: 'quantity', label: 'Quantity' },
  { key: 'weight', label: 'Unit Wt (kg)' }
]

const newLocation = ref<Location>({
  name: '',
  city: '',
  address: ''
})

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'city', label: 'City' },
  { key: 'address', label: 'Address' },
  { key: 'actions', label: 'Actions', align: 'right' as const }
]

const loadData = async () => {
  loading.value = true
  try {
    const [locRes, prodRes] = await Promise.all([
      getLocations(),
      getProducts()
    ])
    locations.value = locRes.locations || []
    products.value = prodRes.items || []
  } catch (err) {
    console.error('Failed to load data', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})

const submitLocation = async () => {
  if (!newLocation.value.name || !newLocation.value.city) return
  
  try {
    await createLocation(newLocation.value)
    showModal.value = false
    newLocation.value = { name: '', city: '', address: '' }
    loadData()
  } catch (err) {
    console.error('Failed to create location', err)
    alert('Failed to create location')
  }
}

const openAddProductModal = (uuid?: string) => {
  if (!uuid) return
  selectedLocationUUID.value = uuid
  showAddProductModal.value = true
}

const submitAddProduct = async () => {
  if (!selectedLocationUUID.value || !addProductData.value.productUUID || addProductData.value.quantity <= 0) return
  
  try {
    await addProductsToLocation(selectedLocationUUID.value, [
      { productUUID: addProductData.value.productUUID, quantity: addProductData.value.quantity }
    ])
    showAddProductModal.value = false
    addProductData.value = { productUUID: '', quantity: 1 }
    alert('Product added successfully')
  } catch (err) {
    console.error('Failed to add product', err)
    alert('Failed to add product to location')
  }
}

const openInventoryModal = async (uuid?: string) => {
  if (!uuid) return
  showInventoryModal.value = true
  inventoryLoading.value = true
  try {
    const res = await getLocationContents(uuid)
    inventoryData.value = res.items || []
  } catch (err) {
    console.error('Failed to load inventory', err)
    inventoryData.value = []
  } finally {
    inventoryLoading.value = false
  }
}
</script>

<template>
  <div class="locations-view">
    <div class="page-header flex justify-between items-center mb-6">
      <div>
        <h2 class="mb-2">Distribution Centers</h2>
        <p class="text-secondary">Manage your geographic logistics network.</p>
      </div>
      <Button @click="showModal = true">Add Location</Button>
    </div>

    <DataTable 
      :columns="columns" 
      :data="locations" 
      :loading="loading"
    >
      <template #cell-actions="{ row }">
        <div class="flex gap-2 justify-end">
          <Button variant="ghost" size="sm" @click="openAddProductModal(row.locationUUID)">Add Product</Button>
          <Button variant="ghost" size="sm" @click="openInventoryModal(row.locationUUID)">View Inventory</Button>
        </div>
      </template>
    </DataTable>

    <!-- Simple Modal Overlay -->
    <div v-if="showModal" class="modal-overlay">
      <Card class="modal-content">
        <h3 class="mb-4">Add New Location</h3>
        
        <div class="form-group">
          <label>Location Name</label>
          <input v-model="newLocation.name" type="text" class="form-input" placeholder="e.g. Central Hub" />
        </div>
        
        <div class="form-group">
          <label>City</label>
          <input v-model="newLocation.city" type="text" class="form-input" placeholder="e.g. New York" />
        </div>
        
        <div class="form-group">
          <label>Address</label>
          <input v-model="newLocation.address" type="text" class="form-input" placeholder="Full street address" />
        </div>

        <div class="modal-actions flex justify-between gap-4 mt-6">
          <Button variant="secondary" @click="showModal = false">Cancel</Button>
          <Button @click="submitLocation">Create Location</Button>
        </div>
      </Card>
    </div>

    <!-- Add Product Modal Overlay -->
    <div v-if="showAddProductModal" class="modal-overlay">
      <Card class="modal-content">
        <h3 class="mb-4">Add Product to Location</h3>
        
        <div class="form-group">
          <label>Select Product</label>
          <select v-model="addProductData.productUUID" class="form-input">
            <option disabled value="">Please select one</option>
            <option v-for="product in products" :key="product.productUUID" :value="product.productUUID">
              {{ product.name }}
            </option>
          </select>
        </div>
        
        <div class="form-group">
          <label>Quantity</label>
          <input v-model.number="addProductData.quantity" type="number" min="1" class="form-input" />
        </div>

        <div class="modal-actions flex justify-between gap-4 mt-6">
          <Button variant="secondary" @click="showAddProductModal = false">Cancel</Button>
          <Button @click="submitAddProduct">Add Product</Button>
        </div>
      </Card>
    </div>

    <!-- Inventory Modal Overlay -->
    <div v-if="showInventoryModal" class="modal-overlay">
      <Card class="modal-content inventory-modal">
        <h3 class="mb-4">Location Inventory</h3>
        
        <DataTable 
          :columns="inventoryColumns" 
          :data="inventoryData" 
          :loading="inventoryLoading"
        />

        <div class="modal-actions flex justify-end mt-6">
          <Button variant="secondary" @click="showInventoryModal = false">Close</Button>
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

.inventory-modal {
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
