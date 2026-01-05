<template>
  <div class="stores-page">
    <div class="page-header">
      <h2>门店管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增门店
      </el-button>
    </div>
    
    <!-- 搜索筛选 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="门店名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入门店名称"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="营业中" value="operating" />
            <el-option label="停业" value="closed" />
            <el-option label="关闭" value="shutdown" />
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
    
    <!-- 门店列表 -->
    <el-card class="table-card">
      <el-table
        :data="storeList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="门店名称" min-width="120" />
        
        <el-table-column prop="address" label="门店地址" min-width="200" />
        
        <el-table-column prop="phone" label="联系电话" width="120" />
        
        <el-table-column label="营业时间" width="150">
          <template #default="{ row }">
            <span v-if="row.business_hours_start && row.business_hours_end">
              {{ row.business_hours_start }} - {{ row.business_hours_end }}
            </span>
            <span v-else>-</span>
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
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
              删除
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getStores, deleteStore } from '@/api/stores'
import dayjs from 'dayjs'

const router = useRouter()

const loading = ref(false)
const storeList = ref([])
const selectedStores = ref([])

const searchForm = reactive({
  name: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 获取门店列表
const fetchStoreList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const response = await getStores(params)
    // 后端返回格式: { code: 200, data: { list: [], pagination: { total, page, page_size } } }
    if (response.data && response.data.list) {
      storeList.value = response.data.list
      pagination.total = response.data.pagination?.total || 0
    } else {
      // 兼容旧格式
      storeList.value = response.data?.stores || []
      pagination.total = response.data?.total || 0
    }
  } catch (error) {
    console.error('获取门店列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    'operating': '营业中',
    'closed': '停业',
    'shutdown': '关闭'
  }
  return statusMap[status] || status
}

// 获取状态类型
const getStatusType = (status) => {
  if (status === 'operating') return 'success'
  if (status === 'closed') return 'warning'
  if (status === 'shutdown') return 'danger'
  return 'info'
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchStoreList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    name: '',
    status: ''
  })
  handleSearch()
}

// 新增门店
const handleCreate = () => {
  router.push('/stores/create')
}

// 编辑门店
const handleEdit = (row) => {
  router.push(`/stores/${row.id}/edit`)
}

// 删除门店
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除门店"${row.name}"吗？删除后不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteStore(row.id)
    ElMessage.success('删除成功')
    fetchStoreList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除门店失败:', error)
    }
  }
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedStores.value = selection
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchStoreList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchStoreList()
}

onMounted(() => {
  fetchStoreList()
})
</script>

<style lang="scss" scoped>
.stores-page {
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
