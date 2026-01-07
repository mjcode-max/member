<template>
  <div class="dashboard-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="store-info">
          <h1 class="store-name">{{ storeInfo.name || '加载中...' }}</h1>
          <p class="store-address" v-if="storeInfo.address">{{ storeInfo.address }}</p>
          <div class="manager-info" v-if="userInfo.username">
            <i class="manager-icon">👨‍💼</i>
            <span>店长：{{ userInfo.username }}</span>
          </div>
        </div>
        <div class="header-actions">
          <div class="status-indicator" :class="getStatusClass(storeInfo.status)">
            {{ getStatusText(storeInfo.status) }}
          </div>
        </div>
      </div>
    </div>

    <!-- 今日数据概览 -->
    <div class="overview-section">
      <div class="overview-grid">
        <div class="overview-item bookings-item">
          <div class="overview-icon">📅</div>
          <div class="overview-content">
            <div class="overview-number">{{ todayData.bookings }}</div>
            <div class="overview-label">今日预约</div>
          </div>
          <div class="overview-trend positive">
            <i class="trend-icon">📈</i>
            <span>+12%</span>
          </div>
        </div>
        
        <div class="overview-item members-item">
          <div class="overview-icon">💎</div>
          <div class="overview-content">
            <div class="overview-number">{{ todayData.newMembers }}</div>
            <div class="overview-label">新增会员</div>
          </div>
          <div class="overview-trend positive">
            <i class="trend-icon">📈</i>
            <span>+5</span>
          </div>
        </div>
        
        <div class="overview-item revenue-item">
          <div class="overview-icon">💰</div>
          <div class="overview-content">
            <div class="overview-number">¥{{ todayData.revenue }}</div>
            <div class="overview-label">今日收入</div>
          </div>
          <div class="overview-trend positive">
            <i class="trend-icon">📈</i>
            <span>+8%</span>
          </div>
        </div>
        
        <div class="overview-item rate-item">
          <div class="overview-icon">📊</div>
          <div class="overview-content">
            <div class="overview-number">{{ todayData.arrivalRate }}%</div>
            <div class="overview-label">到店率</div>
          </div>
          <div class="overview-trend neutral">
            <i class="trend-icon">➡️</i>
            <span>持平</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 员工状态 -->
    <div class="staff-section">
      <div class="section-header">
        <!-- <i class="section-icon">👥</i> -->
        <span class="section-title">👥 员工状态</span>
        <span class="section-more" @click="$router.push('/staff')">管理 ›</span>
      </div>
      <div class="staff-status-grid">
        <div class="status-summary">
          <div class="summary-item active-summary">
            <div class="summary-number">{{ staffStats.active }}</div>
            <div class="summary-label">在岗</div>
          </div>
          <div class="summary-item rest-summary">
            <div class="summary-number">{{ staffStats.rest }}</div>
            <div class="summary-label">休息</div>
          </div>
          <div class="summary-item total-summary">
            <div class="summary-number">{{ staffStats.total }}</div>
            <div class="summary-label">总计</div>
          </div>
        </div>
        <div class="staff-list-preview">
          <div 
            v-for="staff in staffList.slice(0, 3)"
            :key="staff.id"
            class="staff-preview-item"
          >
            <div class="staff-avatar" :class="getWorkStatusClass(staff.work_status)">
              {{ (staff.username || staff.name || '').charAt(0) }}
            </div>
            <div class="staff-info">
              <div class="staff-name">{{ staff.username || staff.name || '未知' }}</div>
              <div class="staff-status" :class="getWorkStatusClass(staff.work_status)">
                {{ getWorkStatusText(staff.work_status) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 核心功能 -->
    <div class="functions-section">
      <div class="section-header">
        <!-- <i class="section-icon">⚡</i> -->
        <span class="section-title">⚡ 管理功能</span>
      </div>
      <div class="functions-grid">
        <div class="function-card" @click="$router.push('/staff')">
          <div class="function-icon staff-icon">
            <i>👥</i>
          </div>
          <div class="function-content">
            <h3 class="function-title">员工管理</h3>
            <p class="function-desc">管理美甲师工作状态</p>
          </div>
          <div class="function-arrow">›</div>
        </div>
        
        <div class="function-card" @click="$router.push('/members')">
          <div class="function-icon member-icon">
            <i>💎</i>
          </div>
          <div class="function-content">
            <h3 class="function-title">会员管理</h3>
            <p class="function-desc">创建和管理会员信息</p>
          </div>
          <div class="function-arrow">›</div>
        </div>
      </div>
    </div>

    <!-- 最新会员 -->
    <div class="recent-members-section">
      <div class="section-header">
        <!-- <i class="section-icon">👑</i> -->
        <span class="section-title">👑 最新会员</span>
        <span class="section-more" @click="$router.push('/members')">查看全部 ›</span>
      </div>
      <div class="members-preview">
        <div 
          v-for="member in recentMembers"
          :key="member.id"
          class="member-preview-item"
          @click="$router.push(`/members/${member.id}`)"
        >
          <div class="member-avatar">
            {{ member.name.charAt(0) }}
          </div>
          <div class="member-info">
            <div class="member-name">{{ member.name }}</div>
            <div class="member-package">{{ member.package_name }}</div>
            <div class="member-date">{{ formatDate(member.created_at) }}</div>
          </div>
          <div class="member-status" :class="member.status">
            {{ member.status === 'active' ? '有效' : '过期' }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getStoreDashboard } from '@/api/dashboard'
import { getStoreById } from '@/api/stores'
import { getStaffList } from '@/api/staff'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

// 从store获取用户信息
const userInfo = computed(() => userStore.userInfo || {})

// 店铺信息
const storeInfo = reactive({
  id: null,
  name: '',
  address: '',
  phone: '',
  contact_person: '',
  status: '',
  business_hours_start: '',
  business_hours_end: '',
  deposit_amount: 0
})

const todayData = reactive({
  bookings: 25,
  newMembers: 3,
  revenue: '3,280',
  arrivalRate: 92
})

const staffStats = reactive({
  active: 0,
  rest: 0,
  total: 0
})

const staffList = ref([])

const recentMembers = ref([
  { 
    id: 1, 
    name: '刘小姐', 
    package_name: '高级会员', 
    status: 'active',
    created_at: '2023-10-15'
  },
  { 
    id: 2, 
    name: '王女士', 
    package_name: 'VIP会员', 
    status: 'active',
    created_at: '2023-10-14'
  },
  { 
    id: 3, 
    name: '张先生', 
    package_name: '基础会员', 
    status: 'active',
    created_at: '2023-10-13'
  }
])

// 获取门店状态文本
const getStatusText = (status) => {
  const statusMap = {
    operating: '营业中',
    closed: '停业',
    shutdown: '已关闭'
  }
  return statusMap[status] || '未知'
}

// 获取门店状态样式类
const getStatusClass = (status) => {
  return {
    'status-operating': status === 'operating',
    'status-closed': status === 'closed',
    'status-shutdown': status === 'shutdown'
  }
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('MM-DD')
}

// 获取店铺详情（仅用于首页头部显示）
const fetchStoreInfo = async () => {
  // 从用户信息中获取store_id
  const storeId = userInfo.value.store_id
  
  if (!storeId) {
    console.error('用户未关联门店')
    return
  }
  
  try {
    const response = await getStoreById(storeId)
    if (response.data) {
      // 只更新首页头部需要的基本信息
      storeInfo.name = response.data.name || ''
      storeInfo.address = response.data.address || ''
      storeInfo.status = response.data.status || ''
    }
  } catch (error) {
    console.error('获取店铺详情失败:', error)
  }
}

// 获取仪表板数据
const fetchDashboardData = async () => {
  try {
    const response = await getStoreDashboard()
    if (response.data) {
      Object.assign(todayData, response.data)
    }
  } catch (error) {
    console.error('获取仪表板数据失败:', error)
  }
}

// 获取员工列表
const fetchStaffList = async () => {
  try {
    const response = await getStaffList()
    // 处理分页响应格式
    if (response && response.data) {
      let staffData = []
      if (response.data.list && Array.isArray(response.data.list)) {
        staffData = response.data.list
      } else if (Array.isArray(response.data)) {
        staffData = response.data
      }
      
      // 更新员工列表
      staffList.value = staffData
      
      // 计算统计数据
      const active = staffData.filter(s => {
        const status = s.work_status
        return status === 'working' || status === 'active'
      }).length
      const rest = staffData.filter(s => {
        const status = s.work_status
        return status === 'rest'
      }).length
      
      staffStats.active = active
      staffStats.rest = rest
      staffStats.total = staffData.length
    } else {
      staffList.value = []
      staffStats.active = 0
      staffStats.rest = 0
      staffStats.total = 0
    }
  } catch (error) {
    console.error('获取员工列表失败:', error)
    staffList.value = []
    staffStats.active = 0
    staffStats.rest = 0
    staffStats.total = 0
  }
}

// 获取工作状态文本
const getWorkStatusText = (status) => {
  if (!status) return '未知'
  const statusMap = {
    working: '在岗',
    rest: '休息',
    offline: '离岗',
    active: '在岗' // 兼容旧数据
  }
  return statusMap[status] || '未知'
}

// 获取工作状态样式类
const getWorkStatusClass = (status) => {
  if (!status) return 'offline'
  // 兼容旧数据
  if (status === 'active') return 'active'
  if (status === 'working') return 'active'
  if (status === 'rest') return 'rest'
  return 'offline'
}

onMounted(() => {
  // 先获取店铺信息
  fetchStoreInfo()
  // 再获取仪表板数据
  fetchDashboardData()
  // 获取员工列表
  fetchStaffList()
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
  background: linear-gradient(135deg, #ff6b6b 0%, #ffa726 100%);
  padding: 24px 20px;
  color: white;

}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.store-info {
  flex: 1;
}

.store-name {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 6px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.store-address {
  font-size: 14px;
  opacity: 0.9;
  margin: 0 0 8px 0;
}

.manager-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  opacity: 0.8;
}

.manager-icon {
  font-size: 16px;
}

.status-indicator {
  padding: 8px 16px;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 600;
  
  &.status-operating {
    background: rgba(82, 196, 26, 0.2);
    color: #73d13d;
  }
  
  &.status-closed {
    background: rgba(250, 173, 20, 0.2);
    color: #ffc53d;
  }
  
  &.status-shutdown {
    background: rgba(255, 77, 79, 0.2);
    color: #ff7875;
  }
}

// 数据概览
.overview-section {
  padding: 20px 16px;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.overview-item {
  background: white;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  position: relative;
  overflow: hidden;
  transition: transform 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
  }
  
  &.bookings-item {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
  }
  
  &.members-item {
    background: linear-gradient(135deg, #f093fb, #f5576c);
    color: white;
  }
  
  &.revenue-item {
    background: linear-gradient(135deg, #4facfe, #00f2fe);
    color: white;
  }
  
  &.rate-item {
    background: linear-gradient(135deg, #43e97b, #38f9d7);
    color: white;
  }
}

.overview-icon {
  font-size: 28px;
}

.overview-content {
  flex: 1;
}

.overview-number {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.overview-label {
  font-size: 14px;
  opacity: 0.9;
}

.overview-trend {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 2px;
  
  &.positive {
    background: rgba(82, 196, 26, 0.2);
    color: #73d13d;
  }
  
  &.neutral {
    background: rgba(255, 255, 255, 0.2);
    color: rgba(255, 255, 255, 0.8);
  }
}

.trend-icon {
  font-size: 8px;
}

// 员工状态
.staff-section {
  padding: 0 16px 20px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.center-header {
  justify-content: center;
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
  color: #ff6b6b;
  cursor: pointer;
  
  &:hover {
    color: #ffa726;
  }
}

.staff-status-grid {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.status-summary {
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.summary-item {
  text-align: center;
  
  &.active-summary {
    color: #52c41a;
  }
  
  &.rest-summary {
    color: #faad14;
  }
  
  &.total-summary {
    color: #667eea;
  }
}

.summary-number {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.summary-label {
  font-size: 12px;
  opacity: 0.8;
}

.staff-list-preview {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.staff-preview-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}

.staff-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: white;
  
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

.staff-info {
  flex: 1;
}

.staff-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.staff-status {
  font-size: 12px;
  
  &.active {
    color: #52c41a;
  }
  
  &.rest {
    color: #faad14;
  }
  
  &.offline {
    color: #999;
  }
}

// 功能区域
.functions-section {
  padding: 0 16px 20px;
}

.functions-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.function-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  }
}

.function-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-right: 16px;
  
  &.staff-icon {
    background: linear-gradient(135deg, #ff6b6b, #ffa726);
    color: white;
  }
  
  &.member-icon {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
  }
}

.function-content {
  flex: 1;
}

.function-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.function-desc {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.function-arrow {
  font-size: 18px;
  color: #ccc;
  margin-left: 8px;
}

// 最新会员
.recent-members-section {
  padding: 0 16px 20px;
}

.members-preview {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.member-preview-item {
  display: flex;
  align-items: center;
  gap: 12px;
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

.member-avatar {
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

.member-info {
  flex: 1;
}

.member-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 2px;
}

.member-package {
  font-size: 13px;
  color: #667eea;
  margin-bottom: 2px;
}

.member-date {
  font-size: 12px;
  color: #999;
}

.member-status {
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  
  &.active {
    background: #f6ffed;
    color: #52c41a;
  }
  
  &.expired {
    background: #fff2f0;
    color: #ff4d4f;
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

.overview-section,
.staff-section,
.functions-section,
.recent-members-section {
  animation: slideInUp 0.6s ease-out;
}

.overview-section { animation-delay: 0.1s; }
.staff-section { animation-delay: 0.2s; }
.functions-section { animation-delay: 0.3s; }
.recent-members-section { animation-delay: 0.4s; }

// 响应式设计
@media (max-width: 375px) {
  .overview-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  
  .overview-item {
    padding: 16px;
  }
  
  .overview-number {
    font-size: 20px;
  }
  
  .store-name {
    font-size: 20px;
  }
}
</style>