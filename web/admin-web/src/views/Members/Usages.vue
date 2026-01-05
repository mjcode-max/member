<template>
  <div class="member-usages-page">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>使用记录 - {{ memberName }}</h2>
      <el-button type="primary" @click="handleCreateUsage">
        <el-icon><Plus /></el-icon>
        新增使用记录
      </el-button>
    </div>
    
    <el-card class="usages-card" v-loading="loading">
      <el-table :data="usageList" stripe>
        <el-table-column prop="service_item" label="服务项目" min-width="150" />
        <el-table-column prop="package_name" label="套餐名称" width="150" />
        <el-table-column prop="store_name" label="使用门店" width="150" />
        <el-table-column prop="technician_name" label="美甲师" width="120">
          <template #default="{ row }">
            {{ row.technician_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="usage_date" label="使用日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.usage_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="200">
          <template #default="{ row }">
            {{ row.remark || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button 
              type="danger" 
              size="small" 
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-if="usageList.length === 0 && !loading" class="empty-state">
        <el-empty description="暂无使用记录" />
      </div>
    </el-card>
    
    <!-- 新增使用记录对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="新增使用记录"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="服务项目" prop="service_item">
          <el-input
            v-model="form.service_item"
            placeholder="请输入服务项目，如：美甲-单色"
            maxlength="100"
          />
        </el-form-item>
        
        <el-form-item label="使用门店" prop="store_id">
          <el-select
            v-model="form.store_id"
            placeholder="请选择门店"
            style="width: 100%"
          >
            <el-option
              v-for="store in storeList"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="美甲师" prop="technician_id">
          <el-select
            v-model="form.technician_id"
            placeholder="请选择美甲师（可选）"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="technician in technicianList"
              :key="technician.id"
              :label="technician.username"
              :value="technician.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="使用日期" prop="usage_date">
          <el-date-picker
            v-model="form.usage_date"
            type="date"
            placeholder="选择使用日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Plus } from '@element-plus/icons-vue'
import { getMemberById, getMemberUsages, createUsage, deleteUsage } from '@/api/members'
import { getStores } from '@/api/stores'
import { getUsers } from '@/api/users'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const submitting = ref(false)
const usageList = ref([])
const storeList = ref([])
const technicianList = ref([])
const memberName = ref('')
const dialogVisible = ref(false)

const formRef = ref(null)
const form = reactive({
  service_item: '',
  store_id: null,
  technician_id: null,
  usage_date: dayjs().format('YYYY-MM-DD'),
  remark: ''
})

const rules = {
  service_item: [
    { required: true, message: '请输入服务项目', trigger: 'blur' }
  ],
  store_id: [
    { required: true, message: '请选择使用门店', trigger: 'change' }
  ],
  usage_date: [
    { required: true, message: '请选择使用日期', trigger: 'change' }
  ]
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

// 格式化日期时间
const formatDateTime = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 获取会员信息
const fetchMemberInfo = async () => {
  const memberId = route.params.id
  if (!memberId) {
    ElMessage.error('会员ID不存在')
    router.back()
    return
  }
  
  try {
    const response = await getMemberById(memberId)
    if (response.data) {
      memberName.value = response.data.name || '未知会员'
    }
  } catch (error) {
    console.error('获取会员信息失败:', error)
  }
}

// 获取使用记录列表
const fetchUsageList = async () => {
  const memberId = route.params.id
  if (!memberId) return
  
  loading.value = true
  try {
    const response = await getMemberUsages(memberId)
    usageList.value = response.data || []
  } catch (error) {
    console.error('获取使用记录失败:', error)
    ElMessage.error('获取使用记录失败')
    usageList.value = []
  } finally {
    loading.value = false
  }
}

// 获取门店列表
const fetchStoreList = async () => {
  try {
    const response = await getStores({ page: 1, page_size: 1000 })
    if (response.code === 0 && response.data) {
      storeList.value = response.data.list || []
    }
  } catch (error) {
    console.error('获取门店列表失败:', error)
  }
}

// 获取美甲师列表
const fetchTechnicianList = async () => {
  try {
    const response = await getUsers({ role: 'technician', page: 1, page_size: 1000 })
    if (response.code === 0 && response.data) {
      technicianList.value = response.data.list || []
    }
  } catch (error) {
    console.error('获取美甲师列表失败:', error)
  }
}

// 返回
const goBack = () => {
  router.back()
}

// 新增使用记录
const handleCreateUsage = () => {
  dialogVisible.value = true
}

// 重置表单
const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    service_item: '',
    store_id: null,
    technician_id: null,
    usage_date: dayjs().format('YYYY-MM-DD'),
    remark: ''
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    const memberId = route.params.id
    submitting.value = true
    
    try {
      const data = {
        service_item: form.service_item,
        store_id: form.store_id,
        technician_id: form.technician_id || undefined,
        usage_date: form.usage_date,
        remark: form.remark || undefined
      }
      
      const submitData = {
        service_item: form.service_item,
        store_id: form.store_id,
        technician_id: form.technician_id || undefined,
        usage_date: form.usage_date ? dayjs(form.usage_date).startOf('day').toISOString() : undefined,
        remark: form.remark || undefined
      }
      
      await createUsage(memberId, submitData)
      ElMessage.success('使用记录创建成功')
      dialogVisible.value = false
      fetchUsageList()
    } catch (error) {
      console.error('创建使用记录失败:', error)
      ElMessage.error('创建使用记录失败')
    } finally {
      submitting.value = false
    }
  })
}

// 删除使用记录
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除这条使用记录吗？`,
      '删除使用记录',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteUsage(row.id)
    ElMessage.success('删除成功')
    fetchUsageList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除使用记录失败:', error)
      ElMessage.error('删除使用记录失败')
    }
  }
}

onMounted(() => {
  fetchMemberInfo()
  fetchUsageList()
  fetchStoreList()
  fetchTechnicianList()
})
</script>

<style lang="scss" scoped>
.member-usages-page {
  .page-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    
    h2 {
      color: #333;
      margin: 0;
      flex: 1;
    }
  }
  
  .usages-card {
    .empty-state {
      padding: 40px 0;
    }
  }
}
</style>

