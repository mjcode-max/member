<template>
  <div class="bookings-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1 class="page-title">我的预约</h1>
          <p class="page-subtitle">查看和管理您的预约订单</p>
        </div>
        <div class="header-actions">
          <div class="new-booking-btn" @click="$router.push('/booking/create')">
            <i class="booking-icon">📅</i>
            <span>新预约</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 状态筛选 -->
    <div class="filter-section">
      <div class="filter-tabs">
        <div 
          v-for="tab in statusTabs"
          :key="tab.value"
          class="filter-tab"
          :class="{ active: activeTab === tab.value }"
          @click="changeTab(tab.value)"
        >
          <i class="tab-icon">{{ tab.icon }}</i>
          <span class="tab-text">{{ tab.label }}</span>
          <span class="tab-count" v-if="getTabCount(tab.value) > 0">
            {{ getTabCount(tab.value) }}
          </span>
        </div>
      </div>
    </div>

    <!-- 预约列表 -->
    <div class="bookings-container">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <div class="bookings-list">
          <div 
            v-for="booking in filteredBookings"
            :key="booking.id"
            class="booking-card"
            :class="booking.booking_status"
            @click="handleBookingClick(booking)"
          >
            <div class="booking-header">
              <div class="booking-info">
                <div class="booking-number">
                  订单号：{{ booking.order_no }}
                </div>
                <div class="booking-status" :class="booking.booking_status">
                  {{ getBookingStatusText(booking.booking_status) }}
                </div>
              </div>
              <div class="booking-time">
                {{ formatDateTime(booking.booking_date, booking.time_slot) }}
              </div>
            </div>

            <div class="booking-content">
              <div class="service-info">
                <div class="service-icon">
                  {{ booking.service_type === 'manicure' ? '💅' : '👁️' }}
                </div>
                <div class="service-details">
                  <h4 class="service-name">
                    {{ booking.service_type === 'manicure' ? '美甲服务' : '美睫服务' }}
                  </h4>
                  <p class="store-name">{{ booking.store?.name }}</p>
                  <p class="store-address">{{ booking.store?.address }}</p>
                </div>
              </div>

              <div class="booking-meta">
                <div class="meta-item">
                  <i class="meta-icon">👤</i>
                  <span>{{ booking.customer_name }}</span>
                </div>
                <div class="meta-item">
                  <i class="meta-icon">📱</i>
                  <span>{{ formatPhone(booking.customer_phone) }}</span>
                </div>
                <div class="meta-item" v-if="booking.remark">
                  <i class="meta-icon">📝</i>
                  <span>{{ booking.remark }}</span>
                </div>
              </div>

              <div class="payment-info">
                <div class="payment-item">
                  <span class="payment-label">押金：</span>
                  <span class="payment-amount">¥{{ (booking.deposit_amount / 100).toFixed(2) }}</span>
                  <span class="payment-status" :class="booking.payment_status">
                    {{ getPaymentStatusText(booking.payment_status) }}
                  </span>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="booking-actions" v-if="showActions(booking)">
                <button 
                  v-if="canCancel(booking)"
                  class="action-btn cancel-btn"
                  @click.stop="cancelBooking(booking)"
                >
                  <i class="btn-icon">❌</i>
                  取消预约
                </button>
                <button 
                  v-if="canPay(booking)"
                  class="action-btn pay-btn"
                  @click.stop="payDeposit(booking)"
                >
                  <i class="btn-icon">💳</i>
                  支付押金
                </button>
                <button 
                  v-if="canRebook(booking)"
                  class="action-btn rebook-btn"
                  @click.stop="rebookService(booking)"
                >
                  <i class="btn-icon">🔄</i>
                  再次预约
                </button>
              </div>
            </div>
          </div>
        </div>
      </van-pull-refresh>

      <!-- 空状态 -->
      <div v-if="filteredBookings.length === 0" class="empty-state">
        <div class="empty-icon">📅</div>
        <div class="empty-title">暂无{{ getTabLabel(activeTab) }}</div>
        <div class="empty-desc">
          <span v-if="activeTab === 'all'">您还没有任何预约记录</span>
          <span v-else>当前状态下没有相关预约</span>
        </div>
        <button class="empty-action" @click="$router.push('/booking/create')">
          立即预约
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import PullRefresh from 'vant/es/pull-refresh'
import { getMyBookings, cancelBooking as cancelBookingApi } from '@/api/bookings'
import dayjs from 'dayjs'

const router = useRouter()

const activeTab = ref('all')
const refreshing = ref(false)
const bookingList = ref([])

// 状态标签页
const statusTabs = [
  { value: 'all', label: '全部', icon: '📋' },
  { value: 'pending', label: '待支付', icon: '💳' },
  { value: 'confirmed', label: '已确认', icon: '✅' },
  { value: 'completed', label: '已完成', icon: '🎉' },
  { value: 'cancelled', label: '已取消', icon: '❌' }
]

// 筛选后的预约列表
const filteredBookings = computed(() => {
  if (activeTab.value === 'all') {
    return bookingList.value
  }
  return bookingList.value.filter(booking => booking.booking_status === activeTab.value)
})

// 获取预约列表
const fetchBookings = async () => {
  try {
    const response = await getMyBookings()
    bookingList.value = response.data.bookings || []
  } catch (error) {
    console.error('获取预约列表失败:', error)
  }
}

// 切换标签页
const changeTab = (tab) => {
  activeTab.value = tab
}

// 获取标签数量
const getTabCount = (status) => {
  if (status === 'all') return bookingList.value.length
  return bookingList.value.filter(b => b.booking_status === status).length
}

// 获取标签文本
const getTabLabel = (status) => {
  const tab = statusTabs.find(t => t.value === status)
  return tab ? tab.label : '预约'
}

// 获取预约状态文本
const getBookingStatusText = (status) => {
  const statusMap = {
    pending: '待支付',
    confirmed: '已确认',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || '未知'
}

// 获取支付状态文本
const getPaymentStatusText = (status) => {
  const statusMap = {
    pending: '待支付',
    paid: '已支付',
    refunded: '已退款',
    failed: '支付失败'
  }
  return statusMap[status] || '未知'
}

// 格式化日期时间
const formatDateTime = (date, timeSlot) => {
  return `${dayjs(date).format('MM月DD日')} ${timeSlot}`
}

// 格式化手机号
const formatPhone = (phone) => {
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 是否显示操作按钮
const showActions = (booking) => {
  return canCancel(booking) || canPay(booking) || canRebook(booking)
}

// 是否可以取消
const canCancel = (booking) => {
  if (booking.booking_status === 'cancelled' || booking.booking_status === 'completed') {
    return false
  }
  
  // 检查是否提前3小时
  const serviceTime = dayjs(`${booking.booking_date} ${booking.time_slot.split('-')[0]}`)
  const now = dayjs()
  const hoursDiff = serviceTime.diff(now, 'hour')
  
  return hoursDiff >= 3
}

// 是否可以支付
const canPay = (booking) => {
  return booking.booking_status === 'pending' && booking.payment_status === 'pending'
}

// 是否可以再次预约
const canRebook = (booking) => {
  return booking.booking_status === 'completed' || booking.booking_status === 'cancelled'
}

// 处理预约点击
const handleBookingClick = (booking) => {
  router.push(`/booking/${booking.id}`)
}

// 取消预约
const cancelBooking = async (booking) => {
  try {
    await showConfirmDialog({
      title: '确认取消',
      message: '确定要取消这个预约吗？押金将原路退回。'
    })

    await cancelBookingApi(booking.id)
    showToast('预约已取消，押金将在3-5个工作日内退回')
    fetchBookings()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消预约失败:', error)
      showToast('取消失败，请重试')
    }
  }
}

// 支付押金
const payDeposit = (booking) => {
  router.push(`/booking/${booking.id}/payment`)
}

// 再次预约
const rebookService = (booking) => {
  router.push(`/booking/create?service=${booking.service_type}&store=${booking.store_id}`)
}

// 下拉刷新
const onRefresh = async () => {
  refreshing.value = true
  await fetchBookings()
  refreshing.value = false
}

onMounted(() => {
  fetchBookings()
})
</script>

<style lang="scss" scoped>
.bookings-page {
  background: #ffffff;
  min-height: 100vh;
}

// 页面头部
.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  color: white;
  position: relative;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -15px;
    left: 0;
    right: 0;
    height: 15px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 0 0 15px 15px;
  }
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.page-subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

.new-booking-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.3);
  }
}

// 筛选区域
.filter-section {
  padding: 16px;
}

.filter-tabs {
  display: flex;
  background: white;
  border-radius: 12px;
  padding: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  overflow-x: auto;
  
  &::-webkit-scrollbar {
    display: none;
  }
}

.filter-tab {
  flex: 1;
  min-width: 60px;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  white-space: nowrap;
  
  &.active {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    transform: scale(1.05);
  }
}

.tab-icon {
  font-size: 16px;
  margin-bottom: 4px;
}

.tab-text {
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 2px;
}

.tab-count {
  font-size: 10px;
  background: rgba(255, 255, 255, 0.2);
  padding: 2px 6px;
  border-radius: 8px;
  min-width: 16px;
  text-align: center;
  
  .filter-tab.active & {
    background: rgba(255, 255, 255, 0.3);
  }
  
  .filter-tab:not(.active) & {
    background: #f0f0f0;
    color: #999;
  }
}

// 预约列表
.bookings-container {
  padding: 0 16px;
}

.bookings-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.booking-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.3s ease;
  border-left: 4px solid transparent;
  
  &.pending {
    border-left-color: #fa8c16;
  }
  
  &.confirmed {
    border-left-color: #1890ff;
  }
  
  &.completed {
    border-left-color: #52c41a;
  }
  
  &.cancelled {
    border-left-color: #ff4d4f;
    opacity: 0.7;
  }
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  }
}

.booking-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.booking-info {
  flex: 1;
}

.booking-number {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.booking-status {
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  display: inline-block;
  
  &.pending {
    background: #fff7e6;
    color: #fa8c16;
  }
  
  &.confirmed {
    background: #e6f7ff;
    color: #1890ff;
  }
  
  &.completed {
    background: #f6ffed;
    color: #52c41a;
  }
  
  &.cancelled {
    background: #fff2f0;
    color: #ff4d4f;
  }
}

.booking-time {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  text-align: right;
}

.booking-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.service-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.service-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.service-details {
  flex: 1;
}

.service-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 4px 0;
}

.store-name {
  font-size: 14px;
  color: #666;
  margin: 0 0 2px 0;
}

.store-address {
  font-size: 12px;
  color: #999;
  margin: 0;
}

.booking-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}

.meta-icon {
  font-size: 14px;
  width: 16px;
  text-align: center;
}

.payment-info {
  padding: 12px;
  background: #f0f2ff;
  border-radius: 8px;
}

.payment-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.payment-label {
  font-size: 14px;
  color: #666;
}

.payment-amount {
  font-size: 16px;
  font-weight: 600;
  color: #ff6b6b;
}

.payment-status {
  padding: 2px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  
  &.pending {
    background: #fff7e6;
    color: #fa8c16;
  }
  
  &.paid {
    background: #e6f7ff;
    color: #1890ff;
  }
  
  &.refunded {
    background: #f6ffed;
    color: #52c41a;
  }
}

.booking-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.action-btn {
  flex: 1;
  min-width: 80px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  
  &.cancel-btn {
    background: #fff2f0;
    color: #ff4d4f;
    
    &:hover {
      background: #ff4d4f;
      color: white;
    }
  }
  
  &.pay-btn {
    background: linear-gradient(135deg, #ff6b6b, #ffa726);
    color: white;
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
    }
  }
  
  &.rebook-btn {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
    }
  }
}

.btn-icon {
  font-size: 12px;
}

// 空状态
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 8px;
  color: #666;
}

.empty-desc {
  font-size: 14px;
  margin-bottom: 24px;
  line-height: 1.4;
}

.empty-action {
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
  }
}

// 响应式设计
@media (max-width: 375px) {
  .booking-card {
    padding: 16px;
  }
  
  .service-icon {
    width: 40px;
    height: 40px;
    font-size: 20px;
  }
  
  .booking-actions {
    flex-direction: column;
  }
  
  .action-btn {
    min-width: auto;
  }
}

// 动画效果
@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.booking-card {
  animation: slideInUp 0.4s ease-out;
}

.booking-card:nth-child(1) { animation-delay: 0.1s; }
.booking-card:nth-child(2) { animation-delay: 0.2s; }
.booking-card:nth-child(3) { animation-delay: 0.3s; }
.booking-card:nth-child(4) { animation-delay: 0.4s; }
.booking-card:nth-child(5) { animation-delay: 0.5s; }
</style>
