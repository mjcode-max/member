<template>
  <div class="store-info-page">
    <!-- 页面头部 -->
    <van-nav-bar
      title="门店信息"
      left-text="返回"
      left-arrow
      @click-left="$router.back()"
      class="custom-nav"
    />
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <van-loading type="spinner" size="24px" />
      <span class="loading-text">加载中...</span>
    </div>
    
    <!-- 门店信息 -->
    <div v-else-if="storeInfo.name" class="content-container">
      <!-- 门店基本信息卡片 -->
      <div class="info-card">
        <div class="card-header">
          <div class="store-name-large">{{ storeInfo.name }}</div>
          <div class="status-badge" :class="getStatusClass(storeInfo.status)">
            {{ getStatusText(storeInfo.status) }}
          </div>
        </div>
        
        <div class="info-item" v-if="storeInfo.address">
          <div class="info-label">
            <i class="info-icon">📍</i>
            <span>门店地址</span>
          </div>
          <div class="info-value">{{ storeInfo.address }}</div>
        </div>
        
        <div class="info-item" v-if="storeInfo.phone">
          <div class="info-label">
            <i class="info-icon">📞</i>
            <span>联系电话</span>
          </div>
          <div class="info-value">
            <a :href="`tel:${storeInfo.phone}`" class="phone-link">{{ storeInfo.phone }}</a>
          </div>
        </div>
        
        <div class="info-item" v-if="storeInfo.contact_person">
          <div class="info-label">
            <i class="info-icon">👤</i>
            <span>联系人</span>
          </div>
          <div class="info-value">{{ storeInfo.contact_person }}</div>
        </div>
      </div>
      
      <!-- 营业信息卡片 -->
      <div class="info-card" v-if="storeInfo.business_hours_start || storeInfo.deposit_amount > 0">
        <div class="card-title">营业信息</div>
        
        <div class="info-item" v-if="storeInfo.business_hours_start && storeInfo.business_hours_end">
          <div class="info-label">
            <i class="info-icon">🕐</i>
            <span>营业时间</span>
          </div>
          <div class="info-value">{{ storeInfo.business_hours_start }} - {{ storeInfo.business_hours_end }}</div>
        </div>
        
        <div class="info-item" v-if="storeInfo.deposit_amount > 0">
          <div class="info-label">
            <i class="info-icon">💰</i>
            <span>押金金额</span>
          </div>
          <div class="info-value">¥{{ storeInfo.deposit_amount }}</div>
        </div>
      </div>
      
      <!-- 操作按钮 -->
      <!-- <div class="action-section">
        <van-button 
          type="primary" 
          block 
          round
          @click="handleEdit"
          class="edit-btn"
        >
          编辑门店信息
        </van-button>
      </div> -->
    </div>
    
    <!-- 错误状态 -->
    <div v-else class="error-container">
      <div class="error-icon">⚠️</div>
      <div class="error-text">获取门店信息失败</div>
      <van-button size="small" type="primary" @click="fetchStoreInfo">重试</van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getStoreById } from '@/api/stores'
import { showFailToast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)

// 店铺信息
const storeInfo = reactive({
  id: null,
  name: '',
  address: '',
  phone: '',
  contact_person: '',
  status: '',
  business_hours_start: '',
  business_hours_end: '',
  deposit_amount: 0
})

// 获取门店状态文本
const getStatusText = (status) => {
  const statusMap = {
    operating: '营业中',
    closed: '停业',
    shutdown: '已关闭'
  }
  return statusMap[status] || '未知'
}

// 获取门店状态样式类
const getStatusClass = (status) => {
  return {
    'status-operating': status === 'operating',
    'status-closed': status === 'closed',
    'status-shutdown': status === 'shutdown'
  }
}

// 获取店铺详情
const fetchStoreInfo = async () => {
  // 从用户信息中获取store_id
  const storeId = userStore.userInfo?.store_id
  
  if (!storeId) {
    showFailToast('用户未关联门店，请联系管理员')
    return
  }
  
  loading.value = true
  try {
    const response = await getStoreById(storeId)
    if (response.data) {
      // 更新店铺信息
      Object.assign(storeInfo, {
        id: response.data.id,
        name: response.data.name || '',
        address: response.data.address || '',
        phone: response.data.phone || '',
        contact_person: response.data.contact_person || '',
        status: response.data.status || '',
        business_hours_start: response.data.business_hours_start || '',
        business_hours_end: response.data.business_hours_end || '',
        deposit_amount: response.data.deposit_amount || 0
      })
    }
  } catch (error) {
    console.error('获取店铺详情失败:', error)
    showFailToast('获取店铺信息失败')
  } finally {
    loading.value = false
  }
}

// 编辑门店信息
const handleEdit = () => {
  // 跳转到编辑页面（如果后续需要实现编辑功能）
  showFailToast('编辑功能开发中...')
  // router.push(`/store/edit/${storeInfo.id}`)
}

onMounted(() => {
  fetchStoreInfo()
})
</script>

<style lang="scss" scoped>
.store-info-page {
  background: linear-gradient(180deg, #f8f9ff 0%, #f0f2ff 100%);
  min-height: 100vh;
  padding-bottom: 20px;
}

.custom-nav {
  background: linear-gradient(135deg, #ff6b6b 0%, #ffa726 100%);
  
  :deep(.van-nav-bar__title) {
    color: white;
    font-weight: 600;
  }
  
  :deep(.van-nav-bar__text) {
    color: white;
  }
  
  :deep(.van-nav-bar__arrow) {
    color: white;
  }
}

// 加载状态
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  gap: 12px;
}

.loading-text {
  font-size: 14px;
  color: #666;
}

// 内容容器
.content-container {
  padding: 20px 16px;
}

// 信息卡片
.info-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.store-name-large {
  font-size: 22px;
  font-weight: 700;
  color: #333;
  flex: 1;
}

.status-badge {
  padding: 6px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  
  &.status-operating {
    background: #e6f7ff;
    color: #1890ff;
  }
  
  &.status-closed {
    background: #fff7e6;
    color: #fa8c16;
  }
  
  &.status-shutdown {
    background: #f5f5f5;
    color: #8c8c8c;
  }
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 12px 0;
  border-bottom: 1px solid #f8f8f8;
  
  &:last-child {
    border-bottom: none;
  }
}

.info-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
  min-width: 100px;
}

.info-icon {
  font-size: 16px;
}

.info-value {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  text-align: right;
  flex: 1;
  margin-left: 16px;
}

.phone-link {
  color: #1890ff;
  text-decoration: none;
  
  &:active {
    opacity: 0.7;
  }
}

// 操作按钮
.action-section {
  padding: 20px 0;
}

.edit-btn {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

// 错误状态
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  gap: 16px;
}

.error-icon {
  font-size: 48px;
  opacity: 0.5;
}

.error-text {
  font-size: 16px;
  color: #666;
  margin-bottom: 8px;
}

// 响应式设计
@media (max-width: 375px) {
  .info-card {
    padding: 16px;
  }
  
  .store-name-large {
    font-size: 20px;
  }
  
  .info-item {
    flex-direction: column;
    gap: 8px;
  }
  
  .info-value {
    text-align: left;
    margin-left: 0;
  }
}
</style>

