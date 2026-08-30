const TOKEN_KEY = 'timecd_token'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

async function request(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) }
  const token = getToken()
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch(`/api${path}`, { ...options, headers })
  let body = null
  try {
    body = await res.json()
  } catch {
    /* non-JSON response */
  }
  if (res.status === 401) {
    clearToken()
    if (location.hash !== '#/enter') location.hash = '#/enter'
  }
  if (!res.ok) {
    const err = new Error((body && body.error) || `请求失败(${res.status})`)
    err.status = res.status
    throw err
  }
  return body
}

export const api = {
  enter: (password) => request('/enter', { method: 'POST', body: JSON.stringify({ password }) }),
  me: () => request('/me'),
  list: () => request('/countdowns'),
  create: (payload) => request('/countdowns', { method: 'POST', body: JSON.stringify(payload) }),
  update: (id, payload) => request(`/countdowns/${id}`, { method: 'PUT', body: JSON.stringify(payload) }),
  remove: (id) => request(`/countdowns/${id}`, { method: 'DELETE' }),
}
