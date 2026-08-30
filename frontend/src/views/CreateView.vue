<template>
  <div class="create-page">
    <header class="page-header">
      <button class="icon-btn" @click="$router.back()">←</button>
      <div class="page-title">新建倒计时</div>
      <span style="width: 38px"></span>
    </header>

    <main class="form-wrap">
      <label class="field-label">标题</label>
      <input
        v-model="title"
        class="text-input"
        type="text"
        maxlength="100"
        placeholder="例如:元旦、生日、项目上线…"
      />

      <label class="field-label">目标时间(精确到秒)</label>
      <input v-model="target" class="text-input" type="datetime-local" step="1" :min="minLocal" />

      <div v-if="title && target" class="preview">
        <p>预览</p>
        <div class="preview-main">
          <span class="p-title">距离「{{ title }}」还有</span>
          <span class="p-time">{{ previewText }}</span>
        </div>
      </div>

      <button class="btn btn-primary save-btn" :disabled="saving || !canSave" @click="save">
        {{ saving ? '保存中…' : '保存' }}
      </button>
    </main>

    <transition name="fade">
      <div v-if="error" class="error-banner">{{ error }}</div>
    </transition>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const title = ref('')
const target = ref('')
const saving = ref(false)
const error = ref('')

const now = new Date()
const minLocal = computed(() => {
  const d = new Date(now.getTime() + 60000)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}T${p(d.getHours())}:${p(d.getMinutes())}`
})

const canSave = computed(() => title.value.trim().length > 0 && !!target.value)

const previewText = computed(() => {
  const t = new Date(target.value).getTime()
  const diff = t - Date.now()
  if (diff <= 0) return '已到时间'
  const s = Math.floor(diff / 1000)
  const d = Math.floor(s / 86400)
  const h = Math.floor((s % 86400) / 3600)
  const m = Math.floor((s % 3600) / 60)
  const sec = s % 60
  const parts = []
  if (d) parts.push(`${d}天`)
  if (h) parts.push(`${h}时`)
  if (m) parts.push(`${m}分`)
  parts.push(`${sec}秒`)
  return parts.join(' ')
})

async function save() {
  error.value = ''
  const t = new Date(target.value)
  if (isNaN(t.getTime())) {
    error.value = '请选择有效的目标时间'
    return
  }
  if (t.getTime() <= Date.now()) {
    error.value = '目标时间必须晚于当前时间'
    return
  }
  saving.value = true
  try {
    await api.create({ title: title.value.trim(), target_time: t.toISOString() })
    router.push('/')
  } catch (e) {
    error.value = e.message
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.create-page {
  min-height: 100dvh;
}

.form-wrap {
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-label {
  font-size: 13px;
  color: var(--text-dim);
  margin: 8px 2px 2px;
}

.text-input {
  height: 50px;
  background: var(--bg-2);
  border: 1px solid var(--card-border);
  border-radius: 12px;
  padding: 0 14px;
  color: var(--text);
  font-size: 15px;
  outline: none;
}

.text-input:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.18);
}

.preview {
  margin-top: 18px;
  background: linear-gradient(160deg, var(--card), #131c30);
  border: 1px dashed var(--card-border);
  border-radius: var(--radius);
  padding: 16px;
}

.preview p {
  font-size: 11px;
  color: var(--text-dim);
  margin-bottom: 10px;
}

.preview-main {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.p-title {
  font-size: 14px;
  word-break: break-all;
}

.p-time {
  font-family: Consolas, monospace;
  font-size: 22px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
  color: var(--accent);
}

.save-btn {
  margin-top: 22px;
  width: 100%;
}

.error-banner {
  margin: 0 20px;
  text-align: center;
  color: #fff;
  background: rgba(248, 113, 113, 0.16);
  border: 1px solid rgba(248, 113, 113, 0.4);
  padding: 10px 14px;
  border-radius: 10px;
  font-size: 13px;
}
</style>
