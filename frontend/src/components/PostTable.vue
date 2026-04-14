<script setup lang="ts">
import PostTableRow from './PostTableRow.vue'

defineProps<{
  posts: any[]
  selectedIds: string[]
}>()

const emit = defineEmits(['toggleSelect'])
</script>

<template>
  <div class="post-table-container">
    <table class="post-table">
      <thead>
        <tr>
          <th class="col-check"><input type="checkbox" disabled /></th>
          <th>ID BÀI VIẾT</th>
          <th>NGUỒN (PAGE/PROFILE)</th>
          <th>NỘI DUNG TÓM TẮT</th>
          <th>THỜI GIAN ĐĂNG</th>
          <th>TƯƠNG TÁC</th>
          <th>MEDIA</th>
          <th>TRẠNG THÁI</th>
        </tr>
      </thead>
      <tbody>
        <PostTableRow 
          v-for="post in posts" 
          :key="post.id"
          :post="post"
          :is-checked="selectedIds.includes(post.id)"
          @toggle="emit('toggleSelect', post.id)"
        />
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.post-table-container {
  overflow-x: auto;
  border: 1px solid var(--c-border);
  background: var(--c-surface);
  /* border handling so it connects seamlessly to filter logic above */
}

.post-table {
  width: 100%;
  border-collapse: collapse;
}

.post-table th {
  padding: 16px 24px;
  text-align: left;
  font-size: 11px;
  font-weight: 700;
  color: var(--c-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background-color: #f9fafb;
  border-bottom: 1px solid var(--c-border);
}

.col-check { width: 40px; padding-right: 8px !important; }
</style>
