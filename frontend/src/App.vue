<script setup>
</script>

<template>
  <nav class="top-nav">
    <div class="nav-inner">
      <span class="brand">⚡ 高并发秒杀平台</span>
      <div class="nav-links">
        <RouterLink to="/user">用户端</RouterLink>
        <RouterLink to="/manager">管理端</RouterLink>
      </div>
    </div>
  </nav>
  <main class="page">
    <RouterView></RouterView>
  </main>
</template>

<style>
.top-nav {
  background: linear-gradient(135deg, #1f2d3d 0%, #2c3e50 100%);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}
.nav-inner {
  max-width: 1100px;
  margin: 0 auto;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}
.brand { color: #fff; font-size: 16px; font-weight: 600; letter-spacing: 1px; }
.nav-links { display: flex; gap: 8px; }
.nav-links a {
  color: #c0c4cc;
  text-decoration: none;
  padding: 8px 18px;
  border-radius: 6px;
  font-size: 14px;
  transition: all .2s;
}
.nav-links a:hover { color: #fff; background: rgba(255, 255, 255, 0.12); }
.nav-links a.router-link-active { color: #fff; background: #409eff; }
.page { max-width: 1100px; margin: 0 auto; padding: 28px 24px; }
</style>
