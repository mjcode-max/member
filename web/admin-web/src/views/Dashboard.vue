<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h2>数据看板</h2>
      <p>实时监控系统运营数据</p>
    </div>
    
    <!-- 数据统计卡片 -->
    <el-row :gutter="20" class="stats-cards" v-loading="loading" element-loading-text="正在加载数据...">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon booking">
              <el-icon><Calendar /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ dashboardData.totalBookings }}</div>
              <div class="stat-label">总预约数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon member">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ dashboardData.totalMembers }}</div>
              <div class="stat-label">会员总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon store">
              <el-icon><Shop /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ dashboardData.totalStores }}</div>
              <div class="stat-label">门店数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 门店数据表格 -->
    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span>门店数据</span>
        </div>
      </template>
      
      <el-table :data="dashboardData.storeStats" stripe>
        <el-table-column prop="store_name" label="门店名称" />
        <el-table-column prop="booking_count" label="预约数量" />
        <el-table-column prop="member_count" label="会员数量" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="viewStoreDetail(row.store_id)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { getAdminDashboard } from '@/api/reports'

const router = useRouter()

const dashboardData = ref({
  totalBookings: 0,
  totalMembers: 0,
  totalStores: 0,
  storeStats: []
})

const loading = ref(false)


// 格式化金额
const formatMoney = (amount) => {
  return (amount / 100).toFixed(2)
}

// 获取数据看板数据
const fetchDashboardData = async () => {
  loading.value = true
  try {
    console.log('开始获取Dashboard数据...')
    const response = await getAdminDashboard()
    console.log('Dashboard API Response:', response)
    console.log('Response data:', response?.data)
    
    // 确保数据结构正确
    if (response && response.data) {
      const data = response.data
      console.log('处理数据:', {
        totalBookings: data.total_bookings || data.totalBookings,
        totalMembers: data.total_members || data.totalMembers,
        totalStores: data.total_stores || data.totalStores,
        storeStatsLength: (data.store_stats || data.storeStats)?.length,
        storeStatsData: data.store_stats || data.storeStats
      })
      
      dashboardData.value = {
        totalBookings: data.total_bookings || data.totalBookings || 0,
        totalMembers: data.total_members || data.totalMembers || 0,
        totalStores: data.total_stores || data.totalStores || 0,
        storeStats: data.store_stats || data.storeStats || []
      }
      
      console.log('设置后的dashboardData:', dashboardData.value)
    } else {
      console.warn('API 返回数据格式异常:', response)
      // 如果API没有返回totalStores，尝试从storeStats计算
      if (response?.data?.storeStats) {
        dashboardData.value.totalStores = response.data.storeStats.length
        console.log('从storeStats计算门店数量:', dashboardData.value.totalStores)
      }
    }
    
    await nextTick()
  } catch (error) {
    console.error('获取数据看板失败:', error)
    console.error('错误详情:', error.response?.data || error.message)
    
    // 设置默认数据，避免页面崩溃
    dashboardData.value = {
      totalBookings: 0,
      totalMembers: 0,
      totalStores: 0,
      storeStats: []
    }
  } finally {
    loading.value = false
  }
}


// 查看门店详情
const viewStoreDetail = (storeId) => {
  router.push(`/stores/${storeId}/edit`)
}

onMounted(() => {
  fetchDashboardData()
})
</script>

<style lang="scss" scoped>
.dashboard {
  .dashboard-header {
    margin-bottom: 20px;
    
    h2 {
      color: #333;
      margin-bottom: 8px;
    }
    
    p {
      color: #666;
      font-size: 14px;
    }
  }
  
  .stats-cards {
    margin-bottom: 20px;
    
    .stat-card {
      .stat-content {
        display: flex;
        align-items: center;
        
        .stat-icon {
          width: 60px;
          height: 60px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 16px;
          font-size: 24px;
          color: white;
          
          &.booking {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          }
          
          &.member {
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
          }
          
          &.revenue {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
          }
          
          &.store {
            background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
          }
        }
        
        .stat-info {
          .stat-value {
            font-size: 24px;
            font-weight: bold;
            color: #333;
            margin-bottom: 4px;
          }
          
          .stat-label {
            font-size: 14px;
            color: #666;
          }
        }
      }
    }
  }
  
  .charts-section {
    margin-bottom: 20px;
    
    .chart-card {
      .chart-container {
        height: 300px;
      }
    }
  }
  
  .table-card {
    .card-header {
      font-weight: bold;
      color: #333;
    }
  }
}
</style>
