<template>
  <div class="members-page">
    <div class="page-header">
      <h2>会员管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增会员
      </el-button>
    </div>
    
    <!-- 搜索筛选 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="会员姓名">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入会员姓名"
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
        
        <el-form-item label="会员状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="有效" value="active" />
            <el-option label="已过期" value="expired" />
            <el-option label="已停用" value="inactive" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="套餐类型">
          <el-select
            v-model="searchForm.package_type"
            placeholder="请选择套餐"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="基础会员" value="basic" />
            <el-option label="高级会员" value="premium" />
            <el-option label="VIP会员" value="vip" />
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
    
    <!-- 会员列表 -->
    <el-card class="table-card">
      <el-table
        :data="memberList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="会员姓名" width="100" />
        
        <el-table-column prop="phone" label="手机号" width="120" />
        
        <el-table-column prop="store.name" label="所属门店" width="120" />
        
        <el-table-column prop="package_name" label="套餐名称" width="120" />
        
        <el-table-column prop="total_times" label="总次数" width="80" />
        
        <el-table-column prop="remaining_times" label="剩余次数" width="100">
          <template #default="{ row }">
            <el-tag :type="row.remaining_times > 0 ? 'success' : 'danger'">
              {{ row.remaining_times }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="valid_from" label="有效期开始" width="120">
          <template #default="{ row }">
            {{ formatDate(row.valid_from) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="valid_to" label="有效期结束" width="120">
          <template #default="{ row }">
            {{ formatDate(row.valid_to) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="160">
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
            <el-button type="warning" size="small" @click="handleConsumption(row)">
              消费记录
            </el-button>
            <el-button type="info" size="small" @click="handleFace(row)">
              人脸管理
            </el-button>
            <el-button 
              v-if="row.status === 'active'"
              type="danger" 
              size="small" 
              @click="handleDisable(row)"
            >
              停用
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getStores } from '@/api/stores'
import { getMembers, disableMember, enableMember } from '@/api/members'
import dayjs from 'dayjs'

const loading = ref(false)
const memberList = ref([])
const storeList = ref([])
const selectedMembers = ref([])

const searchForm = reactive({
  name: '',
  phone: '',
  store_id: '',
  status: '',
  package_type: ''
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

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    active: 'success',
    expired: 'warning',
    inactive: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const textMap = {
    active: '有效',
    expired: '已过期',
    inactive: '已停用'
  }
  return textMap[status] || '未知'
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

// 获取会员列表
const fetchMemberList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const response = await getMembers(params)
    memberList.value = response.data.members || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取会员列表失败:', error)
    ElMessage.error('获取会员列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchMemberList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    name: '',
    phone: '',
    store_id: '',
    status: '',
    package_type: ''
  })
  handleSearch()
}

// 新增会员
const handleCreate = () => {
  ElMessage.info('跳转到新增会员页面')
}

// 查看详情
const handleView = (row) => {
  ElMessage.info(`查看会员详情：${row.name}`)
}

// 编辑会员
const handleEdit = (row) => {
  ElMessage.info(`编辑会员：${row.name}`)
}

// 消费记录
const handleConsumption = (row) => {
  ElMessage.info(`查看消费记录：${row.name}`)
}

// 人脸管理
const handleFace = (row) => {
  ElMessage.info(`管理人脸数据：${row.name}`)
}

// 停用会员
const handleDisable = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要停用会员"${row.name}"吗？`,
      '停用会员',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await disableMember(row.id)
    ElMessage.success('会员停用成功')
    fetchMemberList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('停用会员失败:', error)
      ElMessage.error('停用会员失败')
    }
  }
}

// 启用会员
const handleEnable = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要启用会员"${row.name}"吗？`,
      '启用会员',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await enableMember(row.id)
    ElMessage.success('会员启用成功')
    fetchMemberList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('启用会员失败:', error)
      ElMessage.error('启用会员失败')
    }
  }
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedMembers.value = selection
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchMemberList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchMemberList()
}

onMounted(() => {
  fetchStoreList()
  fetchMemberList()
})
</script>

<style lang="scss" scoped>
.members-page {
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
