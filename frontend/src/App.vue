<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const name = ref('')
const message = ref('')
const currentLanguage = ref<'vi' | 'en'>('vi')

const languages = computed(() => [
  { value: 'vi', label: 'Tiếng Việt' },
  { value: 'en', label: 'English' },
])

function changeLanguage() {
  locale.value = currentLanguage.value
}

async function greet() {
  const target = name.value.trim() || 'Developer'
  message.value = t('message.ready', { name: target })
}
</script>

<template>
  <main class="app-shell">
    <section class="card">
      <h1>{{ t('app.title') }}</h1>
      <p class="subtitle">{{ t('app.subtitle') }}</p>

      <div class="controls-row">
        <label>{{ t('form.language') }}</label>
        <select v-model="currentLanguage" @change="changeLanguage">
          <option v-for="item in languages" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </div>

      <div class="form-row">
        <input v-model="name" type="text" :placeholder="t('form.placeholder')" />
        <button @click="greet">{{ t('form.run') }}</button>
      </div>

      <p v-if="message" class="result">{{ message }}</p>
    </section>
  </main>
</template>
