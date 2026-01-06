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
          <div class="step-text">完成创建</div>
        </div>
      </div>
    </div>

    <div class="content-container">
      <!-- 步骤1：基本信息 -->
      <div v-if="currentStep === 1" class="step-content">
        <van-form @submit="handleNextStep">
          <!-- 会员信息 -->
          <div class="form-section">
            <div class="section-title">会员信息</div>
            <van-field
              v-model="memberForm.name"
              label="会员姓名"
              placeholder="请输入会员姓名"
              required
              :rules="[{ required: true, message: '请输入会员姓名' }]"
            />
            <van-field
              v-model="memberForm.phone"
              label="手机号"
              type="tel"
              placeholder="请输入手机号"
              required
              :rules="[
                { required: true, message: '请输入手机号' },
                { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' }
              ]"
            />
            <van-field
              v-model="memberForm.store_name"
              label="所属门店"
              placeholder="请选择门店"
              readonly
              required
              is-link
              @click="showStorePicker = true"
              :rules="[{ required: true, message: '请选择门店' }]"
            />
            <van-field
              v-model="memberForm.service_type"
              label="服务类型"
              placeholder="请选择服务类型"
              readonly
              required
              is-link
              @click="showServiceTypePicker = true"
              :rules="[{ required: true, message: '请选择服务类型' }]"
            />
            <van-field
              v-model="memberForm.package_name"
              label="套餐名称"
              placeholder="请输入套餐名称"
              required
              :rules="[{ required: true, message: '请输入套餐名称' }]"
            />
            <van-field
              v-model="memberForm.package_price"
              label="套餐价格"
              type="number"
              placeholder="0.00"
              required
              :rules="[{ required: true, message: '请输入套餐价格' }]"
            >
              <template #button>
                <div class="number-stepper">
                  <button type="button" class="stepper-btn" @click.stop="decreasePrice">-</button>
                  <button type="button" class="stepper-btn" @click.stop="increasePrice">+</button>
                </div>
              </template>
            </van-field>
            <van-field
              v-model="memberForm.purchase_amount"
              label="购买金额"
              type="number"
              placeholder="0.00"
            >
              <template #button>
                <div class="number-stepper">
                  <button type="button" class="stepper-btn" @click.stop="decreasePurchase">-</button>
                  <button type="button" class="stepper-btn" @click.stop="increasePurchase">+</button>
                </div>
              </template>
            </van-field>
            <van-field
              :model-value="statusDisplayText"
              label="会员状态"
              placeholder="有效"
              readonly
              is-link
              @click="showStatusPicker = true"
            />
          </div>

          <!-- 有效期设置 -->
          <div class="form-section">
            <div class="section-title">有效期设置</div>
            <van-field
              v-model="memberForm.valid_from"
              label="有效期开始"
              placeholder="选择开始日期"
              readonly
              required
              is-link
              @click="showStartDatePicker = true"
              :rules="[{ required: true, message: '请选择有效期开始日期' }]"
            />
            <van-field
              v-model="memberForm.valid_to"
              label="有效期结束"
              placeholder="选择结束日期"
              readonly
              required
              is-link
              @click="showEndDatePicker = true"
              :rules="[{ required: true, message: '请选择有效期结束日期' }]"
            />
            <van-field
              v-model="memberForm.fixed_duration"
              label="固定时长 (天)"
              type="number"
              placeholder="请输入固定时长天数"
            >
              <template #button>
                <div class="number-stepper">
                  <button type="button" class="stepper-btn" @click.stop="decreaseDuration">-</button>
                  <button type="button" class="stepper-btn" @click.stop="increaseDuration">+</button>
                </div>
              </template>
            </van-field>
            <div class="form-hint">
              输入固定时长后,系统会自动计算结束日期;选择开始/结束日期后,系统会自动计算固定时长
            </div>
          </div>

          <!-- 备注 -->
          <div class="form-section">
            <van-field
              v-model="memberForm.remarks"
              label="备注"
              type="textarea"
              placeholder="请输入备注"
              maxlength="500"
              show-word-limit
              rows="3"
            />
          </div>

          <!-- 下一步按钮 -->
          <div class="form-actions">
            <van-button
              round
              type="primary"
              native-type="submit"
              block
              :disabled="!canNextStep"
              class="next-btn"
            >
              下一步
            </van-button>
          </div>
        </van-form>
      </div>

      <!-- 步骤2：人脸录入 -->
      <div v-if="currentStep === 2" class="step-content">
        <div class="face-capture-section">
          <div class="section-header">
            <i class="section-icon">📸</i>
            <span class="section-title">人脸信息录入</span>
          </div>
          
          <!-- 摄像头预览区域 -->
          <div v-if="!faceImageCaptured && !cameraError" class="camera-container">
            <video
              ref="videoRef"
              class="camera-video"
              autoplay
              playsinline
              muted
            ></video>
            <canvas
              ref="canvasRef"
              class="capture-canvas"
              style="display: none;"
            ></canvas>
            <div class="camera-hint">请正对摄像头，确保光线充足</div>
          </div>

          <!-- 拍照预览区域 -->
          <div v-if="faceImageCaptured" class="face-preview-container">
            <img :src="faceImagePreview" alt="人脸预览" class="face-preview-image" />
            <div class="preview-hint">✓ 已录入人脸</div>
          </div>

          <!-- 错误提示 -->
          <div v-if="cameraError" class="camera-error">
            <div class="error-icon">⚠️</div>
            <div class="error-text">{{ cameraError }}</div>
            <van-button type="primary" size="small" @click="retryCamera">重试</van-button>
          </div>

          <!-- 操作按钮 -->
          <div class="face-actions">
            <van-button
              v-if="!faceImageCaptured && !cameraError"
              type="primary"
              size="large"
              block
              :loading="isCapturing"
              @click="capturePhoto"
            >
              拍照
            </van-button>
            <van-button
              v-if="faceImageCaptured"
              type="default"
              size="large"
              block
              @click="retakePhoto"
            >
              重新拍照
            </van-button>
            <van-button
              v-if="faceImageCaptured"
              type="primary"
              size="large"
              block
              @click="goToNextStep"
              class="next-step-btn"
            >
              下一步
            </van-button>
          </div>
        </div>
      </div>

      <!-- 步骤3：完成创建 -->
      <div v-if="currentStep === 3" class="step-content">
        <div class="review-section">
          <div class="section-header">
            <i class="section-icon">✓</i>
            <span class="section-title">确认信息</span>
          </div>
          <div class="review-content">
            <div class="review-item">
              <span class="review-label">会员姓名：</span>
              <span class="review-value">{{ memberForm.name }}</span>
            </div>
            <div class="review-item">
              <span class="review-label">手机号：</span>
              <span class="review-value">{{ memberForm.phone }}</span>
            </div>
            <div class="review-item">
              <span class="review-label">所属门店：</span>
              <span class="review-value">{{ memberForm.store_name }}</span>
            </div>
            <div class="review-item">
              <span class="review-label">服务类型：</span>
              <span class="review-value">{{ getServiceTypeText(memberForm.service_type) }}</span>
            </div>
            <div class="review-item">
              <span class="review-label">套餐名称：</span>
              <span class="review-value">{{ memberForm.package_name }}</span>
            </div>
            <div class="review-item">
              <span class="review-label">套餐价格：</span>
              <span class="review-value">¥{{ memberForm.package_price || '0.00' }}</span>
            </div>
            <div class="review-item">
              <span class="review-label">有效期：</span>
              <span class="review-value">{{ memberForm.valid_from }} 至 {{ memberForm.valid_to }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部操作按钮 -->
    <div class="bottom-actions">
      <van-button
        v-if="currentStep > 1"
        class="action-btn prev-btn"
        @click="prevStep"
      >
        上一步
      </van-button>
      <van-button
        v-if="currentStep === 2 && !faceImageCaptured"
        class="action-btn create-btn"
        type="primary"
        @click="createMemberWithoutFace"
      >
        创建
      </van-button>
      <van-button
        v-if="currentStep === 3"
        class="action-btn submit-btn"
        type="primary"
        @click="submitMember"
        :loading="submitting"
        :disabled="!canSubmit"
      >
        {{ submitting ? '创建中...' : '创建会员' }}
      </van-button>
    </div>

    <!-- 门店选择器 -->
    <van-popup v-model:show="showStorePicker" position="bottom" round>
      <van-picker
        :columns="storeOptions"
        @confirm="onStoreConfirm"
        @cancel="showStorePicker = false"
      />
    </van-popup>

    <!-- 服务类型选择器 -->
    <van-popup v-model:show="showServiceTypePicker" position="bottom" round>
      <van-picker
        :columns="serviceTypeOptions"
        @confirm="onServiceTypeConfirm"
        @cancel="showServiceTypePicker = false"
      />
    </van-popup>

    <!-- 会员状态选择器 -->
    <van-popup v-model:show="showStatusPicker" position="bottom" round>
      <van-picker
        :columns="statusOptions"
        @confirm="onStatusConfirm"
        @cancel="showStatusPicker = false"
      />
    </van-popup>

    <!-- 日期选择器 -->
    <van-popup v-model:show="showStartDatePicker" position="bottom" round>
      <van-date-picker
        v-model="startDate"
        @confirm="onStartDateConfirm"
        @cancel="showStartDatePicker = false"
      />
    </van-popup>

    <van-popup v-model:show="showEndDatePicker" position="bottom" round>
      <van-date-picker
        v-model="endDate"
        @confirm="onEndDateConfirm"
        @cancel="showEndDatePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showSuccessToast, showFailToast } from 'vant'
import { createMember } from '@/api/members'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const currentStep = ref(1)
const submitting = ref(false)
const showStorePicker = ref(false)
const showServiceTypePicker = ref(false)
const showStatusPicker = ref(false)
const showStartDatePicker = ref(false)
const showEndDatePicker = ref(false)

const startDate = ref([dayjs().format('YYYY'), dayjs().format('MM'), dayjs().format('DD')])
const endDate = ref([dayjs().add(365, 'day').format('YYYY'), dayjs().add(365, 'day').format('MM'), dayjs().add(365, 'day').format('DD')])

const memberForm = reactive({
  name: '',
  phone: '',
  store_id: null,
  store_name: '',
  service_type: '',
  package_name: '',
  package_price: '0.00',
  purchase_amount: '0.00',
  status: 'active',
  valid_from: dayjs().format('YYYY-MM-DD'),
  valid_to: dayjs().add(365, 'day').format('YYYY-MM-DD'),
  fixed_duration: '',
  remarks: '',
  face_image: ''
})

// 摄像头相关
const videoRef = ref(null)
const canvasRef = ref(null)
const faceImageCaptured = ref(false)
const faceImagePreview = ref('')
const faceImageFile = ref(null)
const isCapturing = ref(false)
const cameraError = ref('')
let mediaStream = null

// 门店选项（从用户信息中获取）
const storeOptions = computed(() => {
  const storeId = userStore.userInfo?.store_id
  const storeName = userStore.userInfo?.store_name || '当前门店'
  if (storeId) {
    return [{ text: storeName, value: storeId }]
  }
  return []
})

// 服务类型选项
const serviceTypeOptions = [
  { text: '美甲', value: 'nail' },
  { text: '美睫', value: 'eyelash' },
  { text: '组合', value: 'combo' }
]

// 会员状态选项
const statusOptions = [
  { text: '有效', value: 'active' },
  { text: '过期', value: 'expired' },
  { text: '停用', value: 'inactive' }
]

// 获取服务类型文本
const getServiceTypeText = (value) => {
  const option = serviceTypeOptions.find(opt => opt.value === value)
  return option ? option.text : value
}

// 获取会员状态文本
const getStatusText = (value) => {
  const option = statusOptions.find(opt => opt.value === value)
  return option ? option.text : value
}

// 会员状态显示文本（计算属性）
const statusDisplayText = computed(() => {
  return getStatusText(memberForm.status) || '有效'
})

// 检查是否可以进入下一步
const canNextStep = computed(() => {
  return memberForm.name &&
    memberForm.phone &&
    memberForm.store_id &&
    memberForm.service_type &&
    memberForm.package_name &&
    memberForm.package_price &&
    memberForm.valid_from &&
    memberForm.valid_to
})

// 检查是否可以提交
const canSubmit = computed(() => {
  return canNextStep.value && currentStep.value === 3
})

// 数字步进器方法
const decreasePrice = () => {
  const price = parseFloat(memberForm.package_price) || 0
  if (price > 0) {
    memberForm.package_price = (price - 1).toFixed(2)
  }
}

const increasePrice = () => {
  const price = parseFloat(memberForm.package_price) || 0
  memberForm.package_price = (price + 1).toFixed(2)
}

const decreasePurchase = () => {
  const amount = parseFloat(memberForm.purchase_amount) || 0
  if (amount > 0) {
    memberForm.purchase_amount = (amount - 1).toFixed(2)
  }
}

const increasePurchase = () => {
  const amount = parseFloat(memberForm.purchase_amount) || 0
  memberForm.purchase_amount = (amount + 1).toFixed(2)
}

const decreaseDuration = () => {
  const duration = parseInt(memberForm.fixed_duration) || 0
  if (duration > 0) {
    memberForm.fixed_duration = (duration - 1).toString()
    calculateEndDate()
  }
}

const increaseDuration = () => {
  const duration = parseInt(memberForm.fixed_duration) || 0
  memberForm.fixed_duration = (duration + 1).toString()
  calculateEndDate()
}

// 计算结束日期
const calculateEndDate = () => {
  if (memberForm.fixed_duration && memberForm.valid_from) {
    const days = parseInt(memberForm.fixed_duration)
    if (days > 0) {
      memberForm.valid_to = dayjs(memberForm.valid_from).add(days, 'day').format('YYYY-MM-DD')
      endDate.value = [
        dayjs(memberForm.valid_to).format('YYYY'),
        dayjs(memberForm.valid_to).format('MM'),
        dayjs(memberForm.valid_to).format('DD')
      ]
    }
  }
}

// 计算固定时长
const calculateDuration = () => {
  if (memberForm.valid_from && memberForm.valid_to) {
    const days = dayjs(memberForm.valid_to).diff(dayjs(memberForm.valid_from), 'day')
    if (days > 0) {
      memberForm.fixed_duration = days.toString()
    }
  }
}

// 选择器确认事件
const onStoreConfirm = ({ selectedOptions }) => {
  if (selectedOptions.length > 0) {
    memberForm.store_id = selectedOptions[0].value
    memberForm.store_name = selectedOptions[0].text
  }
  showStorePicker.value = false
}

const onServiceTypeConfirm = ({ selectedOptions }) => {
  if (selectedOptions.length > 0) {
    memberForm.service_type = selectedOptions[0].value
  }
  showServiceTypePicker.value = false
}

const onStatusConfirm = ({ selectedOptions }) => {
  if (selectedOptions.length > 0) {
    memberForm.status = selectedOptions[0].value
  }
  showStatusPicker.value = false
}

const onStartDateConfirm = ({ selectedValues }) => {
  const date = selectedValues.join('-')
  memberForm.valid_from = date
  startDate.value = selectedValues
  calculateDuration()
  showStartDatePicker.value = false
}

const onEndDateConfirm = ({ selectedValues }) => {
  const date = selectedValues.join('-')
  memberForm.valid_to = date
  endDate.value = selectedValues
  calculateDuration()
  showEndDatePicker.value = false
}

// 下一步
const handleNextStep = () => {
  if (canNextStep.value) {
    currentStep.value = 2
  }
}

// 上一步
const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

// 启动摄像头
const startCamera = async () => {
  try {
    cameraError.value = ''
    mediaStream = await navigator.mediaDevices.getUserMedia({
      video: {
        facingMode: 'user', // 前置摄像头
        width: { ideal: 640 },
        height: { ideal: 480 }
      }
    })
    if (videoRef.value) {
      videoRef.value.srcObject = mediaStream
    }
  } catch (error) {
    console.error('启动摄像头失败:', error)
    cameraError.value = '无法访问摄像头，请检查权限设置'
    showFailToast('启动摄像头失败，请检查权限设置')
  }
}

// 停止摄像头
const stopCamera = () => {
  if (mediaStream) {
    mediaStream.getTracks().forEach(track => track.stop())
    mediaStream = null
  }
  if (videoRef.value) {
    videoRef.value.srcObject = null
  }
}

// 拍照
const capturePhoto = async () => {
  if (!videoRef.value || !canvasRef.value) return
  
  isCapturing.value = true
  try {
    const canvas = canvasRef.value
    const video = videoRef.value
    const context = canvas.getContext('2d')
    
    // 设置画布尺寸
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    
    // 绘制当前视频帧
    context.drawImage(video, 0, 0, canvas.width, canvas.height)
    
    // 转换为blob
    canvas.toBlob((blob) => {
      if (blob) {
        // 创建File对象
        faceImageFile.value = new File([blob], 'face.jpg', { type: 'image/jpeg' })
        
        // 创建预览URL
        faceImagePreview.value = URL.createObjectURL(blob)
        faceImageCaptured.value = true
        
        // 停止摄像头
        stopCamera()
        
        showSuccessToast('拍照成功')
      }
    }, 'image/jpeg', 0.8)
    
  } catch (error) {
    console.error('拍照失败:', error)
    showFailToast('拍照失败，请重试')
  } finally {
    isCapturing.value = false
  }
}

// 重新拍照
const retakePhoto = () => {
  faceImageCaptured.value = false
  faceImagePreview.value = ''
  faceImageFile.value = null
  if (faceImagePreview.value) {
    URL.revokeObjectURL(faceImagePreview.value)
  }
  startCamera()
}

// 重试摄像头
const retryCamera = () => {
  cameraError.value = ''
  startCamera()
}

// 不录入人脸直接创建会员
const createMemberWithoutFace = () => {
  stopCamera()
  // 直接提交创建会员（不包含人脸图片）
  submitMember()
}

// 检查是否可以不录入人脸直接创建
const canCreateWithoutFace = computed(() => {
  // 如果没有拍照，允许不录入人脸直接创建
  return !faceImageCaptured.value
})

// 进入下一步（完成人脸录入后）
const goToNextStep = () => {
  if (faceImageCaptured.value) {
    stopCamera()
    currentStep.value = 3
  }
}


// 监听步骤变化，进入步骤2时启动摄像头
watch(currentStep, (newStep) => {
  if (newStep === 2 && !faceImageCaptured.value) {
    startCamera()
  } else if (newStep !== 2) {
    stopCamera()
  }
})

// 提交创建会员
const submitMember = async () => {
  if (!canSubmit.value) {
    showFailToast('请完善必填信息')
    return
  }

  submitting.value = true
  try {
    // 使用FormData上传
    const formData = new FormData()
    
    // 添加表单字段
    formData.append('name', memberForm.name)
    formData.append('phone', memberForm.phone)
    formData.append('store_id', memberForm.store_id.toString())
    formData.append('service_type', memberForm.service_type)
    formData.append('package_name', memberForm.package_name)
    formData.append('price', parseFloat(memberForm.package_price).toString())
    formData.append('purchase_amount', (parseFloat(memberForm.purchase_amount) || 0).toString())
    formData.append('status', memberForm.status)
    if (memberForm.remarks) {
      formData.append('description', memberForm.remarks)
    }
    
    // 处理有效期：使用YYYY-MM-DD格式
    if (memberForm.valid_from) {
      formData.append('valid_from', memberForm.valid_from)
    }
    
    if (memberForm.valid_to) {
      formData.append('valid_to', memberForm.valid_to)
    }
    
    // 如果提供了固定时长，也发送
    if (memberForm.fixed_duration) {
      formData.append('validity_duration', memberForm.fixed_duration)
    }

    // 如果有人脸图片，添加到FormData
    if (faceImageFile.value) {
      formData.append('face_image', faceImageFile.value)
    }

    const response = await createMember(formData)
    
    if (response.code === 200 || response.code === 0) {
      showSuccessToast('创建会员成功')
      router.push('/members')
    } else {
      showFailToast(response.message || '创建会员失败')
    }
  } catch (error) {
    console.error('创建会员失败:', error)
    showFailToast(error.response?.data?.message || '创建会员失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

// 初始化
onMounted(() => {
  // 自动填充门店信息
  const storeId = userStore.userInfo?.store_id
  const storeName = userStore.userInfo?.store_name || '当前门店'
  if (storeId) {
    memberForm.store_id = storeId
    memberForm.store_name = storeName
  }
})

// 组件卸载时清理
onUnmounted(() => {
  stopCamera()
  if (faceImagePreview.value) {
    URL.revokeObjectURL(faceImagePreview.value)
  }
})
</script>

<style lang="scss" scoped>
.member-create-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 80px;
}

.custom-nav {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  :deep(.van-nav-bar__title) {
    color: white;
  }
  
  :deep(.van-nav-bar__text),
  :deep(.van-nav-bar__arrow) {
    color: white;
  }
}

// 步骤指示器
.steps-section {
  background: white;
  padding: 20px;
  margin-bottom: 12px;
}

.steps-container {
  display: flex;
  align-items: center;
  justify-content: center;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  
  .step-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #e0e0e0;
    color: #999;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 600;
    transition: all 0.3s;
  }
  
  .step-text {
    margin-top: 8px;
    font-size: 12px;
    color: #999;
  }
  
  &.active {
    .step-icon {
      background: linear-gradient(135deg, #667eea, #764ba2);
      color: white;
    }
    .step-text {
      color: #667eea;
      font-weight: 600;
    }
  }
  
  &.completed {
    .step-icon {
      background: #52c41a;
      color: white;
    }
  }
}

.step-line {
  width: 60px;
  height: 2px;
  background: #e0e0e0;
  margin: 0 12px;
  margin-top: -20px;
}

.content-container {
  padding: 0 16px;
}

.step-content {
  animation: fadeIn 0.3s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 表单样式
.form-section {
  background: white;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.van-field) {
  padding: 12px 0;
  
  .van-field__label {
    width: 90px;
    font-size: 14px;
    color: #333;
  }
  
  .van-field__control {
    font-size: 14px;
  }
  
  .van-field__required-mark {
    color: #ff4d4f;
  }
}

.number-stepper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stepper-btn {
  width: 28px;
  height: 28px;
  border: 1px solid #e0e0e0;
  background: white;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: #666;
  cursor: pointer;
  
  &:active {
    background: #f5f5f5;
  }
}

.form-hint {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
  line-height: 1.5;
}

.form-actions {
  margin-top: 24px;
  margin-bottom: 24px;
}

.next-btn {
  height: 48px;
  font-size: 16px;
  font-weight: 600;
}

// 人脸录入
.face-capture-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 24px;
  
  .section-icon {
    font-size: 24px;
  }
  
  .section-title {
    font-size: 18px;
    font-weight: 600;
    color: #333;
  }
}

.face-placeholder {
  padding: 60px 20px;
  
  .placeholder-icon {
    font-size: 64px;
    margin-bottom: 16px;
  }
  
  .placeholder-text {
    font-size: 16px;
    color: #666;
    margin-bottom: 8px;
  }
  
  .placeholder-hint {
    font-size: 14px;
    color: #999;
  }
}

// 摄像头相关样式
.camera-container {
  margin-bottom: 24px;
  position: relative;
  background: #000;
  border-radius: 12px;
  overflow: hidden;
  
  .camera-video {
    width: 100%;
    max-width: 100%;
    display: block;
    object-fit: cover;
  }
  
  .camera-hint {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: linear-gradient(to top, rgba(0,0,0,0.7), transparent);
    color: white;
    padding: 16px;
    font-size: 14px;
    text-align: center;
  }
}

.face-preview-container {
  margin-bottom: 24px;
  text-align: center;
  
  .face-preview-image {
    width: 100%;
    max-width: 400px;
    border-radius: 12px;
    box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  }
  
  .preview-hint {
    margin-top: 12px;
    color: #07c160;
    font-size: 14px;
    font-weight: 500;
  }
}

.camera-error {
  padding: 40px 20px;
  text-align: center;
  
  .error-icon {
    font-size: 48px;
    margin-bottom: 16px;
  }
  
  .error-text {
    font-size: 14px;
    color: #666;
    margin-bottom: 16px;
  }
}

.face-actions {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  
  .next-step-btn {
    margin-top: 8px;
  }
}

// 确认信息
.review-section {
  background: white;
  border-radius: 12px;
  padding: 16px;
}

.review-content {
  .review-item {
    display: flex;
    padding: 12px 0;
    border-bottom: 1px solid #f0f0f0;
    
    &:last-child {
      border-bottom: none;
    }
    
    .review-label {
      width: 100px;
      font-size: 14px;
      color: #666;
    }
    
    .review-value {
      flex: 1;
      font-size: 14px;
      color: #333;
      font-weight: 500;
    }
  }
}

// 底部操作按钮
.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  padding: 12px 16px;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  display: flex;
  gap: 12px;
  z-index: 100;
}

.action-btn {
  flex: 1;
  height: 44px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  
  &.prev-btn {
    background: #f5f5f5;
    color: #666;
  }
  
  &.skip-btn {
    background: #f5f5f5;
    color: #666;
  }
  
  &.submit-btn {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
  }
}
</style>
