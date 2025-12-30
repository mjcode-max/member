<template>
  <div class="bookings-page">
    <div class="page-header">
      <h2>预约管理</h2>
    </div>
    
    <!-- 搜索筛选 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="订单号">
          <el-input
            v-model="searchForm.order_no"
            placeholder="请输入订单号"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        
        <el-form-item label="客户姓名">
          <el-input
            v-model="searchForm.customer_name"
            placeholder="请输入客户姓名"
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
        
        <el-form-item label="预约状态">
          <el-select
            v-model="searchForm.booking_status"
            placeholder="请选择状态"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="待支付" value="pending" />
            <el-option label="已确认" value="confirmed" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="预约日期">
          <el-date-picker
            v-model="searchForm.booking_date"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="handleSearch"
          />
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
    
    <!-- 预约列表 -->
    <el-card class="table-card">
      <el-table
        :data="bookingList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="order_no" label="订单号" width="140" />
        
        <el-table-column prop="customer_name" label="客户姓名" width="100" />
        
        <el-table-column prop="customer_phone" label="联系电话" width="120" />
        
        <el-table-column prop="store.name" label="门店" width="120" />
        
        <el-table-column prop="service_type" label="服务类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.service_type === 'manicure' ? 'primary' : 'success'">
              {{ row.service_type === 'manicure' ? '美甲' : '美睫' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="booking_date" label="预约日期" width="120">
          <template #default="{ row }">
            {{ formatBookingDate(row.booking_date) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="time_slot" label="时间段" width="120" />
        
        <el-table-column prop="booking_status" label="预约状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.booking_status)">
              {{ getStatusText(row.booking_status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="payment_status" label="支付状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getPaymentStatusType(row.payment_status)">
              {{ getPaymentStatusText(row.payment_status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="deposit_amount" label="押金金额" width="100">
          <template #default="{ row }">
            ¥{{ (row.deposit_amount / 100).toFixed(2) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="350" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-button 
              v-if="row.booking_status === 'pending'"
              type="success" 
              size="small" 
              @click="handleConfirm(row)"
            >
              确认
            </el-button>
            <el-button 
              v-if="row.booking_status === 'confirmed'"
              type="warning" 
              size="small" 
              @click="handleComplete(row)"
            >
              完成
            </el-button>
            <el-button 
              v-if="['pending', 'confirmed'].includes(row.booking_status)"
              type="danger" 
              size="small" 
              @click="handleCancel(row)"
            >
              取消
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
import { getBookings, confirmBooking, completeBooking, cancelBooking } from '@/api/bookings'
import dayjs from 'dayjs'

const loading = ref(false)
const bookingList = ref([])
const storeList = ref([])
const selectedBookings = ref([])

const searchForm = reactive({
  order_no: '',
  customer_name: '',
  store_id: '',
  booking_status: '',
  booking_date: ''
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

// 格式化预约日期
const formatBookingDate = (date) => {
  if (!date) return '-'
  // 如果是ISO格式的时间戳，先解析
  if (date.includes('T') || date.includes('+')) {
    return dayjs(date).format('YYYY-MM-DD')
  }
  // 如果已经是YYYY-MM-DD格式，直接返回
  return date
}

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    confirmed: 'primary',
    completed: 'success',
    cancelled: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const textMap = {
    pending: '待支付',
    confirmed: '已确认',
    completed: '已完成',
    cancelled: '已取消'
  }
  return textMap[status] || '未知'
}

// 获取支付状态类型
const getPaymentStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    paid: 'success',
    refunded: 'info',
    failed: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取支付状态文本
const getPaymentStatusText = (status) => {
  const textMap = {
    pending: '待支付',
    paid: '已支付',
    refunded: '已退款',
    failed: '支付失败'
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

// 获取预约列表
const fetchBookingList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const response = await getBookings(params)
    bookingList.value = response.data.bookings || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取预约列表失败:', error)
    ElMessage.error('获取预约列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchBookingList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    order_no: '',
    customer_name: '',
    store_id: '',
    booking_status: '',
    booking_date: ''
  })
  handleSearch()
}

// 查看详情
const handleView = (row) => {
  ElMessage.info(`查看预约详情：${row.order_no}`)
}

// 确认预约
const handleConfirm = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要确认预约"${row.order_no}"吗？`,
      '确认预约',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await confirmBooking(row.id)
    ElMessage.success('预约确认成功')
    fetchBookingList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('确认预约失败:', error)
      ElMessage.error('确认预约失败')
    }
  }
}

// 完成预约
const handleComplete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要完成预约"${row.order_no}"吗？`,
      '完成预约',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await completeBooking(row.id)
    ElMessage.success('预约完成成功')
    fetchBookingList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('完成预约失败:', error)
      ElMessage.error('完成预约失败')
    }
  }
}

// 取消预约
const handleCancel = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要取消预约"${row.order_no}"吗？`,
      '取消预约',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await cancelBooking(row.id)
    ElMessage.success('预约取消成功')
    fetchBookingList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消预约失败:', error)
      ElMessage.error('取消预约失败')
    }
  }
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedBookings.value = selection
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchBookingList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchBookingList()
}

onMounted(() => {
  fetchStoreList()
  fetchBookingList()
})
</script>

<style lang="scss" scoped>
.bookings-page {
  .page-header {
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
