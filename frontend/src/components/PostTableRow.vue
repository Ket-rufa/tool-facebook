<script setup lang="ts">
import { inject } from 'vue'
import PostStatusBadge from './PostStatusBadge.vue'

defineProps<{
  post: any
  isChecked: boolean
}>()

const emit = defineEmits(['toggle'])
const navigate = inject('navigate') as (route: string) => void
</script>

<template>
  <tr class="post-row" :class="{ 'selected': isChecked }">
    <td class="col-check">
      <input 
        type="checkbox" 
        class="custom-checkbox" 
        :checked="isChecked" 
        @change="emit('toggle')"
      />
    </td>
    <td class="col-id"><a href="#" @click.prevent="navigate('post-detail')" class="id-link">{{ post.id }}</a></td>
    <td class="col-source">
      <div class="source-flex">
        <img v-if="post.source.avatarImg" :src="post.source.avatarImg" class="src-avatar" alt="avatar" />
        <div v-else class="src-avatar-text" :style="{ backgroundColor: post.source.bgColor, color: post.source.textColor }">
          {{ post.source.avatar }}
        </div>
        <div class="src-info">
          <div class="src-name">{{ post.source.name }}</div>
          <div class="src-type">{{ post.source.type }}</div>
        </div>
      </div>
    </td>
    <td class="col-content">
      <p class="content-text">{{ post.content }}</p>
    </td>
    <td class="col-time">
      <div class="time-text">{{ post.time }}</div>
    </td>
    <td class="col-engagement">
      <div class="eng-row">
        <span class="eng-item">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" stroke="none"><path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path></svg>
          {{ post.engagement.likes }}
        </span>
        <span class="eng-item">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" stroke="none"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
          {{ post.engagement.comments }}
        </span>
      </div>
    </td>
    <td class="col-media">
      <div v-if="post.media.count > 0" class="media-box">
        <span>{{ post.media.count }}</span>
        <svg v-if="post.media.type === 'video'" width="12" height="12" viewBox="0 0 24 24" fill="currentColor" stroke="none"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
        <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="currentColor" stroke="none"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
      </div>
      <div v-else class="media-none">&mdash;</div>
    </td>
    <td class="col-status">
      <PostStatusBadge :status="post.status" />
    </td>
  </tr>
</template>

<style scoped>
.post-row {
  border-bottom: 1px solid var(--c-border);
  transition: background-color 0.2s;
  background: var(--c-surface);
}

.post-row:hover {
  background: #f9fafb;
}

.post-row.selected {
  background: #eff6ff;
}

td {
  padding: 16px 24px;
  vertical-align: top;
}

.col-check { width: 40px; padding-right: 8px; }
.custom-checkbox {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: var(--c-primary);
}

.col-id { font-size: 13px; color: var(--c-text-muted); font-family: monospace; }
.id-link { color: var(--c-primary); text-decoration: none; font-weight: 700; }
.id-link:hover { text-decoration: underline; }

.source-flex { display: flex; gap: 12px; align-items: flex-start; }
.src-avatar { width: 32px; height: 32px; border-radius: 4px; object-fit: cover; }
.src-avatar-text {
  width: 32px; height: 32px; border-radius: 4px; display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 14px;
}
.src-info { display: flex; flex-direction: column; gap: 4px; }
.src-name { font-size: 13px; font-weight: 700; color: var(--c-text-title); }
.src-type { font-size: 11px; color: var(--c-text-muted); }

.col-content { width: 30%; }
.content-text { margin: 0; font-size: 13px; color: var(--c-text-main); line-height: 1.5; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }

.col-time { width: 100px; }
.time-text { font-size: 12px; color: var(--c-text-muted); white-space: pre-line; line-height: 1.4; }

.col-engagement { min-width: 120px; }
.eng-row { display: flex; gap: 12px; }
.eng-item { display: flex; align-items: center; gap: 6px; font-size: 12px; color: var(--c-text-main); font-weight: 500; }
.eng-item svg { color: var(--c-primary); } /* Tinting icon with primary, similar to fb like */

.media-box {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
  color: var(--c-text-main);
}
.media-none { color: var(--c-text-muted); }
</style>
