<template>
  <div class="scanner-page">
    <!-- 页面头部 -->
    <van-nav-bar
      title="扫码核销"
      left-text="返回"
      left-arrow
      @click-left="$router.back()"
      class="custom-nav"
    />

    <!-- 扫码区域 -->
    <div class="scanner-container">
      <div class="scanner-header">
        <h2 class="scanner-title">扫描会员码</h2>
        <p class="scanner-desc">请将会员码对准扫描框</p>
      </div>

      <div class="camera-container">
        <video
          ref="videoRef"
          class="camera-video"
          autoplay
          muted
          playsinline
        ></video>
        
        <!-- 扫描框 -->
        <div class="scan-frame">
          <div class="scan-corner corner-tl"></div>
          <div class="scan-corner corner-tr"></div>
          <div class="scan-corner corner-bl"></div>
          <div class="scan-corner corner-br"></div>
          <div class="scan-line"></div>
        </div>
        
        <!-- 扫描提示 -->
        <div class="scan-tips">
          <div class="tips-text">请将二维码放入框内</div>
        </div>
      </div>

      <div class="scanner-actions">
        <button class="action-btn torch-btn" @click="toggleTorch">
          <i class="btn-icon">🔦</i>
          {{ torchOn ? '关闭闪光灯' : '打开闪光灯' }}
        </button>
        <button class="action-btn manual-btn" @click="showManualInput = true">
          <i class="btn-icon">⌨️</i>
          手动输入
        </button>
      </div>
    </div>

    <!-- 扫描历史 -->
    <div class="history-section">
      <div class="section-header">
        <i class="section-icon">🕒</i>
        <span class="section-title">最近扫描</span>
        <span class="section-clear" @click="clearHistory">清空</span>
      </div>
      <div class="history-list">
        <div 
          v-for="record in scanHistory"
          :key="record.id"
          class="history-item"
          @click="handleHistoryClick(record)"
        >
          <div class="history-avatar">
            {{ record.member_name.charAt(0) }}
          </div>
          <div class="history-info">
            <div class="history-name">{{ record.member_name }}</div>
            <div class="history-service">{{ getServiceText(record.service_type) }}</div>
            <div class="history-time">{{ formatDateTime(record.scan_time) }}</div>
          </div>
          <div class="history-status">
            <i class="status-icon">✅</i>
          </div>
        </div>
      </div>
    </div>

    <!-- 手动输入弹窗 -->
    <van-popup 
      v-model:show="showManualInput" 
      position="center"
      round
      :style="{ width: '90%' }"
    >
      <div class="manual-dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">手动输入会员码</h3>
          <div class="dialog-close" @click="showManualInput = false">×</div>
        </div>
        <div class="dialog-content">
          <div class="input-group">
            <label class="input-label">会员码</label>
            <input 
              v-model="manualCode"
              type="text"
              placeholder="请输入会员码"
              class="code-input"
              @keyup.enter="verifyManualCode"
            />
          </div>
          <div class="service-selection">
            <label class="input-label">服务类型</label>
            <div class="service-options">
              <div 
                class="service-option"
                :class="{ active: selectedService === 'manicure' }"
                @click="selectedService = 'manicure'"
              >
                <i class="service-icon">💅</i>
                <span>美甲</span>
              </div>
              <div 
                class="service-option"
                :class="{ active: selectedService === 'eyelash' }"
                @click="selectedService = 'eyelash'"
              >
                <i class="service-icon">👁️</i>
                <span>美睫</span>
              </div>
            </div>
          </div>
        </div>
        <div class="dialog-actions">
          <button class="dialog-btn cancel-btn" @click="showManualInput = false">
            取消
          </button>
          <button class="dialog-btn verify-btn" @click="verifyManualCode">
            验证核销
          </button>
        </div>
      </div>
    </van-popup>

    <!-- 核销成功弹窗 -->
    <van-popup 
      v-model:show="showSuccessDialog" 
      position="center"
      round
      :style="{ width: '90%' }"
    >
      <div class="success-dialog">
        <div class="success-icon">🎉</div>
        <h3 class="success-title">核销成功</h3>
        <div class="success-content">
          <div class="member-info">
            <div class="info-item">
              <span class="info-label">会员姓名：</span>
              <span class="info-value">{{ verifiedMember.name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">服务类型：</span>
              <span class="info-value">{{ getServiceText(selectedService) }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">剩余次数：</span>
              <span class="info-value">{{ verifiedMember.remaining_times }}次</span>
            </div>
          </div>
        </div>
        <div class="success-actions">
          <button class="success-btn" @click="closeSuccessDialog">
            完成
          </button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { verifyMemberCode } from '@/api/members'
import dayjs from 'dayjs'

const router = useRouter()

const videoRef = ref(null)
const showManualInput = ref(false)
const showSuccessDialog = ref(false)
const torchOn = ref(false)
const manualCode = ref('')
const selectedService = ref('manicure')

const verifiedMember = reactive({
  name: '',
  remaining_times: 0
})

const scanHistory = ref([
  {
    id: 1,
    member_name: '张小姐',
    service_type: 'manicure',
    scan_time: '2023-10-15 14:30:00'
  },
  {
    id: 2,
    member_name: '王女士',
    service_type: 'eyelash',
    scan_time: '2023-10-15 13:15:00'
  },
  {
    id: 3,
    member_name: '刘小姐',
    service_type: 'manicure',
    scan_time: '2023-10-15 11:45:00'
  }
])

let stream = null
let scanTimer = null

// 启动摄像头
const startCamera = async () => {
  try {
    stream = await navigator.mediaDevices.getUserMedia({
      video: { 
        facingMode: 'environment', // 后置摄像头
        width: { ideal: 1280 },
        height: { ideal: 720 }
      }
    })
    if (videoRef.value) {
      videoRef.value.srcObject = stream
      startScanning()
    }
  } catch (error) {
    console.error('启动摄像头失败:', error)
    showToast('启动摄像头失败，请检查权限设置')
  }
}

// 开始扫描
const startScanning = () => {
  // 这里应该集成二维码扫描库，如 jsQR
  // 简化处理，使用定时器模拟扫描
  scanTimer = setInterval(() => {
    // 模拟扫描到二维码
    // 实际实现需要使用 jsQR 或类似库
  }, 100)
}

// 停止摄像头
const stopCamera = () => {
  if (stream) {
    stream.getTracks().forEach(track => track.stop())
    stream = null
  }
  if (scanTimer) {
    clearInterval(scanTimer)
    scanTimer = null
  }
}

// 切换闪光灯
const toggleTorch = async () => {
  if (!stream) return
  
  try {
    const track = stream.getVideoTracks()[0]
    const capabilities = track.getCapabilities()
    
    if (capabilities.torch) {
      await track.applyConstraints({
        advanced: [{ torch: !torchOn.value }]
      })
      torchOn.value = !torchOn.value
      showToast(torchOn.value ? '闪光灯已打开' : '闪光灯已关闭')
    } else {
      showToast('设备不支持闪光灯')
    }
  } catch (error) {
    console.error('切换闪光灯失败:', error)
    showToast('闪光灯操作失败')
  }
}

// 手动验证会员码
const verifyManualCode = async () => {
  if (!manualCode.value) {
    showToast('请输入会员码')
    return
  }

  if (!selectedService.value) {
    showToast('请选择服务类型')
    return
  }

  try {
    const response = await verifyMemberCode({
      member_code: manualCode.value,
      service_type: selectedService.value
    })

    if (response.code === 200) {
      Object.assign(verifiedMember, response.data.member)
      
      // 添加到扫描历史
      scanHistory.value.unshift({
        id: Date.now(),
        member_name: verifiedMember.name,
        service_type: selectedService.value,
        scan_time: new Date().toISOString()
      })

      showManualInput.value = false
      showSuccessDialog.value = true
      manualCode.value = ''
      
      showToast('会员码验证成功')
    }
  } catch (error) {
    console.error('验证会员码失败:', error)
    showToast('验证失败，请检查会员码是否正确')
  }
}

// 关闭成功弹窗
const closeSuccessDialog = () => {
  showSuccessDialog.value = false
  Object.assign(verifiedMember, { name: '', remaining_times: 0 })
}

// 处理历史记录点击
const handleHistoryClick = (record) => {
  showToast(`${record.member_name} - ${getServiceText(record.service_type)}`)
}

// 清空历史
const clearHistory = () => {
  scanHistory.value = []
  showToast('历史记录已清空')
}

// 获取服务类型文本
const getServiceText = (type) => {
  return type === 'manicure' ? '美甲服务' : '美睫服务'
}

// 格式化日期时间
const formatDateTime = (datetime) => {
  return dayjs(datetime).format('MM-DD HH:mm')
}

onMounted(() => {
  startCamera()
})

onUnmounted(() => {
  stopCamera()
})
</script>

<style lang="scss" scoped>
.scanner-page {
  background: linear-gradient(180deg, #1a1a1a 0%, #2d2d2d 100%);
  min-height: 100vh;
  color: white;
}

// 自定义导航栏
.custom-nav {
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(10px);
  
  :deep(.van-nav-bar__title) {
    color: white;
    font-weight: 600;
  }
  
  :deep(.van-nav-bar__text) {
    color: rgba(255, 255, 255, 0.8);
  }
  
  :deep(.van-icon) {
    color: white;
  }
}

// 扫码容器
.scanner-container {
  padding: 20px;
  text-align: center;
}

.scanner-header {
  margin-bottom: 30px;
}

.scanner-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 8px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
}

.scanner-desc {
  font-size: 16px;
  opacity: 0.8;
  margin: 0;
}

// 摄像头容器
.camera-container {
  position: relative;
  width: 280px;
  height: 280px;
  margin: 0 auto 30px;
  border-radius: 20px;
  overflow: hidden;
  background: #000;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.camera-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

// 扫描框
.scan-frame {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 200px;
  height: 200px;
  pointer-events: none;
}

.scan-corner {
  position: absolute;
  width: 20px;
  height: 20px;
  border: 3px solid #ff6b6b;
  
  &.corner-tl {
    top: 0;
    left: 0;
    border-right: none;
    border-bottom: none;
  }
  
  &.corner-tr {
    top: 0;
    right: 0;
    border-left: none;
    border-bottom: none;
  }
  
  &.corner-bl {
    bottom: 0;
    left: 0;
    border-right: none;
    border-top: none;
  }
  
  &.corner-br {
    bottom: 0;
    right: 0;
    border-left: none;
    border-top: none;
  }
}

.scan-line {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #ff6b6b, transparent);
  animation: scan-line 2s linear infinite;
}

.scan-tips {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.7);
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  backdrop-filter: blur(10px);
}

// 扫码操作
.scanner-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

.action-btn {
  padding: 12px 20px;
  border: none;
  border-radius: 24px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  
  &.torch-btn {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: 2px solid rgba(255, 255, 255, 0.3);
    
    &:hover {
      background: rgba(255, 255, 255, 0.3);
    }
  }
  
  &.manual-btn {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
    }
  }
}

.btn-icon {
  font-size: 16px;
}

// 扫描历史
.history-section {
  padding: 20px;
  background: white;
  color: #333;
  border-radius: 20px 20px 0 0;
  margin-top: 20px;
  min-height: 300px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
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

.section-clear {
  font-size: 14px;
  color: #999;
  cursor: pointer;
  
  &:hover {
    color: #ff6b6b;
  }
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: #e9ecef;
    transform: translateY(-1px);
  }
}

.history-avatar {
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

.history-info {
  flex: 1;
}

.history-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 2px;
}

.history-service {
  font-size: 13px;
  color: #666;
  margin-bottom: 2px;
}

.history-time {
  font-size: 12px;
  color: #999;
}

.history-status {
  display: flex;
  align-items: center;
}

.status-icon {
  font-size: 16px;
  color: #52c41a;
}

// 手动输入弹窗
.manual-dialog,
.success-dialog {
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
}

.dialog-content {
  padding: 20px;
}

.input-group,
.service-selection {
  margin-bottom: 20px;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.input-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.code-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #f0f0f0;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: #667eea;
  }
}

.service-options {
  display: flex;
  gap: 12px;
}

.service-option {
  flex: 1;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 2px solid transparent;
  
  &.active {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    border-color: #667eea;
  }
}

.service-icon {
  display: block;
  font-size: 24px;
  margin-bottom: 8px;
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
  
  &.cancel-btn {
    background: #f5f5f5;
    color: #666;
  }
  
  &.verify-btn {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
  }
}

// 成功弹窗
.success-dialog {
  text-align: center;
}

.success-icon {
  font-size: 64px;
  margin: 20px 0;
}

.success-title {
  font-size: 20px;
  font-weight: 700;
  color: #333;
  margin-bottom: 20px;
}

.success-content {
  padding: 0 20px 20px;
}

.member-info {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 16px;
  text-align: left;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  
  &:last-child {
    margin-bottom: 0;
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

.success-actions {
  padding: 20px;
  border-top: 1px solid #f0f0f0;
}

.success-btn {
  width: 100%;
  height: 44px;
  background: linear-gradient(135deg, #52c41a, #73d13d);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
}

// 动画效果
@keyframes scan-line {
  0% { top: 0; }
  100% { top: 100%; }
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

.history-section {
  animation: slideInUp 0.6s ease-out;
}

.history-item {
  animation: slideInUp 0.4s ease-out;
}

.history-item:nth-child(1) { animation-delay: 0.1s; }
.history-item:nth-child(2) { animation-delay: 0.2s; }
.history-item:nth-child(3) { animation-delay: 0.3s; }

// 响应式设计
@media (max-width: 375px) {
  .camera-container {
    width: 240px;
    height: 240px;
  }
  
  .scan-frame {
    width: 160px;
    height: 160px;
  }
  
  .scanner-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .action-btn {
    width: 200px;
  }
}
</style>
