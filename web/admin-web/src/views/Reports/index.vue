<template>
  <div class="reports-page">
    <div class="page-header">
      <h2>报表中心</h2>
    </div>
    
    <!-- 报表筛选 -->
    <el-card class="filter-card">
      <el-form :model="filterForm" inline>
        <el-form-item label="报表类型">
          <el-select
            v-model="filterForm.type"
            placeholder="请选择报表类型"
            style="width: 150px"
            @change="handleFilterChange"
          >
            <el-option label="预约报表" value="booking" />
            <el-option label="会员报表" value="member" />
            <el-option label="收入报表" value="revenue" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="开始日期">
          <el-date-picker
            v-model="filterForm.start_date"
            type="date"
            placeholder="选择开始日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="handleFilterChange"
          />
        </el-form-item>
        
        <el-form-item label="结束日期">
          <el-date-picker
            v-model="filterForm.end_date"
            type="date"
            placeholder="选择结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="handleFilterChange"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch" :loading="loading">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
          <el-button type="success" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 报表内容 -->
    <el-card class="report-card">
      <template #header>
        <div class="card-header">
          <span>{{ getReportTitle() }}</span>
          <el-tag type="info">{{ reportData.report_type }}</el-tag>
        </div>
      </template>
      
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>
      
      <div v-else-if="reportData.data && reportData.data.length > 0">
        <!-- 预约报表 -->
        <div v-if="filterForm.type === 'booking'">
          <el-table :data="reportData.data" stripe>
            <el-table-column prop="date" label="日期" />
            <el-table-column prop="store_name" label="门店名称" />
            <el-table-column prop="booking_count" label="预约数量" />
            <el-table-column prop="completed_count" label="完成数量" />
            <el-table-column prop="cancelled_count" label="取消数量" />
            <el-table-column prop="completion_rate" label="完成率">
              <template #default="{ row }">
                {{ (row.completion_rate * 100).toFixed(1) }}%
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <!-- 会员报表 -->
        <div v-else-if="filterForm.type === 'member'">
          <el-table :data="reportData.data" stripe>
            <el-table-column prop="date" label="日期" />
            <el-table-column prop="store_name" label="门店名称" />
            <el-table-column prop="new_members" label="新增会员" />
            <el-table-column prop="active_members" label="活跃会员" />
            <el-table-column prop="consumption_count" label="消费次数" />
            <el-table-column prop="total_amount" label="消费金额">
              <template #default="{ row }">
                ¥{{ (row.total_amount / 100).toFixed(2) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <!-- 收入报表 -->
        <div v-else-if="filterForm.type === 'revenue'">
          <el-table :data="reportData.data" stripe>
            <el-table-column prop="date" label="日期" />
            <el-table-column prop="store_name" label="门店名称" />
            <el-table-column prop="booking_revenue" label="预约收入">
              <template #default="{ row }">
                ¥{{ (row.booking_revenue / 100).toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column prop="member_revenue" label="会员收入">
              <template #default="{ row }">
                ¥{{ (row.member_revenue / 100).toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column prop="total_revenue" label="总收入">
              <template #default="{ row }">
                ¥{{ (row.total_revenue / 100).toFixed(2) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
      
      <el-empty v-else description="暂无数据" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getReports } from '@/api/reports'
import dayjs from 'dayjs'

const loading = ref(false)
const reportData = ref({})

const filterForm = reactive({
  type: 'booking',
  start_date: dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
  end_date: dayjs().format('YYYY-MM-DD')
})

// 获取报表标题
const getReportTitle = () => {
  const typeMap = {
    booking: '预约报表',
    member: '会员报表',
    revenue: '收入报表'
  }
  return typeMap[filterForm.type] || '报表'
}

// 获取报表数据
const fetchReportData = async () => {
  if (!filterForm.type) {
    ElMessage.warning('请选择报表类型')
    return
  }
  
  loading.value = true
  try {
    const params = {
      type: filterForm.type,
      start_date: filterForm.start_date,
      end_date: filterForm.end_date
    }
    
    const response = await getReports(params)
    reportData.value = response.data
  } catch (error) {
    console.error('获取报表数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 筛选条件变化
const handleFilterChange = () => {
  // 可以在这里添加实时查询逻辑
}

// 搜索
const handleSearch = () => {
  fetchReportData()
}

// 重置
const handleReset = () => {
  Object.assign(filterForm, {
    type: 'booking',
    start_date: dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
    end_date: dayjs().format('YYYY-MM-DD')
  })
  fetchReportData()
}

// 导出
const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

onMounted(() => {
  fetchReportData()
})
</script>

<style lang="scss" scoped>
.reports-page {
  .page-header {
    margin-bottom: 20px;
    
    h2 {
      color: #333;
      margin: 0;
    }
  }
  
  .filter-card {
    margin-bottom: 20px;
  }
  
  .report-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: bold;
      color: #333;
    }
    
    .loading-container {
      padding: 20px;
    }
  }
}
</style>
