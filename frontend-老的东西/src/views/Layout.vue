<template>
  <ElContainer class="dashboard-container">
    <ElHeader class="dashboard-header">
      <div class="header-content">
        <div class="header-left">
          <h2>今天吃什么呢？</h2>
          <div class="notification-wrapper">
            <!-- <ElButton type="primary" @click="handleGoToProfile" plain>
              <el-icon><User /></el-icon>
              个人信息
            </ElButton> -->
            <div
              v-if="hasPendingRequests && userType === 1"
              class="notification-dot"
            ></div>
          </div>
        </div>
        <ElButton type="danger" @click="handleLogout">退出登录</ElButton>
      </div>
    </ElHeader>

    <ElContainer style="padding-top: 80px">
      <ElAside width="200px" style="padding-top: 0px">
        <ElMenu class="dashboard-menu" router :default-active="$route.path">
          <ElMenuItem index="/profile">
            <el-icon><List /></el-icon>
            <span>个人资料</span>
          </ElMenuItem>
          <ElMenuItem index="/postrecipe">
            <el-icon><List /></el-icon>
            <span>发布食谱</span>
          </ElMenuItem>
          <ElMenuItem index="/findrecipe">
            <el-icon><List /></el-icon>
            <span>随便看看</span>
          </ElMenuItem>
          <ElMenuItem index="/myrecipe">
            <el-icon><List /></el-icon>
            <span>我的食谱</span>
          </ElMenuItem>
          <!-- 其他菜单项 -->
        </ElMenu>
      </ElAside>

      <ElMain>
        <router-view></router-view>
      </ElMain>
    </ElContainer>
  </ElContainer>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import {
  ElContainer,
  ElAside,
  ElMain,
  ElHeader,
  ElMenu,
  ElMenuItem,
  ElButton,
} from "element-plus";
import { User, List } from "@element-plus/icons-vue";

const router = useRouter();
const uid = ref(localStorage.getItem("uid"));
const userType = ref(1);
const hasPendingRequests = ref(false);

const handleLogout = () => {
  localStorage.removeItem("jwt");
  localStorage.removeItem("uid");
  router.push("/login");
};

const handleGoToProfile = () => {
  router.push("/profile");
};

onMounted(async () => {
  try {
    const token = localStorage.getItem("jwt");
    if (!token) throw new Error("未登录");

    // 获取用户信息
    const response = await fetch(`http://localhost:3000/user/info`, {
      method: "GET",
      headers: { Authorization: token },
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || "获取用户信息失败");
    }
    const userData = (await response.json()).data;
    userType.value = userData.user_type;
  } catch (error) {
    console.error(error);
    router.push("/login");
  }
});
</script>

<style scoped>
.dashboard-container {
  height: 100vh;
}

.dashboard-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
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
}

.dashboard-menu {
  height: 100%;
  border-right: none;
}
</style>
