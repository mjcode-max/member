<template>
  <div class="member-code-page">
    <!-- 页面头部 -->
    <van-nav-bar
      title="会员码"
      left-text="返回"
      left-arrow
      @click-left="$router.back()"
    />
    
    <!-- 人脸验证区域 -->
    <div v-if="!isVerified" class="face-verify-section">
      <div class="verify-title">人脸验证</div>
      <div class="verify-desc">请正对手机，确保光线充足</div>
      
      <!-- 摄像头区域 -->
      <div class="camera-container">
        <video
          ref="videoRef"
          class="camera-video"
          autoplay
          muted
          playsinline
        ></video>
        <canvas
          ref="canvasRef"
          class="capture-canvas"
          style="display: none;"
        ></canvas>
      </div>
      
      <!-- 操作按钮 -->
      <div class="verify-actions">
        <van-button
          type="primary"
          size="large"
          :loading="isCapturing"
          @click="capturePhoto"
        >
          拍照验证
        </van-button>
        <van-button
          type="default"
          size="large"
          @click="retryCapture"
        >
          重新拍照
        </van-button>
      </div>
    </div>
    
    <!-- 会员码显示区域 -->
    <div v-else class="member-code-section">
      <div class="code-title">动态会员码</div>
      <div class="code-desc">请向店员出示此码</div>
      
      <!-- 二维码 -->
      <div class="qrcode-container">
        <div class="qrcode" ref="qrcodeRef"></div>
        <div class="countdown">
          {{ countdown }}秒后刷新
        </div>
      </div>
      
      <!-- 会员信息 -->
      <div class="member-info">
        <div class="info-item">
          <span class="label">会员姓名：</span>
          <span class="value">{{ memberInfo.name }}</span>
        </div>
        <div class="info-item">
          <span class="label">剩余次数：</span>
          <span class="value">{{ memberInfo.remaining_times }}</span>
        </div>
        <div class="info-item">
          <span class="label">有效期至：</span>
          <span class="value">{{ formatDate(memberInfo.valid_to) }}</span>
        </div>
      </div>
      
      <!-- 操作按钮 -->
      <div class="code-actions">
        <van-button
          type="primary"
          size="large"
          @click="refreshCode"
        >
          刷新会员码
        </van-button>
        <van-button
          type="default"
          size="large"
          @click="finish"
        >
          完成
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showDialog } from 'vant'
import QRCode from 'qrcode'
import { generateMemberCode } from '@/api/members'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()

const videoRef = ref(null)
const canvasRef = ref(null)
const qrcodeRef = ref(null)

const isVerified = ref(false)
const isCapturing = ref(false)
const countdown = ref(10)
const memberCode = ref('')
const memberInfo = ref({})

let stream = null
let countdownTimer = null
let refreshTimer = null

// 启动摄像头
const startCamera = async () => {
  try {
    stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user' }
    })
    videoRef.value.srcObject = stream
  } catch (error) {
    console.error('启动摄像头失败:', error)
    showToast('启动摄像头失败，请检查权限设置')
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
    
    // 获取图片数据
    const imageData = canvas.toDataURL('image/jpeg', 0.8)
    
    // 停止摄像头
    stopCamera()
    
    // 调用人脸验证API
    await verifyFace(imageData)
    
  } catch (error) {
    console.error('拍照失败:', error)
    showToast('拍照失败，请重试')
  } finally {
    isCapturing.value = false
  }
}

// 人脸验证
const verifyFace = async (imageData) => {
  try {
    const memberId = route.params.id
    const response = await generateMemberCode({
      member_id: memberId,
      face_image: imageData
    })
    
    if (response.code === 200) {
      memberCode.value = response.data.member_code
      memberInfo.value = response.data.member
      isVerified.value = true
      
      // 生成二维码
      generateQRCode()
      
      // 启动倒计时
      startCountdown()
      
      showToast('验证成功')
    } else {
      showToast(response.message || '验证失败')
    }
  } catch (error) {
    console.error('人脸验证失败:', error)
    showToast('验证失败，请重试')
  }
}

// 生成二维码
const generateQRCode = async () => {
  try {
    const qrCodeDataURL = await QRCode.toDataURL(memberCode.value, {
      width: 200,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    })
    
    if (qrcodeRef.value) {
      qrcodeRef.value.innerHTML = `<img src="${qrCodeDataURL}" alt="会员码" />`
    }
  } catch (error) {
    console.error('生成二维码失败:', error)
  }
}

// 启动倒计时
const startCountdown = () => {
  countdown.value = 10
  
  // 清除之前的定时器
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
  
  countdownTimer = setInterval(() => {
    countdown.value--
    
    if (countdown.value <= 0) {
      refreshCode()
    }
  }, 1000)
}

// 刷新会员码
const refreshCode = async () => {
  try {
    const memberId = route.params.id
    const response = await generateMemberCode({
      member_id: memberId,
      face_image: '' // 重新验证时不需要人脸图片
    })
    
    if (response.code === 200) {
      memberCode.value = response.data.member_code
      await generateQRCode()
      startCountdown()
      showToast('会员码已刷新')
    }
  } catch (error) {
    console.error('刷新会员码失败:', error)
    showToast('刷新失败，请重试')
  }
}

// 重新拍照
const retryCapture = () => {
  startCamera()
}

// 完成
const finish = () => {
  showDialog({
    title: '确认',
    message: '确定要退出会员码页面吗？',
  }).then(() => {
    router.back()
  }).catch(() => {
    // 用户取消
  })
}

// 停止摄像头
const stopCamera = () => {
  if (stream) {
    stream.getTracks().forEach(track => track.stop())
    stream = null
  }
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

// 清理定时器
const cleanup = () => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
  if (refreshTimer) {
    clearTimeout(refreshTimer)
    refreshTimer = null
  }
  stopCamera()
}

onMounted(() => {
  startCamera()
})

onUnmounted(() => {
  cleanup()
})
</script>

<style lang="scss" scoped>
.member-code-page {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
  position: relative;
  overflow-x: hidden;
}

// 人脸验证部分
.face-verify-section {
  padding: 40px 20px;
  text-align: center;
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border-radius: 0 0 30px 30px;
  }
  
  > * {
    position: relative;
    z-index: 1;
  }
}

.verify-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.verify-desc {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 40px;
  line-height: 1.5;
}

.camera-container {
  position: relative;
  width: 300px;
  height: 300px;
  margin: 0 auto 40px;
  border-radius: 20px;
  overflow: hidden;
  background: #000;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  border: 4px solid rgba(255, 255, 255, 0.2);
  
  &::before {
    content: '';
    position: absolute;
    top: 10px;
    left: 10px;
    right: 10px;
    bottom: 10px;
    border: 2px dashed rgba(255, 255, 255, 0.5);
    border-radius: 16px;
    z-index: 2;
    pointer-events: none;
  }
  
  &::after {
    content: '请将人脸对准框内';
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    color: rgba(255, 255, 255, 0.8);
    font-size: 12px;
    z-index: 3;
    background: rgba(0, 0, 0, 0.5);
    padding: 4px 12px;
    border-radius: 12px;
    backdrop-filter: blur(10px);
  }
}

.camera-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transform: scaleX(-1); // 镜像翻转
}

.verify-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
  
  .van-button {
    min-width: 120px;
    height: 48px;
    border-radius: 24px;
    font-weight: 600;
    font-size: 16px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
    
    &--primary {
      background: linear-gradient(135deg, #ff6b6b, #ffa726);
      border: none;
    }
    
    &--default {
      background: rgba(255, 255, 255, 0.2);
      color: #fff;
      border: 2px solid rgba(255, 255, 255, 0.3);
      backdrop-filter: blur(10px);
    }
  }
}

// 会员码展示部分
.member-code-section {
  padding: 30px 20px;
  text-align: center;
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    top: -20px;
    left: 20px;
    right: 20px;
    bottom: 20px;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(15px);
    border-radius: 30px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  }
  
  > * {
    position: relative;
    z-index: 1;
  }
}

.code-title {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  
  &::before {
    content: '✨';
    font-size: 20px;
  }
}

.code-desc {
  font-size: 15px;
  color: #666;
  margin-bottom: 30px;
  opacity: 0.8;
}

.qrcode-container {
  margin-bottom: 30px;
  position: relative;
}

.qrcode {
  width: 220px;
  height: 220px;
  margin: 0 auto 20px;
  background: #fff;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  border: 3px solid transparent;
  background-clip: padding-box;
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    inset: -3px;
    background: linear-gradient(135deg, #ff6b6b, #ffa726, #667eea, #764ba2);
    border-radius: 23px;
    z-index: -1;
  }
  
  img {
    width: 190px;
    height: 190px;
    border-radius: 12px;
  }
}

.countdown {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #ff6b6b, #ffa726);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 10px;
  
  &::before {
    content: '⏰ ';
    -webkit-text-fill-color: #ff6b6b;
  }
}

.member-info {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 30px;
  text-align: left;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    background: rgba(255, 255, 255, 0.5);
    margin: 0 -12px;
    padding: 16px 12px;
    border-radius: 12px;
  }
}

.label {
  font-size: 15px;
  color: #666;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
  
  &::before {
    content: '•';
    color: #ff6b6b;
    font-weight: bold;
  }
}

.value {
  font-size: 15px;
  color: #333;
  font-weight: 600;
}

.code-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
  padding-bottom: 20px;
  
  .van-button {
    min-width: 140px;
    height: 48px;
    border-radius: 24px;
    font-weight: 600;
    font-size: 16px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
    
    &--primary {
      background: linear-gradient(135deg, #667eea, #764ba2);
      border: none;
    }
    
    &--default {
      background: rgba(255, 255, 255, 0.8);
      color: #666;
      border: 2px solid rgba(0, 0, 0, 0.1);
      backdrop-filter: blur(10px);
    }
  }
}

// 响应式设计
@media (max-width: 375px) {
  .camera-container {
    width: 260px;
    height: 260px;
  }
  
  .qrcode {
    width: 180px;
    height: 180px;
    
    img {
      width: 150px;
      height: 150px;
    }
  }
  
  .verify-actions,
  .code-actions {
    flex-direction: column;
    align-items: center;
    
    .van-button {
      width: 200px;
    }
  }
}

// 动画效果
@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.qrcode {
  animation: pulse 2s ease-in-out infinite;
}

.member-code-section > * {
  animation: fadeInUp 0.6s ease-out;
}

.member-code-section > *:nth-child(1) { animation-delay: 0.1s; }
.member-code-section > *:nth-child(2) { animation-delay: 0.2s; }
.member-code-section > *:nth-child(3) { animation-delay: 0.3s; }
.member-code-section > *:nth-child(4) { animation-delay: 0.4s; }

.countdown::before {
  animation: rotate 2s linear infinite;
  display: inline-block;
}

// 加载状态
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
  
  .loading-content {
    background: rgba(255, 255, 255, 0.95);
    padding: 30px;
    border-radius: 20px;
    text-align: center;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
    
    .loading-text {
      margin-top: 16px;
      font-size: 16px;
      color: #333;
      font-weight: 500;
    }
  }
}
</style>
