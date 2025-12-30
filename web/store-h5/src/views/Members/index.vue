<template>
  <div class="members-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1 class="page-title">会员管理</h1>
          <p class="page-subtitle">管理门店会员信息</p>
        </div>
        <div class="header-actions">
          <div class="add-button" @click="$router.push('/members/create')">
            <i class="add-icon">➕</i>
          </div>
        </div>
      </div>
    </div>

    <!-- 会员统计 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-number">{{ totalMembers }}</div>
          <div class="stat-label">总会员</div>
        </div>
        <div class="stat-item active-stat">
          <div class="stat-number">{{ activeMembers }}</div>
          <div class="stat-label">有效会员</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ newMembersToday }}</div>
          <div class="stat-label">今日新增</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ memberRevenue }}</div>
          <div class="stat-label">会员收入</div>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-section">
      <div class="search-container">
        <div class="search-input-container">
          <i class="search-icon">🔍</i>
          <input 
            v-model="searchKeyword"
            type="text"
            placeholder="搜索会员姓名或手机号"
            class="search-input"
            @input="handleSearch"
          />
        </div>
        <div class="filter-button" @click="showFilterDialog = true">
          <i class="filter-icon">🔽</i>
        </div>
      </div>
    </div>

    <!-- 会员列表 -->
    <div class="members-container">
      <div class="members-list">
        <div 
          v-for="member in memberList"
          :key="member.id"
          class="member-card"
          @click="handleMemberClick(member)"
        >
          <div class="member-header">
            <div class="member-info">
              <div class="member-avatar" :class="member.status">
                {{ member.name.charAt(0) }}
                <div class="member-level">{{ getMemberLevel(member.package_name) }}</div>
              </div>
              <div class="member-details">
                <h4 class="member-name">{{ member.name }}</h4>
                <p class="member-phone">{{ formatPhone(member.phone) }}</p>
                <div class="member-meta">
                  <span class="member-package">{{ member.package_name }}</span>
                  <span class="member-times">剩余{{ member.remaining_times }}次</span>
                </div>
              </div>
            </div>
            <div class="member-status" :class="member.status">
              {{ getMemberStatusText(member.status) }}
            </div>
          </div>

          <div class="member-content">
            <div class="member-progress">
              <div class="progress-info">
                <span class="progress-label">套餐进度</span>
                <span class="progress-text">
                  {{ member.total_times - member.remaining_times }}/{{ member.total_times }}
                </span>
              </div>
              <div class="progress-bar">
                <div 
                  class="progress-fill"
                  :style="{ width: getProgressPercent(member) + '%' }"
                ></div>
              </div>
            </div>

            <div class="member-validity">
              <div class="validity-item">
                <i class="validity-icon">📅</i>
                <span class="validity-text">
                  有效期至 {{ formatDate(member.valid_to) }}
                </span>
                <span class="validity-days" :class="getValidityClass(member.valid_to)">
                  {{ getValidityDays(member.valid_to) }}
                </span>
              </div>
            </div>

            <div class="member-actions">
              <button class="action-btn face-btn" @click.stop="uploadFace(member)">
                <i class="btn-icon">📸</i>
                录入人脸
              </button>
              <button class="action-btn edit-btn" @click.stop="editMember(member)">
                <i class="btn-icon">✏️</i>
                编辑
              </button>
              <button class="action-btn history-btn" @click.stop="viewHistory(member)">
                <i class="btn-icon">📋</i>
                消费记录
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 筛选弹窗 -->
    <van-popup 
      v-model:show="showFilterDialog" 
      position="bottom" 
      round
      :style="{ height: '60%' }"
    >
      <div class="filter-dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">筛选条件</h3>
          <div class="dialog-close" @click="showFilterDialog = false">×</div>
        </div>
        <div class="filter-content">
          <div class="filter-group">
            <div class="filter-label">会员状态</div>
            <div class="filter-options">
              <div 
                v-for="status in statusOptions"
                :key="status.value"
                class="filter-option"
                :class="{ active: filterStatus === status.value }"
                @click="filterStatus = status.value"
              >
                {{ status.label }}
              </div>
            </div>
          </div>
          <div class="filter-group">
            <div class="filter-label">套餐类型</div>
            <div class="filter-options">
              <div 
                v-for="pkg in packageOptions"
                :key="pkg.value"
                class="filter-option"
                :class="{ active: filterPackage === pkg.value }"
                @click="filterPackage = pkg.value"
              >
                {{ pkg.label }}
              </div>
            </div>
          </div>
        </div>
        <div class="filter-actions">
          <button class="filter-btn reset-btn" @click="resetFilter">
            重置
          </button>
          <button class="filter-btn apply-btn" @click="applyFilter">
            应用筛选
          </button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { getMembers } from '@/api/members'
import dayjs from 'dayjs'

const router = useRouter()

const searchKeyword = ref('')
const showFilterDialog = ref(false)
const memberList = ref([])
const filterStatus = ref('all')
const filterPackage = ref('all')

// 统计数据
const totalMembers = computed(() => memberList.value.length)
const activeMembers = computed(() => memberList.value.filter(m => m.status === 'active').length)
const newMembersToday = ref(3)
const memberRevenue = ref('¥12,800')

// 筛选选项
const statusOptions = [
  { value: 'all', label: '全部' },
  { value: 'active', label: '有效' },
  { value: 'expired', label: '过期' },
  { value: 'inactive', label: '停用' }
]

const packageOptions = [
  { value: 'all', label: '全部套餐' },
  { value: 'basic', label: '基础会员' },
  { value: 'premium', label: '高级会员' },
  { value: 'vip', label: 'VIP会员' }
]

// 获取会员列表
const fetchMembersList = async () => {
  try {
    const response = await getMembers({
      keyword: searchKeyword.value,
      status: filterStatus.value,
      package: filterPackage.value
    })
    memberList.value = response.data.members || []
  } catch (error) {
    console.error('获取会员列表失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  fetchMembersList()
}

// 获取会员等级
const getMemberLevel = (packageName) => {
  if (packageName.includes('VIP')) return 'V'
  if (packageName.includes('高级')) return 'P'
  return 'B'
}

// 获取会员状态文本
const getMemberStatusText = (status) => {
  const statusMap = {
    active: '有效',
    expired: '过期',
    inactive: '停用'
  }
  return statusMap[status] || '未知'
}

// 获取进度百分比
const getProgressPercent = (member) => {
  if (member.total_times === 0) return 0
  return Math.round(((member.total_times - member.remaining_times) / member.total_times) * 100)
}

// 获取有效期天数
const getValidityDays = (validTo) => {
  const days = dayjs(validTo).diff(dayjs(), 'day')
  if (days < 0) return '已过期'
  if (days === 0) return '今日到期'
  if (days <= 7) return `${days}天后到期`
  return `还有${days}天`
}

// 获取有效期样式类
const getValidityClass = (validTo) => {
  const days = dayjs(validTo).diff(dayjs(), 'day')
  if (days < 0) return 'expired'
  if (days <= 7) return 'warning'
  return 'normal'
}

// 格式化手机号
const formatPhone = (phone) => {
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

// 处理会员点击
const handleMemberClick = (member) => {
  router.push(`/members/${member.id}`)
}

// 上传人脸
const uploadFace = (member) => {
  router.push(`/members/${member.id}/face`)
}

// 编辑会员
const editMember = (member) => {
  router.push(`/members/${member.id}/edit`)
}

// 查看消费记录
const viewHistory = (member) => {
  router.push(`/members/${member.id}/history`)
}

// 重置筛选
const resetFilter = () => {
  filterStatus.value = 'all'
  filterPackage.value = 'all'
  searchKeyword.value = ''
}

// 应用筛选
const applyFilter = () => {
  showFilterDialog.value = false
  fetchMembersList()
}

onMounted(() => {
  fetchMembersList()
})
</script>

<style lang="scss" scoped>
.members-page {
  background: linear-gradient(180deg, #f8f9ff 0%, #f0f2ff 100%);
  min-height: 100vh;
  padding-bottom: 20px;
}

// 页面头部
.page-header {
  background: linear-gradient(135deg, #ff6b6b 0%, #ffa726 100%);
  padding: 20px;
  color: white;
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

.add-button {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.1);
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
  
  &.active-stat {
    background: linear-gradient(135deg, #52c41a, #73d13d);
    color: white;
  }
}

.stat-number {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  opacity: 0.8;
}

// 搜索区域
.search-section {
  padding: 0 16px 16px;
}

.search-container {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-input-container {
  flex: 1;
  display: flex;
  align-items: center;
  background: white;
  border-radius: 24px;
  padding: 0 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.search-icon {
  font-size: 16px;
  margin-right: 8px;
  color: #999;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  padding: 12px 0;
  font-size: 14px;
  
  &::placeholder {
    color: #999;
  }
}

.filter-button {
  width: 40px;
  height: 40px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.2s ease;
  
  &:hover {
    transform: scale(1.1);
  }
}

// 会员列表
.members-container {
  padding: 0 16px;
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.member-card {
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

.member-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.member-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.member-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 600;
  color: white;
  position: relative;
  
  &.active {
    background: linear-gradient(135deg, #667eea, #764ba2);
  }
  
  &.expired {
    background: linear-gradient(135deg, #999, #666);
  }
  
  &.inactive {
    background: linear-gradient(135deg, #ff4d4f, #ff7875);
  }
}

.member-level {
  position: absolute;
  bottom: -4px;
  right: -4px;
  width: 20px;
  height: 20px;
  background: #ffa726;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: bold;
  color: white;
  border: 2px solid white;
}

.member-details {
  flex: 1;
}

.member-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 4px 0;
}

.member-phone {
  font-size: 13px;
  color: #999;
  margin: 0 0 4px 0;
}

.member-meta {
  display: flex;
  gap: 12px;
}

.member-package {
  font-size: 12px;
  color: #667eea;
  background: #f0f2ff;
  padding: 2px 8px;
  border-radius: 8px;
}

.member-times {
  font-size: 12px;
  color: #52c41a;
  background: #f6ffed;
  padding: 2px 8px;
  border-radius: 8px;
}

.member-status {
  padding: 6px 12px;
  border-radius: 12px;
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
  
  &.inactive {
    background: #f5f5f5;
    color: #999;
  }
}

// 会员内容
.member-content {
  border-top: 1px solid #f0f0f0;
  padding-top: 16px;
}

.member-progress {
  margin-bottom: 16px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.progress-label {
  font-size: 14px;
  color: #666;
}

.progress-text {
  font-size: 12px;
  color: #999;
}

.progress-bar {
  height: 6px;
  background: #f0f0f0;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.member-validity {
  margin-bottom: 16px;
}

.validity-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.validity-icon {
  font-size: 14px;
}

.validity-text {
  font-size: 14px;
  color: #666;
  flex: 1;
}

.validity-days {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 8px;
  
  &.normal {
    background: #f6ffed;
    color: #52c41a;
  }
  
  &.warning {
    background: #fff7e6;
    color: #fa8c16;
  }
  
  &.expired {
    background: #fff2f0;
    color: #ff4d4f;
  }
}

.member-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  flex: 1;
  padding: 8px 12px;
  border: none;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  
  &.face-btn {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
  }
  
  &.edit-btn {
    background: #f0f0f0;
    color: #666;
  }
  
  &.history-btn {
    background: #f0f0f0;
    color: #666;
  }
  
  &:hover {
    transform: translateY(-1px);
  }
}

.btn-icon {
  font-size: 12px;
}

// 筛选弹窗
.filter-dialog {
  background: white;
  border-radius: 16px 16px 0 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.dialog-close {
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
}

.filter-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.filter-group {
  margin-bottom: 24px;
}

.filter-label {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.filter-option {
  padding: 8px 16px;
  background: #f5f5f5;
  border-radius: 20px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &.active {
    background: linear-gradient(135deg, #ff6b6b, #ffa726);
    color: white;
  }
}

.filter-actions {
  display: flex;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #f0f0f0;
}

.filter-btn {
  flex: 1;
  height: 44px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &.reset-btn {
    background: #f5f5f5;
    color: #666;
  }
  
  &.apply-btn {
    background: linear-gradient(135deg, #ff6b6b, #ffa726);
    color: white;
  }
}

// 响应式设计
@media (max-width: 375px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }
  
  .member-card {
    padding: 16px;
  }
  
  .member-actions {
    flex-direction: column;
    gap: 6px;
  }
  
  .action-btn {
    padding: 6px 8px;
    font-size: 11px;
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

.member-card {
  animation: slideInUp 0.4s ease-out;
}

.member-card:nth-child(1) { animation-delay: 0.1s; }
.member-card:nth-child(2) { animation-delay: 0.2s; }
.member-card:nth-child(3) { animation-delay: 0.3s; }
.member-card:nth-child(4) { animation-delay: 0.4s; }
.member-card:nth-child(5) { animation-delay: 0.5s; }
</style>