<template>
  <div class="user-center-container">
    <el-row :gutter="20">
      <!-- 左侧：用户头像与信息 -->
      <el-col :span="6" :xs="24">
        <el-card shadow="hover" class="user-profile-card">
          <div class="avatar-wrapper">
            <el-avatar
              :size="100"
              src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png "
            />
          </div>
          <div class="user-info-text">
            <h3>{{ userInfo.name }}</h3>
            <p>ID: {{ uid }}</p>
          </div>
          <el-button type="primary" @click="showResetPwdDialog = true">修改密码</el-button>
        </el-card>
      </el-col>

      <!-- 右侧：用户信息展示 -->
      <el-col :span="18" :xs="24">
        <el-card shadow="hover" class="user-info-card">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
            </div>
          </template>
          <div class="user-info-list">
            <el-row :gutter="20" v-for="(item, index) in userInfoFields" :key="index">
              <el-col :span="8" class="label-col">
                <span class="info-label">{{ item.label }}</span>
              </el-col>
              <el-col :span="16" class="value-col">
                <span class="info-text">{{ userInfo[item.key] || '未设置' }}</span>
              </el-col>
            </el-row>
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
import { ref } from 'vue';
import { ElMessage } from 'element-plus';

// 用户信息模拟数据
const uid = ref(localStorage.getItem('uid') || '12345');
const userInfo = ref({
  name: '张三',
  username: 'zhangsan',
  email: 'zhangsan@example.com',
  phone: '13800001111',
  gender: '男',
  address: '北京市朝阳区',
});

// 用户信息字段
const userInfoFields = [
  { label: '用户名', key: 'username' },
  { label: '邮箱', key: 'email' },
  { label: '手机号', key: 'phone' },
  { label: '性别', key: 'gender' },
  { label: '地址', key: 'address' },
];

// 控制修改密码弹窗
const showResetPwdDialog = ref(false);

// 密码表单数据
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
});

// 表单验证规则
const rules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在6到20个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: any) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入密码不一致'));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    },
  ],
};

// 提交修改密码
const handleResetPassword = () => {
  // 这里可以调用 API 请求
  ElMessage.success('密码修改成功');
  showResetPwdDialog.value = false;
};
</script>

<style scoped>
.user-center-container {
  padding: 20px;
  max-width: 1400px; /* 最大宽度 */
  margin: 0 auto; /* 居中 */
  width: calc(100% - 40px); /* 留出左右间距 */
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-center-container {
    padding: 15px;
    margin: 0 20px;
  }
}

/* 用户信息卡片 */
.user-profile-card {
  text-align: center;
  padding: 20px;
}

.avatar-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 10px;
}

.user-info-text {
  text-align: center;
  margin-bottom: 10px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 用户信息列表 */
.user-info-card {
  margin-top: 20px;
}

.user-info-list {
  padding: 20px;
}

.info-label {
  font-weight: bold;
  color: #666;
}

.info-text {
  display: block;
  font-size: 14px;
  color: #666;
}

/* 表格布局 */
.user-info-list .el-row {
  margin-bottom: 10px;
}

.label-col {
  text-align: right;
}

.value-col {
  text-align: left;
}

/* 小屏幕时的样式 */
@media (max-width: 768px) {
  .user-profile-card,
  .user-info-card {
    margin: 0 auto;
  }

  .user-info-card {
    margin-top: 20px;
  }
}
</style>