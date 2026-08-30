<template>
  <div class="card" :class="{ ended }">
    <button class="del" title="删除" @click="$emit('delete')">✕</button>

    <div class="title" :class="{ dim: ended }">
      <span v-if="!ended">距离「{{ item.title }}」还有</span>
      <span v-else>「{{ item.title }}」已结束</span>
    </div>

    <div v-if="!ended" class="digits">
      <div class="unit">
        <span class="num">{{ pad2(r.days) }}</span>
        <span class="lbl">天</span>
      </div>
      <span class="sep">:</span>
      <div class="unit">
        <span class="num">{{ pad2(r.hours) }}</span>
        <span class="lbl">时</span>
      </div>
      <span class="sep">:</span>
      <div class="unit">
        <span class="num">{{ pad2(r.minutes) }}</span>
        <span class="lbl">分</span>
      </div>
      <span class="sep">:</span>
      <div class="unit">
        <span class="num">{{ pad2(r.seconds) }}</span>
        <span class="lbl">秒</span>
      </div>
    </div>

    <div v-else class="ended-badge">⏹ 已完成</div>

    <div class="target">
      {{ formatTarget }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { pad2, splitRemaining } from '../utils/countdown'

const props = defineProps({
  item: { type: Object, required: true },
  nowMs: { type: Number, required: true },
})
defineEmits(['delete'])

const targetMs = computed(() => new Date(props.item.target_time).getTime())
const r = computed(() => splitRemaining(targetMs.value - props.nowMs))
const ended = computed(() => r.value.ended)

const formatTarget = computed(() => {
  const d = new Date(targetMs.value)
  const p = (n) => String(n).padStart(2, '0')
  return `目标 ${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
})
</script>

<style scoped>
.card {
  position: relative;
  background: linear-gradient(160deg, var(--card), #131c30);
  border: 1px solid var(--card-border);
  border-radius: var(--radius);
  padding: 18px 18px 14px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.35);
  overflow: hidden;
}

.card::before {
  content: '';
  position: absolute;
  inset: 0 0 auto 0;
  height: 3px;
  background: linear-gradient(90deg, var(--accent), var(--accent-2));
  opacity: 0.85;
}

.card.ended {
  opacity: 0.55;
  filter: saturate(0.4);
}

.del {
  position: absolute;
  top: 10px;
  right: 12px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  color: var(--text-dim);
  background: rgba(255, 255, 255, 0.05);
  font-size: 13px;
  z-index: 2;
}

.del:active {
  color: var(--danger);
  background: rgba(248, 113, 113, 0.15);
}

.title {
  font-size: 14px;
  color: var(--text);
  margin: 4px 34px 14px 0;
  word-break: break-all;
  line-height: 1.5;
}

.title.dim {
  color: var(--text-dim);
}

.digits {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  justify-content: center;
}

.unit {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.num {
  font-family: 'SF Mono', 'JetBrains Mono', Consolas, 'Courier New', monospace;
  font-size: 34px;
  font-weight: 800;
  font-variant-numeric: tabular-nums;
  color: #fff;
  background: linear-gradient(180deg, #ffffff, #9fd8ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  text-shadow: 0 0 18px rgba(56, 189, 248, 0.5);
  line-height: 1.1;
}

.lbl {
  font-size: 11px;
  color: var(--text-dim);
}

.sep {
  font-size: 28px;
  font-weight: 800;
  color: var(--accent);
  opacity: 0.7;
  line-height: 1.1;
  margin-bottom: 14px;
}

.ended-badge {
  text-align: center;
  font-size: 16px;
  font-weight: 700;
  color: var(--text-dim);
  padding: 6px 0;
}

.target {
  margin-top: 12px;
  text-align: center;
  font-size: 11px;
  color: var(--text-dim);
  font-variant-numeric: tabular-nums;
}
</style>
