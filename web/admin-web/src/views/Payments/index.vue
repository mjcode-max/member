<template>
  <div class="payments-page">
    <div class="page-header">
      <h2>支付管理</h2>
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
        
        <el-form-item label="支付方式">
          <el-select
            v-model="searchForm.payment_method"
            placeholder="请选择支付方式"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="微信支付" value="wechat" />
            <el-option label="支付宝" value="alipay" />
            <el-option label="现金" value="cash" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="支付状态">
          <el-select
            v-model="searchForm.payment_status"
            placeholder="请选择支付状态"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option label="待支付" value="pending" />
            <el-option label="已支付" value="paid" />
            <el-option label="已退款" value="refunded" />
            <el-option label="支付失败" value="failed" />
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
        
        <el-form-item label="支付时间">
          <el-date-picker
            v-model="searchForm.payment_date"
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
          <el-button type="success" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 支付列表 -->
    <el-card class="table-card">
      <el-table
        :data="paymentList"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="order_no" label="订单号" width="140" />
        
        <el-table-column prop="customer_name" label="客户姓名" width="100" />
        
        <el-table-column prop="store.name" label="门店" width="120" />
        
        <el-table-column prop="payment_method" label="支付方式" width="100">
          <template #default="{ row }">
            <el-tag :type="getPaymentMethodType(row.payment_method)">
              {{ getPaymentMethodText(row.payment_method) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="amount" label="支付金额" width="120">
          <template #default="{ row }">
            ¥{{ (row.amount / 100).toFixed(2) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="payment_status" label="支付状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getPaymentStatusType(row.payment_status)">
              {{ getPaymentStatusText(row.payment_status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="transaction_id" label="交易号" width="160" />
        
        <el-table-column prop="paid_at" label="支付时间" width="160">
          <template #default="{ row }">
            {{ row.paid_at ? formatDateTime(row.paid_at) : '-' }}
          </template>
        </el-table-column>
        
        <el-table-column prop="refund_amount" label="退款金额" width="120">
          <template #default="{ row }">
            {{ row.refund_amount ? `¥${(row.refund_amount / 100).toFixed(2)}` : '-' }}
          </template>
        </el-table-column>
        
        <el-table-column prop="refunded_at" label="退款时间" width="160">
          <template #default="{ row }">
            {{ row.refunded_at ? formatDateTime(row.refunded_at) : '-' }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-button 
              v-if="row.payment_status === 'paid' && !row.refund_amount"
              type="warning" 
              size="small" 
              @click="handleRefund(row)"
            >
              退款
            </el-button>
            <el-button 
              v-if="row.payment_status === 'failed'"
              type="success" 
              size="small" 
              @click="handleRetry(row)"
            >
              重试
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
import { getPayments, refundPayment, retryPayment } from '@/api/payments'
import dayjs from 'dayjs'

const loading = ref(false)
const paymentList = ref([])
const storeList = ref([])
const selectedPayments = ref([])

const searchForm = reactive({
  order_no: '',
  payment_method: '',
  payment_status: '',
  store_id: '',
  payment_date: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 格式化日期时间
const formatDateTime = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 获取支付方式类型
const getPaymentMethodType = (method) => {
  const typeMap = {
    wechat: 'success',
    alipay: 'primary',
    cash: 'warning'
  }
  return typeMap[method] || 'info'
}

// 获取支付方式文本
const getPaymentMethodText = (method) => {
  const textMap = {
    wechat: '微信支付',
    alipay: '支付宝',
    cash: '现金'
  }
  return textMap[method] || '未知'
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

// 获取支付列表
const fetchPaymentList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const response = await getPayments(params)
    paymentList.value = response.data.payments || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取支付列表失败:', error)
    ElMessage.error('获取支付列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchPaymentList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    order_no: '',
    payment_method: '',
    payment_status: '',
    store_id: '',
    payment_date: ''
  })
  handleSearch()
}

// 查看详情
const handleView = (row) => {
  ElMessage.info(`查看支付详情：${row.order_no}`)
}

// 退款
const handleRefund = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要退款订单"${row.order_no}"吗？退款金额：¥${(row.amount / 100).toFixed(2)}`,
      '确认退款',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await refundPayment(row.id)
    ElMessage.success('退款申请已提交')
    fetchPaymentList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('退款失败:', error)
      ElMessage.error('退款失败')
    }
  }
}

// 重试支付
const handleRetry = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要重试支付订单"${row.order_no}"吗？`,
      '重试支付',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await retryPayment(row.id)
    ElMessage.success('支付重试已提交')
    fetchPaymentList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重试支付失败:', error)
      ElMessage.error('重试支付失败')
    }
  }
}

// 导出
const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedPayments.value = selection
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchPaymentList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchPaymentList()
}

onMounted(() => {
  fetchStoreList()
  fetchPaymentList()
})
</script>

<style lang="scss" scoped>
.payments-page {
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
