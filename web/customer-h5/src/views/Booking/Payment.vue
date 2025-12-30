<template>
  <div class="payment-page">
    <!-- 页面头部 -->
    <van-nav-bar
      title="支付押金"
      left-text="返回"
      left-arrow
      @click-left="$router.back()"
      class="custom-nav"
    />

    <!-- 预约信息确认 -->
    <div class="booking-info-section">
      <div class="info-card">
        <div class="info-header">
          <i class="info-icon">📋</i>
          <span class="info-title">预约信息确认</span>
        </div>
        <div class="info-content">
          <div class="info-item">
            <span class="info-label">门店</span>
            <span class="info-value">{{ bookingInfo.store_name }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">服务</span>
            <span class="info-value">{{ getServiceText(bookingInfo.service_type) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">时间</span>
            <span class="info-value">{{ formatDateTime() }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">客户</span>
            <span class="info-value">{{ bookingInfo.customer_name }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 支付信息 -->
    <div class="payment-info-section">
      <div class="payment-card">
        <div class="payment-header">
          <i class="payment-icon">💳</i>
          <span class="payment-title">支付信息</span>
        </div>
        <div class="payment-details">
          <div class="payment-item">
            <span class="payment-label">押金金额</span>
            <span class="payment-amount">¥10.00</span>
          </div>
          <div class="payment-note">
            <i class="note-icon">💡</i>
            <span>服务完成后押金将原路退回</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 支付方式 -->
    <div class="payment-method-section">
      <div class="method-card">
        <div class="method-header">
          <i class="method-icon">💰</i>
          <span class="method-title">支付方式</span>
        </div>
        <div class="payment-methods">
          <div class="method-item active">
            <div class="method-logo">
              <i class="wechat-icon">💚</i>
            </div>
            <div class="method-info">
              <div class="method-name">微信支付</div>
              <div class="method-desc">安全便捷的支付方式</div>
            </div>
            <div class="method-check">
              <i class="check-icon">✓</i>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 重要提示 -->
    <div class="tips-section">
      <div class="tips-card">
        <div class="tips-header">
          <i class="tips-icon">⚠️</i>
          <span class="tips-title">重要提示</span>
        </div>
        <div class="tips-content">
          <div class="tip-item">
            <i class="tip-icon">⏰</i>
            <span>请在15分钟内完成支付，否则订单将自动取消</span>
          </div>
          <div class="tip-item">
            <i class="tip-icon">📅</i>
            <span>可在预约时间前3小时免费取消</span>
          </div>
          <div class="tip-item">
            <i class="tip-icon">💰</i>
            <span>服务完成后押金将在3-5个工作日内退回</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部支付按钮 -->
    <div class="payment-footer">
      <div class="payment-summary">
        <div class="summary-text">
          <span class="summary-label">应付金额</span>
          <span class="summary-amount">¥10.00</span>
        </div>
      </div>
      <button 
        class="pay-button"
        @click="handlePayment"
        :disabled="paying"
      >
        <i class="pay-icon">{{ paying ? '⏳' : '💳' }}</i>
        <span>{{ paying ? '支付中...' : '立即支付' }}</span>
      </button>
    </div>

    <!-- 支付结果弹窗 -->
    <van-popup 
      v-model:show="showPaymentResult" 
      position="center"
      round
      :style="{ width: '90%' }"
      :close-on-click-overlay="false"
    >
      <div class="payment-result-dialog">
        <div class="result-icon" :class="paymentResult.success ? 'success' : 'failed'">
          {{ paymentResult.success ? '🎉' : '😞' }}
        </div>
        <h3 class="result-title">
          {{ paymentResult.success ? '支付成功' : '支付失败' }}
        </h3>
        <div class="result-content" v-if="paymentResult.success">
          <p class="result-text">预约已确认，请按时到店享受服务</p>
          <div class="booking-number">
            订单号：{{ bookingInfo.order_no }}
          </div>
        </div>
        <div class="result-content" v-else>
          <p class="result-text">{{ paymentResult.message }}</p>
        </div>
        <div class="result-actions">
          <button 
            v-if="paymentResult.success"
            class="result-btn success-btn"
            @click="goToOrders"
          >
            查看订单
          </button>
          <button 
            v-else
            class="result-btn retry-btn"
            @click="retryPayment"
          >
            重新支付
          </button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { createPayment } from '@/api/payments'
import { getBookingById } from '@/api/bookings'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()

const paying = ref(false)
const showPaymentResult = ref(false)

const bookingInfo = reactive({
  id: '',
  order_no: '',
  store_name: '',
  service_type: '',
  booking_date: '',
  time_slot: '',
  customer_name: ''
})

const paymentResult = reactive({
  success: false,
  message: ''
})

// 获取预约信息
const fetchBookingInfo = async () => {
  try {
    const bookingId = route.params.id
    const response = await getBookingById(bookingId)
    Object.assign(bookingInfo, response.data)
  } catch (error) {
    console.error('获取预约信息失败:', error)
    showToast('获取预约信息失败')
    router.back()
  }
}

// 处理支付
const handlePayment = async () => {
  paying.value = true
  
  try {
    const response = await createPayment(bookingInfo.id, {
      amount: 1000 // 10元 = 1000分
    })
    
    if (response.data.payment_url) {
      // 跳转到微信支付
      window.location.href = response.data.payment_url
    } else {
      // 模拟支付成功
      setTimeout(() => {
        paymentResult.success = true
        paymentResult.message = ''
        showPaymentResult.value = true
        paying.value = false
      }, 2000)
    }
  } catch (error) {
    console.error('支付失败:', error)
    paymentResult.success = false
    paymentResult.message = error.message || '支付失败，请重试'
    showPaymentResult.value = true
    paying.value = false
  }
}

// 重新支付
const retryPayment = () => {
  showPaymentResult.value = false
  handlePayment()
}

// 查看订单
const goToOrders = () => {
  router.push('/booking')
}

// 获取服务类型文本
const getServiceText = (type) => {
  return type === 'manicure' ? '美甲服务' : '美睫服务'
}

// 格式化日期时间
const formatDateTime = () => {
  return `${dayjs(bookingInfo.booking_date).format('MM月DD日')} ${bookingInfo.time_slot}`
}

onMounted(() => {
  fetchBookingInfo()
})
</script>

<style lang="scss" scoped>
.payment-page {
  background: linear-gradient(180deg, #f8f9ff 0%, #f0f2ff 100%);
  min-height: 100vh;
  padding-bottom: 100px;
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
}

// 预约信息
.booking-info-section {
  padding: 20px 16px;
}

.info-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 16px;
}

.info-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.info-icon {
  font-size: 18px;
  margin-right: 8px;
}

.info-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f8f8f8;
  
  &:last-child {
    border-bottom: none;
  }
}

.info-label {
  font-size: 14px;
  color: #666;
}

.info-value {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

// 支付信息
.payment-info-section {
  padding: 0 16px 16px;
}

.payment-card {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
}

.payment-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.payment-icon {
  font-size: 18px;
  margin-right: 8px;
}

.payment-title {
  font-size: 16px;
  font-weight: 600;
}

.payment-details {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.payment-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.payment-label {
  font-size: 14px;
  opacity: 0.9;
}

.payment-amount {
  font-size: 24px;
  font-weight: 700;
  color: #ffa726;
}

.payment-note {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  font-size: 14px;
  opacity: 0.9;
}

.note-icon {
  font-size: 16px;
}

// 支付方式
.payment-method-section {
  padding: 0 16px 16px;
}

.method-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.method-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.method-icon {
  font-size: 18px;
  margin-right: 8px;
}

.method-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.payment-methods {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.method-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 2px solid transparent;
  
  &.active {
    background: #e6f7ff;
    border-color: #1890ff;
  }
}

.method-logo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  background: #07c160;
}

.wechat-icon {
  color: white;
}

.method-info {
  flex: 1;
}

.method-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 2px;
}

.method-desc {
  font-size: 13px;
  color: #666;
}

.method-check {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #1890ff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

// 提示信息
.tips-section {
  padding: 0 16px 16px;
}

.tips-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.tips-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.tips-icon {
  font-size: 18px;
  margin-right: 8px;
  color: #fa8c16;
}

.tips-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.tips-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tip-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 14px;
  color: #666;
  line-height: 1.4;
}

.tip-icon {
  font-size: 14px;
  margin-top: 2px;
  color: #fa8c16;
}

// 底部支付区域
.payment-footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.payment-summary {
  flex: 1;
}

.summary-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.summary-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 2px;
}

.summary-amount {
  font-size: 20px;
  font-weight: 700;
  color: #ff6b6b;
}

.pay-button {
  padding: 14px 32px;
  background: linear-gradient(135deg, #07c160, #38f9d7);
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 8px 24px rgba(7, 193, 96, 0.3);
  
  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 12px 32px rgba(7, 193, 96, 0.4);
  }
  
  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
}

.pay-icon {
  font-size: 18px;
}

// 支付结果弹窗
.payment-result-dialog {
  background: white;
  border-radius: 16px;
  padding: 40px 24px;
  text-align: center;
  max-width: 400px;
  width: 100%;
}

.result-icon {
  font-size: 64px;
  margin-bottom: 16px;
  
  &.success {
    animation: bounce 1s ease-out;
  }
  
  &.failed {
    animation: shake 0.5s ease-out;
  }
}

.result-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 16px;
  
  .payment-result-dialog .result-icon.success + & {
    color: #52c41a;
  }
  
  .payment-result-dialog .result-icon.failed + & {
    color: #ff4d4f;
  }
}

.result-content {
  margin-bottom: 24px;
}

.result-text {
  font-size: 16px;
  color: #666;
  margin: 0 0 12px 0;
  line-height: 1.4;
}

.booking-number {
  font-size: 14px;
  color: #999;
  background: #f8f9fa;
  padding: 8px 12px;
  border-radius: 8px;
}

.result-actions {
  display: flex;
  gap: 12px;
}

.result-btn {
  flex: 1;
  height: 44px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &.success-btn {
    background: linear-gradient(135deg, #52c41a, #73d13d);
    color: white;
  }
  
  &.retry-btn {
    background: linear-gradient(135deg, #ff6b6b, #ffa726);
    color: white;
  }
}

// 动画效果
@keyframes bounce {
  0%, 20%, 53%, 80%, 100% {
    transform: translate3d(0, 0, 0);
  }
  40%, 43% {
    transform: translate3d(0, -10px, 0);
  }
  70% {
    transform: translate3d(0, -5px, 0);
  }
  90% {
    transform: translate3d(0, -2px, 0);
  }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
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

.booking-info-section,
.payment-info-section,
.payment-method-section,
.tips-section {
  animation: slideInUp 0.6s ease-out;
}

.booking-info-section { animation-delay: 0.1s; }
.payment-info-section { animation-delay: 0.2s; }
.payment-method-section { animation-delay: 0.3s; }
.tips-section { animation-delay: 0.4s; }

// 响应式设计
@media (max-width: 375px) {
  .payment-footer {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .pay-button {
    width: 100%;
    justify-content: center;
  }
  
  .result-actions {
    flex-direction: column;
  }
}
</style>
