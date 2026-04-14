<script setup lang="ts">
defineProps<{
  posts: any[]
}>()
</script>

<template>
  <div class="latest-posts">
    <div class="header">
      <h3 class="title">BÀI VIẾT MỚI NHẤT</h3>
    </div>
    <div class="post-list">
      <div v-for="post in posts" :key="post.id" class="post-item">
        <img :src="post.image" class="thumb" alt="thumb" />
        <div class="post-info">
          <h4 class="post-title">{{ post.title }}</h4>
          <div class="post-meta">{{ post.time }} • {{ post.engagement }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.latest-posts {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.header {
  margin-bottom: 16px;
}

.title {
  margin: 0;
  font-size: 12px;
  font-weight: 600;
  color: var(--c-text-title);
  letter-spacing: 0.5px;
}

.post-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.post-item {
  display: flex;
  gap: 12px;
  align-items: center;
}

.thumb {
  width: 60px;
  height: 48px;
  border-radius: var(--radius-sm);
  object-fit: cover;
  background: #f3f4f6;
  flex-shrink: 0;
}

.post-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.post-title {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-meta {
  font-size: 11px;
  color: var(--c-text-muted);
}
</style>
