<template>
  <div class="profile-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>个人资料</span>
        </div>
      </template>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        v-loading="loading"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        
        <el-form-item label="角色">
          <el-input :value="getRoleName(form.role)" disabled />
        </el-form-item>
        
        <el-form-item label="状态">
          <el-input :value="form.status === 'active' ? '激活' : '禁用'" disabled />
        </el-form-item>
        
        <el-form-item label="新密码">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="留空则不修改密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="确认密码">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="修改密码时需要确认"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            保存
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getCurrentUser } from '@/api/auth'
import { updateUser } from '@/api/users'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)
const submitting = ref(false)

const form = reactive({
  username: '',
  phone: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: '',
  status: ''
})

const originalForm = reactive({
  username: '',
  phone: '',
  email: '',
  role: '',
  status: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  phone: [
    {
      validator: (rule, value, callback) => {
        if (value && !/^1[3-9]\d{9}$/.test(value)) {
          callback(new Error('请输入正确的手机号'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  confirmPassword: [
    {
      validator: (rule, value, callback) => {
        if (form.password && form.password !== value) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取角色名称
const getRoleName = (role) => {
  const roleMap = {
    admin: '管理员',
    store_manager: '店长',
    technician: '美甲师',
    customer: '顾客'
  }
  return roleMap[role] || role
}

// 获取当前用户信息
const fetchUserInfo = async () => {
  loading.value = true
  try {
    const response = await getCurrentUser()
    if (response.data && response.data.user) {
      const user = response.data.user
      Object.assign(form, {
        username: user.username || '',
        phone: user.phone || '',
        email: user.email || '',
        password: '',
        confirmPassword: '',
        role: user.role || '',
        status: user.status || 'active'
      })
      
      // 保存原始数据用于重置
      Object.assign(originalForm, {
        username: user.username || '',
        phone: user.phone || '',
        email: user.email || '',
        role: user.role || '',
        status: user.status || 'active'
      })
      
      // 更新 store 中的用户信息
      userStore.userInfo = user
      localStorage.setItem('userInfo', JSON.stringify(user))
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      const userId = userStore.userInfo.id
      const data = {
        username: form.username,
        phone: form.phone || undefined,
        email: form.email || undefined
      }
      
      // 如果填写了密码，则更新密码
      if (form.password) {
        data.password = form.password
      }
      
      const response = await updateUser(userId, data)
      if (response.data) {
        // 更新 store 中的用户信息
        userStore.userInfo = response.data
        localStorage.setItem('userInfo', JSON.stringify(response.data))
      }
      
      ElMessage.success('更新个人资料成功')
      
      // 重置密码字段
      form.password = ''
      form.confirmPassword = ''
      
      // 更新原始数据
      Object.assign(originalForm, {
        username: form.username,
        phone: form.phone,
        email: form.email
      })
    } catch (error) {
      console.error('更新个人资料失败:', error)
      ElMessage.error(error.message || '更新个人资料失败')
    } finally {
      submitting.value = false
    }
  })
}

// 重置表单
const handleReset = () => {
  Object.assign(form, {
    username: originalForm.username,
    phone: originalForm.phone,
    email: originalForm.email,
    password: '',
    confirmPassword: ''
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

onMounted(() => {
  // 如果 store 中有用户信息，先使用，然后从服务器获取最新信息
  if (userStore.userInfo && userStore.userInfo.id) {
    const user = userStore.userInfo
    Object.assign(form, {
      username: user.username || '',
      phone: user.phone || '',
      email: user.email || '',
      role: user.role || '',
      status: user.status || 'active'
    })
    Object.assign(originalForm, {
      username: user.username || '',
      phone: user.phone || '',
      email: user.email || '',
      role: user.role || '',
      status: user.status || 'active'
    })
  }
  
  fetchUserInfo()
})
</script>

<style lang="scss" scoped>
.profile-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 16px;
    font-weight: 500;
  }
}
</style>

