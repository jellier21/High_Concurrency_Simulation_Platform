<script setup>
import { ref, onMounted } from 'vue'
import { createActivity, getActivities } from '../../api'


const startTime = ref('')
const endTime = ref('')
const activities = ref([])
const submitting = ref(false)

// datetime-local 的值("2026-10-12T15:04")转毫秒时间戳
function toMs(dt) {
  return new Date(dt).getTime().toString()
}

async function submit() {
  if (!startTime.value || !endTime.value) {
    alert('请先选择开始和结束时间')
    return
  }
  if (new Date(startTime.value).getTime() >= new Date(endTime.value).getTime()) {
    alert('结束时间必须晚于开始时间')
    return
  }
  submitting.value = true
  try {
    const res = await createActivity(toMs(startTime.value), toMs(endTime.value))
    alert(res.message || '创建成功')
    await loadActivities()
  } finally {
    submitting.value = false
  }
}

async function loadActivities() {
  const res = await getActivities()
  activities.value = res.data || []
}

onMounted(loadActivities)
</script>

<template>
  <div class="card">
    <h3 class="card-title">设置秒杀时间</h3>
    <div class="form-row">
      <input type="datetime-local" v-model="startTime" class="input">
      <input type="datetime-local" v-model="endTime" class="input">
      <button class="btn" :disabled="submitting" @click="submit">
        {{ submitting ? '提交中...' : '提交' }}
      </button>
    </div>
  </div>

  <div class="card">
    <h3 class="card-title">秒杀活动列表</h3>
    <div class="empty" v-if="activities.length === 0">
      <span class="empty-icon">🗓️</span>
      <span>暂无秒杀活动，请先创建活动</span>
    </div>
    <ul class="list" v-else>
      <li v-for="a in activities" :key="a.msl_uid">
        <div class="list-main">
          <span class="tag">活动编号</span>
          <code>{{ a.msl_uid }}</code>
        </div>
        <div class="list-sub">开始 {{ a.start_time }} · 结束 {{ a.end_time }}</div>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.list { list-style: none; }
.list li {
  padding: 14px 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
  margin-bottom: 10px;
  background: #fafbfc;
}
.list-main { display: flex; align-items: center; gap: 10px; margin-bottom: 6px; }
.list-main code { color: var(--primary); font-size: 13px; }
.tag { background: var(--primary); color: #fff; font-size: 12px; padding: 2px 8px; border-radius: 4px; }
.list-sub { font-size: 13px; color: var(--text-light); }
</style>
