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
  let timeout: ReturnType<typeof setTimeout>
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

// 时区相关工具函数
export const formatDateTimeForPicker = (date: Date): string => {
  // 将本地时间转换为YYYY-MM-DD HH:mm:ss格式（用于时间选择器）
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

export const formatDateTimeForAPI = (date: Date): string => {
  // 将本地时间转换为带时区的ISO格式字符串（用于API请求）
  return date.toISOString()
}

export const parseDateTimeFromLocal = (dateTimeString: string): Date => {
  // 解析本地时间字符串并创建Date对象
  // 格式：YYYY-MM-DD HH:mm:ss
  const [datePart, timePart] = dateTimeString.split(' ')
  const [year, month, day] = datePart.split('-').map(Number)
  const [hours, minutes, seconds] = timePart.split(':').map(Number)
  
  // 创建Date对象（使用本地时间）
  const date = new Date(year, month - 1, day, hours, minutes, seconds || 0)
  
  return date
}

export const formatDateTimeForDisplay = (dateString: string | null): string => {
  if (!dateString) return ''
  
  // 解析时间字符串，假设是UTC时间
  const date = new Date(dateString)
  
  if (isNaN(date.getTime())) {
    return ''
  }
  
  // 转换为本地时间字符串
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

export const getDefaultScheduleTime = (): string => {
  // 返回1小时后的本地时间
  const now = new Date()
  now.setHours(now.getHours() + 1)
  return formatDateTimeForPicker(now)
}

export const isValidScheduleTime = (timeString: string): boolean => {
  if (!timeString) return false
  
  const selectedTime = new Date(timeString)
  const now = new Date()
  
  // 选择的时间必须比当前时间晚1分钟
  return selectedTime.getTime() > now.getTime() + 60000
}