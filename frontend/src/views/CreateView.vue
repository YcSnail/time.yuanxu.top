<template>
  <div class="create-page">
    <header class="page-header">
      <button class="icon-btn" @click="$router.back()">←</button>
      <div class="page-title">新建倒计时</div>
      <span style="width: 38px"></span>
    </header>

    <main class="form-wrap">
      <van-field
        v-model="title"
        label="标题"
        placeholder="例如:元旦、生日、项目上线…"
        maxlength="100"
        clearable
        class="cd-field"
      />

      <van-field
        :model-value="dateText"
        readonly
        is-link
        label="日期"
        placeholder="选择日期"
        @click="showDate = true"
        class="cd-field"
      />

      <van-field
        :model-value="timeText"
        readonly
        is-link
        label="时间"
        placeholder="选择时间"
        @click="showTime = true"
        class="cd-field"
      />

      <div v-if="title && dateText && timeText" class="preview">
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

    <!-- 日期选择 -->
    <van-popup v-model:show="showDate" position="bottom" round>
      <van-date-picker
        v-model="dateValue"
        :min-date="minDate"
        :max-date="maxDate"
        title="选择日期"
        @confirm="onDateConfirm"
        @cancel="showDate = false"
      />
    </van-popup>

    <!-- 时间选择(精确到秒) -->
    <van-popup v-model:show="showTime" position="bottom" round>
      <van-time-picker
        v-model="timeValue"
        title="选择时间"
        @confirm="onTimeConfirm"
        @cancel="showTime = false"
      />
    </van-popup>

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
const saving = ref(false)
const error = ref('')

/* ---- 日期/时间选择 ---- */
const showDate = ref(false)
const showTime = ref(false)

const minDate = new Date()
const maxDate = new Date(2099, 11, 31)

// Vant DatePicker 的值是 [年, 月, 日] 数组
const today = new Date()
const dateValue = ref([String(today.getFullYear()), String(today.getMonth() + 1).padStart(2, '0'), String(today.getDate()).padStart(2, '0')])
// Vant TimePicker 的值是 "HH:mm:ss" 字符串
const timeValue = ref(
  `${String((today.getHours() + 1) % 24).padStart(2, '0')}:${String(today.getMinutes()).padStart(2, '0')}:00`,
)

const dateText = computed(() => dateValue.value.join('-'))
const timeText = computed(() => timeValue.value)

function onDateConfirm({ selectedValues }) {
  dateValue.value = selectedValues
  showDate.value = false
}

function onTimeConfirm({ selectedValues }) {
  // selectedValues: ['时','分','秒']
  timeValue.value = selectedValues.join(':')
  showTime.value = false
}

/* ---- 提交 ---- */
const canSave = computed(() => title.value.trim().length > 0 && !!dateText.value && !!timeText.value)

const previewText = computed(() => {
  const t = new Date(`${dateText.value} ${timeText.value}`).getTime()
  if (isNaN(t)) return ''
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
  const t = new Date(`${dateText.value} ${timeText.value}`)
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
  padding: 18px 16px;
}

.cd-field {
  background: var(--bg-2);
  border: 1px solid var(--card-border);
  border-radius: 12px;
  margin-bottom: 12px;
  overflow: hidden;
}

.cd-field :deep(.van-field__label) {
  color: var(--text-dim);
  width: 64px;
}

.cd-field :deep(.van-field__control) {
  color: var(--text);
}

.preview {
  margin-top: 8px;
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
  margin: 0 16px;
  text-align: center;
  color: #fff;
  background: rgba(248, 113, 113, 0.16);
  border: 1px solid rgba(248, 113, 113, 0.4);
  padding: 10px 14px;
  border-radius: 10px;
  font-size: 13px;
}
</style>
