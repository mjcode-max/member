<template>
  <div class="dashboard-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="user-info">
          <div class="user-avatar" :class="userStatus.work_status">
            {{ (userInfo.username || userInfo.name || 'U').charAt(0) }}
            <div class="status-indicator" :class="userStatus.work_status"></div>
          </div>
          <div class="user-details">
            <h2 class="user-name">{{ userInfo.username || userInfo.name || '用户' }}</h2>
            <p class="user-role">美甲师</p>
            <div class="work-status" :class="userStatus.work_status">
              {{ getWorkStatusText(userStatus.work_status) }}
            </div>
          </div>
        </div>
        <div class="status-toggle" @click="toggleWorkStatus">
          <i class="toggle-icon">{{ userStatus.work_status === 'working' ? '😴' : '💼' }}</i>
        </div>
      </div>
    </div>

    <!-- 今日统计 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-icon">📅</div>
          <div class="stat-content">
            <div class="stat-number">{{ todayStats.total }}</div>
            <div class="stat-label">今日预约</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon">✅</div>
          <div class="stat-content">
            <div class="stat-number">{{ todayStats.completed }}</div>
            <div class="stat-label">已完成</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon">⏰</div>
          <div class="stat-content">
            <div class="stat-number">{{ todayStats.pending }}</div>
            <div class="stat-label">待服务</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon">💰</div>
          <div class="stat-content">
            <div class="stat-number">¥{{ todayStats.earnings }}</div>
            <div class="stat-label">今日收入</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 核心功能 -->
    <div class="functions-section">
      <div class="section-header">
        <i class="section-icon">⚡</i>
        <span class="section-title">核心功能</span>
      </div>
      <div class="functions-grid">
        <div class="function-item primary-function" @click="$router.push('/bookings')">
          <div class="function-icon">
            <i>📋</i>
          </div>
          <div class="function-content">
            <h3 class="function-title">我的日程</h3>
            <p class="function-desc">查看今日预约安排</p>
            <div class="function-badge" v-if="todayStats.pending > 0">
              {{ todayStats.pending }}
            </div>
          </div>
          <div class="function-arrow">›</div>
        </div>
        
        <div class="function-item primary-function" @click="$router.push('/scanner')">
          <div class="function-icon scanner-icon">
            <i>📱</i>
          </div>
          <div class="function-content">
            <h3 class="function-title">扫码核销</h3>
            <p class="function-desc">扫描会员码进行核销</p>
          </div>
          <div class="function-arrow">›</div>
        </div>
      </div>
    </div>

    <!-- 今日预约预览 -->
    <div class="preview-section">
      <div class="section-header">
        <i class="section-icon">📅</i>
        <span class="section-title">今日预约</span>
        <span class="section-more" @click="$router.push('/bookings')">查看全部 ›</span>
      </div>
      <div class="bookings-preview">
        <div 
          v-for="booking in todayBookings.slice(0, 3)"
          :key="booking.id"
          class="booking-preview-item"
          @click="handleBookingClick(booking)"
        >
          <div class="booking-time">{{ booking.time_slot }}</div>
          <div class="booking-info">
            <div class="customer-name">{{ booking.customer_name }}</div>
            <div class="service-type">{{ getServiceText(booking.service_type) }}</div>
          </div>
          <div class="booking-status" :class="booking.booking_status">
            {{ getStatusText(booking.booking_status) }}
          </div>
        </div>
        
        <div v-if="todayBookings.length === 0" class="no-bookings">
          <div class="no-bookings-icon">😊</div>
          <div class="no-bookings-text">今日暂无预约</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { useUserStore } from '@/stores/user'
import { updateMyWorkStatus, getTodayBookings } from '@/api/staff'
import { getCurrentUser } from '@/api/auth'

const router = useRouter()
const userStore = useUserStore()

const userInfo = reactive({
  username: '',
  name: '',
  phone: ''
})

const userStatus = reactive({
  work_status: 'working' // working, rest, offline
})

const todayStats = reactive({
  total: 8,
  completed: 3,
  pending: 5,
  earnings: '280'
})

const todayBookings = ref([])

// 获取工作状态文本
const getWorkStatusText = (status) => {
  const statusMap = {
    working: '在岗中',
    rest: '休息中',
    offline: '离岗'
  }
  return statusMap[status] || '未知'
}

// 切换工作状态
const toggleWorkStatus = async () => {
  const newStatus = userStatus.work_status === 'working' ? 'rest' : 'working'
  const statusText = newStatus === 'working' ? '在岗' : '休息'
  
  try {
    await showConfirmDialog({
      title: '确认状态切换',
      message: `确定要切换为${statusText}状态吗？`
    })

    // 获取当前用户ID
    const userId = userStore.userInfo?.id
    if (!userId) {
      showToast('无法获取用户ID')
      return
    }
    await updateMyWorkStatus({ work_status: newStatus }, userId)
    userStatus.work_status = newStatus
    // 更新store中的用户信息
    if (userStore.userInfo) {
      userStore.userInfo.work_status = newStatus
      localStorage.setItem('userInfo', JSON.stringify(userStore.userInfo))
    }
    showToast(`已切换为${statusText}状态`)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('更新工作状态失败:', error)
      showToast('状态切换失败')
    }
  }
}

// 处理预约点击
const handleBookingClick = (booking) => {
  router.push(`/bookings/${booking.id}`)
}

// 获取服务类型文本
const getServiceText = (type) => {
  return type === 'manicure' ? '美甲' : '美睫'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    pending: '待确认',
    confirmed: '待服务',
    completed: '已完成'
  }
  return statusMap[status] || '未知'
}

// 获取今日预约
const fetchTodayBookings = async () => {
  try {
    const response = await getTodayBookings()
    todayBookings.value = response.data || []
  } catch (error) {
    console.error('获取今日预约失败:', error)
  }
}

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const user = userStore.userInfo
    if (user && Object.keys(user).length > 0) {
      Object.assign(userInfo, {
        username: user.username || user.name || '',
        name: user.name || user.username || '',
        phone: user.phone || ''
      })
      userStatus.work_status = user.work_status || 'working'
    } else {
      // 如果本地没有用户信息，尝试从接口获取
      const response = await getCurrentUser()
      if (response && response.data) {
        const userData = response.data.user || response.data
        Object.assign(userInfo, {
          username: userData.username || userData.name || '',
          name: userData.name || userData.username || '',
          phone: userData.phone || ''
        })
        userStatus.work_status = userData.work_status || 'working'
        // 更新store
        userStore.userInfo = userData
      }
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

onMounted(() => {
  fetchUserInfo()
  fetchTodayBookings()
})
</script>

<style lang="scss" scoped>
.dashboard-page {
  background: linear-gradient(180deg, #f8f9ff 0%, #f0f2ff 100%);
  min-height: 100vh;
  padding-bottom: 20px;
}

// 页面头部
.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30px 20px;
  color: white;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: white;
  position: relative;
  border: 3px solid rgba(255, 255, 255, 0.3);
  
  &.working {
    background: linear-gradient(135deg, #52c41a, #73d13d);
    border-color: rgba(82, 196, 26, 0.5);
  }
  
  &.rest {
    background: linear-gradient(135deg, #faad14, #ffc53d);
    border-color: rgba(250, 173, 20, 0.5);
  }
  
  &.offline {
    background: linear-gradient(135deg, #999, #bbb);
    border-color: rgba(153, 153, 153, 0.5);
  }
}

.status-indicator {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 3px solid white;
  
  &.working {
    background: #52c41a;
    animation: pulse-green 2s infinite;
  }
  
  &.rest {
    background: #faad14;
    animation: pulse-orange 2s infinite;
  }
  
  &.offline {
    background: #999;
  }
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.user-role {
  font-size: 16px;
  opacity: 0.9;
  margin: 0 0 8px 0;
}

.work-status {
  padding: 6px 12px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  display: inline-block;
  
  &.working {
    background: rgba(82, 196, 26, 0.2);
    color: #73d13d;
  }
  
  &.rest {
    background: rgba(250, 173, 20, 0.2);
    color: #ffc53d;
  }
  
  &.offline {
    background: rgba(153, 153, 153, 0.2);
    color: #999;
  }
}

.status-toggle {
  width: 50px;
  height: 50px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  
  &:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.1);
  }
}

.toggle-icon {
  font-size: 24px;
}

// 统计区域
.stats-section {
  padding: 20px 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-item {
  background: white;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: transform 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
  }
}

.stat-icon {
  font-size: 24px;
}

.stat-content {
  flex: 1;
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

// 功能区域
.functions-section {
  padding: 0 16px 20px;
}

.section-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.section-icon {
  font-size: 18px;
  margin-right: 8px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.section-more {
  font-size: 14px;
  color: #667eea;
  cursor: pointer;
  margin-left: auto;
}

.functions-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.function-item {
  background: white;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  position: relative;
  overflow: hidden;
  
  &.primary-function {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 40px rgba(102, 126, 234, 0.3);
    }
  }
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  }
}

.function-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  margin-right: 16px;
  
  .primary-function & {
    background: rgba(255, 255, 255, 0.2);
  }
  
  .function-item:not(.primary-function) & {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
  }
  
  &.scanner-icon {
    animation: pulse 2s infinite;
  }
}

.function-content {
  flex: 1;
}

.function-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 4px;
}

.function-desc {
  font-size: 14px;
  opacity: 0.8;
  margin: 0;
}

.function-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  background: #ff6b6b;
  color: white;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.function-arrow {
  font-size: 20px;
  opacity: 0.6;
  margin-left: 8px;
}

// 预约预览
.preview-section {
  padding: 0 16px 20px;
}

.bookings-preview {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.booking-preview-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f8f8f8;
  cursor: pointer;
  transition: background 0.2s ease;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    background: #f8f9ff;
  }
}

.booking-time {
  font-size: 16px;
  font-weight: 600;
  color: #667eea;
  margin-right: 16px;
  min-width: 80px;
}

.booking-info {
  flex: 1;
}

.customer-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 2px;
}

.service-type {
  font-size: 13px;
  color: #666;
}

.booking-status {
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  
  &.confirmed {
    background: #e6f7ff;
    color: #1890ff;
  }
  
  &.pending {
    background: #fff7e6;
    color: #fa8c16;
  }
}

// 无预约状态
.no-bookings {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.no-bookings-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.no-bookings-text {
  font-size: 16px;
}

// 动画效果
@keyframes pulse-green {
  0%, 100% { 
    box-shadow: 0 0 0 0 rgba(82, 196, 26, 0.7);
  }
  50% { 
    box-shadow: 0 0 0 8px rgba(82, 196, 26, 0);
  }
}

@keyframes pulse-orange {
  0%, 100% { 
    box-shadow: 0 0 0 0 rgba(250, 173, 20, 0.7);
  }
  50% { 
    box-shadow: 0 0 0 8px rgba(250, 173, 20, 0);
  }
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

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

.stats-section,
.functions-section,
.preview-section {
  animation: slideInUp 0.6s ease-out;
}

.stats-section { animation-delay: 0.1s; }
.functions-section { animation-delay: 0.2s; }
.preview-section { animation-delay: 0.3s; }

// 响应式设计
@media (max-width: 375px) {
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  
  .user-avatar {
    width: 50px;
    height: 50px;
    font-size: 20px;
  }
  
  .user-name {
    font-size: 20px;
  }
  
  .function-icon {
    width: 48px;
    height: 48px;
    font-size: 24px;
  }
}
</style>
