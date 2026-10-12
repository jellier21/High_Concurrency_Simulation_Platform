// 网关层接口统一封装（相对路径走 Vite proxy → localhost:8080）
async function request(url, options = {}) {
    const res = await fetch(url, options)
    return res.json()
}

// 1. 创建秒杀活动（start/end 为毫秒时间戳字符串）
export function createActivity(start, end) {
    return request('/api/activity/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: `start=${encodeURIComponent(start)}&end=${encodeURIComponent(end)}`
    })
}

// 2. 获取全部活动列表
export function getActivities() {
    return request('/api/activity/get')
}

// 3. 创建订单任务（模拟秒杀并发）
export function createOrderTask(mslNo, count) {
    return request('/api/order/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: `msl_no=${encodeURIComponent(mslNo)}&count=${encodeURIComponent(count)}`
    })
}

// 4. 查询任务列表及状态
export function getTaskList() {
    return request('/api/order/status/task')
}

// 5. 按活动编号查询订单列表
export function getOrdersByMslNo(mslNo) {
    return request(`/api/order/status/order?msl_no=${encodeURIComponent(mslNo)}`)
}
