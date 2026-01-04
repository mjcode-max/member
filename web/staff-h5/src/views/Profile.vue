<template>
  <div class="profile-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="profile-info">
          <div class="avatar-container">
            <div class="user-avatar" :class="userInfo.work_status">
              {{ (userInfo.username || userInfo.name || 'U').charAt(0) }}
              <div class="status-ring" :class="userInfo.work_status"></div>
            </div>
            <div class="work-status-badge" :class="userInfo.work_status">
              {{ getWorkStatusText(userInfo.work_status) }}
            </div>
          </div>
          <div class="user-details">
            <h2 class="user-name">{{ userInfo.username || userInfo.name || '用户' }}</h2>
            <p class="user-role">专业美甲师</p>
            <div class="user-meta">
              <span class="join-date" v-if="userInfo.created_at">入职时间：{{ formatDate(userInfo.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 工作状态切换 -->
    <div class="status-section">
      <div class="status-card">
        <div class="status-header">
          <i class="status-icon">💼</i>
          <span class="status-title">工作状态</span>
        </div>
        <div class="status-toggle">
          <div 
            class="toggle-option"
            :class="{ active: userInfo.work_status === 'working' || userInfo.work_status === 'active' }"
            @click="updateWorkStatus('working')"
          >
            <div class="option-icon">🟢</div>
            <div class="option-content">
              <div class="option-title">在岗工作</div>
              <div class="option-desc">接受新的预约服务</div>
            </div>
          </div>
          <div 
            class="toggle-option"
            :class="{ active: userInfo.work_status === 'rest' }"
            @click="updateWorkStatus('rest')"
          >
            <div class="option-icon">😴</div>
            <div class="option-content">
              <div class="option-title">休息状态</div>
              <div class="option-desc">暂停接受新预约</div>
            </div>
          </div>
        </div>
        <div class="status-note">
          <i class="note-icon">💡</i>
          <span>切换为休息状态将影响门店的可预约数量</span>
        </div>
      </div>
    </div>

    <!-- 今日统计 -->
    <div class="stats-section">
      <div class="section-header">
        <i class="section-icon">📊</i>
        <span class="section-title">今日数据</span>
      </div>
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-icon">📅</div>
          <div class="stat-content">
            <div class="stat-number">{{ todayStats.bookings }}</div>
            <div class="stat-label">预约数</div>
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
          <div class="stat-icon">💰</div>
          <div class="stat-content">
            <div class="stat-number">¥{{ todayStats.earnings }}</div>
            <div class="stat-label">收入</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon">⭐</div>
          <div class="stat-content">
            <div class="stat-number">{{ todayStats.rating }}</div>
            <div class="stat-label">评分</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 快捷功能 -->
    <div class="actions-section">
      <div class="section-header">
        <i class="section-icon">⚡</i>
        <span class="section-title">快捷功能</span>
      </div>
      <div class="actions-grid">
        <div class="action-item" @click="$router.push('/bookings')">
          <div class="action-icon action-icon-1">
            <i>📋</i>
          </div>
          <div class="action-text">我的日程</div>
          <div class="action-desc">查看今日预约</div>
        </div>
        <div class="action-item" @click="openScanner">
          <div class="action-icon action-icon-2">
            <i>📱</i>
          </div>
          <div class="action-text">扫码核销</div>
          <div class="action-desc">扫描会员码</div>
        </div>
        <div class="action-item" @click="$router.push('/schedule')">
          <div class="action-icon action-icon-3">
            <i>📅</i>
          </div>
          <div class="action-text">排班表</div>
          <div class="action-desc">查看工作安排</div>
        </div>
        <div class="action-item" @click="$router.push('/settings')">
          <div class="action-icon action-icon-4">
            <i>⚙️</i>
          </div>
          <div class="action-text">设置</div>
          <div class="action-desc">个人设置</div>
        </div>
      </div>
    </div>

    <!-- 最近服务记录 -->
    <div class="recent-section">
      <div class="section-header">
        <i class="section-icon">🕒</i>
        <span class="section-title">最近服务</span>
        <span class="section-more" @click="$router.push('/bookings')">查看全部 ›</span>
      </div>
      <div class="recent-list">
        <div 
          v-for="record in recentServices"
          :key="record.id"
          class="recent-item"
        >
          <div class="recent-avatar">
            {{ record.customer_name.charAt(0) }}
          </div>
          <div class="recent-info">
            <div class="recent-customer">{{ record.customer_name }}</div>
            <div class="recent-service">{{ getServiceText(record.service_type) }}</div>
            <div class="recent-time">{{ formatDateTime(record.completed_time) }}</div>
          </div>
          <div class="recent-status">
            <i class="status-check">✅</i>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onActivated } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { useUserStore } from '@/stores/user'
import { updateMyWorkStatus } from '@/api/staff'
import { getCurrentUser } from '@/api/auth'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const userInfo = reactive({
  username: '',
  name: '',
  work_status: 'working', // working, rest, offline
  created_at: ''
})

const todayStats = reactive({
  bookings: 8,
  completed: 5,
  earnings: '280',
  rating: '4.9'
})

const recentServices = ref([
  {
    id: 1,
    customer_name: '张小姐',
    service_type: 'manicure',
    completed_time: '2023-10-15 14:30:00'
  },
  {
    id: 2,
    customer_name: '王女士',
    service_type: 'eyelash',
    completed_time: '2023-10-15 13:00:00'
  },
  {
    id: 3,
    customer_name: '刘小姐',
    service_type: 'manicure',
    completed_time: '2023-10-15 11:30:00'
  }
])

// 获取工作状态文本
const getWorkStatusText = (status) => {
  const statusMap = {
    working: '在岗中',
    rest: '休息中',
    offline: '离岗',
    active: '在岗中' // 兼容旧数据
  }
  return statusMap[status] || '未知'
}

// 更新工作状态
const updateWorkStatus = async (status) => {
  if (userInfo.work_status === status) return

  try {
    const statusText = getWorkStatusText(status)
    const message = status === 'working' ? '确定要设置为在岗状态吗？' : '确定要设置为休息状态吗？'
    const note = status === 'rest' ? '休息状态下将不会接收新的预约' : '在岗状态下可以接收新的预约'
    
    await showConfirmDialog({
      title: '确认状态切换',
      message: `${message}\n\n${note}`
    })

    // 获取当前用户ID
    const userId = userStore.userInfo?.id
    if (!userId) {
      showToast('无法获取用户ID')
      return
    }
    await updateMyWorkStatus({ work_status: status }, userId)
    userInfo.work_status = status
    // 更新store中的用户信息
    if (userStore.userInfo) {
      userStore.userInfo.work_status = status
      localStorage.setItem('userInfo', JSON.stringify(userStore.userInfo))
    }
    // 重新获取用户信息以确保数据同步
    await fetchUserInfo()
    
    showToast(`已切换为${statusText}状态`)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('更新工作状态失败:', error)
      showToast('状态更新失败，请重试')
    }
  }
}

// 打开扫码器
const openScanner = () => {
  router.push('/scanner')
}

// 获取服务类型文本
const getServiceText = (type) => {
  return type === 'manicure' ? '美甲服务' : '美睫服务'
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY年MM月DD日')
}

// 格式化日期时间
const formatDateTime = (datetime) => {
  return dayjs(datetime).format('MM-DD HH:mm')
}

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const response = await getCurrentUser()
    if (response && response.data) {
      const userData = response.data.user || response.data
      Object.assign(userInfo, {
        username: userData.username || userData.name || '',
        name: userData.name || userData.username || '',
        work_status: userData.work_status || 'working',
        created_at: userData.created_at || ''
      })
      // 更新store
      if (userStore.userInfo) {
        Object.assign(userStore.userInfo, userData)
      } else {
        userStore.userInfo = userData
      }
      localStorage.setItem('userInfo', JSON.stringify(userStore.userInfo))
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    // 如果接口失败，尝试使用本地存储的数据
    const user = userStore.userInfo
    if (user && Object.keys(user).length > 0) {
      Object.assign(userInfo, {
        username: user.username || user.name || '',
        name: user.name || user.username || '',
        work_status: user.work_status || 'working',
        created_at: user.created_at || ''
      })
    }
  }
}

onMounted(async () => {
  // 先尝试从本地获取
  const user = userStore.userInfo
  if (user && Object.keys(user).length > 0) {
    Object.assign(userInfo, {
      username: user.username || user.name || '',
      name: user.name || user.username || '',
      work_status: user.work_status || 'working',
      created_at: user.created_at || ''
    })
  }
  // 然后从接口获取最新数据
  await fetchUserInfo()
})

// 页面激活时刷新用户信息（从其他页面返回时）
onActivated(() => {
  fetchUserInfo()
})
</script>

<style lang="scss" scoped>
.profile-page {
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

.profile-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar-container {
  position: relative;
  text-align: center;
}

.user-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 700;
  color: white;
  position: relative;
  border: 4px solid rgba(255, 255, 255, 0.3);
  
  &.working,
  &.active {
    background: linear-gradient(135deg, #52c41a, #73d13d);
  }
  
  &.rest {
    background: linear-gradient(135deg, #faad14, #ffc53d);
  }
  
  &.offline {
    background: linear-gradient(135deg, #999, #bbb);
  }
}

.status-ring {
  position: absolute;
  top: -4px;
  left: -4px;
  right: -4px;
  bottom: -4px;
  border-radius: 50%;
  border: 3px solid transparent;
  
  &.working,
  &.active {
    border-color: #52c41a;
    animation: pulse-green 2s infinite;
  }
  
  &.rest {
    border-color: #faad14;
    animation: pulse-orange 2s infinite;
  }
  
  &.offline {
    border-color: #999;
  }
}

.work-status-badge {
  margin-top: 8px;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  
  &.working,
  &.active {
    background: rgba(82, 196, 26, 0.2);
    color: #52c41a;
  }
  
  &.rest {
    background: rgba(250, 140, 22, 0.2);
    color: #faad14;
  }
  
  &.offline {
    background: rgba(153, 153, 153, 0.2);
    color: #999;
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

.user-meta {
  font-size: 14px;
  opacity: 0.8;
}

// 工作状态区域
.status-section {
  padding: 20px 16px;
}

.status-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.status-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.status-icon {
  font-size: 20px;
  margin-right: 8px;
}

.status-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.status-toggle {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.toggle-option {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  
  &.active,
  &.working {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    border-color: #667eea;
    transform: scale(1.02);
  }
  
  &:hover:not(.active) {
    background: #e9ecef;
  }
}

.option-icon {
  font-size: 24px;
}

.option-content {
  flex: 1;
}

.option-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
}

.option-desc {
  font-size: 14px;
  opacity: 0.8;
}

.status-note {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #fff7e6;
  border-radius: 8px;
  font-size: 14px;
  color: #fa8c16;
}

.note-icon {
  font-size: 16px;
}

// 通用区块样式
.stats-section,
.actions-section,
.recent-section {
  padding: 0 16px 20px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
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
}

// 统计网格
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-item {
  background: white;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
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

// 快捷功能
.actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.action-item {
  background: white;
  border-radius: 16px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  }
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
  font-size: 20px;
  
  &.action-icon-1 {
    background: linear-gradient(135deg, #667eea, #764ba2);
  }
  
  &.action-icon-2 {
    background: linear-gradient(135deg, #ff6b6b, #ffa726);
  }
  
  &.action-icon-3 {
    background: linear-gradient(135deg, #4facfe, #00f2fe);
  }
  
  &.action-icon-4 {
    background: linear-gradient(135deg, #43e97b, #38f9d7);
  }
}

.action-text {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.action-desc {
  font-size: 12px;
  color: #999;
}

// 最近服务
.recent-list {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #f8f8f8;
  transition: background 0.2s ease;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    background: #f8f9ff;
  }
}

.recent-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

.recent-info {
  flex: 1;
}

.recent-customer {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.recent-service {
  font-size: 13px;
  color: #666;
  margin-bottom: 2px;
}

.recent-time {
  font-size: 12px;
  color: #999;
}

.recent-status {
  display: flex;
  align-items: center;
}

.status-check {
  font-size: 16px;
  color: #52c41a;
}

// 动画效果
@keyframes pulse-green {
  0%, 100% { 
    border-color: #52c41a;
    opacity: 1;
  }
  50% { 
    border-color: rgba(82, 196, 26, 0.5);
    opacity: 0.7;
  }
}

@keyframes pulse-orange {
  0%, 100% { 
    border-color: #faad14;
    opacity: 1;
  }
  50% { 
    border-color: rgba(250, 173, 20, 0.5);
    opacity: 0.7;
  }
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

.status-section,
.stats-section,
.actions-section,
.recent-section {
  animation: slideInUp 0.6s ease-out;
}

.status-section { animation-delay: 0.1s; }
.stats-section { animation-delay: 0.2s; }
.actions-section { animation-delay: 0.3s; }
.recent-section { animation-delay: 0.4s; }

// 响应式设计
@media (max-width: 375px) {
  .profile-info {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
  
  .stats-grid,
  .actions-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  
  .user-avatar {
    width: 60px;
    height: 60px;
    font-size: 24px;
  }
  
  .user-name {
    font-size: 20px;
  }
}
</style>