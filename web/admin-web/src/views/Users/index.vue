<template>
  <div class="users-page">
    <div class="page-header">
      <h2>用户管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增用户
      </el-button>
    </div>
    
    <!-- 搜索筛选 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="用户名">
          <el-input
            v-model="searchForm.username"
            placeholder="请输入用户名"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        
        <el-form-item label="手机号">
          <el-input
            v-model="searchForm.phone"
            placeholder="请输入手机号"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        
        <el-form-item label="角色">
          <el-select
            v-model="searchForm.role"
            placeholder="请选择角色"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="总后台" value="admin" />
            <el-option label="店长" value="store_manager" />
            <el-option label="美甲师" value="technician" />
            <el-option label="顾客" value="customer" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="激活" value="active" />
            <el-option label="禁用" value="inactive" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="门店">
          <el-select
            v-model="searchForm.store_id"
            placeholder="请选择门店"
            clearable
            style="width: 200px"
            @change="handleSearch"
          >
            <el-option
              v-for="store in storeList"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 用户列表 -->
    <el-card class="table-card">
      <el-table
        :data="userList"
        v-loading="loading"
        stripe
        empty-text="暂无数据"
      >
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column prop="username" label="用户名" width="120" />
        
        <el-table-column prop="phone" label="手机号" width="120" />
        
        <el-table-column prop="email" label="邮箱" width="180" />
        
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="store_id" label="所属门店" width="120">
          <template #default="{ row }">
            {{ getStoreName(row.store_id) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="work_status" label="工作状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.work_status" :type="getWorkStatusType(row.work_status)">
              {{ getWorkStatusText(row.work_status) }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button 
              v-if="row.role === 'technician'"
              type="warning" 
              size="small" 
              @click="handleUpdateWorkStatus(row)"
            >
              工作状态
            </el-button>
            <el-button 
              v-if="row.status === 'active'"
              type="danger" 
              size="small" 
              @click="handleDisable(row)"
            >
              禁用
            </el-button>
            <el-button 
              v-if="row.status === 'inactive'"
              type="success" 
              size="small" 
              @click="handleEnable(row)"
            >
              启用
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 工作状态对话框 -->
    <el-dialog v-model="workStatusDialogVisible" title="更新工作状态" width="400px">
      <el-form :model="workStatusForm" label-width="100px">
        <el-form-item label="工作状态">
          <el-select v-model="workStatusForm.work_status" style="width: 100%">
            <el-option label="在岗" value="working" />
            <el-option label="休息" value="rest" />
            <el-option label="离岗" value="offline" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="workStatusDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleConfirmWorkStatus">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onActivated } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUsers, updateUserStatus, updateWorkStatus } from '@/api/users'
import dayjs from 'dayjs'

const router = useRouter()

const loading = ref(false)
const userList = ref([])
const storeList = ref([])

const searchForm = reactive({
  username: '',
  phone: '',
  role: '',
  status: '',
  store_id: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const workStatusDialogVisible = ref(false)
const workStatusForm = reactive({
  user_id: null,
  work_status: ''
})

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 获取角色类型
const getRoleType = (role) => {
  const typeMap = {
    admin: 'danger',
    store_manager: 'warning',
    technician: 'primary',
    customer: 'success'
  }
  return typeMap[role] || 'info'
}

// 获取角色文本
const getRoleText = (role) => {
  const textMap = {
    admin: '总后台',
    store_manager: '店长',
    technician: '美甲师',
    customer: '顾客'
  }
  return textMap[role] || '未知'
}

// 获取工作状态类型
const getWorkStatusType = (status) => {
  const typeMap = {
    working: 'success',
    rest: 'warning',
    offline: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取工作状态文本
const getWorkStatusText = (status) => {
  const textMap = {
    working: '在岗',
    rest: '休息',
    offline: '离岗'
  }
  return textMap[status] || '未知'
}

// 获取门店名称
const getStoreName = (storeId) => {
  if (!storeId) return '-'
  const store = storeList.value.find(s => s.id === storeId)
  return store ? store.name : '-'
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

// 获取用户列表
const fetchUserList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      role: searchForm.role || undefined,
      status: searchForm.status || undefined,
      store_id: searchForm.store_id || undefined,
      username: searchForm.username || undefined,
      phone: searchForm.phone || undefined
    }
    
    const response = await getUsers(params)
    console.log('用户列表API响应:', response)
    
    // 响应拦截器返回的是 { code, message, data }
    // 对于分页接口，data 的结构是 { list: [...], pagination: {...} }
    if (response && response.data) {
      // 检查是否是分页响应格式
      if (response.data.list && Array.isArray(response.data.list)) {
        userList.value = response.data.list
        pagination.total = response.data.pagination?.total || 0
        console.log('解析后的用户列表:', userList.value)
        console.log('解析后的总数:', pagination.total)
      } else if (Array.isArray(response.data)) {
        // 如果不是分页格式，直接是数组
        userList.value = response.data
        pagination.total = response.data.length
      } else {
        console.warn('响应数据格式不正确:', response.data)
        userList.value = []
        pagination.total = 0
      }
    } else {
      console.warn('响应为空或没有data字段:', response)
      userList.value = []
      pagination.total = 0
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
    userList.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchUserList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    username: '',
    phone: '',
    role: '',
    status: '',
    store_id: ''
  })
  handleSearch()
}

// 新增用户
const handleCreate = () => {
  router.push('/users/create')
}

// 编辑用户
const handleEdit = (row) => {
  router.push(`/users/${row.id}/edit`)
}

// 更新工作状态
const handleUpdateWorkStatus = (row) => {
  workStatusForm.user_id = row.id
  workStatusForm.work_status = row.work_status || 'working'
  workStatusDialogVisible.value = true
}

// 确认更新工作状态
const handleConfirmWorkStatus = async () => {
  try {
    await updateWorkStatus(workStatusForm.user_id, workStatusForm.work_status)
    ElMessage.success('更新工作状态成功')
    workStatusDialogVisible.value = false
    fetchUserList()
  } catch (error) {
    console.error('更新工作状态失败:', error)
    ElMessage.error('更新工作状态失败')
  }
}

// 禁用用户
const handleDisable = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要禁用用户"${row.username || row.phone}"吗？`,
      '禁用用户',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await updateUserStatus(row.id, 'inactive')
    ElMessage.success('用户禁用成功')
    fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('禁用用户失败:', error)
      ElMessage.error('禁用用户失败')
    }
  }
}

// 启用用户
const handleEnable = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要启用用户"${row.username || row.phone}"吗？`,
      '启用用户',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await updateUserStatus(row.id, 'active')
    ElMessage.success('用户启用成功')
    fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('启用用户失败:', error)
      ElMessage.error('启用用户失败')
    }
  }
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchUserList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchUserList()
}

onMounted(() => {
  fetchStoreList()
  fetchUserList()
})

// 页面激活时刷新数据（从其他页面返回时）
onActivated(() => {
  fetchUserList()
})
</script>

<style lang="scss" scoped>
.users-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      color: #333;
      margin: 0;
    }
  }
  
  .search-card {
    margin-bottom: 20px;
  }
  
  .table-card {
    .pagination-container {
      margin-top: 20px;
      display: flex;
      justify-content: flex-end;
    }
  }
}
</style>

