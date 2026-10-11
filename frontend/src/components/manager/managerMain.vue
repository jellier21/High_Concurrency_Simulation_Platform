<script setup>

</script>

<template>
  <div class="container">
    <h3 class="useNotice">设置秒杀时间</h3>
    <div class="setServiceTimeContent">
      <input type="datetime-local">
      <button>提交</button>
    </div>
  </div>
  <div class="SetServiceList">
    <h3>秒杀活动列表</h3>
    <div class="s_kills_List">
      占位
    </div>
  </div>
</template>

<style scoped>
/* 整体页面样式 */
body {
  font-family: "Microsoft YaHei", Arial, sans-serif;
  background-color: #f2f2f2;
  margin: 0;
  padding: 20px;
}

/* 设置秒杀时间容器 */
.container {
  width: 500px;
  margin: 0 auto 20px;
  background-color: #fff;
  padding: 20px 25px;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 标题 */
.useNotice {
  font-size: 18px;
  color: #333;
  margin-bottom: 15px;
}

/* 设置时间区域 */
.setServiceTimeContent {
  display: flex;
  align-items: center;
  gap: 15px;
}

/* 时间选择器 */
.setServiceTimeContent input[type="datetime-local"] {
  flex: 1;
  padding: 8px 10px;
  font-size: 14px;
  border: 1px solid #ccc;
  border-radius: 4px;
  outline: none;
}

.setServiceTimeContent input[type="datetime-local"]:focus {
  border-color: #409eff;
}

/* 提交按钮 */
.setServiceTimeContent button {
  padding: 8px 20px;
  font-size: 14px;
  color: #fff;
  background-color: #409eff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.setServiceTimeContent button:hover {
  background-color: #66b1ff;
}

/* 秒杀活动列表容器 */
.SetServiceList {
  width: 500px;
  margin: 0 auto;
  background-color: #fff;
  padding: 20px 25px;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.SetServiceList h3 {
  font-size: 18px;
  color: #333;
  margin-bottom: 15px;
}

/* 秒杀列表内容 */
.s_kills_List {
  min-height: 120px;
  border: 1px dashed #ccc;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 14px;
}
</style>