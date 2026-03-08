<script setup lang="ts">
import { computed } from 'vue';

export interface Column {
  key: string;
  label: string;
  sortable?: boolean;
  align?: 'left' | 'center' | 'right';
  format?: (val: any) => string;
}

const props = defineProps<{
  columns: Column[];
  data: any[];
  loading?: boolean;
}>();

const alignClass = (align?: string) => {
  if (align === 'center') return 'text-center';
  if (align === 'right') return 'text-right';
  return 'text-left';
};
</script>

<template>
  <div class="table-container glass-panel">
    <table class="data-table">
      <thead>
        <tr>
          <th 
            v-for="col in columns" 
            :key="col.key"
            :class="alignClass(col.align)"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading">
          <td :colspan="columns.length" class="text-center py-8">
            <span class="loading-state">Loading data...</span>
          </td>
        </tr>
        <tr v-else-if="!data || data.length === 0">
           <td :colspan="columns.length" class="text-center py-8 text-muted">
             No records found.
           </td>
        </tr>
        <tr v-else v-for="(row, idx) in data" :key="idx" class="table-row">
          <td 
            v-for="col in columns" 
            :key="col.key"
            :class="alignClass(col.align)"
          >
            <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
               {{ col.format ? col.format(row[col.key]) : row[col.key] }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.table-container {
  width: 100%;
  overflow-x: auto;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

th {
  padding: var(--spacing-4) var(--spacing-6);
  font-size: var(--font-size-xs);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
  color: var(--color-text-muted);
  border-bottom: 1px solid var(--color-border);
  background: rgba(15, 23, 42, 0.4); /* slight dark tint */
}

td {
  padding: var(--spacing-4) var(--spacing-6);
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.table-row {
  transition: background-color var(--transition-fast);
}

.table-row:hover {
  background-color: var(--color-bg-surface-hover);
}

.table-row:last-child td {
  border-bottom: none;
}

.text-left { text-align: left; }
.text-center { text-align: center; }
.text-right { text-align: right; }

.py-8 {
  padding-top: var(--spacing-8);
  padding-bottom: var(--spacing-8);
}

.text-muted {
  color: var(--color-text-muted);
}

.loading-state {
  display: inline-block;
  color: var(--color-brand-primary);
  animation: pulse 1.5s infinite ease-in-out;
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}
</style>
