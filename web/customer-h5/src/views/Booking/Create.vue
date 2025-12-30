<template>
  <div class="booking-create-page">
    <!-- 页面头部 -->
    <van-nav-bar
      title="预约服务"
      left-text="返回"
      left-arrow
      @click-left="$router.back()"
      class="custom-nav"
    />
    
    <!-- 进度指示器 -->
    <div class="progress-container">
      <div class="progress-steps">
        <div class="step active">
          <div class="step-icon">1</div>
          <div class="step-text">选择服务</div>
        </div>
        <div class="step-line"></div>
        <div class="step" :class="{ active: form.store_id }">
          <div class="step-icon">2</div>
          <div class="step-text">选择时间</div>
        </div>
        <div class="step-line"></div>
        <div class="step" :class="{ active: form.customer_name }">
          <div class="step-icon">3</div>
          <div class="step-text">填写信息</div>
        </div>
      </div>
    </div>

    <div class="content-container">
      <!-- 门店选择卡片 -->
      <div class="section-card">
        <div class="section-header">
          <i class="icon-store">🏪</i>
          <span class="section-title">选择门店</span>
        </div>
        <div class="store-selector" @click="showStorePicker = true">
          <div v-if="selectedStore" class="selected-store">
            <div class="store-name">{{ selectedStore.name }}</div>
            <div class="store-address">{{ selectedStore.address }}</div>
          </div>
          <div v-else class="placeholder">请选择您要预约的门店</div>
          <i class="arrow-right">›</i>
        </div>
      </div>

      <!-- 服务类型卡片 -->
      <div class="section-card">
        <div class="section-header">
          <i class="icon-service">💄</i>
          <span class="section-title">服务类型</span>
        </div>
        <div class="service-types">
          <div 
            class="service-type-item"
            :class="{ active: form.service_type === 'manicure' }"
            @click="form.service_type = 'manicure'"
          >
            <div class="service-icon">💅</div>
            <div class="service-name">美甲</div>
            <div class="service-desc">精致美甲服务</div>
          </div>
          <div 
            class="service-type-item"
            :class="{ active: form.service_type === 'eyelash' }"
            @click="form.service_type = 'eyelash'"
          >
            <div class="service-icon">👁️</div>
            <div class="service-name">美睫</div>
            <div class="service-desc">专业美睫服务</div>
          </div>
        </div>
      </div>

      <!-- 预约时间卡片 -->
      <div class="section-card">
        <div class="section-header">
          <i class="icon-time">📅</i>
          <span class="section-title">预约时间</span>
        </div>
        
        <!-- 日期选择 -->
        <div class="time-selector">
          <div class="time-item" @click="showDatePicker = true">
            <div class="time-label">预约日期</div>
            <div class="time-value">
              {{ form.booking_date || '选择日期' }}
            </div>
            <i class="arrow-right">›</i>
          </div>
        </div>

        <!-- 时间段选择 -->
        <div class="time-slots-container" v-if="form.booking_date">
          <div class="time-slots-header">
            <span class="slots-title">选择时段</span>
            <span class="slots-note">数字为可预约人数</span>
          </div>
          
          <!-- 加载状态 -->
          <div v-if="loadingTimeSlots" class="loading-container">
            <van-loading type="spinner" size="24px" />
            <span class="loading-text">正在加载可用时段...</span>
          </div>
          
          <!-- 错误状态 -->
          <div v-else-if="timeSlotsError" class="error-container">
            <i class="error-icon">⚠️</i>
            <span class="error-text">{{ timeSlotsError }}</span>
            <van-button size="small" type="primary" @click="fetchAvailableSlots">重试</van-button>
          </div>
          
          <!-- 无可用时段 -->
          <div v-else-if="availableTimeSlots.length === 0" class="empty-container">
            <i class="empty-icon">📅</i>
            <span class="empty-text">该日期暂无可用时段</span>
            <span class="empty-desc">请选择其他日期或联系门店</span>
          </div>
          
          <!-- 可用时段列表 -->
          <div v-else class="time-slots-grid">
            <div 
              v-for="slot in availableTimeSlots"
              :key="slot.time_slot"
              class="time-slot-item"
              :class="{ 
                selected: form.time_slot === slot.time_slot,
                disabled: slot.available_count === 0
              }"
              @click="selectTimeSlot(slot)"
            >
              <div class="slot-time">{{ slot.time_slot }}</div>
              <div class="slot-availability">
                <span class="available-count">{{ slot.available_count }}</span>
                <span class="total-count">/{{ slot.total_count }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 客户信息卡片 -->
      <div class="section-card">
        <div class="section-header">
          <i class="icon-user">👤</i>
          <span class="section-title">客户信息</span>
        </div>
        <div class="form-fields">
          <div class="form-field">
            <input 
              v-model="form.customer_name"
              placeholder="请输入您的姓名"
              class="form-input"
            />
          </div>
          <div class="form-field">
            <input 
              v-model="form.customer_phone"
              placeholder="请输入手机号"
              type="tel"
              class="form-input"
            />
          </div>
        </div>
      </div>

      <!-- 预约信息预览 -->
      <div class="preview-card" v-if="canSubmit">
        <div class="preview-header">
          <i class="icon-preview">📋</i>
          <span class="preview-title">预约信息</span>
        </div>
        <div class="preview-content">
          <div class="preview-item">
            <span class="label">门店：</span>
            <span class="value">{{ selectedStore?.name }}</span>
          </div>
          <div class="preview-item">
            <span class="label">服务：</span>
            <span class="value">{{ form.service_type === 'manicure' ? '美甲' : '美睫' }}</span>
          </div>
          <div class="preview-item">
            <span class="label">时间：</span>
            <span class="value">{{ form.booking_date }} {{ form.time_slot }}</span>
          </div>
          <div class="preview-item">
            <span class="label">客户：</span>
            <span class="value">{{ form.customer_name }} {{ form.customer_phone }}</span>
          </div>
          <div class="deposit-info">
            <span class="deposit-label">押金：</span>
            <span class="deposit-amount">¥10</span>
            <span class="deposit-desc">（服务完成后原路退回）</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 弹窗组件 -->
    <van-popup 
      v-model:show="showStorePicker" 
      position="bottom" 
      round
      :style="{ height: '60%' }"
    >
      <div class="popup-header">
        <div class="popup-title">选择门店</div>
        <div class="popup-close" @click="showStorePicker = false">×</div>
      </div>
      <van-picker
        :columns="storeOptions"
        @confirm="onStoreConfirm"
        @cancel="showStorePicker = false"
      />
    </van-popup>
    
    <van-calendar
      v-model:show="showDatePicker"
      :min-date="minDate"
      :max-date="maxDate"
      @confirm="onDateConfirm"
    />
    
    <van-popup 
      v-model:show="showTimePicker" 
      position="bottom" 
      round
      :style="{ height: '60%' }"
    >
      <div class="popup-header">
        <div class="popup-title">选择时段</div>
        <div class="popup-close" @click="showTimePicker = false">×</div>
      </div>
      <van-picker
        :columns="timeSlotOptions"
        @confirm="onTimeConfirm"
        @cancel="showTimePicker = false"
      />
    </van-popup>
    
    <!-- 底部提交按钮 -->
    <div class="submit-container">
      <van-button
        round
        block
        type="primary"
        size="large"
        :loading="loading"
        :disabled="!canSubmit"
        @click="handleSubmit"
        class="submit-button"
      >
        <template v-if="loading">
          预约中...
        </template>
        <template v-else>
          确认预约并支付 ¥10
        </template>
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { getStores } from '@/api/stores'
import { getAvailableSlots, createBooking } from '@/api/bookings'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const loadingTimeSlots = ref(false)
const timeSlotsError = ref('')
const showStorePicker = ref(false)
const showDatePicker = ref(false)
const showTimePicker = ref(false)
const selectedDate = ref(new Date())
const selectedStore = ref(null)
const storeList = ref([])
const availableTimeSlots = ref([])

const form = reactive({
  store_id: '',
  service_type: 'manicure',
  booking_date: '',
  time_slot: '',
  customer_name: '',
  customer_phone: ''
})

const minDate = new Date()
const maxDate = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000) // 30天后

const phoneRules = [
  { required: true, message: '请输入手机号' },
  { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' }
]

const storeOptions = computed(() => {
  return storeList.value.map(store => ({
    text: store.name,
    value: store.id
  }))
})

const timeSlotOptions = computed(() => {
  return availableSlots.value.map(slot => ({
    text: slot,
    value: slot
  }))
})

const canSubmit = computed(() => {
  return form.store_id && form.booking_date && form.time_slot && 
         form.customer_name && form.customer_phone
})

// 获取门店列表
const fetchStores = async () => {
  try {
    const response = await getStores()
    storeList.value = response.data.stores || []
    
    // 如果URL中有storeId参数，自动选择该门店
    const storeId = route.query.storeId
    if (storeId) {
      const store = storeList.value.find(s => s.id == storeId)
      if (store) {
        selectedStore.value = store
        form.store_id = store.id
      }
    }
  } catch (error) {
    console.error('获取门店列表失败:', error)
  }
}

// 获取可用时间段
const fetchAvailableSlots = async () => {
  if (!form.store_id || !form.booking_date || !form.service_type) {
    availableTimeSlots.value = []
    timeSlotsError.value = ''
    return
  }
  
  loadingTimeSlots.value = true
  timeSlotsError.value = ''
  
  try {
    const params = {
      store_id: form.store_id,
      date: form.booking_date,
      service_type: form.service_type
    }
    
    console.log('正在获取可用时间段，参数:', params)
    const response = await getAvailableSlots(params)
    console.log('获取可用时间段响应:', response)
    
    availableTimeSlots.value = response.data.time_slots || []
    
    if (availableTimeSlots.value.length === 0) {
      console.log('该日期暂无可用时段')
    }
  } catch (error) {
    console.error('获取可用时间段失败:', error)
    timeSlotsError.value = error.message || '获取可用时段失败，请重试'
    availableTimeSlots.value = []
    showToast.fail(timeSlotsError.value)
  } finally {
    loadingTimeSlots.value = false
  }
}

// 选择时间段
const selectTimeSlot = (slot) => {
  if (slot.available_count === 0) {
    showToast('该时段已满，请选择其他时间')
    return
  }
  form.time_slot = slot.time_slot
}

// 门店确认
const onStoreConfirm = ({ selectedOptions }) => {
  const store = storeList.value.find(s => s.id === selectedOptions[0].value)
  selectedStore.value = store
  form.store_id = store.id
  showStorePicker.value = false
  
  // 重新获取可用时间段
  fetchAvailableSlots()
}

// 日期确认
const onDateConfirm = (value) => {
  console.log('Calendar确认的值:', value)
  
  // van-calendar返回的是Date对象
  selectedDate.value = value
  form.booking_date = dayjs(value).format('YYYY-MM-DD')
  showDatePicker.value = false
  
  // 重新获取可用时间段
  fetchAvailableSlots()
}

// 时间确认
const onTimeConfirm = ({ selectedOptions }) => {
  form.time_slot = selectedOptions[0].value
  showTimePicker.value = false
}

// 提交预约
const handleSubmit = async () => {
  loading.value = true
  try {
    await createBooking(form)
    showToast.success('预约创建成功')
    router.push('/booking')
  } catch (error) {
    console.error('创建预约失败:', error)
  } finally {
    loading.value = false
  }
}

// 监听服务类型变化
watch(() => form.service_type, () => {
  fetchAvailableSlots()
})

onMounted(() => {
  fetchStores()
})
</script>

<style lang="scss" scoped>
.booking-create-page {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
  padding-bottom: 120px;
}

// 自定义导航栏
.custom-nav {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  
  :deep(.van-nav-bar__title) {
    color: #333;
    font-weight: 600;
    font-size: 18px;
  }
  
  :deep(.van-nav-bar__text) {
    color: #666;
  }
}

// 进度指示器
.progress-container {
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.progress-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  max-width: 300px;
  margin: 0 auto;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: rgba(255, 255, 255, 0.6);
  transition: all 0.3s ease;
  
  &.active {
    color: #fff;
    
    .step-icon {
      background: linear-gradient(135deg, #ff6b6b, #ffa726);
      color: #fff;
      transform: scale(1.1);
    }
  }
}

.step-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 8px;
  transition: all 0.3s ease;
}

.step-text {
  font-size: 12px;
  text-align: center;
  white-space: nowrap;
}

.step-line {
  width: 40px;
  height: 2px;
  background: rgba(255, 255, 255, 0.3);
  margin: 0 10px;
}

// 内容容器
.content-container {
  padding: 20px 16px;
}

// 卡片样式
.section-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  
  &:active {
    transform: translateY(1px);
  }
}

.section-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-left: 8px;
}

// 门店选择器
.store-selector {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: #e9ecef;
  }
}

.selected-store {
  flex: 1;
}

.store-name {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.store-address {
  font-size: 13px;
  color: #666;
}

.placeholder {
  color: #999;
  font-size: 14px;
}

.arrow-right {
  font-size: 18px;
  color: #ccc;
  font-weight: bold;
}

// 服务类型
.service-types {
  display: flex;
  gap: 12px;
}

.service-type-item {
  flex: 1;
  padding: 20px 12px;
  background: #f8f9fa;
  border-radius: 12px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  
  &.active {
    background: linear-gradient(135deg, #ff6b6b, #ffa726);
    color: #fff;
    border-color: #fff;
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(255, 107, 107, 0.3);
  }
  
  &:not(.active):hover {
    background: #e9ecef;
    transform: translateY(-1px);
  }
}

.service-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.service-name {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 4px;
}

.service-desc {
  font-size: 12px;
  opacity: 0.8;
}

// 时间选择器
.time-selector {
  margin-bottom: 12px;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.time-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: #e9ecef;
  }
}

.time-label {
  font-size: 14px;
  color: #666;
}

.time-value {
  font-size: 15px;
  color: #333;
  font-weight: 500;
}

// 表单字段
.form-fields {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-field {
  position: relative;
}

.form-input {
  width: 100%;
  padding: 16px;
  border: 2px solid #e9ecef;
  border-radius: 12px;
  font-size: 15px;
  background: #fff;
  transition: border-color 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: #ff6b6b;
    box-shadow: 0 0 0 3px rgba(255, 107, 107, 0.1);
  }
  
  &::placeholder {
    color: #999;
  }
}

// 预约信息预览
.preview-card {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  border-radius: 16px;
  padding: 20px;
  margin-top: 8px;
}

.preview-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.preview-title {
  font-size: 16px;
  font-weight: 600;
  margin-left: 8px;
}

.preview-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  
  &:last-child {
    border-bottom: none;
  }
}

.label {
  font-size: 14px;
  opacity: 0.9;
}

.value {
  font-size: 14px;
  font-weight: 500;
}

.deposit-info {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  margin-top: 8px;
}

.deposit-label {
  font-size: 14px;
  margin-right: 8px;
}

.deposit-amount {
  font-size: 18px;
  font-weight: bold;
  color: #ffa726;
  margin-right: 8px;
}

.deposit-desc {
  font-size: 12px;
  opacity: 0.8;
}

// 底部提交按钮
.submit-container {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.submit-button {
  background: linear-gradient(135deg, #ff6b6b, #ffa726);
  border: none;
  box-shadow: 0 8px 24px rgba(255, 107, 107, 0.3);
  
  &:disabled {
    background: #ccc;
    box-shadow: none;
  }
  
  :deep(.van-button__text) {
    font-weight: 600;
    font-size: 16px;
  }
}

// 弹窗样式
.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
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
  .progress-steps {
    max-width: 250px;
  }
  
  .step-text {
    font-size: 11px;
  }
  
  .service-types {
    flex-direction: column;
    gap: 8px;
  }
  
  .service-type-item {
    padding: 16px 12px;
  }
}

// 动画效果
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.section-card {
  animation: slideIn 0.3s ease-out;
}

.section-card:nth-child(1) { animation-delay: 0.1s; }
.section-card:nth-child(2) { animation-delay: 0.2s; }
.section-card:nth-child(3) { animation-delay: 0.3s; }
.section-card:nth-child(4) { animation-delay: 0.4s; }

// 时间段容器
.time-slots-container {
  margin-top: 16px;
}

.time-slots-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.slots-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.slots-note {
  font-size: 12px;
  color: #999;
}

.time-slots-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.time-slot-item {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 12px 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  
  &.selected {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    border-color: #667eea;
    transform: scale(1.05);
  }
  
  &.disabled {
    background: #f0f0f0;
    color: #ccc;
    cursor: not-allowed;
    
    .slot-availability {
      color: #ccc;
    }
  }
  
  &:hover:not(.disabled):not(.selected) {
    background: #e9ecef;
    transform: translateY(-2px);
  }
}

.slot-time {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 4px;
}

.slot-availability {
  font-size: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
}

.available-count {
  font-weight: 700;
  color: #52c41a;
  
  .time-slot-item.selected & {
    color: rgba(255, 255, 255, 0.9);
  }
  
  .time-slot-item.disabled & {
    color: #ff4d4f;
  }
}

.total-count {
  opacity: 0.7;
}

// 加载状态
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  gap: 12px;
}

.loading-text {
  font-size: 14px;
  color: #666;
}

// 错误状态
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  gap: 12px;
  background: #fff2f0;
  border: 1px solid #ffccc7;
  border-radius: 12px;
}

.error-icon {
  font-size: 24px;
}

.error-text {
  font-size: 14px;
  color: #ff4d4f;
  text-align: center;
}

// 空状态
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  gap: 8px;
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 12px;
}

.empty-icon {
  font-size: 24px;
}

.empty-text {
  font-size: 14px;
  color: #52c41a;
  font-weight: 500;
}

.empty-desc {
  font-size: 12px;
  color: #8c8c8c;
  text-align: center;
}
</style>
