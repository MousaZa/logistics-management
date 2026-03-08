<script setup lang="ts">
// Simple wrapping container to maintain consistent panel styling
</script>

<template>
  <div class="card glass-panel">
    <slot></slot>
  </div>
</template>

<style scoped>
.card {
  border-radius: var(--radius-lg);
  padding: var(--spacing-6);
  box-shadow: var(--shadow-md);
  transition: transform var(--transition-fast), box-shadow var(--transition-fast);
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}
</style>
