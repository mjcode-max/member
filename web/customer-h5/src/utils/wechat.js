// 微信登录相关工具函数
import { WECHAT_CONFIG } from '@/config/wechat'

// 检查是否在微信浏览器中
export const isWechatBrowser = () => {
  const ua = navigator.userAgent.toLowerCase()
  return ua.includes('micromessenger')
}

// 跳转到微信授权页面获取code
export const redirectToWechatAuth = () => {
  const appId = WECHAT_CONFIG.appId
  
  // 检查AppID是否已配置
  if (!appId || appId === 'your_wechat_app_id') {
    throw new Error('微信AppID未配置，请联系管理员')
  }
  
  // 构建redirect_uri
  let redirectUri = WECHAT_CONFIG.redirectUri || (window.location.origin + window.location.pathname)
  
  // 如果使用配置的redirectUri，直接使用；否则处理当前URL
  if (!WECHAT_CONFIG.redirectUri) {
    // 移除非标准端口号（微信不允许redirect_uri包含非标准端口）
    const urlObj = new URL(redirectUri)
    if (urlObj.port && urlObj.port !== '80' && urlObj.port !== '443') {
      if ((urlObj.protocol === 'http:' && urlObj.port !== '80') || 
          (urlObj.protocol === 'https:' && urlObj.port !== '443')) {
        redirectUri = `${urlObj.protocol}//${urlObj.hostname}${urlObj.pathname}`
      }
    }
    
    // 确保redirect_uri指向登录页面
    if (!redirectUri.endsWith('/login')) {
      redirectUri = redirectUri.replace(/\/[^/]*$/, '/login')
    }
  }
  
  const encodedRedirectUri = encodeURIComponent(redirectUri)
  const scope = 'snsapi_base' // 静默授权，只需要openid
  const state = 'STATE'
  
  const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appId}&redirect_uri=${encodedRedirectUri}&response_type=code&scope=${scope}&state=${state}#wechat_redirect`
  
  window.location.href = authUrl
}
