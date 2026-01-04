<template>
  <div class="user-create-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>创建用户</span>
          <el-button @click="handleCancel">返回</el-button>
        </div>
      </template>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
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
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="总后台" value="admin" />
            <el-option label="店长" value="store_manager" />
            <el-option label="美甲师" value="technician" />
            <el-option label="顾客" value="customer" />
          </el-select>
        </el-form-item>
        
        <el-form-item
          v-if="form.role === 'store_manager' || form.role === 'technician'"
          label="所属门店"
          prop="store_id"
        >
          <el-select v-model="form.store_id" placeholder="请选择门店" style="width: 100%">
            <el-option
              v-for="store in storeList"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="active">激活</el-radio>
            <el-radio label="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            创建
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createUser } from '@/api/users'

const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)
const storeList = ref([])

const form = reactive({
  username: '',
  phone: '',
  email: '',
  password: '',
  role: '',
  store_id: null,
  status: 'active'
})

const rules = {
  username: [
    {
      validator: (rule, value, callback) => {
        if (form.role && form.role !== 'customer' && !value) {
          callback(new Error('员工必须输入用户名'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  phone: [
    {
      validator: (rule, value, callback) => {
        if (form.role === 'customer' && !value) {
          callback(new Error('顾客必须输入手机号'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  password: [
    {
      validator: (rule, value, callback) => {
        if (form.role && form.role !== 'customer' && !value) {
          callback(new Error('员工必须输入密码'))
        } else if (value && value.length < 6) {
          callback(new Error('密码长度不能少于6位'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  store_id: [
    {
      validator: (rule, value, callback) => {
        if ((form.role === 'store_manager' || form.role === 'technician') && !value) {
          callback(new Error('店长和美甲师必须选择门店'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
}

// 假的门店数据（用于测试）
const mockStores = [
  { id: 1, name: '总店' },
  { id: 2, name: '分店A' },
  { id: 3, name: '分店B' },
  { id: 4, name: '分店C' },
  { id: 5, name: '分店D' }
]

// 获取门店列表（使用假数据）
const fetchStoreList = () => {
  // 直接使用假数据，用于测试
  storeList.value = mockStores
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      const data = {
        username: form.username,
        phone: form.phone || undefined,
        email: form.email || undefined,
        password: form.password,
        role: form.role,
        store_id: form.store_id || undefined,
        status: form.status
      }
      
      await createUser(data)
      ElMessage.success('创建用户成功')
      router.push('/users')
    } catch (error) {
      console.error('创建用户失败:', error)
      ElMessage.error(error.message || '创建用户失败')
    } finally {
      submitting.value = false
    }
  })
}

// 取消
const handleCancel = () => {
  router.push('/users')
}

onMounted(() => {
  fetchStoreList()
})
</script>

<style lang="scss" scoped>
.user-create-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>

