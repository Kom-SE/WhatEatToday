<template>
  <div class="basic-layout">
    <button @click="login">登录傻逼</button>
    <div v-if="response">
      <p>message：{{ response.message }}</p>
      <p>token：{{ response.token }}</p>
    </div>
    <div v-if="error" style="color: red">
      <p>错误：{{ error }}</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import axios from 'axios'
import { ref } from 'vue'

// 响应式数据
interface ResponseData {
  message: string
  token: string
  uid: number
}
const response = ref<ResponseData | null>(null)
const error = ref<string | null>(null)
// 方法
const login = async () => {
  try {
    error.value = null // 清除之前的错误信息
    // 这里添加你的请求逻辑
    const requestData = {
      username: 'fuck',
      password: '123456',
    }
    console.log('发送请求数据:', requestData)
    const result = await axios.post('http://127.0.0.1:3000/auth/login', requestData)

    response.value = result.data
    console.log('注册成功:', result.data)
  } catch (err: any) {
    console.error('注册失败:', err)
    console.error('错误响应:', err.response?.data)
    error.value = err.response?.data?.error || '请求失败'
  }
}
</script>

<style scoped></style>
