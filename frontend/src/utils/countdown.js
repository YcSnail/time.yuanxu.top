/**
 * Split a remaining duration (ms) into its components.
 * Non-negative; a finished countdown returns all zeros and ended: true.
 */
export function splitRemaining(ms) {
  const total = Math.max(0, Math.floor(ms))
  const days = Math.floor(total / 86400000)
  const hours = Math.floor((total % 86400000) / 3600000)
  const minutes = Math.floor((total % 3600000) / 60000)
  const seconds = Math.floor((total % 60000) / 1000)
  const millis = total % 1000
  return { days, hours, minutes, seconds, millis, ended: ms <= 0 }
}

export function pad2(n) {
  return String(n).padStart(2, '0')
}

export function pad3(n) {
  return String(n).padStart(3, '0')
}
