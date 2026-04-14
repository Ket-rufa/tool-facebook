<script setup lang="ts">
defineProps<{
  columns: { key: string; label: string; width?: string }[]
  data: any[]
  selectedId?: string | number
}>()

const emit = defineEmits<{
  (e: 'row-click', row: any): void
}>()
</script>

<template>
  <div class="data-table-container">
    <table class="data-table">
      <thead>
        <tr>
          <th 
            v-for="col in columns" 
            :key="col.key"
            :style="{ width: col.width }"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr 
          v-for="row in data" 
          :key="row.id"
          :class="{ selected: row.id === selectedId }"
          @click="emit('row-click', row)"
        >
          <td v-for="col in columns" :key="col.key">
            <slot :name="col.key" :row="row" :value="row[col.key]">
              {{ row[col.key] }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.data-table-container {
  overflow-x: auto;
  border: 1px solid var(--c-border);
  border-radius: 0 0 var(--radius-md) var(--radius-md);
  background: var(--c-surface);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 16px 24px;
  text-align: left;
  border-bottom: 1px solid var(--c-border);
}

.data-table th {
  font-size: 12px;
  font-weight: 600;
  color: var(--c-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background-color: #f9fafb;
}

.data-table tbody tr {
  cursor: pointer;
  transition: background-color 0.2s;
}

.data-table tbody tr:hover {
  background-color: #f9fafb;
}

.data-table tbody tr.selected {
  background-color: #eff6ff;
}

.data-table tbody tr:last-child td {
  border-bottom: none;
}
</style>
