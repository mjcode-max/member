<template>
  <div class="member-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1 class="page-title">我的会员</h1>
          <p class="page-subtitle">享受专属会员服务</p>
        </div>
      </div>
    </div>

    <!-- 会员卡片 -->
    <div class="member-container" v-if="memberInfo">
      <div class="member-card">
        <div class="card-background">
          <div class="card-pattern"></div>
        </div>
        <div class="card-content">
          <div class="card-header">
            <div class="member-level">{{ getMemberLevel() }}</div>
            <div class="member-status" :class="memberInfo.status">
              {{ getMemberStatusText() }}
            </div>
          </div>
          <div class="card-body">
            <h2 class="member-name">{{ memberInfo.name }}</h2>
            <p class="member-phone">{{ formatPhone(memberInfo.phone) }}</p>
            <div class="member-package">{{ memberInfo.package_name }}</div>
          </div>
          <div class="card-footer">
            <div class="validity-info">
              <i class="validity-icon">📅</i>
              <span>有效期至 {{ formatDate(memberInfo.valid_to) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 会员统计 -->
      <div class="member-stats">
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-icon">🎯</div>
            <div class="stat-content">
              <div class="stat-number">{{ memberInfo.remaining_times }}</div>
              <div class="stat-label">剩余次数</div>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-icon">📅</div>
            <div class="stat-content">
              <div class="stat-number">{{ getValidityDays() }}</div>
              <div class="stat-label">剩余天数</div>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-icon">🏆</div>
            <div class="stat-content">
              <div class="stat-number">{{ consumptionCount }}</div>
              <div class="stat-label">已消费</div>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-icon">💎</div>
            <div class="stat-content">
              <div class="stat-number">{{ memberInfo.total_times }}</div>
              <div class="stat-label">总次数</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 出示会员码按钮 -->
      <div class="member-actions">
        <button 
          class="show-code-btn"
          @click="showMemberCode"
          :disabled="!canShowCode"
        >
          <i class="code-icon">📱</i>
          <span>出示会员码</span>
          <div class="btn-glow"></div>
        </button>
        <div class="action-note" v-if="!canShowCode">
          <i class="note-icon">⚠️</i>
          <span>{{ getDisabledReason() }}</span>
        </div>
      </div>

      <!-- 消费记录 -->
      <div class="consumption-section">
        <div class="section-header">
          <i class="section-icon">📋</i>
          <span class="section-title">消费记录</span>
        </div>
        <div class="consumption-list">
          <div 
            v-for="record in consumptionList"
            :key="record.id"
            class="consumption-item"
          >
            <div class="consumption-info">
              <div class="consumption-service">
                <i class="service-icon">{{ record.service_type === 'manicure' ? '💅' : '👁️' }}</i>
                <span class="service-name">{{ getServiceText(record.service_type) }}</span>
              </div>
              <div class="consumption-time">{{ formatDateTime(record.consumed_at) }}</div>
              <div class="consumption-staff">服务技师：{{ record.staff_name || '未知' }}</div>
            </div>
            <div class="consumption-status">
              <i class="status-icon">✅</i>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 非会员状态 -->
    <div class="non-member-container" v-else>
      <div class="non-member-card">
        <div class="non-member-icon">💎</div>
        <h3 class="non-member-title">您还不是会员</h3>
        <p class="non-member-desc">请联系门店工作人员办理会员卡</p>
        <div class="member-benefits">
          <div class="benefit-item">
            <i class="benefit-icon">✨</i>
            <span>专享会员价格</span>
          </div>
          <div class="benefit-item">
            <i class="benefit-icon">🎁</i>
            <span>生日特权礼遇</span>
          </div>
          <div class="benefit-item">
            <i class="benefit-icon">⚡</i>
            <span>优先预约服务</span>
          </div>
        </div>
        <button class="contact-btn" @click="handleContact">
          <i class="contact-icon">📞</i>
          联系门店
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { getMemberInfo } from '@/api/members'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const memberInfo = ref(null)
const consumptionList = ref([])
const consumptionCount = computed(() => consumptionList.value.length)

// 检查是否可以出示会员码
const canShowCode = computed(() => {
  if (!memberInfo.value) return false
  if (memberInfo.value.status !== 'active') return false
  if (memberInfo.value.remaining_times <= 0) return false
  if (dayjs().isAfter(dayjs(memberInfo.value.valid_to))) return false
  return true
})

// 获取会员等级
const getMemberLevel = () => {
  if (!memberInfo.value) return ''
  const packageName = memberInfo.value.package_name
  if (packageName.includes('VIP')) return 'VIP会员'
  if (packageName.includes('高级')) return '高级会员'
  return '基础会员'
}

// 获取会员状态文本
const getMemberStatusText = () => {
  if (!memberInfo.value) return ''
  const statusMap = {
    active: '有效',
    expired: '已过期',
    inactive: '已停用'
  }
  return statusMap[memberInfo.value.status] || '未知'
}

// 获取剩余天数
const getValidityDays = () => {
  if (!memberInfo.value) return 0
  const days = dayjs(memberInfo.value.valid_to).diff(dayjs(), 'day')
  return Math.max(0, days)
}

// 获取禁用原因
const getDisabledReason = () => {
  if (!memberInfo.value) return '非会员用户'
  if (memberInfo.value.status !== 'active') return '会员状态异常'
  if (memberInfo.value.remaining_times <= 0) return '次数已用完'
  if (dayjs().isAfter(dayjs(memberInfo.value.valid_to))) return '会员已过期'
  return ''
}

// 出示会员码
const showMemberCode = () => {
  if (!canShowCode.value) {
    showToast(getDisabledReason())
    return
  }
  router.push(`/member/code/${memberInfo.value.id}`)
}

// 获取服务类型文本
const getServiceText = (type) => {
  return type === 'manicure' ? '美甲服务' : '美睫服务'
}

// 格式化手机号
const formatPhone = (phone) => {
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY年MM月DD日')
}

// 格式化日期时间
const formatDateTime = (datetime) => {
  return dayjs(datetime).format('MM月DD日 HH:mm')
}

// 联系门店
const handleContact = () => {
  showToast('请联系门店：400-123-4567')
}

// 获取会员信息
const fetchMemberInfo = async () => {
  try {
    // 这里应该根据用户手机号获取会员信息
    const phone = userStore.getUserPhone() // 假设有这个方法
    if (phone) {
      const response = await getMemberInfo(phone)
      memberInfo.value = response.data.member
      consumptionList.value = response.data.consumptions || []
    }
  } catch (error) {
    console.error('获取会员信息失败:', error)
    // 如果获取失败，说明不是会员
    memberInfo.value = null
  }
}

onMounted(() => {
  fetchMemberInfo()
})
</script>

<style lang="scss" scoped>
.member-page {
  background: #ffffff;
  min-height: 100vh;
  padding-bottom: 0;
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
    bottom: -20px;
    left: 0;
    right: 0;
    height: 20px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 0 0 20px 20px;
  }
}

.header-content {
  text-align: center;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.page-subtitle {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
}

// 会员容器
.member-container {
  padding: 20px 16px;
}

// 会员卡片
.member-card {
  position: relative;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 20px;
  overflow: hidden;
  box-shadow: 0 15px 35px rgba(102, 126, 234, 0.3);
}

.card-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  opacity: 0.1;
}

.card-pattern {
  position: absolute;
  top: -50%;
  right: -20%;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.3) 0%, transparent 70%);
  border-radius: 50%;
  animation: float 6s ease-in-out infinite;
}

.card-content {
  position: relative;
  z-index: 1;
  color: white;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.member-level {
  font-size: 14px;
  font-weight: 600;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  backdrop-filter: blur(10px);
}

.member-status {
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  
  &.active {
    background: rgba(82, 196, 26, 0.2);
    color: #73d13d;
  }
  
  &.expired {
    background: rgba(255, 77, 79, 0.2);
    color: #ff7875;
  }
}

.card-body {
  margin-bottom: 20px;
}

.member-name {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 8px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.member-phone {
  font-size: 16px;
  opacity: 0.9;
  margin: 0 0 12px 0;
}

.member-package {
  font-size: 14px;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 12px;
  display: inline-block;
}

.card-footer {
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  padding-top: 16px;
}

.validity-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.validity-icon {
  font-size: 16px;
}

// 会员统计
.member-stats {
  margin-bottom: 24px;
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

// 会员操作
.member-actions {
  margin-bottom: 24px;
  text-align: center;
}

.show-code-btn {
  position: relative;
  width: 100%;
  height: 64px;
  background: linear-gradient(135deg, #ff6b6b, #ffa726);
  border: none;
  border-radius: 32px;
  color: white;
  font-size: 18px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  overflow: hidden;
  
  &:hover:not(:disabled) {
    transform: translateY(-4px);
    box-shadow: 0 12px 40px rgba(255, 107, 107, 0.4);
  }
  
  &:disabled {
    background: #ccc;
    cursor: not-allowed;
  }
}

.code-icon {
  font-size: 24px;
}

.btn-glow {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  animation: shine 3s infinite;
}

.action-note {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 12px;
  font-size: 14px;
  color: #ff4d4f;
}

.note-icon {
  font-size: 16px;
}

// 消费记录
.consumption-section {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
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

.consumption-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.consumption-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  transition: all 0.2s ease;
  
  &:hover {
    background: #e9ecef;
  }
}

.consumption-info {
  flex: 1;
}

.consumption-service {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.service-icon {
  font-size: 18px;
}

.service-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.consumption-time {
  font-size: 13px;
  color: #666;
  margin-bottom: 2px;
}

.consumption-staff {
  font-size: 12px;
  color: #999;
}

.consumption-status {
  display: flex;
  align-items: center;
}

.status-icon {
  font-size: 18px;
  color: #52c41a;
}

// 非会员状态
.non-member-container {
  padding: 40px 16px;
  text-align: center;
}

.non-member-card {
  background: white;
  border-radius: 20px;
  padding: 40px 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.non-member-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.non-member-title {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin-bottom: 12px;
}

.non-member-desc {
  font-size: 16px;
  color: #666;
  margin: 0 0 24px 0;
  line-height: 1.5;
}

.member-benefits {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 32px;
}

.benefit-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 15px;
  color: #666;
}

.benefit-icon {
  font-size: 18px;
  color: #667eea;
}

.contact-btn {
  padding: 14px 32px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin: 0 auto;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
  }
}

.contact-icon {
  font-size: 16px;
}

// 动画效果
@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-10px) rotate(5deg); }
}

@keyframes shine {
  0% { left: -100%; }
  50%, 100% { left: 100%; }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.member-card,
.member-stats,
.member-actions,
.consumption-section {
  animation: slideInUp 0.6s ease-out;
}

.member-card { animation-delay: 0.1s; }
.member-stats { animation-delay: 0.2s; }
.member-actions { animation-delay: 0.3s; }
.consumption-section { animation-delay: 0.4s; }

// 响应式设计
@media (max-width: 375px) {
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  
  .member-card {
    padding: 20px;
  }
  
  .member-name {
    font-size: 20px;
  }
  
  .show-code-btn {
    height: 56px;
    font-size: 16px;
  }
}
</style>