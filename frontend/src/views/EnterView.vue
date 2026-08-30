<template>
  <div class="enter-page">
    <div class="enter-hero">
      <div class="enter-logo">⏳</div>
      <h1>我的倒计时</h1>
      <p>为重要的时刻倒数<br />最近优先展示,精确到秒</p>
    </div>

    <div class="enter-card">
      <label class="field-label" for="pwd">输入你的密码</label>
      <div class="pwd-box">
        <input
          id="pwd"
          v-model="password"
          :type="show ? 'text' : 'password'"
          placeholder="包含大小写字母和数字"
          autocomplete="off"
          maxlength="64"
          @keyup.enter="submit"
        />
        <button type="button" class="pwd-eye" @click="show = !show">
          {{ show ? '🙈' : '👁️' }}
        </button>
      </div>

      <ul class="checks">
        <li :class="{ ok: rules.length }">长度至少 6 位</li>
        <li :class="{ ok: rules.upper }">包含大写字母 A-Z</li>
        <li :class="{ ok: rules.lower }">包含小写字母 a-z</li>
        <li :class="{ ok: rules.digit }">包含数字 0-9</li>
      </ul>

      <p class="hint">同一密码对应同一账号 · 首次输入自动创建,否则直接登录</p>

      <button class="btn btn-primary enter-btn" :disabled="loading" @click="submit">
        {{ loading ? '进入中…' : '进入' }}
      </button>
    </div>

    <transition name="fade">
      <div v-if="error" class="error-banner">{{ error }}</div>
    </transition>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, setToken } from '../api'

const router = useRouter()
const password = ref('')
const show = ref(false)
const loading = ref(false)
const error = ref('')

const rules = computed(() => ({
  length: password.value.length >= 6,
  upper: /[A-Z]/.test(password.value),
  lower: /[a-z]/.test(password.value),
  digit: /[0-9]/.test(password.value),
}))

const allPass = computed(() => Object.values(rules.value).every(Boolean))

async function submit() {
  error.value = ''
  if (!password.value) {
    error.value = '请输入密码'
    return
  }
  if (!allPass.value) {
    error.value = '密码太简单:必须同时包含大小写字母和数字,长度至少 6 位'
    return
  }
  loading.value = true
  try {
    const res = await api.enter(password.value)
    setToken(res.token)
    router.push('/')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.enter-page {
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 30px 24px 40px;
}

.enter-hero {
  text-align: center;
  margin-bottom: 30px;
}

.enter-logo {
  font-size: 52px;
  margin-bottom: 10px;
  filter: drop-shadow(0 6px 18px rgba(56, 189, 248, 0.45));
}

.enter-hero h1 {
  font-size: 26px;
  font-weight: 800;
  letter-spacing: 2px;
  background: linear-gradient(120deg, var(--accent), var(--accent-2));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.enter-hero p {
  margin-top: 10px;
  color: var(--text-dim);
  font-size: 13px;
  line-height: 1.7;
}

.enter-card {
  background: var(--card);
  border: 1px solid var(--card-border);
  border-radius: 20px;
  padding: 24px 20px 20px;
  box-shadow: 0 18px 50px rgba(0, 0, 0, 0.4);
}

.field-label {
  display: block;
  font-size: 13px;
  color: var(--text-dim);
  margin-bottom: 8px;
}

.pwd-box {
  display: flex;
  align-items: center;
  background: var(--bg-2);
  border: 1px solid var(--card-border);
  border-radius: 12px;
  padding: 0 12px;
}

.pwd-box:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.18);
}

.pwd-box input {
  flex: 1;
  height: 48px;
  background: transparent;
  border: none;
  outline: none;
  color: var(--text);
  font-size: 15px;
  letter-spacing: 1px;
}

.pwd-eye {
  font-size: 16px;
  padding: 6px;
}

.checks {
  list-style: none;
  margin: 14px 2px 4px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px 12px;
}

.checks li {
  font-size: 12px;
  color: var(--text-dim);
  padding-left: 18px;
  position: relative;
}

.checks li::before {
  content: '✕';
  position: absolute;
  left: 0;
  color: var(--danger);
  font-weight: 700;
}

.checks li.ok {
  color: var(--ok);
}

.checks li.ok::before {
  content: '✓';
  color: var(--ok);
}

.hint {
  margin: 12px 2px 16px;
  font-size: 11px;
  color: var(--text-dim);
  line-height: 1.6;
}

.enter-btn {
  width: 100%;
}

.error-banner {
  margin-top: 16px;
  text-align: center;
  color: #fff;
  background: rgba(248, 113, 113, 0.16);
  border: 1px solid rgba(248, 113, 113, 0.4);
  padding: 10px 14px;
  border-radius: 10px;
  font-size: 13px;
}
</style>
