<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Card from '../components/common/Card.vue'
import { getLocations, getProducts } from '../services/inventory'
import { getOrders } from '../services/orders'

const stats = ref({
  locations: 0,
  products: 0,
  ordersPending: 0,
  ordersCompleted: 0
})

const loading = ref(true)

onMounted(async () => {
  try {
    const [locs, prods, ords] = await Promise.all([
      getLocations().catch(() => ({ locations: [] })),
      getProducts().catch(() => ({ items: [] })),
      getOrders().catch(() => ({ orders: [] }))
    ])

    stats.value.locations = locs.locations?.length || 0
    stats.value.products = prods.items?.length || 0
    
    const orders = ords.orders || []
    stats.value.ordersPending = orders.filter((o: any) => o.status === 'pending').length
    stats.value.ordersCompleted = orders.filter((o: any) => o.status === 'completed').length
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="dashboard">
    <div class="stats-grid">
      <Card class="stat-card">
        <div class="stat-icon locations-icon">L</div>
        <div class="stat-info">
          <h3>Total Locations</h3>
          <p class="stat-value" :class="{ 'skeleton-text': loading }">{{ loading ? '0' : stats.locations }}</p>
        </div>
      </Card>
      
      <Card class="stat-card">
        <div class="stat-icon products-icon">P</div>
        <div class="stat-info">
           <h3>Total Products</h3>
           <p class="stat-value" :class="{ 'skeleton-text': loading }">{{ loading ? '0' : stats.products }}</p>
        </div>
      </Card>

      <Card class="stat-card">
        <div class="stat-icon orders-pending-icon">O</div>
        <div class="stat-info">
           <h3>Pending Orders</h3>
           <p class="stat-value" :class="{ 'skeleton-text': loading }">{{ loading ? '0' : stats.ordersPending }}</p>
        </div>
      </Card>

      <Card class="stat-card">
        <div class="stat-icon orders-completed-icon">C</div>
        <div class="stat-info">
           <h3>Completed Orders</h3>
           <p class="stat-value" :class="{ 'skeleton-text': loading }">{{ loading ? '0' : stats.ordersCompleted }}</p>
        </div>
      </Card>
    </div>

    <div class="dashboard-content mt-8">
      <Card>
        <h2>Welcome to Logistics Management</h2>
        <p class="text-secondary mt-4">
          Manage your inventory across multiple locations, track product stocks, and process outgoing orders efficiently. 
          Use the sidebar navigation to explore different modules.
        </p>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: var(--spacing-6);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xl);
  font-weight: 700;
  color: white;
  flex-shrink: 0;
}

.locations-icon { background: linear-gradient(135deg, #8b5cf6, #6d28d9); }
.products-icon { background: linear-gradient(135deg, #10b981, #047857); }
.orders-pending-icon { background: linear-gradient(135deg, #f59e0b, #b45309); }
.orders-completed-icon { background: linear-gradient(135deg, #3b82f6, #1d4ed8); }

.stat-info h3 {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  font-weight: 500;
}

.stat-value {
  margin: 0;
  font-size: var(--font-size-2xl);
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.mt-4 { margin-top: var(--spacing-4); }
.mt-8 { margin-top: var(--spacing-8); }
.text-secondary { color: var(--color-text-secondary); line-height: 1.6; }

.skeleton-text {
  color: transparent !important;
  background-color: var(--color-border);
  border-radius: var(--radius-sm);
  animation: pulse 1.5s infinite ease-in-out;
}
</style>
