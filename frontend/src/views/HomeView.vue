<template>
  <div class="home-page">
    <header class="page-header">
      <div class="page-title">
        ⏳ 我的倒计时
        <small>最近的排在最前</small>
      </div>
      <button class="logout-btn" @click="logout">退出</button>
    </header>

    <main class="list-wrap">
      <transition-group name="list" tag="div" class="cd-list">
        <CountdownCard
          v-for="item in sorted"
          :key="item.id"
          :item="item"
          :now-ms="nowMs"
          @edit="edit(item)"
          @delete="remove(item)"
        />
      </transition-group>

      <div v-if="!loading && sorted.length === 0" class="empty-state">
        <div class="big">🕰️</div>
        <p>还没有倒计时</p>
        <p style="margin-top: 6px; font-size: 12px">点击下方 + 添加第一个吧</p>
      </div>
    </main>

    <button class="fab" @click="$router.push('/create')">＋</button>

    <transition name="fade">
      <div v-if="toast" class="toast">{{ toast }}</div>
    </transition>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, clearToken } from '../api'
import CountdownCard from '../components/CountdownCard.vue'

const router = useRouter()
const items = ref([])
const loading = ref(true)
const toast = ref('')
const nowMs = ref(Date.now())

let timer = 0
function tick() {
  nowMs.value = Date.now()
  timer = setInterval(() => {
    nowMs.value = Date.now()
  }, 1000)
}

/* 最近的排在最前;已结束的沉底 */
const sorted = computed(() => {
  const now = nowMs.value
  const upcoming = items.value
    .filter((it) => new Date(it.target_time).getTime() > now)
    .sort((a, b) => new Date(a.target_time) - new Date(b.target_time))
  const ended = items.value
    .filter((it) => new Date(it.target_time).getTime() <= now)
    .sort((a, b) => new Date(b.target_time) - new Date(a.target_time))
  return [...upcoming, ...ended]
})

function flashToast(msg) {
  toast.value = msg
  setTimeout(() => (toast.value = ''), 2200)
}

async function load() {
  loading.value = true
  try {
    const res = await api.list()
    items.value = res.items || []
  } catch (e) {
    flashToast(e.message)
  } finally {
    loading.value = false
  }
}

function edit(item) {
  router.push({ path: '/create', query: { id: item.id } })
}

async function remove(item) {
  if (!window.confirm(`删除「${item.title}」这个倒计时?`)) return
  try {
    await api.remove(item.id)
    items.value = items.value.filter((it) => it.id !== item.id)
    flashToast('已删除')
  } catch (e) {
    flashToast(e.message)
  }
}

function logout() {
  clearToken()
  router.push('/enter')
}

onMounted(() => {
  load()
  tick()
})

onBeforeUnmount(() => clearInterval(timer))
</script>

<style scoped>
.home-page {
  min-height: 100dvh;
  padding-bottom: 96px;
}

.list-wrap {
  padding: 14px 16px 20px;
}

.cd-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.logout-btn {
  height: 32px;
  padding: 0 14px;
  border-radius: 16px;
  font-size: 13px;
  color: var(--text-dim);
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid var(--card-border);
}

.logout-btn:active {
  color: var(--danger);
  border-color: rgba(248, 113, 113, 0.5);
  background: rgba(248, 113, 113, 0.1);
}

.fab {
  position: fixed;
  right: max(20px, calc(50vw - 195px));
  bottom: 30px;
  width: 58px;
  height: 58px;
  border-radius: 50%;
  font-size: 30px;
  font-weight: 300;
  color: #fff;
  background: linear-gradient(135deg, var(--accent), var(--accent-2));
  box-shadow: 0 10px 30px rgba(56, 189, 248, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.1s ease;
}

.fab:active {
  transform: scale(0.92);
}

.list-enter-active {
  transition: all 0.3s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateY(14px);
}

.list-leave-active {
  transition: all 0.25s ease;
  position: absolute;
  width: 100%;
}

.list-leave-to {
  opacity: 0;
  transform: translateX(40px);
}
</style>
