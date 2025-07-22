<template>
  <ElContainer class="dashboard-container">
    <ElHeader class="dashboard-header">
      <div class="header-content">
        <div class="header-left">
          <h2>今天吃什么呢~</h2>
          <div class="notification-wrapper">
            <div v-if="hasPendingRequests && userType === 1" class="notification-dot"></div>
          </div>
        </div>
        <div class="header-right">
          <ElButton type="danger" @click="handleLogout">退出登录</ElButton>
        </div>
      </div>
    </ElHeader>

    <ElContainer style="padding-top:80px">
      <ElAside width="200px" style="padding-top:0px">
        <ElMenu
          class="dashboard-menu"
          router
          :default-active="$route.path"
        >
          <ElMenuItem index="/profile">
            <el-icon><List /></el-icon>
            <span>个人资料</span>
          </ElMenuItem>
          <ElMenuItem index="/postrecipe">
            <el-icon><List /></el-icon>
            <span>发布食谱</span>
          </ElMenuItem>
          </ElMenu>
      </ElAside>

      <ElMain>
        <router-view></router-view>
      </ElMain>
    </ElContainer>
  </ElContainer>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElContainer, ElAside, ElMain, ElHeader, ElMenu, ElMenuItem, ElButton } from 'element-plus'
import { User, List } from '@element-plus/icons-vue' // Keep User icon if you uncomment the button later

const router = useRouter()
const uid = ref(localStorage.getItem('uid'))
const userType = ref(1) // Example user type, update based on actual data
const hasPendingRequests = ref(false) // Example notification state

const handleLogout = () => {
  localStorage.removeItem('jwt')
  localStorage.removeItem('uid')
  router.push('/login')
}

const handleGoToProfile = () => {
  router.push('/profile')
}

onMounted(async () => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('未登录')

    // Get user info
    const response = await fetch(`http://localhost:3000/user/info`, {
      method: 'GET',
      headers: { Authorization: token },
    })
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || '获取用户信息失败')
    }
    const userData = (await response.json()).data
    userType.value = userData.user_type
    // Example: Check for pending requests based on userType or other data
    // hasPendingRequests.value = userType.value === 1 && userData.some_pending_flag;
  } catch (error) {
    console.error('Failed to fetch user info or token missing:', error)
    router.push('/login')
  }
})
</script>

<style scoped>
.dashboard-container {
  height: 100vh;
}

.dashboard-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0 20px;
  /* Use fixed positioning to keep the header at the top when scrolling */
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000; /* Ensure it stays above other content */
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05); /* Add a subtle shadow */
}

.header-content {
  display: flex;
  justify-content: space-between; /* Pushes left and right content to edges */
  align-items: center; /* Vertically centers items */
  height: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

/* Optional: style for the title */
.header-left h2 {
  margin: 0; /* Remove default margin */
  font-size: 24px;
  color: #333;
}

.notification-wrapper {
  position: relative;
  display: inline-block;
}

.notification-dot {
  position: absolute;
  top: -5px;
  right: -5px;
  width: 10px;
  height: 10px;
  background-color: #f56c6c;
  border-radius: 50%;
  box-shadow: 0 0 0 2px #fff; /* Add a white border for visibility */
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px; /* Space between buttons if multiple */
}

 .dashboard-menu :deep(.el-menu-item) {
  font-size: 16px;
  padding: 12px 20px; /* 增加内边距使菜单项更大 */
}

.dashboard-menu :deep(.el-menu-item span) {
  font-size: 16px;
  font-weight: 500;
}

.dashboard-menu :deep(.el-icon) {
  font-size: 18px;
  margin-right: 10px;
}

/* Ensure the main content starts below the fixed header */
.ElContainer {
  padding-top: 80px; /* Adjust this to match your header's height */
}

/* Specific styling for the ElMain to adjust for fixed header and sticky elements */
.ElMain {
  padding: 20px; /* Add some padding to the main content area */
}
</style>