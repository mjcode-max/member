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
            <el-form-item label="纬度" prop="latitude">
              <el-input-number
                v-model="form.latitude"
                :precision="6"
                :step="0.000001"
                placeholder="请输入纬度"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          
          <el-col :span="12">
            <el-form-item label="经度" prop="longitude">
              <el-input-number
                v-model="form.longitude"
                :precision="6"
                :step="0.000001"
                placeholder="请输入经度"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="营业时间" prop="business_hours">
          <el-input
            v-model="form.business_hours"
            placeholder="请输入营业时间，如：09:00-21:00"
            maxlength="20"
          />
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createStore } from '@/api/stores'

const router = useRouter()

const formRef = ref()
const loading = ref(false)

const form = reactive({
  name: '',
  address: '',
  phone: '',
  latitude: null,
  longitude: null,
  business_hours: '09:00-21:00'
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
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  latitude: [
    { required: true, message: '请输入纬度', trigger: 'blur' }
  ],
  longitude: [
    { required: true, message: '请输入经度', trigger: 'blur' }
  ],
  business_hours: [
    { required: true, message: '请输入营业时间', trigger: 'blur' },
    { pattern: /^\d{2}:\d{2}-\d{2}:\d{2}$/, message: '请输入正确的营业时间格式，如：09:00-21:00', trigger: 'blur' }
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

// 返回
const goBack = () => {
  router.push('/stores')
}
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
