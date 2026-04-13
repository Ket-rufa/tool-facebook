import { createI18n } from 'vue-i18n'
import en from './locales/en'
import vi from './locales/vi'

export const messages = {
  en,
  vi,
}

const i18n = createI18n({
  legacy: false,
  locale: 'vi',
  fallbackLocale: 'en',
  messages,
})

export default i18n
