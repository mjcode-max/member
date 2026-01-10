<template>
  <div class="template-form-page">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>{{ isEdit ? '编辑时段模板' : '创建时段模板' }}</h2>
    </div>
    
    <el-card class="form-card" v-loading="loading">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入模板名称"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="时间段" prop="time_slots">
          <div class="time-slots-container">
            <div
              v-for="(slot, index) in form.time_slots"
              :key="index"
              class="time-slot-item"
            >
              <el-time-picker
                v-model="slot.start_time"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="开始时间"
                style="width: 150px; margin-right: 10px"
                @change="validateTimeSlot(index)"
              />
              <span style="margin-right: 10px">至</span>
              <el-time-picker
                v-model="slot.end_time"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="结束时间"
                style="width: 150px; margin-right: 10px"
                @change="validateTimeSlot(index)"
              />
              <el-button
                type="danger"
                link
                size="small"
                @click="removeTimeSlot(index)"
              >
                删除
              </el-button>
            </div>
            
            <el-button
              type="primary"
              @click="addTimeSlot"
              style="width: 100%; margin-top: 10px"
            >
              + 添加时间段
            </el-button>
          </div>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '保存修改' : '创建模板' }}
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
import { getTemplateById, createTemplate, updateTemplate } from '@/api/templates'

const router = useRouter()
const route = useRoute()

const formRef = ref()
const loading = ref(false)
const submitting = ref(false)

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  name: '',
  status: 'active',
  time_slots: []
})

const rules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
    { min: 2, max: 100, message: '模板名称长度在2到100个字符', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  time_slots: [
    { required: true, message: '请至少添加一个时间段', trigger: 'change' },
    {
      validator: (rule, value, callback) => {
        if (!value || value.length === 0) {
          callback(new Error('请至少添加一个时间段'))
          return
        }
        
        // 检查每个时间段是否完整
        for (let i = 0; i < value.length; i++) {
          const slot = value[i]
          if (!slot.start_time || !slot.end_time) {
            callback(new Error(`时间段 ${i + 1} 的开始时间和结束时间不能为空`))
            return
          }
          
          if (slot.start_time >= slot.end_time) {
            callback(new Error(`时间段 ${i + 1} 的开始时间必须早于结束时间`))
            return
          }
        }
        
        // 检查时间段是否重叠
        for (let i = 0; i < value.length; i++) {
          for (let j = i + 1; j < value.length; j++) {
            const slot1 = value[i]
            const slot2 = value[j]
            if (slot1.start_time < slot2.end_time && slot1.end_time > slot2.start_time) {
              callback(new Error('时间段不能重叠'))
              return
            }
          }
        }
        
        callback()
      },
      trigger: 'change'
    }
  ]
}

// 添加时间段
const addTimeSlot = () => {
  form.time_slots.push({
    start_time: '09:00',
    end_time: '10:00'
  })
}

// 删除时间段
const removeTimeSlot = (index) => {
  form.time_slots.splice(index, 1)
}

// 验证时间段
const validateTimeSlot = (index) => {
  // 触发表单验证
  if (formRef.value) {
    formRef.value.validateField('time_slots')
  }
}

// 获取模板详情
const fetchTemplateDetail = async () => {
  const templateId = route.params.id
  if (!templateId) return
  
  loading.value = true
  try {
    const response = await getTemplateById(templateId)
    const templateData = response.data || response
    
    Object.assign(form, {
      name: templateData.name || '',
      status: templateData.status || 'active',
      time_slots: templateData.time_slots || []
    })
  } catch (error) {
    console.error('获取模板详情失败:', error)
    ElMessage.error('获取模板详情失败')
    router.push('/templates')
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
    
    const templateData = {
      name: form.name,
      status: form.status,
      time_slots: form.time_slots
    }
    
    if (isEdit.value) {
      await updateTemplate(route.params.id, templateData)
      ElMessage.success('模板更新成功')
    } else {
      await createTemplate(templateData)
      ElMessage.success('模板创建成功')
    }
    
    router.push('/templates')
  } catch (error) {
    console.error('保存模板失败:', error)
  } finally {
    submitting.value = false
  }
}

// 返回
const goBack = () => {
  router.push('/templates')
}

onMounted(() => {
  if (isEdit.value) {
    fetchTemplateDetail()
  } else {
    // 默认添加一个时间段
    addTimeSlot()
  }
})
</script>

<style lang="scss" scoped>
.template-form-page {
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
  }
  
  .time-slots-container {
    .time-slot-item {
      display: flex;
      align-items: center;
      margin-bottom: 10px;
      padding: 10px;
      border: 1px solid #dcdfe6;
      border-radius: 4px;
      background-color: #f5f7fa;
    }
  }
}
</style>
