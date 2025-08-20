import { useI18n } from 'vue-i18n'

export const formatDate = (dateString: string | undefined | null): string => {
  const { t, d } = useI18n()

  if (!dateString) return t('time.date_not_available')
  
  // Handle "YYYY-MM-DD HH:MM:SS" format by replacing space with T
  const formattedDateString = dateString.replace(' ', 'T');
  const date = new Date(formattedDateString);

  if (isNaN(date.getTime())) {
    // Try parsing without replacement if the first attempt fails
    const originalDate = new Date(dateString);
    if (isNaN(originalDate.getTime())) {
      return t('time.date_not_available'); // Return placeholder if still invalid
    }
    date.setTime(originalDate.getTime());
  }

  const now = new Date()
  const diffTime = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    const diffHours = Math.floor(diffTime / (1000 * 60 * 60))
    if (diffHours === 0) {
      const diffMinutes = Math.floor(diffTime / (1000 * 60))
      return diffMinutes <= 0 ? t('time.just_now') : t('time.minutes_ago', { n: diffMinutes })
    }
    return t('time.hours_ago', { n: diffHours })
  } else if (diffDays === 1) {
    return t('time.yesterday')
  } else if (diffDays < 7) {
    return t('time.days_ago', { n: diffDays })
  } else {
    return d(date, 'long')
  }
}

export const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text
  return text.slice(0, maxLength) + '...'
}

export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout
  return (...args: Parameters<T>) => {
    clearTimeout(timeout)
    timeout = setTimeout(() => func(...args), wait)
  }
}

export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

export const generateId = (): string => {
  return Math.random().toString(36).substr(2, 9)
}

export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch (err) {
    console.error('Failed to copy text: ', err)
    return false
  }
}

export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

export const extractTags = (content: string): string[] => {
  const tagRegex = /#(\w+)/g
  const matches = content.match(tagRegex)
  return matches ? matches.map(tag => tag.slice(1)) : []
}