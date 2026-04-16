<script setup lang="ts">
import { ref, provide } from 'vue'
import AppShell from './components/AppShell.vue'
import Overview from './pages/Overview.vue'
import Jobs from './pages/Jobs.vue'
import PostScraper from './pages/PostScraper.vue'
import Settings from './pages/Settings.vue'
import ActionLikeTest from './pages/ActionLikeTest.vue'
import ActionCommentTest from './pages/ActionCommentTest.vue'
import ActionPostTest from './pages/ActionPostTest.vue'
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
    <PostScraper v-if="currentRoute === 'post-scraper'" />
    <Settings v-if="currentRoute === 'settings'" />
    <ActionLikeTest v-if="currentRoute === 'action-like-test'" />
    <ActionCommentTest v-if="currentRoute === 'action-comment-test'" />
    <ActionPostTest v-if="currentRoute === 'action-post-test'" />
    <Accounts v-if="currentRoute === 'accounts'" />
  </AppShell>
</template>

<style>
/* Global styles have been moved to style.css */
</style>
