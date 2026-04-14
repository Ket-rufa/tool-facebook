<script setup lang="ts">
import { ref, provide } from 'vue'
import AppShell from './components/AppShell.vue'
import Overview from './pages/Overview.vue'
import Jobs from './pages/Jobs.vue'
import Logs from './pages/Logs.vue'
import Profiles from './pages/Profiles.vue'
import Posts from './pages/Posts.vue'
import PostDetail from './pages/PostDetail.vue'
import ActionCommentTest from './pages/ActionCommentTest.vue'
import ActionPostTest from './pages/ActionPostTest.vue'
import Integrations from './pages/Integrations.vue'
import Settings from './pages/Settings.vue'
import ActionLikeTest from './pages/ActionLikeTest.vue'
import Accounts from './pages/Accounts.vue'

// Basic routing state
const currentRoute = ref('overview')

// Provide a way for SidebarNav to change the route
provide('currentRoute', currentRoute)
provide('navigate', (route: string) => {
  currentRoute.value = route
})
</script>

<template>
  <AppShell>
    <Overview v-if="currentRoute === 'overview'" />
    <Jobs v-if="currentRoute === 'tasks'" />
    <Logs v-if="currentRoute === 'logs'" />
    <Profiles v-if="currentRoute === 'profiles'" />
    <Posts v-if="currentRoute === 'posts'" />
    <PostDetail v-if="currentRoute === 'post-detail'" />
    <ActionCommentTest v-if="currentRoute === 'action-comment-test'" />
    <ActionPostTest v-if="currentRoute === 'action-post-test'" />
    <Integrations v-if="currentRoute === 'integrations'" />
    <Settings v-if="currentRoute === 'settings'" />
    <ActionLikeTest v-if="currentRoute === 'action-like-test'" />
    <Accounts v-if="currentRoute === 'accounts'" />
  </AppShell>
</template>

<style>
/* Global styles have been moved to style.css */
</style>
