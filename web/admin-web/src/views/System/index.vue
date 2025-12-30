<template>
  <div class="system-page">
    <div class="page-header">
      <h2>系统设置</h2>
    </div>
    
    <el-row :gutter="20">
      <!-- 系统配置 -->
      <el-col :span="12">
        <el-card class="config-card">
          <template #header>
            <div class="card-header">
              <span>系统配置</span>
            </div>
          </template>
          
          <el-form :model="systemConfig" label-width="120px">
            <el-form-item label="预约押金">
              <el-input-number
                v-model="systemConfig.booking_deposit_amount"
                :min="0"
                :max="10000"
                :step="100"
                controls-position="right"
              />
              <span class="unit">分</span>
            </el-form-item>
            
            <el-form-item label="支付超时时间">
              <el-input-number
                v-model="systemConfig.payment_timeout_minutes"
                :min="5"
                :max="60"
                controls-position="right"
              />
              <span class="unit">分钟</span>
            </el-form-item>
            
            <el-form-item label="取消限制时间">
              <el-input-number
                v-model="systemConfig.cancel_limit_hours"
                :min="1"
                :max="24"
                controls-position="right"
              />
              <span class="unit">小时</span>
            </el-form-item>
            
            <el-form-item label="时间段长度">
              <el-input-number
                v-model="systemConfig.time_slot_minutes"
                :min="30"
                :max="120"
                :step="30"
                controls-position="right"
              />
              <span class="unit">分钟</span>
            </el-form-item>
            
            <el-form-item label="会员码刷新间隔">
              <el-input-number
                v-model="systemConfig.member_code_refresh_seconds"
                :min="5"
                :max="60"
                controls-position="right"
              />
              <span class="unit">秒</span>
            </el-form-item>
            
            <el-form-item label="人脸相似度阈值">
              <el-input-number
                v-model="systemConfig.face_similarity_threshold"
                :min="0.1"
                :max="1.0"
                :step="0.1"
                :precision="1"
                controls-position="right"
              />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveSystemConfig" :loading="saving">
                保存配置
              </el-button>
              <el-button @click="resetSystemConfig">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      
      <!-- 系统信息 -->
      <el-col :span="12">
        <el-card class="info-card">
          <template #header>
            <div class="card-header">
              <span>系统信息</span>
            </div>
          </template>
          
          <el-descriptions :column="1" border>
            <el-descriptions-item label="系统版本">v1.0.0</el-descriptions-item>
            <el-descriptions-item label="构建时间">2023-10-15 10:30:00</el-descriptions-item>
            <el-descriptions-item label="运行环境">生产环境</el-descriptions-item>
            <el-descriptions-item label="数据库版本">MySQL 8.0</el-descriptions-item>
            <el-descriptions-item label="Redis版本">Redis 6.0</el-descriptions-item>
            <el-descriptions-item label="服务器时间">
              {{ currentTime }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
        
        <!-- 操作日志 -->
        <el-card class="log-card" style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>操作日志</span>
              <el-button size="small" @click="refreshLogs">刷新</el-button>
            </div>
          </template>
          
          <el-timeline>
            <el-timeline-item
              v-for="log in operationLogs"
              :key="log.id"
              :timestamp="log.timestamp"
              :type="log.type"
            >
              {{ log.content }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const saving = ref(false)
const currentTime = ref('')

const systemConfig = reactive({
  booking_deposit_amount: 1000,
  payment_timeout_minutes: 15,
  cancel_limit_hours: 3,
  time_slot_minutes: 60,
  member_code_refresh_seconds: 10,
  face_similarity_threshold: 0.8
})

const operationLogs = ref([
  {
    id: 1,
    timestamp: '2023-10-15 10:30:00',
    type: 'success',
    content: '系统配置已更新'
  },
  {
    id: 2,
    timestamp: '2023-10-15 10:25:00',
    type: 'info',
    content: '用户 admin 登录系统'
  },
  {
    id: 3,
    timestamp: '2023-10-15 10:20:00',
    type: 'warning',
    content: '门店 王府井旗舰店 状态已更新'
  },
  {
    id: 4,
    timestamp: '2023-10-15 10:15:00',
    type: 'success',
    content: '新增门店 朝阳门店'
  }
])

let timeInterval = null

// 更新当前时间
const updateCurrentTime = () => {
  currentTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss')
}

// 保存系统配置
const saveSystemConfig = async () => {
  saving.value = true
  try {
    // 这里应该调用API保存配置
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    
    ElMessage.success('系统配置保存成功')
    
    // 添加到操作日志
    operationLogs.value.unshift({
      id: Date.now(),
      timestamp: dayjs().format('YYYY-MM-DD HH:mm:ss'),
      type: 'success',
      content: '系统配置已更新'
    })
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 重置系统配置
const resetSystemConfig = () => {
  Object.assign(systemConfig, {
    booking_deposit_amount: 1000,
    payment_timeout_minutes: 15,
    cancel_limit_hours: 3,
    time_slot_minutes: 60,
    member_code_refresh_seconds: 10,
    face_similarity_threshold: 0.8
  })
  ElMessage.info('配置已重置')
}

// 刷新日志
const refreshLogs = () => {
  ElMessage.info('日志已刷新')
}

onMounted(() => {
  updateCurrentTime()
  timeInterval = setInterval(updateCurrentTime, 1000)
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style lang="scss" scoped>
.system-page {
  .page-header {
    margin-bottom: 20px;
    
    h2 {
      color: #333;
      margin: 0;
    }
  }
  
  .config-card, .info-card, .log-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: bold;
      color: #333;
    }
  }
  
  .unit {
    margin-left: 8px;
    color: #666;
    font-size: 14px;
  }
}
</style>
