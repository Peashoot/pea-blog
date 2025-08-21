import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import zhCN from './locales/zh-CN.json'

// 获取操作系统语言
const getSystemLanguage = (): string => {
  const systemLang = navigator.language || navigator.languages?.[0] || 'en'
  
  // 检查系统语言是否为中文（包括各种变体）
  if (systemLang.startsWith('zh')) {
    // 如果系统语言是中文简体，返回 zh-CN
    // 其他中文变体（如 zh-TW, zh-HK）也返回 zh-CN 作为默认
    return 'zh-CN'
  }
  
  // 默认返回英文
  return 'en'
}

// 初始化语言设置
const initializeLanguage = (): string => {
  // 检查 localStorage 中是否已有语言设置
  const savedLang = localStorage.getItem('locale')
  if (savedLang) {
    return savedLang
  }
  
  // 如果没有保存的语言设置，检测系统语言并保存到 localStorage
  const systemLang = getSystemLanguage()
  localStorage.setItem('locale', systemLang)
  return systemLang
}

const i18n = createI18n({
  legacy: false,
  locale: initializeLanguage(),
  fallbackLocale: 'en',
  messages: {
    en,
    'zh-CN': zhCN,
  },
})

export default i18n
