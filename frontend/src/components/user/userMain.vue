<script setup>
import { ref, onMounted } from 'vue'
import { getActivities, createOrderTask, getTaskList, getOrdersByMslNo } from '../../api'


const channels = ref([])   // 并发服务时间列表（活动）
const count = ref('')      // 参与人数
const generating = ref(false)
const tasks = ref([])      // 任务列表及状态
const orders = ref([])     // 预查看订单
const selectedMslNo = ref('')

async function loadChannels() {
  const res = await getActivities()
  channels.value = res.data || []
}

async function generate() {
  if (!selectedMslNo.value) { alert('请先选择要模拟秒杀的活动'); return }
  if (!count.value || +count.value <= 0) { alert('请输入正确的参与人数'); return }
  generating.value = true
  try {
    const res = await createOrderTask(selectedMslNo.value, count.value)
    alert(res.message || '任务已创建')
    await loadTasks()
  } finally {
    generating.value = false
  }
}

async function loadTasks() {
  const res = await getTaskList()
  tasks.value = res.data || []
}

async function previewOrders() {
  if (!selectedMslNo.value) { alert('请先选择活动'); return }
  const res = await getOrdersByMslNo(selectedMslNo.value)
  orders.value = res.data || []
}

onMounted(() => { loadChannels(); loadTasks() })
</script>

<template>
  <div class="grid">
    <div class="card">
      <h3 class="card-title">并发服务时间列表</h3>
      <div class="empty" v-if="channels.length === 0">
        <span class="empty-icon">⏱️</span>
        <span>暂无并发服务时间</span>
      </div>
      <ul class="list" v-else>
        <li
            v-for="a in channels"
            :key="a.msl_uid"
            :class="{ active: selectedMslNo === a.msl_uid }"
            @click="selectedMslNo = a.msl_uid"
        >
          <div class="list-main">
            <span class="tag">活动编号</span>
            <code>{{ a.msl_uid }}</code>
          </div>
          <div class="list-sub">{{ a.start_time }} → {{ a.end_time }}</div>
        </li>
      </ul>
    </div>

    <div class="card">
      <h3 class="card-title">用户端模拟秒杀高并发局面</h3>
      <div class="form-row">
        <input type="number" v-model="count" class="input" placeholder="填入并发秒杀参与人数">
        <button class="btn" :disabled="generating" @click="generate">
          {{ generating ? '生成中...' : '生成' }}
        </button>
      </div>

      <h4 class="sub-title">任务列表及状态</h4>
      <div class="empty" v-if="tasks.length === 0">
        <span class="empty-icon">📋</span>
        <span>暂无任务</span>
      </div>
      <ul class="list" v-else>
        <li v-for="t in tasks" :key="t.uuid">
          <div class="list-main">
            <code>{{ t.msl_no }}</code>
            <span class="status" :class="t.status">{{ t.status }}</span>
          </div>
          <div class="list-sub">成功 {{ t.success_int }} / {{ t.total }} · 失败 {{ t.fail_int }}</div>
        </li>
      </ul>

      <h4 class="sub-title">预查看并发递送信息</h4>
      <button class="btn btn-ghost" @click="previewOrders">查询该活动订单</button>
      <div class="empty" v-if="orders.length === 0">
        <span class="empty-icon">📦</span>
        <span>暂无递送信息</span>
      </div>
      <ul class="list" v-else>
        <li v-for="o in orders" :key="o.uuid">
          <div class="list-main"><code>{{ o.uuid }}</code></div>
          <div class="list-sub">{{ o.start_time }}</div>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.sub-title { font-size: 14px; font-weight: 600; margin-bottom: 12px; }
.list { list-style: none; }
.list li {
  padding: 12px 14px;
  border: 1px solid var(--border);
  border-radius: 8px;
  margin-bottom: 10px;
  background: #fafbfc;
  cursor: pointer;
  transition: border-color .2s, box-shadow .2s;
}
.list li:hover { border-color: var(--primary); }
.list li.active { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.15); }
.list-main { display: flex; align-items: center; gap: 10px; margin-bottom: 6px; }
.list-main code { color: var(--primary); font-size: 13px; }
.tag { background: var(--primary); color: #fff; font-size: 12px; padding: 2px 8px; border-radius: 4px; }
.list-sub { font-size: 13px; color: var(--text-light); }
.status { font-size: 12px; padding: 2px 8px; border-radius: 4px; background: #f0f2f5; color: var(--text-light); }
.status.Running { background: #ecf5ff; color: #409eff; }
.status.Success { background: #f0f9eb; color: #67c23a; }
.status.Interrupt { background: #fef0f0; color: #f56c6c; }
.btn-ghost { background: #fff; color: var(--primary); border: 1px solid var(--primary); margin-bottom: 12px; }
.btn-ghost:hover { background: #ecf5ff; }
@media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }
</style>
