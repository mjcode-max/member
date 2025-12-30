<template>
  <div class="bookings-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1 class="page-title">我的日程</h1>
          <p class="page-subtitle">今日预约 {{ todayBookingsCount }} 位客户</p>
        </div>
        <div class="header-actions">
          <div class="date-selector" @click="showDatePicker = true">
            <i class="date-icon">📅</i>
            <span>{{ currentDate }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 快捷统计 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-number">{{ todayBookingsCount }}</div>
          <div class="stat-label">今日预约</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ completedCount }}</div>
          <div class="stat-label">已完成</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ pendingCount }}</div>
          <div class="stat-label">待服务</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ totalEarnings }}</div>
          <div class="stat-label">今日收入</div>
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
          <span class="tab-count">{{ getTabCount(tab.value) }}</span>
        </div>
      </div>
    </div>
    
    <!-- 预约列表 -->
    <div class="bookings-container">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <div class="bookings-list">
          <div 
            v-for="booking in bookingList"
            :key="booking.id"
            class="booking-card"
            @click="handleBookingClick(booking)"
          >
            <div class="booking-header">
              <div class="customer-info">
                <div class="customer-avatar">
                  {{ booking.customer_name.charAt(0) }}
                </div>
                <div class="customer-details">
                  <h4 class="customer-name">{{ booking.customer_name }}</h4>
                  <p class="customer-phone">{{ formatPhone(booking.customer_phone) }}</p>
                </div>
              </div>
              <div class="booking-status" :class="booking.booking_status">
                {{ getBookingStatusText(booking.booking_status) }}
              </div>
            </div>
            
            <div class="booking-content">
              <div class="booking-info">
                <div class="info-item">
                  <i class="info-icon">⏰</i>
                  <span>{{ booking.time_slot }}</span>
                </div>
                <div class="info-item">
                  <i class="info-icon">💄</i>
                  <span>{{ getServiceText(booking.service_type) }}</span>
                </div>
                <div class="info-item" v-if="booking.remark">
                  <i class="info-icon">📝</i>
                  <span>{{ booking.remark }}</span>
                </div>
              </div>
              
              <div class="booking-actions" v-if="booking.booking_status === 'confirmed'">
                <button class="action-btn complete-btn" @click.stop="completeBooking(booking)">
                  <i class="btn-icon">✅</i>
                  完成服务
                </button>
                <button class="action-btn cancel-btn" @click.stop="cancelBooking(booking)">
                  <i class="btn-icon">❌</i>
                  取消预约
                </button>
              </div>
            </div>
          </div>
        </div>
      </van-pull-refresh>
      
      <!-- 空状态 -->
      <div v-if="!loading && bookingList.length === 0" class="empty-state">
        <div class="empty-icon">📅</div>
        <div class="empty-title">暂无{{ getTabLabel(activeTab) }}</div>
        <div class="empty-desc">当前没有相关的预约记录</div>
      </div>
    </div>

    <!-- 日期选择器 -->
    <van-popup 
      v-model:show="showDatePicker" 
      position="bottom" 
      round
      :style="{ height: '50%' }"
    >
      <div class="popup-header">
        <div class="popup-title">选择日期</div>
        <div class="popup-close" @click="showDatePicker = false">×</div>
      </div>
      <van-date-picker
        v-model="selectedDate"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>

    <!-- 扫码核销按钮 -->
    <div class="floating-actions">
      <div class="fab-button" @click="openScanner">
        <i class="fab-icon">📱</i>
        <span class="fab-text">扫码核销</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { showToast } from 'vant'
import { getStaffBookings } from '@/api/bookings'
import dayjs from 'dayjs'

const activeTab = ref('confirmed')
const loading = ref(false)
const refreshing = ref(false)
const finished = ref(false)
const bookingList = ref([])
const showDatePicker = ref(false)
const selectedDate = ref(new Date())
const currentDate = ref(dayjs().format('MM月DD日'))

// 统计数据
const todayBookingsCount = ref(8)
const completedCount = ref(3)
const pendingCount = ref(5)
const totalEarnings = ref('￥280')

// 状态标签页
const statusTabs = [
  { value: 'confirmed', label: '待服务', icon: '⏳' },
  { value: 'completed', label: '已完成', icon: '✅' },
  { value: 'pending', label: '待确认', icon: '⏰' },
  { value: 'cancelled', label: '已取消', icon: '❌' }
]

// 获取预约状态文本
const getBookingStatusText = (status) => {
  const statusMap = {
    pending: '待确认',
    confirmed: '待服务',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || '未知'
}

// 获取服务类型文本
const getServiceText = (type) => {
  return type === 'manicure' ? '美甲服务' : '美睫服务'
}

// 格式化手机号
const formatPhone = (phone) => {
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 获取标签数量
const getTabCount = (status) => {
  return bookingList.value.filter(b => b.booking_status === status).length
}

// 获取标签文本
const getTabLabel = (status) => {
  const tab = statusTabs.find(t => t.value === status)
  return tab ? tab.label : '预约'
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('MM-DD')
}

// 获取预约列表
const fetchBookings = async () => {
  try {
    const params = {
      status: activeTab.value,
      page: 1,
      page_size: 20
    }
    
    const response = await getStaffBookings(params)
    bookingList.value = response.data.bookings || []
  } catch (error) {
    console.error('获取预约列表失败:', error)
  }
}

// 切换标签页
const changeTab = (tab) => {
  activeTab.value = tab
  fetchBookings()
}

// 日期确认
const onDateConfirm = () => {
  currentDate.value = dayjs(selectedDate.value).format('MM月DD日')
  showDatePicker.value = false
  fetchBookings()
}

// 完成预约
const completeBooking = (booking) => {
  showToast(`完成了${booking.customer_name}的预约`)
}

// 取消预约
const cancelBooking = (booking) => {
  showToast(`取消了${booking.customer_name}的预约`)
}

// 打开扫码器
const openScanner = () => {
  showToast('打开扫码核销功能')
}

// 下拉刷新
const onRefresh = async () => {
  refreshing.value = true
  await fetchBookings()
  refreshing.value = false
}

// 上拉加载
const onLoad = async () => {
  loading.value = true
  await fetchBookings()
  loading.value = false
  finished.value = true
}

// 处理预约点击
const handleBookingClick = (booking) => {
  showToast(`点击了预约：${booking.customer_name}`)
}

onMounted(() => {
  fetchBookings()
})
</script>

<style lang="scss" scoped>
.bookings-page {
  background: linear-gradient(180deg, #f8f9ff 0%, #f0f2ff 100%);
  min-height: 100vh;
  padding-bottom: 80px;
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

.date-selector {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.3);
  }
}

// 统计区域
.stats-section {
  padding: 20px 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.stat-item {
  background: white;
  border-radius: 12px;
  padding: 16px 8px;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transition: transform 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
  }
}

.stat-number {
  font-size: 20px;
  font-weight: 700;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

// 筛选区域
.filter-section {
  padding: 0 16px 16px;
}

.filter-tabs {
  display: flex;
  background: white;
  border-radius: 12px;
  padding: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.filter-tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  
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
  border: 1px solid rgba(255, 255, 255, 0.8);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  }
}

.booking-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.customer-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.customer-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
}

.customer-details {
  flex: 1;
}

.customer-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 4px 0;
}

.customer-phone {
  font-size: 13px;
  color: #999;
  margin: 0;
}

.booking-status {
  padding: 6px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  
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

.booking-content {
  border-top: 1px solid #f0f0f0;
  padding-top: 16px;
}

.booking-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #666;
}

.info-icon {
  margin-right: 8px;
  font-size: 14px;
}

.booking-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  flex: 1;
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
  
  &.complete-btn {
    background: linear-gradient(135deg, #52c41a, #73d13d);
    color: white;
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(82, 196, 26, 0.3);
    }
  }
  
  &.cancel-btn {
    background: #f5f5f5;
    color: #999;
    
    &:hover {
      background: #ff4d4f;
      color: white;
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
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 8px;
  color: #666;
}

.empty-desc {
  font-size: 14px;
}

// 悬浮按钮
.floating-actions {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 100;
}

.fab-button {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #ff6b6b, #ffa726);
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  box-shadow: 0 8px 24px rgba(255, 107, 107, 0.4);
  transition: all 0.3s ease;
  
  &:hover {
    transform: scale(1.1);
    box-shadow: 0 12px 32px rgba(255, 107, 107, 0.5);
  }
}

.fab-icon {
  font-size: 20px;
  margin-bottom: 2px;
}

.fab-text {
  font-size: 8px;
  font-weight: 500;
}

// 弹窗样式
.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.popup-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.popup-close {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: #e9ecef;
  }
}

// 响应式设计
@media (max-width: 375px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }
  
  .stat-item {
    padding: 12px 6px;
  }
  
  .stat-number {
    font-size: 18px;
  }
  
  .booking-card {
    padding: 16px;
  }
  
  .customer-avatar {
    width: 36px;
    height: 36px;
    font-size: 14px;
  }
  
  .fab-button {
    width: 56px;
    height: 56px;
    bottom: 16px;
    right: 16px;
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
