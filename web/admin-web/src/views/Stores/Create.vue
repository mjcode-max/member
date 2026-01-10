<template>
  <div class="store-form-page">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>创建门店</h2>
    </div>
    
    <el-card class="form-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        @submit.prevent="handleSubmit"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="门店名称" prop="name">
              <el-input
                v-model="form.name"
                placeholder="请输入门店名称"
                maxlength="50"
                show-word-limit
              />
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input
                v-model="form.phone"
                placeholder="请输入联系电话"
                maxlength="20"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="门店地址" prop="address">
          <el-input
            v-model="form.address"
            placeholder="请输入门店地址"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="联系人" prop="contact_person">
              <el-input
                v-model="form.contact_person"
                placeholder="请输入联系人"
                maxlength="50"
              />
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="押金金额" prop="deposit_amount">
              <el-input-number
                v-model="form.deposit_amount"
                :precision="2"
                :min="0"
                placeholder="请输入押金金额"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="营业开始时间" prop="business_hours_start">
              <el-time-picker
                v-model="form.business_hours_start"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="选择开始时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="营业结束时间" prop="business_hours_end">
              <el-time-picker
                v-model="form.business_hours_end"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="选择结束时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="时段模板" prop="template_id">
          <el-select
            v-model="form.template_id"
            placeholder="请选择时段模板（可选）"
            clearable
            filterable
            style="width: 100%"
            :loading="loadingTemplates"
          >
            <el-option
              v-for="template in templates"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            >
              <span>{{ template.name }}</span>
              <el-tag
                :type="template.status === 'active' ? 'success' : 'info'"
                size="small"
                style="margin-left: 8px"
              >
                {{ template.status === 'active' ? '启用' : '禁用' }}
              </el-tag>
            </el-option>
          </el-select>
          <div style="color: #909399; font-size: 12px; margin-top: 4px">
            选择时段模板后，系统将根据模板自动生成可预约时段
          </div>
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio value="operating">营业中</el-radio>
            <el-radio value="closed">停业</el-radio>
            <el-radio value="shutdown">关闭</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            创建门店
          </el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createStore } from '@/api/stores'
import { getTemplates } from '@/api/templates'

const router = useRouter()
const route = useRoute()

const formRef = ref()
const loading = ref(false)
const loadingTemplates = ref(false)
const templates = ref([])

const form = reactive({
  name: '',
  address: '',
  phone: '',
  contact_person: '',
  business_hours_start: '09:00',
  business_hours_end: '21:00',
  status: 'operating',
  deposit_amount: 0,
  template_id: null
})

const rules = {
  name: [
    { required: true, message: '请输入门店名称', trigger: 'blur' },
    { min: 2, max: 50, message: '门店名称长度在2到50个字符', trigger: 'blur' }
  ],
  address: [
    { required: true, message: '请输入门店地址', trigger: 'blur' },
    { min: 5, max: 200, message: '门店地址长度在5到200个字符', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' }
  ],
  business_hours_start: [
    { required: true, message: '请选择营业开始时间', trigger: 'change' }
  ],
  business_hours_end: [
    { required: true, message: '请选择营业结束时间', trigger: 'change' }
  ]
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    await createStore(form)
    ElMessage.success('门店创建成功')
    router.push('/stores')
  } catch (error) {
    console.error('创建门店失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取时段模板列表
// 注意：创建门店时，模板需要先创建，所以这里可以获取所有模板
// 或者留空，让用户创建门店后再编辑选择模板
const fetchTemplates = async () => {
  loadingTemplates.value = true
  try {
    // 尝试获取所有模板（如果后端支持不传store_id）
    // 如果不支持，则留空，用户可以在创建门店后编辑选择模板
    try {
      const response = await getTemplates({})
      if (response.data) {
        templates.value = Array.isArray(response.data) ? response.data : []
      }
    } catch (error) {
      // 如果获取失败，留空，不显示错误（模板是可选的）
      templates.value = []
    }
  } catch (error) {
    console.error('获取时段模板列表失败:', error)
    templates.value = []
  } finally {
    loadingTemplates.value = false
  }
}

// 返回
const goBack = () => {
  router.push('/stores')
}

onMounted(() => {
  fetchTemplates()
})
</script>

<style lang="scss" scoped>
.store-form-page {
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
    max-width: 800px;
  }
}
</style>
