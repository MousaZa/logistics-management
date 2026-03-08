<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()

const navItems = [
  { name: 'Dashboard', path: '/' },
  { name: 'Locations', path: '/locations' },
  { name: 'Products', path: '/products' },
  { name: 'Orders', path: '/orders' },
]

const isActive = (path: string) => {
  if (path === '/' && route.path !== '/') return false
  return route.path.startsWith(path)
}
</script>

<template>
  <div class="app-layout">
    <!-- Sidebar -->
    <aside class="sidebar glass-panel">
      <div class="sidebar-header">
        <div class="logo">LM</div>
        <h2 class="brand">Logistics</h2>
      </div>
      
      <nav class="sidebar-nav">
        <RouterLink 
          v-for="item in navItems" 
          :key="item.path" 
          :to="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path) }"
        >
          {{ item.name }}
        </RouterLink>
      </nav>
      
      <div class="sidebar-footer">
         <p class="user-info">Mousa Zeydan</p>
         <p class="user-email">mous.zeydan@gmail.com</p>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="main-content">
      <!-- Top header / actions could go here -->
      <header class="top-header">
        <div class="greeting">
          <h1>{{ route.name ? route.name.toString().charAt(0).toUpperCase() + route.name.toString().slice(1) : 'Dashboard' }}</h1>
        </div>
      </header>
      
      <div class="page-container">
        <slot></slot>
      </div>
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  width: 100%;
}

.sidebar {
  width: 260px;
  display: flex;
  flex-direction: column;
  padding: var(--spacing-6) 0;
  border-right: 1px solid var(--color-border);
  z-index: 10;
}

.sidebar-header {
  display: flex;
  align-items: center;
  padding: 0 var(--spacing-6);
  margin-bottom: var(--spacing-8);
}

.logo {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, var(--color-brand-primary), var(--color-brand-primary-hover));
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: var(--font-size-lg);
  color: white;
  margin-right: var(--spacing-3);
  box-shadow: 0 4px 15px rgba(59, 130, 246, 0.4);
}

.brand {
  margin: 0;
  font-size: var(--font-size-xl);
  font-weight: 600;
  letter-spacing: -0.5px;
  background: linear-gradient(90deg, #fff, var(--color-text-secondary));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
  padding: 0 var(--spacing-4);
}

.nav-item {
  padding: var(--spacing-3) var(--spacing-4);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  font-weight: 500;
  transition: all var(--transition-fast);
  display: flex;
  align-items: center;
}

.nav-item:hover {
  background-color: var(--color-bg-surface-hover);
  color: var(--color-text-primary);
  transform: translateX(4px);
}

.nav-item.active {
  background-color: rgba(59, 130, 246, 0.15);
  color: var(--color-brand-primary);
  border-right: 3px solid var(--color-brand-primary);
}

.sidebar-footer {
  padding: var(--spacing-6);
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.user-info {
  margin: 0;
  font-weight: 600;
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
}

.user-email {
  margin: 0;
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden; /* prevents layout shift, child handles scroll */
}

.top-header {
  height: 80px;
  display: flex;
  align-items: center;
  padding: 0 var(--spacing-8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.greeting h1 {
  margin: 0;
  font-size: var(--font-size-2xl);
  font-weight: 600;
}

.page-container {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-8);
}
</style>
