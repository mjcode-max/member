<template>
  <div class="home-page">
    <!-- 状态栏占位 -->
    <div class="status-bar"></div>
    
    <!-- 头部区域 -->
    <div class="header-section">
      <div class="header-content">
        <div class="brand-info">
          <h1 class="brand-title">奈斯美甲</h1>
          <p class="brand-subtitle">专业美甲美睫连锁服务</p>
        </div>
        <div class="header-actions">
          <div class="location-info" @click="handleLocation">
            <i class="location-icon">📍</i>
            <span>北京市</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 轮播图 -->
    <div class="banner-section">
      <van-swipe :autoplay="3000" indicator-color="rgba(255,255,255,0.8)" class="custom-swipe">
        <van-swipe-item>
          <div class="banner-item banner-1">
            <div class="banner-overlay"></div>
            <div class="banner-content">
              <div class="banner-icon">💅</div>
              <h3>专业美甲服务</h3>
              <p>精致工艺，让您的手部更加美丽动人</p>
              <div class="banner-btn" @click="$router.push('/booking/create?service=manicure')">
                立即预约
              </div>
            </div>
          </div>
        </van-swipe-item>
        <van-swipe-item>
          <div class="banner-item banner-2">
            <div class="banner-overlay"></div>
            <div class="banner-content">
              <div class="banner-icon">👁️</div>
              <h3>精致美睫体验</h3>
              <p>专业技师，让您的眼睛更加迷人动感</p>
              <div class="banner-btn" @click="$router.push('/booking/create?service=eyelash')">
                立即预约
              </div>
            </div>
          </div>
        </van-swipe-item>
        <van-swipe-item>
          <div class="banner-item banner-3">
            <div class="banner-overlay"></div>
            <div class="banner-content">
              <div class="banner-icon">🎁</div>
              <h3>会员专享优惠</h3>
              <p>更多专属优惠和贴心服务等您来享</p>
              <div class="banner-btn" @click="$router.push('/member')">
                查看会员
              </div>
            </div>
          </div>
        </van-swipe-item>
      </van-swipe>
    </div>
    
    <!-- 快捷功能 -->
    <div class="quick-actions">
      <div class="section-title">
        <span class="title-text">快捷服务</span>
        <span class="title-desc">一键直达您需要的服务</span>
      </div>
      <div class="action-grid">
        <div class="action-item" @click="$router.push('/booking/create')">
          <div class="action-icon action-icon-1">
            <i>📅</i>
          </div>
          <div class="action-text">立即预约</div>
          <div class="action-desc">快速预约服务</div>
        </div>
        <div class="action-item" @click="$router.push('/member')">
          <div class="action-icon action-icon-2">
            <i>💎</i>
          </div>
          <div class="action-text">我的会员</div>
          <div class="action-desc">查看会员信息</div>
        </div>
        <div class="action-item" @click="showStoreList = true">
          <div class="action-icon action-icon-3">
            <i>🏪</i>
          </div>
          <div class="action-text">门店查询</div>
          <div class="action-desc">查找附近门店</div>
        </div>
        <div class="action-item" @click="handleContact">
          <div class="action-icon action-icon-4">
            <i>📞</i>
          </div>
          <div class="action-text">联系我们</div>
          <div class="action-desc">客服热线咨询</div>
        </div>
      </div>
    </div>
    
    <!-- 门店列表 -->
    <div class="stores-section">
      <div class="section-title">
        <span class="title-text">附近门店</span>
        <span class="title-more" @click="showStoreList = true">查看全部 ›</span>
      </div>
      <div class="stores-list">
        <div 
          v-for="store in nearbyStores" 
          :key="store.id"
          class="store-card"
          @click="handleStoreClick(store)"
        >
          <div class="store-info">
            <div class="store-header">
              <h4 class="store-name">{{ store.name }}</h4>
              <div class="store-status" :class="{ active: store.status === 'active' }">
                {{ store.status === 'active' ? '营业中' : '已停业' }}
              </div>
            </div>
            <p class="store-address">{{ store.address }}</p>
            <div class="store-meta">
              <span class="store-distance">距离 {{ store.distance }}km</span>
              <span class="store-rating">⭐ 4.8</span>
            </div>
          </div>
          <div class="store-arrow">›</div>
        </div>
      </div>
    </div>

    <!-- 服务特色 -->
    <div class="features-section">
      <div class="section-title">
        <span class="title-text">服务特色</span>
        <span class="title-desc">专业品质，贴心服务</span>
      </div>
      <div class="features-grid">
        <div class="feature-item">
          <div class="feature-icon">✨</div>
          <div class="feature-title">专业技师</div>
          <div class="feature-desc">持证上岗，技艺精湛</div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">🛡️</div>
          <div class="feature-title">安全卫生</div>
          <div class="feature-desc">一客一换，严格消毒</div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">💖</div>
          <div class="feature-title">贴心服务</div>
          <div class="feature-desc">用心服务，满意保障</div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">⏰</div>
          <div class="feature-title">便捷预约</div>
          <div class="feature-desc">随时预约，准时服务</div>
        </div>
      </div>
    </div>
    
    <!-- 门店列表弹窗 -->
    <van-popup 
      v-model:show="showStoreList" 
      position="bottom" 
      round
      :style="{ height: '70%' }"
    >
      <div class="store-list-popup">
        <div class="popup-header">
          <div class="popup-title">选择门店</div>
          <div class="popup-close" @click="showStoreList = false">×</div>
        </div>
        <div class="popup-content">
          <div 
            v-for="store in allStores" 
            :key="store.id"
            class="popup-store-item"
            @click="handleStoreSelect(store)"
          >
            <div class="popup-store-info">
              <h4 class="popup-store-name">{{ store.name }}</h4>
              <p class="popup-store-address">{{ store.address }}</p>
            </div>
            <div class="popup-store-arrow">›</div>
          </div>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import Swipe from 'vant/es/swipe'
import SwipeItem from 'vant/es/swipe-item'
import Popup from 'vant/es/popup'
import { getStores } from '@/api/stores'

const router = useRouter()

const showStoreList = ref(false)
const nearbyStores = ref([])
const allStores = ref([])

// 获取门店列表
const fetchStores = async () => {
  try {
    // 获取所有门店（包括营业中、停业、关闭）
    const response = await getStores({
      status: 'operating', // 获取所有门店
      page: 1,
      page_size: 100
    })
    // 后端返回格式：{ code: 0, data: { list: [...], pagination: {...} } }
    if (response.data?.list) {
      allStores.value = response.data.list
    } else if (Array.isArray(response.data)) {
      allStores.value = response.data
    } else if (response.data?.stores) {
      allStores.value = response.data.stores
    } else {
      allStores.value = []
    }
    
    // 模拟附近门店（取前3个）
    nearbyStores.value = allStores.value.slice(0, 3).map(store => ({
      ...store,
      distance: (Math.random() * 5).toFixed(1) // 模拟距离
    }))
  } catch (error) {
    console.error('获取门店列表失败:', error)
  }
}

// 处理门店点击
const handleStoreClick = (store) => {
  router.push(`/booking/create?storeId=${store.id}`)
}

// 处理门店选择
const handleStoreSelect = (store) => {
  showStoreList.value = false
  router.push(`/booking/create?storeId=${store.id}`)
}

// 处理联系我们
const handleContact = () => {
  showToast('客服电话：400-123-4567')
}

// 处理定位
const handleLocation = () => {
  showToast('定位功能开发中...')
}

onMounted(() => {
  fetchStores()
})
</script>

<style lang="scss" scoped>
.home-page {
  background: #ffffff;
  min-height: 100vh;
}

// 状态栏占位
.status-bar {
  height: env(safe-area-inset-top, 20px);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

// 头部区域
.header-section {
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
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 4px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.brand-subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

.location-info {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.3);
  }
}

// 轮播图
.banner-section {
  margin: 0 16px 20px;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  
  .custom-swipe {
    border-radius: 16px;
    overflow: hidden;
  }
  
  .banner-item {
    height: 180px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    position: relative;
    
    &.banner-1 {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }
    
    &.banner-2 {
      background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    }
    
    &.banner-3 {
      background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    }
  }
  
  .banner-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.2);
  }
  
  .banner-content {
    text-align: center;
    position: relative;
    z-index: 1;
    padding: 0 20px;
  }
  
  .banner-icon {
    font-size: 32px;
    margin-bottom: 12px;
  }
  
  .banner-content h3 {
    font-size: 22px;
    margin-bottom: 8px;
    font-weight: 700;
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }
  
  .banner-content p {
    font-size: 14px;
    opacity: 0.95;
    margin-bottom: 16px;
    line-height: 1.4;
  }
  
  .banner-btn {
    display: inline-block;
    padding: 8px 20px;
    background: rgba(255, 255, 255, 0.2);
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 20px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    backdrop-filter: blur(10px);
    
    &:hover {
      background: rgba(255, 255, 255, 0.3);
      transform: translateY(-2px);
    }
  }
}

// 区块标题
.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  padding: 0 4px;
}

.title-text {
  font-size: 18px;
  font-weight: 700;
  color: #333;
}

.title-desc {
  font-size: 12px;
  color: #999;
}

.title-more {
  font-size: 14px;
  color: #667eea;
  cursor: pointer;
  
  &:hover {
    color: #764ba2;
  }
}

// 快捷功能
.quick-actions {
  margin: 0 16px 24px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.action-item {
  background: white;
  border-radius: 16px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.8);
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  }
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
  font-size: 20px;
  
  &.action-icon-1 {
    background: linear-gradient(135deg, #667eea, #764ba2);
  }
  
  &.action-icon-2 {
    background: linear-gradient(135deg, #f093fb, #f5576c);
  }
  
  &.action-icon-3 {
    background: linear-gradient(135deg, #4facfe, #00f2fe);
  }
  
  &.action-icon-4 {
    background: linear-gradient(135deg, #43e97b, #38f9d7);
  }
}

.action-text {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.action-desc {
  font-size: 12px;
  color: #999;
  line-height: 1.3;
}

// 门店列表
.stores-section {
  margin: 0 16px 24px;
}

.stores-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.store-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.8);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  }
}

.store-info {
  flex: 1;
}

.store-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.store-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.store-status {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  background: #f5f5f5;
  color: #999;
  
  &.active {
    background: #e8f5e8;
    color: #52c41a;
  }
}

.store-address {
  font-size: 14px;
  color: #666;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.store-meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.store-distance {
  font-size: 13px;
  color: #999;
}

.store-rating {
  font-size: 13px;
  color: #ffa726;
}

.store-arrow {
  font-size: 18px;
  color: #ccc;
  font-weight: bold;
  margin-left: 12px;
}

// 服务特色
.features-section {
  margin: 0 16px 24px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.feature-item {
  background: white;
  border-radius: 16px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.8);
}

.feature-icon {
  font-size: 24px;
  margin-bottom: 12px;
}

.feature-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 6px;
}

.feature-desc {
  font-size: 12px;
  color: #999;
  line-height: 1.3;
}

// 弹窗样式
.store-list-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: white;
}

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

.popup-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.popup-store-item {
  display: flex;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f8f8f8;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    background: #f8f9ff;
    margin: 0 -20px;
    padding: 16px 20px;
    border-radius: 8px;
  }
}

.popup-store-info {
  flex: 1;
}

.popup-store-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin: 0 0 4px 0;
}

.popup-store-address {
  font-size: 13px;
  color: #666;
  margin: 0;
}

.popup-store-arrow {
  font-size: 16px;
  color: #ccc;
  margin-left: 12px;
}

// 响应式设计
@media (max-width: 375px) {
  .brand-title {
    font-size: 24px;
  }
  
  .action-grid,
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .banner-content h3 {
    font-size: 20px;
  }
  
  .banner-content p {
    font-size: 13px;
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

.quick-actions,
.stores-section,
.features-section {
  animation: slideInUp 0.6s ease-out;
}

.quick-actions { animation-delay: 0.1s; }
.stores-section { animation-delay: 0.2s; }
.features-section { animation-delay: 0.3s; }

.action-item,
.store-card,
.feature-item {
  animation: slideInUp 0.4s ease-out;
}

.action-item:nth-child(1) { animation-delay: 0.1s; }
.action-item:nth-child(2) { animation-delay: 0.2s; }
.action-item:nth-child(3) { animation-delay: 0.3s; }
.action-item:nth-child(4) { animation-delay: 0.4s; }
</style>
