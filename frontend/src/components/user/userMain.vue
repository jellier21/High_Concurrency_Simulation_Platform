<script setup>
import {ref, computed, onMounted, onUnmounted} from 'vue'
import { getActivities, createOrderTask, getTaskList, getOrdersByMslNo } from '../../api'
import { sendOrderToRabit } from '../../api/rabbit.js'

const channels = ref([])
const count = ref('')
const generating = ref(false)
const tasks = ref([])
const orders = ref([])
const selectedMslNo = ref('')
const pollingTimer = ref(null)

// ============ 订单状态统计 ============
const stats = computed(() => {
  const s = { total: orders.value.length, success: 0, waiting: 0, fail: 0 }
  for (const o of orders.value) {
    if (o.status === 'SUCCESS') s.success++
    else if (o.status === 'FAIL' || o.status === 'OVERDUE') s.fail++
    else s.waiting++
  }
  return s
})

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

// 拉取订单列表（轮询也走这里）
async function refreshOrders() {
  if (!selectedMslNo.value) return
  const res = await getOrdersByMslNo(selectedMslNo.value)
  orders.value = res.data || []
  syncOrderTimers(orders.value)
}

function previewOrders() {
  refreshOrders()
  if (pollingTimer.value) clearInterval(pollingTimer.value)
  pollingTimer.value = setInterval(refreshOrders, 4000)
}

onMounted(() => { loadChannels(); loadTasks() })

// ============ 每个订单独立的定时任务 ============
const orderTimers = new Map()   // key=uuid, value=setTimeout id
const sentUuids = new Set()      // 网关确认受理过的 uuid

// 分批发送：每批 50 条，批间隔 100ms；成功才记 sentUuids
function sendOrdersBatch(orderList) {
  const BATCH = 50
  const GAP = 100
  for (let i = 0; i < orderList.length; i += BATCH) {
    setTimeout(() => {
      orderList.slice(i, i + BATCH).forEach(o => {
        sendOrderToRabit(o.uuid)
            .then(() => { sentUuids.add(o.uuid) })
            .catch(() => { /* 失败不记，下轮轮询自动重发 */ })
      })
    }, (i / BATCH) * GAP)
  }
}

// 对账定时器：新增的注册，消失的清除；发过的/已处理的跳过
function syncOrderTimers(orderList) {
  const alive = new Set(orderList.map(o => o.uuid))

  for (const [uuid, timerId] of orderTimers) {
    if (!alive.has(uuid)) {
      clearTimeout(timerId)
      orderTimers.delete(uuid)
    }
  }

  for (const order of orderList) {
    if (orderTimers.has(order.uuid)) continue
    if (sentUuids.has(order.uuid)) continue
    if (order.status && order.status !== 'WAITING') continue

    const target = new Date(order.start_time.replace(' ', 'T')).getTime()
    const delay = target - Date.now()

    if (delay <= 0) {
      sendOrdersBatch([order])
      continue
    }

    const timerId = setTimeout(() => {
      sendOrdersBatch([order])
      orderTimers.delete(order.uuid)
    }, delay)

    orderTimers.set(order.uuid, timerId)
  }
}

onUnmounted(() => {
  for (const timerId of orderTimers.values()) clearTimeout(timerId)
  orderTimers.clear()
  if (pollingTimer.value) clearInterval(pollingTimer.value)
})
</script>

<template>
  <div class="grid">
    <!-- 左：活动列表 -->
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

    <!-- 右：模拟秒杀 -->
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

      <h4 class="sub-title">订单处理进度</h4>
      <button class="btn btn-ghost" @click="previewOrders">查询该活动订单</button>

      <!-- 统计条 -->
      <div class="stats" v-if="orders.length > 0">
        <div class="stat-item">
          <span class="stat-num">{{ stats.total }}</span>
          <span class="stat-label">总数</span>
        </div>
        <div class="stat-item">
          <span class="stat-num num-success">{{ stats.success }}</span>
          <span class="stat-label">成功</span>
        </div>
        <div class="stat-item">
          <span class="stat-num num-waiting">{{ stats.waiting }}</span>
          <span class="stat-label">等待</span>
        </div>
        <div class="stat-item">
          <span class="stat-num num-fail">{{ stats.fail }}</span>
          <span class="stat-label">失败</span>
        </div>
      </div>

      <!-- 订单列表（限高滚动） -->
      <div class="empty" v-if="orders.length === 0">
        <span class="empty-icon">📦</span>
        <span>暂无递送信息</span>
      </div>
      <ul class="list order-scroll" v-else>
        <li v-for="o in orders" :key="o.uuid">
          <div class="list-main">
            <code class="order-uuid">{{ o.uuid }}</code>
            <span class="status" :class="o.status">{{ o.status || 'UNKNOWN' }}</span>
          </div>
          <div class="list-sub">{{ o.start_time }}</div>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.sub-title { font-size: 14px; font-weight: 600; margin: 18px 0 12px; }
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

/* 状态标签 */
.status { font-size: 12px; padding: 2px 10px; border-radius: 10px; background: #f0f2f5; color: var(--text-light); white-space: nowrap; }
.status.Running { background: #ecf5ff; color: #409eff; }
.status.Success, .status.SUCCESS { background: #f0f9eb; color: #67c23a; }
.status.Interrupt, .status.FAIL { background: #fef0f0; color: #f56c6c; }
.status.WAITING { background: #fdf6ec; color: #e6a23c; }
.status.OVERDUE { background: #f4f4f5; color: #909399; }

/* 统计条 */
.stats { display: flex; gap: 12px; margin-bottom: 14px; }
.stat-item {
  flex: 1;
  background: #fafbfc;
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 10px;
  text-align: center;
}
.stat-num { display: block; font-size: 20px; font-weight: 700; color: var(--text); }
.stat-label { font-size: 12px; color: var(--text-light); }
.num-success { color: #67c23a; }
.num-waiting { color: #e6a23c; }
.num-fail { color: #f56c6c; }

/* 订单列表限高滚动 */
.order-scroll { max-height: 420px; overflow-y: auto; padding-right: 4px; }
.order-scroll::-webkit-scrollbar { width: 6px; }
.order-scroll::-webkit-scrollbar-thumb { background: #dcdfe6; border-radius: 3px; }
.order-uuid { word-break: break-all; }

.btn-ghost { background: #fff; color: var(--primary); border: 1px solid var(--primary); margin-bottom: 12px; }
.btn-ghost:hover { background: #ecf5ff; }
@media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }
</style>
