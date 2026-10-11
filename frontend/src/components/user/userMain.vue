<script setup>

</script>

<template>
  <div class="container">
    <div class="Concurrence">
      <h3>并发服务时间列表</h3>
      <div class="channel_lists">
        占位
      </div>
    </div>
    <div class="user-Operator">
      <h2>用户端模拟秒杀高并发局面</h2>
      <div class="SettingArea">
        <input type="number" placeholder="填入并发秒杀参与人数">
        <button>生成</button>
      </div>
      <div class="preView">
        <h3>预查看并发递送信息</h3>
        <div class="preLists">
          占位
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 页面基础样式 */
body {
  font-family: "Microsoft YaHei", Arial, sans-serif;
  background-color: #f2f2f2;
  margin: 0;
  padding: 20px;
}

/* 主容器：左右分栏 */
.container {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  max-width: 1100px;
  margin: 0 auto;
}

/* 左侧 & 右侧公共面板 */
.Concurrence,
.user-Operator {
  background-color: #fff;
  border-radius: 6px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.Concurrence {
  flex: 1;
}

.user-Operator {
  flex: 1;
}

/* 标题 */
h2,
h3 {
  margin: 0 0 15px;
  color: #333;
}

h3 {
  font-size: 16px;
}

h2 {
  font-size: 18px;
}

/* 占位列表 */
.channel_lists,
.preLists {
  min-height: 120px;
  border: 1px dashed #ccc;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 14px;
}

/* 用户端设置区域 */
.SettingArea {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.SettingArea input {
  flex: 1;
  padding: 8px 10px;
  font-size: 14px;
  border: 1px solid #ccc;
  border-radius: 4px;
  outline: none;
}

.SettingArea input:focus {
  border-color: #409eff;
}

.SettingArea button {
  padding: 8px 20px;
  font-size: 14px;
  color: #fff;
  background-color: #409eff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.SettingArea button:hover {
  background-color: #66b1ff;
}

/* 预览区域 */
.preView h3 {
  margin-bottom: 10px;
}
</style>