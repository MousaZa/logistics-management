<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
}>();

const variantClass = computed(() => {
  switch (props.variant) {
    case 'secondary': return 'btn-secondary';
    case 'danger': return 'btn-danger';
    case 'ghost': return 'btn-ghost';
    case 'primary':
    default: return 'btn-primary';
  }
});

const sizeClass = computed(() => {
  switch (props.size) {
    case 'sm': return 'btn-sm';
    case 'lg': return 'btn-lg';
    case 'md':
    default: return 'btn-md';
  }
});
</script>

<template>
  <button 
    class="btn-base" 
    :class="[variantClass, sizeClass]"
    :disabled="disabled"
  >
    <slot></slot>
  </button>
</template>

<style scoped>
.btn-base {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-family-sans);
  font-weight: 500;
  border-radius: var(--radius-md);
  border: none;
  cursor: pointer;
  transition: all var(--transition-normal);
  outline: none;
}

.btn-base:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-base:active:not(:disabled) {
  transform: scale(0.97);
}

/* Variants */
.btn-primary {
  background-color: var(--color-brand-primary);
  color: white;
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.3);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--color-brand-primary-hover);
  box-shadow: 0 0 15px rgba(59, 130, 246, 0.5);
}

.btn-secondary {
  background-color: var(--color-bg-surface-hover);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover:not(:disabled) {
  background-color: var(--color-border);
  color: white;
}

.btn-danger {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--color-accent-danger);
  border: 1px solid rgba(239, 68, 68, 0.5);
}

.btn-danger:hover:not(:disabled) {
  background-color: var(--color-accent-danger);
  color: white;
}

.btn-ghost {
  background-color: transparent;
  color: var(--color-brand-primary);
}

.btn-ghost:hover:not(:disabled) {
  background-color: rgba(59, 130, 246, 0.1);
}

/* Sizes */
.btn-sm {
  padding: var(--spacing-1) var(--spacing-3);
  font-size: var(--font-size-xs);
}

.btn-md {
  padding: var(--spacing-2) var(--spacing-4);
  font-size: var(--font-size-sm);
}

.btn-lg {
  padding: var(--spacing-3) var(--spacing-6);
  font-size: var(--font-size-base);
}
</style>
