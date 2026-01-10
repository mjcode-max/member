<template>
  <div class="member-detail-page">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>会员详情</h2>
    </div>
    
    <el-card class="detail-card" v-loading="loading">
      <div v-if="memberInfo.id" class="member-info">
        <!-- 基本信息 -->
        <el-descriptions title="基本信息" :column="2" border>
          <el-descriptions-item label="会员姓名">
            {{ memberInfo.name }}
          </el-descriptions-item>
          <el-descriptions-item label="手机号">
            {{ memberInfo.phone }}
          </el-descriptions-item>
          <el-descriptions-item label="所属门店">
            {{ getStoreName(memberInfo.store_id) }}
          </el-descriptions-item>
          <el-descriptions-item label="会员状态">
            <el-tag :type="getStatusType(memberInfo.status)">
              {{ getStatusText(memberInfo.status) }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <!-- 套餐信息 -->
        <el-descriptions title="套餐信息" :column="2" border style="margin-top: 20px">
          <el-descriptions-item label="套餐名称">
            {{ memberInfo.package_name }}
          </el-descriptions-item>
          <el-descriptions-item label="服务类型">
            {{ getServiceTypeText(memberInfo.service_type) }}
          </el-descriptions-item>
          <el-descriptions-item label="套餐价格">
            ¥{{ memberInfo.price?.toFixed(2) || '0.00' }}
          </el-descriptions-item>
          <el-descriptions-item label="购买金额">
            ¥{{ memberInfo.purchase_amount?.toFixed(2) || '0.00' }}
          </el-descriptions-item>
          <el-descriptions-item label="已使用次数">
            <el-tag type="info">{{ memberInfo.used_times || 0 }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="购买时间">
            {{ formatDateTime(memberInfo.purchase_time) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <!-- 有效期信息 -->
        <el-descriptions title="有效期信息" :column="2" border style="margin-top: 20px">
          <el-descriptions-item label="有效期开始">
            {{ formatDate(memberInfo.valid_from) }}
          </el-descriptions-item>
          <el-descriptions-item label="有效期结束">
            {{ formatDate(memberInfo.valid_to) }}
          </el-descriptions-item>
          <el-descriptions-item label="有效期时长">
            {{ memberInfo.validity_duration }} 天
          </el-descriptions-item>
          <el-descriptions-item label="剩余天数">
            <el-tag :type="getRemainingDaysType(remainingDays)">
              {{ remainingDays > 0 ? `${remainingDays} 天` : '已过期' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <!-- 备注信息 -->
        <el-descriptions title="备注信息" :column="1" border style="margin-top: 20px" v-if="memberInfo.description">
          <el-descriptions-item label="备注">
            {{ memberInfo.description }}
          </el-descriptions-item>
        </el-descriptions>
        
        <!-- 操作按钮 -->
        <div class="action-buttons" style="margin-top: 30px">
          <el-button type="primary" @click="handleEdit">编辑会员</el-button>
          <el-button type="success" @click="handleViewUsages">查看使用记录</el-button>
          <el-button 
            v-if="memberInfo.status === 'active'"
            type="warning" 
            @click="handleDisable"
          >
            停用会员
          </el-button>
          <el-button 
            v-if="memberInfo.status === 'inactive'"
            type="success" 
            @click="handleEnable"
          >
            启用会员
          </el-button>
        </div>
      </div>
      
      <div v-else class="empty-state">
        <el-empty description="会员信息不存在" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getMemberById, disableMember, enableMember } from '@/api/members'
import { getStores } from '@/api/stores'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const memberInfo = reactive({})
const storeList = ref([])

// 计算剩余天数
const remainingDays = computed(() => {
  if (!memberInfo.valid_to) return 0
  const endDate = dayjs(memberInfo.valid_to)
  const today = dayjs()
  const days = endDate.diff(today, 'day')
  return days > 0 ? days : 0
})

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

// 格式化日期时间
const formatDateTime = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    active: 'success',
    expired: 'warning',
    inactive: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const textMap = {
    active: '有效',
    expired: '已过期',
    inactive: '已停用'
  }
  return textMap[status] || '未知'
}

// 获取服务类型文本
const getServiceTypeText = (serviceType) => {
  const typeMap = {
    nail: '美甲',
    eyelash: '美睫',
    combo: '组合'
  }
  return typeMap[serviceType] || serviceType || '-'
}

// 获取门店名称
const getStoreName = (storeId) => {
  const store = storeList.value.find(s => s.id === storeId)
  return store ? store.name : '-'
}

// 获取剩余天数类型
const getRemainingDaysType = (days) => {
  if (days > 30) return 'success'
  if (days > 7) return 'warning'
  return 'danger'
}

// 获取会员详情
const fetchMemberDetail = async () => {
  const memberId = route.params.id
  if (!memberId) {
    ElMessage.error('会员ID不存在')
    router.back()
    return
  }
  
  loading.value = true
  try {
    const response = await getMemberById(memberId)
    if (response.data) {
      Object.assign(memberInfo, response.data)
    } else {
      ElMessage.error('会员不存在')
      router.back()
    }
  } catch (error) {
    console.error('获取会员详情失败:', error)
    ElMessage.error('获取会员详情失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 获取门店列表
const fetchStoreList = async () => {
  try {
    const response = await getStores({ page: 1, page_size: 1000 })
    if (response.code === 0 && response.data) {
      storeList.value = response.data.list || []
    }
  } catch (error) {
    console.error('获取门店列表失败:', error)
  }
}

// 返回
const goBack = () => {
  router.back()
}

// 编辑会员
const handleEdit = () => {
  router.push(`/members/${memberInfo.id}/edit`)
}

// 查看使用记录
const handleViewUsages = () => {
  router.push(`/members/${memberInfo.id}/usages`)
}

// 停用会员
const handleDisable = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要停用会员"${memberInfo.name}"吗？`,
      '停用会员',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await disableMember(memberInfo.id)
    ElMessage.success('会员停用成功')
    fetchMemberDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('停用会员失败:', error)
      ElMessage.error('停用会员失败')
    }
  }
}

// 启用会员
const handleEnable = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要启用会员"${memberInfo.name}"吗？`,
      '启用会员',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await enableMember(memberInfo.id)
    ElMessage.success('会员启用成功')
    fetchMemberDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('启用会员失败:', error)
      ElMessage.error('启用会员失败')
    }
  }
}

onMounted(() => {
  fetchStoreList()
  fetchMemberDetail()
})
</script>

<style lang="scss" scoped>
.member-detail-page {
  .page-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    
    h2 {
      color: #333;
      margin: 0;
    }
  }
  
  .detail-card {
    .member-info {
      .action-buttons {
        display: flex;
        gap: 12px;
      }
    }
    
    .empty-state {
      padding: 40px 0;
    }
  }
}
</style>

