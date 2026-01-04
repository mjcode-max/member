<template>
  <div class="staff-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1 class="page-title">员工管理</h1>
          <p class="page-subtitle">管理门店美甲师团队</p>
        </div>
        <div class="header-actions">
          <div class="add-button" @click="showAddDialog = true">
            <i class="add-icon">➕</i>
          </div>
        </div>
      </div>
    </div>

    <!-- 员工统计 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-number">{{ totalStaff }}</div>
          <div class="stat-label">总员工</div>
        </div>
        <div class="stat-item active-stat">
          <div class="stat-number">{{ activeStaff }}</div>
          <div class="stat-label">在岗中</div>
        </div>
        <div class="stat-item rest-stat">
          <div class="stat-number">{{ restStaff }}</div>
          <div class="stat-label">休息中</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ todayServices }}</div>
          <div class="stat-label">今日服务</div>
        </div>
      </div>
    </div>

    <!-- 员工列表 -->
    <div class="staff-container">
      <div class="section-header">
        <span class="section-title">员工列表</span>
        <div class="batch-actions">
          <button class="batch-btn" @click="batchSetStatus('active')">
            批量上岗
          </button>
          <button class="batch-btn" @click="batchSetStatus('rest')">
            批量休息
          </button>
        </div>
      </div>

      <div class="staff-list">
        <div 
          v-for="staff in staffList"
          :key="staff.id"
          class="staff-card"
          :class="{ selected: selectedStaff.includes(staff.id) }"
          @click="toggleSelect(staff.id)"
        >
          <div class="staff-header">
            <div class="staff-info">
              <div class="staff-avatar" :class="getWorkStatusClass(staff.work_status)">
                {{ (staff.username || staff.name || 'U').charAt(0) }}
                <div class="status-indicator" :class="getWorkStatusClass(staff.work_status)"></div>
              </div>
              <div class="staff-details">
                <h4 class="staff-name">{{ staff.username || staff.name || '未命名' }}</h4>
                <p class="staff-phone">{{ formatPhone(staff.phone) }}</p>
                <div class="staff-meta">
                  <span class="join-date">入职 {{ formatDate(staff.created_at) }}</span>
                </div>
              </div>
            </div>
            <div class="staff-actions">
              <div class="work-status-toggle">
                <div 
                  class="status-option"
                  :class="{ active: getWorkStatusClass(staff.work_status) === 'working' }"
                  @click.stop="updateWorkStatus(staff.id, 'working')"
                >
                  <i class="status-icon">💼</i>
                  <span>在岗</span>
                </div>
                <div 
                  class="status-option"
                  :class="{ active: getWorkStatusClass(staff.work_status) === 'rest' }"
                  @click.stop="updateWorkStatus(staff.id, 'rest')"
                >
                  <i class="status-icon">😴</i>
                  <span>休息</span>
                </div>
              </div>
            </div>
          </div>

          <div class="staff-stats" v-if="getWorkStatusClass(staff.work_status) === 'working'">
            <div class="stat-row">
              <div class="stat-col">
                <span class="stat-value">{{ staff.today_bookings || 0 }}</span>
                <span class="stat-desc">今日预约</span>
              </div>
              <div class="stat-col">
                <span class="stat-value">{{ staff.completed_today || 0 }}</span>
                <span class="stat-desc">已完成</span>
              </div>
              <div class="stat-col">
                <span class="stat-value">{{ staff.rating || '4.8' }}</span>
                <span class="stat-desc">评分</span>
              </div>
            </div>
          </div>

          <div class="rest-info" v-else>
            <div class="rest-message">
              <i class="rest-icon">💤</i>
              <span>员工休息中，不接受新预约</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加员工弹窗 -->
    <van-popup 
      v-model:show="showAddDialog" 
      position="center"
      round
      :style="{ width: '90%' }"
    >
      <div class="add-dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">添加员工</h3>
          <div class="dialog-close" @click="showAddDialog = false">×</div>
        </div>
        <div class="dialog-content">
          <div class="form-group">
            <label class="form-label">用户名</label>
            <input 
              v-model="newStaff.username"
              type="text"
              placeholder="请输入用户名"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">手机号码</label>
            <input 
              v-model="newStaff.phone"
              type="tel"
              placeholder="请输入手机号码"
              class="form-input"
            />
          </div>
        </div>
        <div class="dialog-actions">
          <button class="dialog-btn cancel-btn" @click="showAddDialog = false">
            取消
          </button>
          <button class="dialog-btn confirm-btn" @click="addStaff">
            确认添加
          </button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onActivated } from 'vue'
import { showToast, showConfirmDialog } from 'vant'
import { getStaffList, createStaff, updateStaffStatus } from '@/api/staff'
import dayjs from 'dayjs'

const showAddDialog = ref(false)
const selectedStaff = ref([])
const staffList = ref([])

const newStaff = reactive({
  username: '',
  phone: '',
  password: '123456'
})

// 统计数据
const totalStaff = computed(() => staffList.value.length)
const activeStaff = computed(() => staffList.value.filter(s => {
  const status = s.work_status
  return status === 'working' || status === 'active'
}).length)
const restStaff = computed(() => staffList.value.filter(s => s.work_status === 'rest').length)
const todayServices = ref(25)

// 获取员工列表
const fetchStaffList = async () => {
  try {
    const response = await getStaffList()
    // 处理分页响应格式
    if (response && response.data) {
      if (response.data.list && Array.isArray(response.data.list)) {
        staffList.value = response.data.list
      } else if (Array.isArray(response.data)) {
        staffList.value = response.data
      } else {
        staffList.value = []
      }
    } else {
      staffList.value = []
    }
  } catch (error) {
    console.error('获取员工列表失败:', error)
    staffList.value = []
  }
}

// 更新工作状态
const updateWorkStatus = async (staffId, status) => {
  try {
    // 后端使用 working/rest/offline，前端需要转换
    const backendStatus = status === 'active' ? 'working' : status
    await updateStaffStatus(staffId, { work_status: backendStatus })
    
    // 更新本地数据
    const staff = staffList.value.find(s => s.id === staffId)
    if (staff) {
      // 后端返回的是 working/rest/offline，前端显示需要转换
      staff.work_status = backendStatus
    }
    
    showToast(`已设置为${getWorkStatusText(backendStatus)}状态`)
    // 重新获取列表以确保数据同步
    await fetchStaffList()
  } catch (error) {
    console.error('更新工作状态失败:', error)
    showToast('更新失败，请重试')
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
  if (status === 'active') return 'working'
  return status
}

// 批量设置状态
const batchSetStatus = async (status) => {
  if (selectedStaff.value.length === 0) {
    showToast('请先选择员工')
    return
  }

  try {
    const statusText = status === 'active' || status === 'working' ? '在岗' : '休息'
    await showConfirmDialog({
      title: '确认操作',
      message: `确定要将选中的${selectedStaff.value.length}名员工设置为${statusText}状态吗？`
    })

    // 批量更新
    const backendStatus = status === 'active' ? 'working' : status
    for (const staffId of selectedStaff.value) {
      await updateWorkStatus(staffId, backendStatus)
    }

    selectedStaff.value = []
    showToast('批量操作成功')
  } catch (error) {
    // 用户取消或操作失败
  }
}

// 切换选择
const toggleSelect = (staffId) => {
  const index = selectedStaff.value.indexOf(staffId)
  if (index > -1) {
    selectedStaff.value.splice(index, 1)
  } else {
    selectedStaff.value.push(staffId)
  }
}

// 添加员工
const addStaff = async () => {
  if (!newStaff.username || !newStaff.phone) {
    showToast('请填写完整信息')
    return
  }

  try {
    const data = {
      username: newStaff.username,
      phone: newStaff.phone,
      password: newStaff.password || '123456'
    }
    await createStaff(data)
    showToast('员工添加成功')
    showAddDialog.value = false
    newStaff.username = ''
    newStaff.phone = ''
    newStaff.password = '123456'
    fetchStaffList()
  } catch (error) {
    console.error('添加员工失败:', error)
    showToast('添加失败，请重试')
  }
}

// 格式化手机号
const formatPhone = (phone) => {
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

onMounted(() => {
  fetchStaffList()
})

// 页面激活时刷新数据（从其他页面返回时）
onActivated(() => {
  fetchStaffList()
})
</script>

<style lang="scss" scoped>
.staff-page {
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

.add-icon {
  font-size: 20px;
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
  
  &.rest-stat {
    background: linear-gradient(135deg, #faad14, #ffc53d);
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

// 员工列表
.staff-container {
  padding: 0 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.batch-actions {
  display: flex;
  gap: 8px;
}

.batch-btn {
  padding: 6px 12px;
  background: #f0f0f0;
  border: none;
  border-radius: 16px;
  font-size: 12px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: #667eea;
    color: white;
  }
}

.staff-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.staff-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  
  &.selected {
    border-color: #ff6b6b;
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(255, 107, 107, 0.2);
  }
  
  &:hover:not(.selected) {
    transform: translateY(-1px);
    box-shadow: 0 6px 24px rgba(0, 0, 0, 0.12);
  }
}

.staff-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.staff-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.staff-avatar {
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

.status-indicator {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 2px solid white;
  
  &.working,
  &.active {
    background: #52c41a;
  }
  
  &.rest {
    background: #faad14;
  }
  
  &.offline {
    background: #999;
  }
}

.staff-details {
  flex: 1;
}

.staff-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 4px 0;
}

.staff-phone {
  font-size: 13px;
  color: #999;
  margin: 0 0 4px 0;
}

.staff-meta {
  display: flex;
  gap: 12px;
}

.join-date {
  font-size: 12px;
  color: #999;
}

.work-status-toggle {
  display: flex;
  background: #f5f5f5;
  border-radius: 20px;
  padding: 2px;
}

.status-option {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: 18px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &.active {
    background: white;
    color: #333;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  &:not(.active) {
    color: #999;
  }
}

.status-icon {
  font-size: 12px;
}

// 员工统计
.staff-stats {
  border-top: 1px solid #f0f0f0;
  padding-top: 16px;
}

.stat-row {
  display: flex;
  justify-content: space-around;
}

.stat-col {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.stat-desc {
  font-size: 12px;
  color: #999;
}

// 休息状态信息
.rest-info {
  border-top: 1px solid #f0f0f0;
  padding-top: 16px;
}

.rest-message {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: #fff7e6;
  border-radius: 8px;
  color: #fa8c16;
  font-size: 14px;
}

.rest-icon {
  font-size: 16px;
}

// 添加员工弹窗
.add-dialog {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  max-width: 400px;
  width: 100%;
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
  transition: all 0.2s ease;
  
  &:hover {
    background: #e9ecef;
  }
}

.dialog-content {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #f0f0f0;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: #ff6b6b;
  }
  
  &::placeholder {
    color: #999;
  }
}

.dialog-actions {
  display: flex;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #f0f0f0;
}

.dialog-btn {
  flex: 1;
  height: 44px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &.cancel-btn {
    background: #f5f5f5;
    color: #666;
    
    &:hover {
      background: #e9ecef;
    }
  }
  
  &.confirm-btn {
    background: linear-gradient(135deg, #ff6b6b, #ffa726);
    color: white;
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(255, 107, 107, 0.3);
    }
  }
}

// 响应式设计
@media (max-width: 375px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }
  
  .staff-card {
    padding: 16px;
  }
  
  .staff-avatar {
    width: 40px;
    height: 40px;
    font-size: 16px;
  }
  
  .work-status-toggle {
    flex-direction: column;
    gap: 4px;
  }
  
  .status-option {
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

.staff-card {
  animation: slideInUp 0.4s ease-out;
}

.staff-card:nth-child(1) { animation-delay: 0.1s; }
.staff-card:nth-child(2) { animation-delay: 0.2s; }
.staff-card:nth-child(3) { animation-delay: 0.3s; }
.staff-card:nth-child(4) { animation-delay: 0.4s; }
.staff-card:nth-child(5) { animation-delay: 0.5s; }
</style>