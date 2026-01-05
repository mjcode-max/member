<template>
  <div class="member-form-page">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>{{ isEdit ? '编辑会员' : '新增会员' }}</h2>
    </div>
    
    <el-card class="form-card" v-loading="loading">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        @submit.prevent="handleSubmit"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="会员姓名" prop="name">
              <el-input
                v-model="form.name"
                placeholder="请输入会员姓名"
                maxlength="50"
              />
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input
                v-model="form.phone"
                placeholder="请输入手机号"
                maxlength="20"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属门店" prop="store_id">
              <el-select
                v-model="form.store_id"
                placeholder="请选择门店"
                style="width: 100%"
              >
                <el-option
                  v-for="store in storeList"
                  :key="store.id"
                  :label="store.name"
                  :value="store.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="服务类型" prop="service_type">
              <el-select
                v-model="form.service_type"
                placeholder="请选择服务类型"
                style="width: 100%"
              >
                <el-option label="美甲" value="nail" />
                <el-option label="美睫" value="eyelash" />
                <el-option label="组合" value="combo" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="套餐名称" prop="package_name">
              <el-input
                v-model="form.package_name"
                placeholder="请输入套餐名称"
                maxlength="100"
              />
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="套餐价格" prop="price">
              <el-input-number
                v-model="form.price"
                :precision="2"
                :min="0"
                placeholder="请输入套餐价格"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="购买金额" prop="purchase_amount">
              <el-input-number
                v-model="form.purchase_amount"
                :precision="2"
                :min="0"
                placeholder="请输入购买金额"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="会员状态" prop="status">
              <el-select
                v-model="form.status"
                placeholder="请选择状态"
                style="width: 100%"
              >
                <el-option label="有效" value="active" />
                <el-option label="已过期" value="expired" />
                <el-option label="已停用" value="inactive" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <!-- 有效期设置 -->
        <el-divider>有效期设置</el-divider>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="有效期开始" prop="valid_from">
              <el-date-picker
                v-model="form.valid_from"
                type="date"
                placeholder="选择开始日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                @change="handleValidFromChange"
              />
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="有效期结束" prop="valid_to">
              <el-date-picker
                v-model="form.valid_to"
                type="date"
                placeholder="选择结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                @change="handleValidToChange"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="固定时长（天）" prop="validity_duration">
              <el-input-number
                v-model="form.validity_duration"
                :min="1"
                placeholder="请输入固定时长天数"
                style="width: 100%"
                @change="handleDurationChange"
              />
              <div class="form-tip">输入固定时长后，系统会自动计算结束日期；选择开始/结束日期后，系统会自动计算固定时长</div>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="备注" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '保存修改' : '创建会员' }}
          </el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getMemberById, createMember, updateMember } from '@/api/members'
import { getStores } from '@/api/stores'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()

const formRef = ref()
const loading = ref(false)
const submitting = ref(false)
const storeList = ref([])
const isUpdatingDuration = ref(false) // 防止循环更新

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  name: '',
  phone: '',
  store_id: null,
  package_name: '',
  service_type: '',
  price: 0,
  purchase_amount: 0,
  status: 'active',
  validity_duration: null,
  valid_from: dayjs().format('YYYY-MM-DD'),
  valid_to: '',
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入会员姓名', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  store_id: [
    { required: true, message: '请选择所属门店', trigger: 'change' }
  ],
  package_name: [
    { required: true, message: '请输入套餐名称', trigger: 'blur' }
  ],
  service_type: [
    { required: true, message: '请选择服务类型', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入套餐价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '套餐价格不能为负数', trigger: 'blur' }
  ],
  purchase_amount: [
    { type: 'number', min: 0, message: '购买金额不能为负数', trigger: 'blur' }
  ],
  valid_from: [
    { required: true, message: '请选择有效期开始日期', trigger: 'change' }
  ],
  valid_to: [
    { required: true, message: '请选择有效期结束日期', trigger: 'change' }
  ]
}

// 处理固定时长变化
const handleDurationChange = (value) => {
  if (isUpdatingDuration.value) return
  if (value && form.valid_from) {
    const startDate = dayjs(form.valid_from)
    form.valid_to = startDate.add(value, 'day').format('YYYY-MM-DD')
  }
}

// 处理有效期开始日期变化
const handleValidFromChange = (value) => {
  if (isUpdatingDuration.value) return
  // 如果开始日期和结束日期都已选择，计算固定时长
  if (value && form.valid_to) {
    const startDate = dayjs(value)
    const endDate = dayjs(form.valid_to)
    if (endDate.isAfter(startDate)) {
      isUpdatingDuration.value = true
      form.validity_duration = endDate.diff(startDate, 'day')
      isUpdatingDuration.value = false
    }
  }
  // 如果固定时长已设置，更新结束日期
  else if (value && form.validity_duration) {
    const startDate = dayjs(value)
    form.valid_to = startDate.add(form.validity_duration, 'day').format('YYYY-MM-DD')
  }
}

// 处理有效期结束日期变化
const handleValidToChange = (value) => {
  if (isUpdatingDuration.value) return
  // 如果开始日期和结束日期都已选择，计算固定时长
  if (value && form.valid_from) {
    const startDate = dayjs(form.valid_from)
    const endDate = dayjs(value)
    if (endDate.isAfter(startDate)) {
      isUpdatingDuration.value = true
      form.validity_duration = endDate.diff(startDate, 'day')
      isUpdatingDuration.value = false
    } else if (endDate.isBefore(startDate)) {
      ElMessage.warning('结束日期不能早于开始日期')
      form.valid_to = ''
    }
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

// 获取会员详情
const fetchMemberDetail = async () => {
  const memberId = route.params.id
  if (!memberId) return
  
  loading.value = true
  try {
    const response = await getMemberById(memberId)
    if (response.data) {
      const data = response.data
      Object.assign(form, {
        name: data.name || '',
        phone: data.phone || '',
        store_id: data.store_id || null,
        package_name: data.package_name || '',
        service_type: data.service_type || '',
        price: data.price || 0,
        purchase_amount: data.purchase_amount || 0,
        status: data.status || 'active',
        validity_duration: data.validity_duration || null,
        valid_from: data.valid_from ? dayjs(data.valid_from).format('YYYY-MM-DD') : '',
        valid_to: data.valid_to ? dayjs(data.valid_to).format('YYYY-MM-DD') : '',
        description: data.description || ''
      })
    }
  } catch (error) {
    console.error('获取会员详情失败:', error)
    ElMessage.error('获取会员详情失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    // 准备提交数据
    const data = {
      name: form.name,
      phone: form.phone,
      store_id: form.store_id,
      package_name: form.package_name,
      service_type: form.service_type,
      price: form.price,
      purchase_amount: form.purchase_amount || form.price,
      status: form.status,
      description: form.description || undefined
    }
    
    // 有效期处理：优先使用日期，如果提供了固定时长则使用固定时长
    if (form.validity_duration && form.valid_from) {
      data.validity_duration = form.validity_duration
      // 将日期字符串转换为RFC3339格式（Gin支持），使用本地时区
      const startDate = dayjs(form.valid_from).startOf('day')
      data.valid_from = startDate.toISOString()
    } else if (form.valid_from && form.valid_to) {
      // 将日期字符串转换为RFC3339格式
      const startDate = dayjs(form.valid_from).startOf('day')
      const endDate = dayjs(form.valid_to).startOf('day')
      data.valid_from = startDate.toISOString()
      data.valid_to = endDate.toISOString()
    }
    
    if (isEdit.value) {
      await updateMember(route.params.id, data)
      ElMessage.success('会员更新成功')
    } else {
      await createMember(data)
      ElMessage.success('会员创建成功')
    }
    
    router.push('/members')
  } catch (error) {
    console.error('保存会员失败:', error)
  } finally {
    submitting.value = false
  }
}

// 返回
const goBack = () => {
  router.push('/members')
}

onMounted(() => {
  fetchStoreList()
  if (isEdit.value) {
    fetchMemberDetail()
  } else {
    // 新建时设置默认值
    form.valid_from = dayjs().format('YYYY-MM-DD')
  }
})
</script>

<style lang="scss" scoped>
.member-form-page {
  .page-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      color: #333;
      margin: 0 0 0 16px;
    }
  }
  
  .form-card {
    max-width: 1000px;
    
    .form-tip {
      font-size: 12px;
      color: #909399;
      margin-top: 4px;
    }
  }
}
</style>

