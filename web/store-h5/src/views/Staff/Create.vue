<template>
  <div class="staff-create-page">
    <!-- 页面头部 -->
    <van-nav-bar
      title="添加员工"
      left-text="返回"
      left-arrow
      @click-left="$router.back()"
    />
    
    <!-- 表单 -->
    <van-form @submit="handleSubmit">
      <van-cell-group inset>
        <van-field
          v-model="form.username"
          name="username"
          label="用户名"
          placeholder="请输入用户名"
          :rules="[{ required: true, message: '请输入用户名' }]"
          left-icon="user-o"
        />
        
        <van-field
          v-model="form.phone"
          name="phone"
          label="手机号码"
          placeholder="请输入手机号码"
          :rules="phoneRules"
          left-icon="phone-o"
        />
      </van-cell-group>
      
      <div class="form-actions">
        <van-button
          round
          block
          type="primary"
          native-type="submit"
          :loading="loading"
          loading-text="创建中..."
        >
          创建员工
        </van-button>
      </div>
    </van-form>
    
    <!-- 提示信息 -->
    <van-cell-group inset style="margin-top: 16px;">
      <van-cell>
        <template #title>
          <div class="tip-content">
            <van-icon name="info-o" color="#1989fa" />
            <span>提示信息</span>
          </div>
        </template>
        <template #label>
          <div class="tip-text">
            • 员工默认密码为：123456<br/>
            • 创建后员工可使用手机号登录<br/>
            • 员工状态可在员工管理页面调整
          </div>
        </template>
      </van-cell>
    </van-cell-group>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { createStaff } from '@/api/staff'

const router = useRouter()

const loading = ref(false)

const form = reactive({
  username: '',
  phone: '',
  password: '123456' // 默认密码
})

const phoneRules = [
  { required: true, message: '请输入手机号码' },
  { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码' }
]

// 提交表单
const handleSubmit = async () => {
  loading.value = true
  try {
    // 如果没有提供密码，使用默认密码
    const data = {
      username: form.username,
      phone: form.phone,
      password: form.password || '123456'
    }
    await createStaff(data)
    showToast.success('员工创建成功')
    router.back()
  } catch (error) {
    console.error('创建员工失败:', error)
    showToast.fail('创建失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.staff-create-page {
  background-color: #f7f8fa;
  min-height: 100vh;
}

.form-actions {
  padding: 16px;
}

.tip-content {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.tip-text {
  margin-top: 8px;
  line-height: 1.6;
  color: #666;
}
</style>
