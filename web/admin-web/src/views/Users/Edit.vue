<template>
  <div class="user-edit-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>编辑用户</span>
          <el-button @click="handleCancel">返回</el-button>
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
        
        <el-form-item label="密码">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="留空则不修改密码"
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
        
        <el-form-item
          v-if="form.role === 'technician'"
          label="工作状态"
          prop="work_status"
        >
          <el-select v-model="form.work_status" placeholder="请选择工作状态" style="width: 100%">
            <el-option label="在岗" value="working" />
            <el-option label="休息" value="rest" />
            <el-option label="离岗" value="offline" />
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
            保存
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getUserById, updateUser } from '@/api/users'
import { getStores } from '@/api/stores'

const router = useRouter()
const route = useRoute()
const formRef = ref(null)
const loading = ref(false)
const submitting = ref(false)
const storeList = ref([])

const form = reactive({
  username: '',
  phone: '',
  email: '',
  password: '',
  role: '',
  store_id: null,
  work_status: null,
  status: 'active'
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
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

// 获取门店列表
const fetchStoreList = async () => {
  try {
    // 获取所有门店，设置较大的 page_size 以获取全部数据
    const response = await getStores({
      page: 1,
      page_size: 1000
    })
    
    // 后端返回格式: { code: 200, data: { list: [], pagination: {...} } }
    if (response.data && response.data.list) {
      storeList.value = response.data.list
    } else if (Array.isArray(response.data)) {
      storeList.value = response.data
    } else {
      storeList.value = []
    }
  } catch (error) {
    console.error('获取门店列表失败:', error)
    ElMessage.error('获取门店列表失败')
    storeList.value = []
  }
}

// 获取用户详情
const fetchUserDetail = async () => {
  loading.value = true
  try {
    const userId = route.params.id
    const response = await getUserById(userId)
    if (response.data) {
      Object.assign(form, {
        username: response.data.username || '',
        phone: response.data.phone || '',
        email: response.data.email || '',
        password: '',
        role: response.data.role || '',
        store_id: response.data.store_id || null,
        work_status: response.data.work_status || null,
        status: response.data.status || 'active'
      })
    }
  } catch (error) {
    console.error('获取用户详情失败:', error)
    ElMessage.error('获取用户详情失败')
    router.push('/users')
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
      const userId = route.params.id
      const data = {
        username: form.username,
        phone: form.phone || undefined,
        email: form.email || undefined,
        password: form.password || undefined,
        role: form.role,
        store_id: form.store_id || undefined,
        work_status: form.work_status || undefined,
        status: form.status
      }
      
      await updateUser(userId, data)
      ElMessage.success('更新用户成功')
      router.push('/users')
    } catch (error) {
      console.error('更新用户失败:', error)
      ElMessage.error(error.message || '更新用户失败')
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
  fetchUserDetail()
})
</script>

<style lang="scss" scoped>
.user-edit-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>

