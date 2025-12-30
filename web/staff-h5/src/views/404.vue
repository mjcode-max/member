<template>
  <div class="not-found">
    <van-empty
      image="error"
      description="页面不存在"
    >
      <van-button
        round
        type="primary"
        @click="goHome"
      >
        返回首页
      </van-button>
    </van-empty>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style lang="scss" scoped>
.not-found {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f7f8fa;
}
</style>
