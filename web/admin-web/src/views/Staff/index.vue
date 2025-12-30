<template>
  <div class="staff-page">
    <div class="page-header">
      <h2>员工管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增员工
      </el-button>
    </div>
    
    <!-- 搜索筛选 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="员工姓名">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入员工姓名"
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
        
        <el-form-item label="职位">
          <el-select
            v-model="searchForm.position"
            placeholder="请选择职位"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="店长" value="manager" />
            <el-option label="技师" value="technician" />
            <el-option label="前台" value="receptionist" />
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
            <el-option label="在职" value="active" />
            <el-option label="离职" value="inactive" />
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
    
    <!-- 员工列表 -->
    <el-card class="table-card">
      <el-table
        :data="staffList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="员工姓名" width="100" />
        
        <el-table-column prop="phone" label="手机号" width="120" />
        
        <el-table-column prop="store.name" label="所属门店" width="120" />
        
        <el-table-column prop="position" label="职位" width="100">
          <template #default="{ row }">
            <el-tag :type="getPositionType(row.position)">
              {{ getPositionText(row.position) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="service_types" label="服务类型" width="120">
          <template #default="{ row }">
            <el-tag 
              v-for="type in row.service_types" 
              :key="type"
              :type="type === 'manicure' ? 'primary' : 'success'"
              size="small"
              style="margin-right: 4px;"
            >
              {{ type === 'manicure' ? '美甲' : '美睫' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="work_schedule" label="工作时间" width="120">
          <template #default="{ row }">
            {{ row.work_schedule || '09:00-18:00' }}
          </template>
        </el-table-column>
        
        <el-table-column prop="salary" label="薪资" width="100">
          <template #default="{ row }">
            ¥{{ row.salary || 0 }}
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '在职' : '离职' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="入职时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="400" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-button type="success" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="warning" size="small" @click="handleSchedule(row)">
              排班
            </el-button>
            <el-button 
              v-if="row.status === 'active'"
              type="danger" 
              size="small" 
              @click="handleResign(row)"
            >
              离职
            </el-button>
            <el-button 
              v-if="row.status === 'inactive'"
              type="success" 
              size="small" 
              @click="handleRehire(row)"
            >
              复职
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getStores } from '@/api/stores'
import { getStaff, resignStaff, rehireStaff } from '@/api/staff'
import dayjs from 'dayjs'

const loading = ref(false)
const staffList = ref([])
const storeList = ref([])
const selectedStaff = ref([])

const searchForm = reactive({
  name: '',
  phone: '',
  store_id: '',
  position: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

// 获取职位类型
const getPositionType = (position) => {
  const typeMap = {
    manager: 'danger',
    technician: 'primary',
    receptionist: 'success'
  }
  return typeMap[position] || 'info'
}

// 获取职位文本
const getPositionText = (position) => {
  const textMap = {
    manager: '店长',
    technician: '技师',
    receptionist: '前台'
  }
  return textMap[position] || '未知'
}

// 获取门店列表
const fetchStoreList = async () => {
  try {
    const response = await getStores()
    storeList.value = response.data.stores || []
  } catch (error) {
    console.error('获取门店列表失败:', error)
  }
}

// 获取员工列表
const fetchStaffList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const response = await getStaff(params)
    staffList.value = response.data.staff || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取员工列表失败:', error)
    ElMessage.error('获取员工列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchStaffList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    name: '',
    phone: '',
    store_id: '',
    position: '',
    status: ''
  })
  handleSearch()
}

// 新增员工
const handleCreate = () => {
  ElMessage.info('跳转到新增员工页面')
}

// 查看详情
const handleView = (row) => {
  ElMessage.info(`查看员工详情：${row.name}`)
}

// 编辑员工
const handleEdit = (row) => {
  ElMessage.info(`编辑员工：${row.name}`)
}

// 排班管理
const handleSchedule = (row) => {
  ElMessage.info(`管理排班：${row.name}`)
}

// 离职
const handleResign = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要让员工"${row.name}"离职吗？`,
      '员工离职',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await resignStaff(row.id)
    ElMessage.success('员工离职成功')
    fetchStaffList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('员工离职失败:', error)
      ElMessage.error('员工离职失败')
    }
  }
}

// 复职
const handleRehire = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要让员工"${row.name}"复职吗？`,
      '员工复职',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await rehireStaff(row.id)
    ElMessage.success('员工复职成功')
    fetchStaffList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('员工复职失败:', error)
      ElMessage.error('员工复职失败')
    }
  }
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedStaff.value = selection
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchStaffList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchStaffList()
}

onMounted(() => {
  fetchStoreList()
  fetchStaffList()
})
</script>

<style lang="scss" scoped>
.staff-page {
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
