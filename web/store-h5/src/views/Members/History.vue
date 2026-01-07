<template>
  <div class="history-page">
    <!-- 页面头部 -->
    <van-nav-bar
      title="使用记录"
      left-arrow
      @click-left="goBack"
      fixed
      placeholder
    >
      <template #right>
        <!-- <div class="add-button" @click="showCreateDialog = true">
          <i class="add-icon">➕</i>
        </div> -->
      </template>
    </van-nav-bar>

    <!-- 会员信息卡片 -->
    <div class="member-info-card" v-if="memberInfo">
      <div class="member-header">
        <div class="member-avatar" :class="memberInfo.status">
          {{ memberInfo.name?.charAt(0) }}
        </div>
        <div class="member-details">
          <h3 class="member-name">{{ memberInfo.name }}</h3>
          <p class="member-phone">{{ formatPhone(memberInfo.phone) }}</p>
          <div class="member-meta">
            <span class="member-package">{{ memberInfo.package_name }}</span>
            <span class="member-times">使用{{ getUsedTimes(memberInfo) }}次</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 使用记录列表 -->
    <div class="usage-list-container">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="fetchUsageList"
      >
        <div class="usage-list">
          <div 
            v-for="usage in usageList"
            :key="usage.id"
            class="usage-card"
          >
            <div class="usage-header">
              <div class="usage-service">
                <i class="service-icon">💅</i>
                <span class="service-name">{{ usage.service_item }}</span>
              </div>
              <div class="usage-date">
                {{ formatDate(usage.usage_date) }}
              </div>
            </div>
            
            <div class="usage-content">
              <div class="usage-info-item">
                <i class="info-icon">🏪</i>
                <span class="info-label">门店：</span>
                <span class="info-value">{{ usage.store_name || '-' }}</span>
              </div>
              <div class="usage-info-item" v-if="usage.technician_name">
                <i class="info-icon">👤</i>
                <span class="info-label">美甲师：</span>
                <span class="info-value">{{ usage.technician_name }}</span>
              </div>
              <div class="usage-info-item" v-if="usage.remark">
                <i class="info-icon">📝</i>
                <span class="info-label">备注：</span>
                <span class="info-value">{{ usage.remark }}</span>
              </div>
            </div>
            
            <div class="usage-footer">
              <span class="usage-time">{{ formatDateTime(usage.created_at) }}</span>
              <button 
                v-if="canDelete"
                class="delete-btn" 
                @click="handleDelete(usage)"
              >
                <i class="delete-icon">🗑️</i>
                删除
              </button>
            </div>
          </div>
        </div>
        
        <div v-if="usageList.length === 0 && !loading" class="empty-state">
          <div class="empty-icon">📋</div>
          <div class="empty-text">暂无使用记录</div>
        </div>
      </van-list>
    </div>

    <!-- 新增使用记录弹窗 -->
    <van-popup 
      v-model:show="showCreateDialog" 
      position="bottom" 
      round
      :style="{ height: '80%' }"
    >
      <div class="create-dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">新增使用记录</h3>
          <div class="dialog-close" @click="showCreateDialog = false">×</div>
        </div>
        
        <div class="dialog-content">
          <van-form @submit="handleSubmit">
            <van-cell-group inset>
              <van-field
                v-model="form.service_item"
                label="服务项目"
                placeholder="请输入服务项目，如：美甲-单色"
                required
                :rules="[{ required: true, message: '请输入服务项目' }]"
              />
              
              <van-field
                v-model="form.store_name"
                label="使用门店"
                placeholder="请选择门店"
                readonly
                is-link
                required
                @click="showStorePicker = true"
                :rules="[{ required: true, message: '请选择门店' }]"
              />
              
              <van-field
                v-model="form.technician_name"
                label="美甲师"
                placeholder="请选择美甲师（可选）"
                readonly
                is-link
                @click="showTechnicianPicker = true"
              />
              
              <van-field
                v-model="form.usage_date"
                label="使用日期"
                placeholder="选择使用日期"
                readonly
                is-link
                required
                @click="showDatePicker = true"
                :rules="[{ required: true, message: '请选择使用日期' }]"
              />
              
              <van-field
                v-model="form.remark"
                label="备注"
                type="textarea"
                placeholder="请输入备注（可选）"
                rows="3"
                autosize
                maxlength="500"
                show-word-limit
              />
            </van-cell-group>
            
            <div class="form-actions">
              <van-button 
                round 
                block 
                type="primary" 
                native-type="submit"
                :loading="submitting"
              >
                确定
              </van-button>
            </div>
          </van-form>
        </div>
      </div>
    </van-popup>

    <!-- 门店选择器 -->
    <van-popup v-model:show="showStorePicker" position="bottom">
      <van-picker
        :columns="storeColumns"
        @confirm="onStoreConfirm"
        @cancel="showStorePicker = false"
      />
    </van-popup>

    <!-- 美甲师选择器 -->
    <van-popup v-model:show="showTechnicianPicker" position="bottom">
      <van-picker
        :columns="technicianColumns"
        @confirm="onTechnicianConfirm"
        @cancel="showTechnicianPicker = false"
      />
    </van-popup>

    <!-- 日期选择器 -->
    <van-popup v-model:show="showDatePicker" position="bottom">
      <van-date-picker
        v-model="currentDate"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getMemberById, getMemberUsages, createUsage, deleteUsage } from '@/api/members'
import { getStores } from '@/api/stores'
import { getStaffList } from '@/api/staff'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const finished = ref(false)
const submitting = ref(false)
const usageList = ref([])
const memberInfo = ref(null)
const storeList = ref([])
const technicianList = ref([])

const showCreateDialog = ref(false)
const showStorePicker = ref(false)
const showTechnicianPicker = ref(false)
const showDatePicker = ref(false)
const currentDate = ref(new Date())

const form = reactive({
  service_item: '',
  store_id: null,
  store_name: '',
  technician_id: null,
  technician_name: '',
  usage_date: dayjs().format('YYYY-MM-DD'),
  remark: ''
})

// 门店选择器列
const storeColumns = computed(() => {
  return storeList.value.map(store => ({
    text: store.name,
    value: store.id
  }))
})

// 美甲师选择器列（包含"无"选项）
const technicianColumns = computed(() => {
  const columns = [{ text: '无', value: null }]
  technicianList.value.forEach(tech => {
    columns.push({
      text: tech.username || tech.name,
      value: tech.id
    })
  })
  return columns
})

// 格式化手机号
const formatPhone = (phone) => {
  if (!phone) return '-'
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

// 格式化日期时间
const formatDateTime = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 计算已使用次数
const getUsedTimes = (member) => {
  if (!member) return 0
  const total = member.total_times || 0
  const remaining = member.remaining_times || 0
  return total - remaining
}

// 检查是否有删除权限（仅管理员）
const canDelete = computed(() => {
  const userRole = userStore.userInfo?.role
  return userRole === 'admin'
})

// 返回
const goBack = () => {
  router.back()
}

// 获取会员信息
const fetchMemberInfo = async () => {
  const memberId = route.params.id
  if (!memberId) {
    showToast('会员ID不存在')
    router.back()
    return
  }
  
  try {
    const response = await getMemberById(memberId)
    if (response.data) {
      memberInfo.value = response.data
    }
  } catch (error) {
    console.error('获取会员信息失败:', error)
    // 如果是401/403错误，request.js已经处理了跳转，这里不需要额外处理
    if (error.response?.status !== 401 && error.response?.status !== 403) {
      showToast('获取会员信息失败')
    }
  }
}

// 获取使用记录列表
const fetchUsageList = async () => {
  const memberId = route.params.id
  if (!memberId || loading.value) return
  
  loading.value = true
  try {
    const response = await getMemberUsages(memberId)
    if (response.data) {
      usageList.value = Array.isArray(response.data) ? response.data : []
    } else {
      usageList.value = []
    }
    finished.value = true
  } catch (error) {
    console.error('获取使用记录失败:', error)
    // 如果是401/403错误，request.js已经处理了跳转，这里不需要额外处理
    if (error.response?.status !== 401 && error.response?.status !== 403) {
      showToast('获取使用记录失败')
    }
    usageList.value = []
    finished.value = true
  } finally {
    loading.value = false
  }
}

// 获取门店列表
const fetchStoreList = async () => {
  try {
    const response = await getStores({ page: 1, page_size: 1000 })
    if (response.data) {
      if (response.data.list) {
        storeList.value = response.data.list
      } else if (Array.isArray(response.data)) {
        storeList.value = response.data
      } else {
        storeList.value = []
      }
    } else {
      storeList.value = []
    }
  } catch (error) {
    console.error('获取门店列表失败:', error)
  }
}

// 获取美甲师列表
const fetchTechnicianList = async () => {
  try {
    const response = await getStaffList()
    if (response.data) {
      if (response.data.list) {
        technicianList.value = response.data.list
      } else if (Array.isArray(response.data)) {
        technicianList.value = response.data
      } else {
        technicianList.value = []
      }
    } else {
      technicianList.value = []
    }
  } catch (error) {
    console.error('获取美甲师列表失败:', error)
  }
}

// 门店选择确认
const onStoreConfirm = ({ selectedOptions }) => {
  const selected = selectedOptions[0]
  form.store_id = selected.value
  form.store_name = selected.text
  showStorePicker.value = false
}

// 美甲师选择确认
const onTechnicianConfirm = ({ selectedOptions }) => {
  const selected = selectedOptions[0]
  form.technician_id = selected.value
  form.technician_name = selected.value ? selected.text : ''
  showTechnicianPicker.value = false
}

// 日期选择确认
const onDateConfirm = () => {
  form.usage_date = dayjs(currentDate.value).format('YYYY-MM-DD')
  showDatePicker.value = false
}

// 重置表单
const resetForm = () => {
  form.service_item = ''
  form.store_id = null
  form.store_name = ''
  form.technician_id = null
  form.technician_name = ''
  form.usage_date = dayjs().format('YYYY-MM-DD')
  form.remark = ''
  currentDate.value = new Date()
}

// 提交表单
const handleSubmit = async () => {
  if (!form.service_item || !form.store_id || !form.usage_date) {
    showToast('请填写完整信息')
    return
  }
  
  const memberId = route.params.id
  submitting.value = true
  
  try {
    const data = {
      service_item: form.service_item,
      store_id: form.store_id,
      technician_id: form.technician_id || undefined,
      usage_date: form.usage_date,
      remark: form.remark || undefined
    }
    
    await createUsage(memberId, data)
    showToast('使用记录创建成功')
    showCreateDialog.value = false
    resetForm()
    // 重新获取列表
    usageList.value = []
    finished.value = false
    fetchUsageList()
    // 重新获取会员信息以更新剩余次数
    fetchMemberInfo()
  } catch (error) {
    console.error('创建使用记录失败:', error)
    showToast('创建使用记录失败')
  } finally {
    submitting.value = false
  }
}

// 删除使用记录
const handleDelete = async (usage) => {
  try {
    await showConfirmDialog({
      title: '删除使用记录',
      message: '确定要删除这条使用记录吗？'
    })
    
    await deleteUsage(usage.id)
    showToast('删除成功')
    // 重新获取列表
    usageList.value = []
    finished.value = false
    fetchUsageList()
    // 重新获取会员信息以更新剩余次数
    fetchMemberInfo()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除使用记录失败:', error)
      showToast('删除使用记录失败')
    }
  }
}

onMounted(() => {
  fetchMemberInfo()
  fetchUsageList()
//   fetchStoreList()
//   fetchTechnicianList()
})
</script>

<style lang="scss" scoped>
.history-page {
  background: linear-gradient(180deg, #f8f9ff 0%, #f0f2ff 100%);
  min-height: 100vh;
  padding-bottom: 20px;
}

.add-button {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.add-icon {
  font-size: 18px;
}

// 会员信息卡片
.member-info-card {
  margin: 16px;
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.member-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.member-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 600;
  color: white;
  
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

.member-details {
  flex: 1;
}

.member-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 4px 0;
}

.member-phone {
  font-size: 14px;
  color: #999;
  margin: 0 0 8px 0;
}

.member-meta {
  display: flex;
  gap: 8px;
}

.member-package {
  font-size: 12px;
  color: #667eea;
  background: #f0f2ff;
  padding: 4px 12px;
  border-radius: 12px;
}

.member-times {
  font-size: 12px;
  color: #52c41a;
  background: #f6ffed;
  padding: 4px 12px;
  border-radius: 12px;
}

// 使用记录列表
.usage-list-container {
  padding: 0 16px;
}

.usage-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.usage-card {
  background: white;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.usage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.usage-service {
  display: flex;
  align-items: center;
  gap: 8px;
}

.service-icon {
  font-size: 20px;
}

.service-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.usage-date {
  font-size: 14px;
  color: #999;
}

.usage-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.usage-info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.info-icon {
  font-size: 16px;
}

.info-label {
  color: #666;
}

.info-value {
  color: #333;
  flex: 1;
}

.usage-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.usage-time {
  font-size: 12px;
  color: #999;
}

.delete-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: #fff2f0;
  color: #ff4d4f;
  border: none;
  border-radius: 8px;
  font-size: 12px;
  cursor: pointer;
}

.delete-icon {
  font-size: 14px;
}

// 空状态
.empty-state {
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 14px;
  color: #999;
}

// 创建弹窗
.create-dialog {
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

.dialog-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.form-actions {
  margin-top: 20px;
  padding: 0 16px 20px;
}

// 响应式设计
@media (max-width: 375px) {
  .usage-card {
    padding: 12px;
  }
  
  .service-name {
    font-size: 14px;
  }
}
</style>

