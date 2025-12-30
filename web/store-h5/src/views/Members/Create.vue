<template>
  <div class="member-create-page">
    <!-- 页面头部 -->
    <van-nav-bar
      title="创建会员"
      left-text="返回"
      left-arrow
      @click-left="$router.back()"
      class="custom-nav"
    />

    <!-- 创建步骤 -->
    <div class="steps-section">
      <div class="steps-container">
        <div class="step" :class="{ active: currentStep >= 1, completed: currentStep > 1 }">
          <div class="step-icon">1</div>
          <div class="step-text">基本信息</div>
        </div>
        <div class="step-line"></div>
        <div class="step" :class="{ active: currentStep >= 2, completed: currentStep > 2 }">
          <div class="step-icon">2</div>
          <div class="step-text">人脸录入</div>
        </div>
        <div class="step-line"></div>
        <div class="step" :class="{ active: currentStep >= 3 }">
          <div class="step-icon">3</div>
          <div class="step-text">套餐选择</div>
        </div>
      </div>
    </div>

    <div class="content-container">
      <!-- 步骤1：基本信息 -->
      <div v-if="currentStep === 1" class="step-content">
        <div class="section-card">
          <div class="section-header">
            <i class="section-icon">👤</i>
            <span class="section-title">会员基本信息</span>
          </div>
          <div class="form-fields">
            <div class="form-group">
              <label class="form-label">会员姓名</label>
              <div class="input-container">
                <i class="input-icon">👨‍💼</i>
                <input
                  v-model="memberForm.name"
                  type="text"
          placeholder="请输入会员姓名"
                  class="form-input"
                  required
                />
              </div>
            </div>
            <div class="form-group">
              <label class="form-label">手机号码</label>
              <div class="input-container">
                <i class="input-icon">📱</i>
                <input
                  v-model="memberForm.phone"
                  type="tel"
          placeholder="请输入手机号码"
                  class="form-input"
                  required
                />
              </div>
            </div>
            <div class="form-group">
              <label class="form-label">有效期设置</label>
              <div class="validity-container">
                <div class="validity-item" @click="showStartDatePicker = true">
                  <div class="validity-label">开始日期</div>
                  <div class="validity-value">
                    {{ memberForm.valid_from || '选择开始日期' }}
                  </div>
                  <i class="arrow-icon">›</i>
                </div>
                <div class="validity-item" @click="showEndDatePicker = true">
                  <div class="validity-label">结束日期</div>
                  <div class="validity-value">
                    {{ memberForm.valid_to || '选择结束日期' }}
                  </div>
                  <i class="arrow-icon">›</i>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 步骤2：人脸录入 -->
      <div v-if="currentStep === 2" class="step-content">
        <div class="section-card">
          <div class="section-header">
            <i class="section-icon">📸</i>
            <span class="section-title">人脸信息录入</span>
          </div>
          <div class="face-capture-container">
            <div class="capture-tips">
              <div class="tips-icon">💡</div>
              <div class="tips-text">
                <p>请让会员正对手机摄像头</p>
                <p>确保光线充足，面部清晰可见</p>
              </div>
            </div>
            
            <div class="camera-area">
              <video
                ref="videoRef"
                class="camera-video"
                autoplay
                muted
                playsinline
                v-show="!capturedImage"
              ></video>
              
              <div v-if="capturedImage" class="captured-preview">
                <img :src="capturedImage" alt="已拍摄的人脸照片" />
                <div class="preview-overlay">
                  <button class="overlay-btn retake-btn" @click="retakePhoto">
                    重新拍摄
                  </button>
                </div>
              </div>
              
              <canvas ref="canvasRef" style="display: none;"></canvas>
            </div>
            
            <div class="capture-actions">
              <button 
                v-if="!capturedImage"
                class="capture-btn"
                @click="capturePhoto"
                :disabled="isCapturing"
              >
                <i class="capture-icon">📷</i>
                {{ isCapturing ? '拍摄中...' : '拍摄照片' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 步骤3：套餐选择 -->
      <div v-if="currentStep === 3" class="step-content">
        <div class="section-card">
          <div class="section-header">
            <i class="section-icon">💎</i>
            <span class="section-title">选择套餐</span>
          </div>
          <div class="packages-grid">
            <div 
              v-for="pkg in packageOptions"
              :key="pkg.value"
              class="package-card"
              :class="{ active: memberForm.package_name === pkg.value }"
              @click="selectPackage(pkg)"
            >
              <div class="package-header">
                <div class="package-icon">{{ pkg.icon }}</div>
                <div class="package-badge" v-if="pkg.popular">推荐</div>
              </div>
              <h3 class="package-title">{{ pkg.title }}</h3>
              <div class="package-price">
                <span class="price-symbol">¥</span>
                <span class="price-amount">{{ pkg.price }}</span>
              </div>
              <div class="package-features">
                <div v-for="feature in pkg.features" :key="feature" class="feature-item">
                  <i class="feature-check">✓</i>
                  <span>{{ feature }}</span>
                </div>
              </div>
            </div>
          </div>
          
          <div class="custom-package" v-if="memberForm.package_name === 'custom'">
            <div class="form-group">
              <label class="form-label">自定义套餐名称</label>
              <div class="input-container">
                <input
                  v-model="memberForm.custom_package_name"
                  type="text"
                  placeholder="请输入套餐名称"
                  class="form-input"
                />
              </div>
            </div>
            <div class="form-group">
              <label class="form-label">套餐金额（元）</label>
              <div class="input-container">
                <input
                  v-model="memberForm.package_amount"
                  type="number"
                  placeholder="请输入套餐金额"
                  class="form-input"
                />
              </div>
            </div>
            <div class="form-group">
              <label class="form-label">服务次数</label>
              <div class="input-container">
                <input
                  v-model="memberForm.total_times"
                  type="number"
                  placeholder="请输入服务次数"
                  class="form-input"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部操作按钮 -->
    <div class="bottom-actions">
      <button 
        v-if="currentStep > 1"
        class="step-btn prev-btn"
        @click="prevStep"
      >
        上一步
      </button>
      <button 
        v-if="currentStep < 3"
        class="step-btn next-btn"
        @click="nextStep"
        :disabled="!canNextStep"
      >
        下一步
      </button>
      <button 
        v-if="currentStep === 3"
        class="step-btn submit-btn"
        @click="submitMember"
        :disabled="!canSubmit"
        :loading="submitting"
      >
        {{ submitting ? '创建中...' : '创建会员' }}
      </button>
    </div>

    <!-- 日期选择器 -->
    <van-popup 
      v-model:show="showStartDatePicker" 
      position="bottom" 
          round
        >
      <div class="popup-header">
        <div class="popup-title">选择开始日期</div>
        <div class="popup-close" @click="showStartDatePicker = false">×</div>
      </div>
      <van-date-picker
        v-model="startDate"
        @confirm="onStartDateConfirm"
        @cancel="showStartDatePicker = false"
      />
    </van-popup>

    <van-popup 
      v-model:show="showEndDatePicker" 
      position="bottom" 
      round
    >
      <div class="popup-header">
        <div class="popup-title">选择结束日期</div>
        <div class="popup-close" @click="showEndDatePicker = false">×</div>
      </div>
      <van-date-picker
        v-model="endDate"
        @confirm="onEndDateConfirm"
        @cancel="showEndDatePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { createMember } from '@/api/members'
import dayjs from 'dayjs'

const router = useRouter()

const currentStep = ref(1)
const submitting = ref(false)
const isCapturing = ref(false)
const capturedImage = ref('')
const videoRef = ref(null)
const canvasRef = ref(null)
const showStartDatePicker = ref(false)
const showEndDatePicker = ref(false)
const startDate = ref(new Date())
const endDate = ref(new Date(Date.now() + 365 * 24 * 60 * 60 * 1000))

let stream = null

const memberForm = reactive({
  name: '',
  phone: '',
  valid_from: '',
  valid_to: '',
  package_name: '',
  custom_package_name: '',
  package_amount: '',
  total_times: '',
  face_image: ''
})

const packageOptions = [
  {
    value: 'basic',
    title: '基础套餐',
    price: '299',
    icon: '🥉',
    features: ['美甲服务 10次', '基础护理', '会员折扣'],
    popular: false
  },
  {
    value: 'premium',
    title: '高级套餐',
    price: '599',
    icon: '🥈',
    features: ['美甲服务 20次', '美睫服务 5次', '专属技师', '优先预约'],
    popular: true
  },
  {
    value: 'vip',
    title: 'VIP套餐',
    price: '999',
    icon: '🥇',
    features: ['不限次数服务', '私人定制', '专属技师', '免费上门'],
    popular: false
  },
  {
    value: 'custom',
    title: '自定义套餐',
    price: '自定义',
    icon: '⚙️',
    features: ['灵活配置', '个性化服务', '自定义价格'],
    popular: false
  }
]

// 检查是否可以进入下一步
const canNextStep = computed(() => {
  if (currentStep.value === 1) {
    return memberForm.name && memberForm.phone && memberForm.valid_from && memberForm.valid_to
  }
  if (currentStep.value === 2) {
    return capturedImage.value
  }
  return false
})

// 检查是否可以提交
const canSubmit = computed(() => {
  return memberForm.package_name && (
    memberForm.package_name !== 'custom' || 
    (memberForm.custom_package_name && memberForm.package_amount && memberForm.total_times)
  )
})

// 下一步
const nextStep = () => {
  if (currentStep.value < 3) {
    currentStep.value++
    if (currentStep.value === 2) {
      startCamera()
    }
  }
}

// 上一步
const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
    if (currentStep.value === 1) {
      stopCamera()
    }
  }
}

// 启动摄像头
const startCamera = async () => {
  try {
    stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user' }
    })
    if (videoRef.value) {
      videoRef.value.srcObject = stream
    }
  } catch (error) {
    console.error('启动摄像头失败:', error)
    showToast('启动摄像头失败，请检查权限设置')
  }
}

// 拍摄照片
const capturePhoto = async () => {
  if (!videoRef.value || !canvasRef.value) return
  
  isCapturing.value = true
  
  try {
    const canvas = canvasRef.value
    const video = videoRef.value
    const context = canvas.getContext('2d')
    
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    context.drawImage(video, 0, 0, canvas.width, canvas.height)
    
    const imageData = canvas.toDataURL('image/jpeg', 0.8)
    capturedImage.value = imageData
    memberForm.face_image = imageData
    
    stopCamera()
    showToast('人脸照片拍摄成功')
    
  } catch (error) {
    console.error('拍照失败:', error)
    showToast('拍照失败，请重试')
  } finally {
    isCapturing.value = false
  }
}

// 重新拍照
const retakePhoto = () => {
  capturedImage.value = ''
  memberForm.face_image = ''
  startCamera()
}

// 停止摄像头
const stopCamera = () => {
  if (stream) {
    stream.getTracks().forEach(track => track.stop())
    stream = null
  }
}

// 选择套餐
const selectPackage = (pkg) => {
  memberForm.package_name = pkg.value
  if (pkg.value !== 'custom') {
    memberForm.package_amount = pkg.price
    memberForm.total_times = pkg.value === 'basic' ? '10' : pkg.value === 'premium' ? '25' : '99'
  }
}

// 日期确认
const onStartDateConfirm = () => {
  memberForm.valid_from = dayjs(startDate.value).format('YYYY-MM-DD')
  showStartDatePicker.value = false
}

const onEndDateConfirm = () => {
  memberForm.valid_to = dayjs(endDate.value).format('YYYY-MM-DD')
  showEndDatePicker.value = false
}

// 提交创建会员
const submitMember = async () => {
  submitting.value = true
  try {
    const data = {
      ...memberForm,
      package_name: memberForm.package_name === 'custom' ? memberForm.custom_package_name : memberForm.package_name,
      package_amount: parseInt(memberForm.package_amount) * 100, // 转换为分
      total_times: parseInt(memberForm.total_times),
      remaining_times: parseInt(memberForm.total_times)
    }

    await createMember(data)
    showToast('会员创建成功')
    router.push('/members')
  } catch (error) {
    console.error('创建会员失败:', error)
    showToast('创建失败，请重试')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  // 设置默认日期
  memberForm.valid_from = dayjs().format('YYYY-MM-DD')
  memberForm.valid_to = dayjs().add(1, 'year').format('YYYY-MM-DD')
})

onUnmounted(() => {
  stopCamera()
})
</script>

<style lang="scss" scoped>
.member-create-page {
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

// 步骤指示器
.steps-section {
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.steps-container {
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
  color: rgba(102, 126, 234, 0.4);
  transition: all 0.3s ease;
  
  &.active {
    color: #667eea;
    
    .step-icon {
      background: linear-gradient(135deg, #667eea, #764ba2);
      color: white;
      transform: scale(1.1);
    }
  }
  
  &.completed {
    color: #52c41a;
    
    .step-icon {
      background: #52c41a;
      color: white;
      
      &::after {
        content: '✓';
        font-size: 12px;
      }
    }
  }
}

.step-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #f0f0f0;
  color: #999;
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
  background: #f0f0f0;
  margin: 0 10px;
}

// 内容容器
.content-container {
  padding: 20px 16px;
}

.step-content {
  animation: slideInUp 0.4s ease-out;
}

// 卡片样式
.section-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.section-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.section-icon {
  font-size: 20px;
  margin-right: 8px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

// 表单样式
.form-fields {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.input-container {
  display: flex;
  align-items: center;
  background: #f8f9fa;
  border-radius: 12px;
  padding: 0 16px;
  border: 2px solid transparent;
  transition: all 0.3s ease;
  
  &:focus-within {
    background: white;
    border-color: #667eea;
    box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
  }
}

.input-icon {
  font-size: 16px;
  margin-right: 12px;
  color: #999;
}

.form-input {
  flex: 1;
  border: none;
  outline: none;
  padding: 16px 0;
  font-size: 16px;
  background: transparent;
  color: #333;
  
  &::placeholder {
    color: #999;
  }
}

// 有效期设置
.validity-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.validity-item {
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

.validity-label {
  font-size: 14px;
  color: #666;
}

.validity-value {
  font-size: 15px;
  color: #333;
  font-weight: 500;
}

.arrow-icon {
  font-size: 16px;
  color: #ccc;
}

// 人脸拍摄
.face-capture-container {
  text-align: center;
}

.capture-tips {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #fff7e6;
  border-radius: 12px;
  margin-bottom: 20px;
}

.tips-icon {
  font-size: 24px;
}

.tips-text {
  flex: 1;
  text-align: left;
  
  p {
    margin: 0;
    font-size: 14px;
    color: #fa8c16;
    line-height: 1.4;
  }
}

.camera-area {
  width: 280px;
  height: 280px;
  margin: 0 auto 20px;
  border-radius: 20px;
  overflow: hidden;
  background: #000;
  position: relative;
  border: 4px solid rgba(102, 126, 234, 0.2);
}

.camera-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transform: scaleX(-1);
}

.captured-preview {
  position: relative;
  width: 100%;
  height: 100%;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.preview-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  
  .captured-preview:hover & {
    opacity: 1;
  }
}

.overlay-btn {
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 20px;
  color: #333;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.capture-actions {
  display: flex;
  justify-content: center;
}

.capture-btn {
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  
  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
  }
  
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.capture-icon {
  font-size: 16px;
}

// 套餐选择
.packages-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.package-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  position: relative;
  
  &.active {
    border-color: #667eea;
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(102, 126, 234, 0.2);
  }
  
  &:hover:not(.active) {
    transform: translateY(-1px);
    box-shadow: 0 6px 24px rgba(0, 0, 0, 0.1);
  }
}

.package-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.package-icon {
  font-size: 32px;
}

.package-badge {
  background: linear-gradient(135deg, #ff6b6b, #ffa726);
  color: white;
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 10px;
  font-weight: 600;
}

.package-title {
  font-size: 18px;
  font-weight: 700;
  color: #333;
  margin-bottom: 8px;
}

.package-price {
  display: flex;
  align-items: baseline;
  margin-bottom: 16px;
}

.price-symbol {
  font-size: 14px;
  color: #ff6b6b;
  font-weight: 600;
}

.price-amount {
  font-size: 24px;
  font-weight: 700;
  color: #ff6b6b;
  margin-left: 4px;
}

.package-features {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.feature-item {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #666;
}

.feature-check {
  color: #52c41a;
  font-weight: bold;
  margin-right: 8px;
}

// 自定义套餐
.custom-package {
  margin-top: 20px;
  padding: 20px;
  background: #f8f9ff;
  border-radius: 12px;
  border: 2px dashed #667eea;
}

// 底部操作
.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
  display: flex;
  gap: 12px;
}

.step-btn {
  flex: 1;
  height: 48px;
  border: none;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &.prev-btn {
    background: #f5f5f5;
    color: #666;
    
    &:hover {
      background: #e9ecef;
    }
  }
  
  &.next-btn {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    
    &:hover:not(:disabled) {
      transform: translateY(-2px);
      box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
    }
  }
  
  &.submit-btn {
    background: linear-gradient(135deg, #52c41a, #73d13d);
    color: white;
    
    &:hover:not(:disabled) {
      transform: translateY(-2px);
      box-shadow: 0 8px 24px rgba(82, 196, 26, 0.3);
    }
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none !important;
  }
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
  .camera-area {
    width: 240px;
    height: 240px;
  }
  
  .packages-grid {
    gap: 8px;
  }
  
  .package-card {
    padding: 16px;
  }
  
  .bottom-actions {
    flex-direction: column;
  }
  
  .step-btn {
    width: 100%;
  }
}

// 动画效果
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
</style>