import FingerprintJS from '@fingerprintjs/fingerprintjs'

export const getFingerprint = async (): Promise<string> => {
  const storedFingerprint = localStorage.getItem('fingerprint')
  if (storedFingerprint) {
    return storedFingerprint
  }

  const fp = await FingerprintJS.load()
  const result = await fp.get()
  const visitorId = result.visitorId
  localStorage.setItem('fingerprint', visitorId)
  return visitorId
}

export const getStoredFingerprint = (): string | null => {
  return localStorage.getItem('fingerprint')
}
