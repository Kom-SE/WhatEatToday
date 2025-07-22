<template>
  <div class="user-center-container">
    <el-row :gutter="20">
      <!-- 左侧：用户头像与信息 -->
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="user-profile-card">
            <el-avatar
              :size="100"
              src="/image/avatar/user/avatar_2_0_20250704_043425.png"
            />
            <h3>{{ userInfo.name }}</h3>
            <p>ID: {{ uid }}</p>
            <el-button type="primary" @click="showResetPwdDialog = true">修改密码</el-button>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：用户信息展示 -->
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
            </div>
          </template>
          <div class="user-info-list">
            <el-form label-position="left" label-width="100px">
              <el-form-item label="用户名">
                <span class="info-text">{{ userInfo.username }}</span>
              </el-form-item>
              <el-form-item label="邮箱">
                <span class="info-text">{{ userInfo.email || '未设置' }}</span>
              </el-form-item>
              <el-form-item label="手机号">
                <span class="info-text">{{ userInfo.phone || '未设置' }}</span>
              </el-form-item>
              <el-form-item label="性别">
                <span class="info-text">{{ userInfo.gender || '未设置' }}</span>
              </el-form-item>
              <el-form-item label="地址">
                <span class="info-text">{{ userInfo.address || '未填写' }}</span>
              </el-form-item>
            </el-form>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 修改密码对话框 -->
    <el-dialog v-model="showResetPwdDialog" title="修改密码" width="40%">
      <el-form ref="formRef" :model="passwordForm" :rules="rules" label-width="100px">
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showResetPwdDialog = false">取消</el-button>
        <el-button type="primary" @click="handleResetPassword">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { onMounted } from 'vue' // 添加 onMounted 导入
import { useRouter } from 'vue-router' // 添加 router 导入

const router = useRouter() // 定义 router 实例

// 用户信息模拟数据
const uid = ref(localStorage.getItem('uid') || '12345')
const userInfo = ref({
  name: '',
  username: '',
  email: '',
  phone: '',
  gender: '',
  address: ''
})
const fetchUserInfo = async () => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await fetch(`http://localhost:3000/user/info`, {
      method: 'GET',
      headers: { Authorization: token },
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || '获取用户信息失败')
    }
    
    const userData = (await response.json()).data
    
    // 更新用户信息
    userInfo.value = {
      name: userData.name || '',
      username: userData.username || '',
      email: userData.email || '',
      phone: userData.phone || '',
      gender: userData.gender || '',
      address: userData.address || ''
    }
    
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
    router.push('/login')
  }
}

// 组件挂载时获取用户信息
onMounted(() => {
  fetchUserInfo()
})

// 控制修改密码弹窗
const showResetPwdDialog = ref(false)

// 密码表单数据
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 表单验证规则
const rules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在6到20个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: any) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 提交修改密码
const handleResetPassword = () => {
  // 这里可以调用 API 请求
  ElMessage.success('密码修改成功')
  showResetPwdDialog.value = false
}
</script>

<style scoped>
.user-center-container {
  padding: 20px;
}

.user-profile-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 20px 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-text {
  display: block;
  font-size: 14px;
  color: #666;
}
</style>