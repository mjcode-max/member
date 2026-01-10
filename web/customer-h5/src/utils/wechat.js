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
  const scope = 'snsapi_userinfo' // 需要用户授权，可以获取用户基本信息（手机号需要特殊权限）
  const state = 'STATE'
  
  const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appId}&redirect_uri=${encodedRedirectUri}&response_type=code&scope=${scope}&state=${state}#wechat_redirect`
  
  console.log('跳转到微信授权页面:', {
    redirectUri,
    encodedRedirectUri,
    scope,
    authUrl
  })
  
  window.location.href = authUrl
}

// 通过JS-SDK获取手机号（需要用户授权）
export const getPhoneNumberBySDK = () => {
  return new Promise((resolve, reject) => {
    // 检查是否在微信环境中
    if (!isWechatBrowser()) {
      reject(new Error('不在微信环境中'))
      return
    }

    // 检查JS-SDK是否加载
    if (typeof wx === 'undefined') {
      reject(new Error('微信JS-SDK未加载，请检查index.html中是否引入了jweixin.js'))
      return
    }

    // 注意：微信H5网页无法直接获取手机号
    // 手机号需要通过小程序接口（getPhoneNumber）获取
    // 或者使用微信开放平台的手机号一键登录功能（需要特殊权限）
    // 这里提供一个获取用户信息的接口，但可能不包含手机号
    wx.getUserInfo({
      success: (res) => {
        // getUserInfo返回的是用户基本信息，不包含手机号
        // 手机号需要通过其他方式获取
        console.warn('getUserInfo无法获取手机号，需要通过其他方式')
        resolve({
          userInfo: res.userInfo,
          phone: '' // 无法直接获取
        })
      },
      fail: (err) => {
        reject(err)
      }
    })
  })
}
