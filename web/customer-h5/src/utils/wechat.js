// 微信相关工具函数

// 检查是否在微信浏览器中
export const isWechatBrowser = () => {
  const ua = navigator.userAgent.toLowerCase()
  return ua.includes('micromessenger')
}

// 获取微信授权码
export const getWechatCode = () => {
  return new Promise((resolve, reject) => {
    // 检查URL中是否已有code
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    
    if (code) {
      resolve(code)
      return
    }

    // 如果没有code，跳转到微信授权页面
    const appId = 'your_wechat_app_id' // 这里需要替换为实际的AppID
    
    // 检查AppID是否已配置
    if (appId === 'your_wechat_app_id') {
      reject(new Error('微信AppID未配置，请联系管理员'))
      return
    }
    
    const redirectUri = encodeURIComponent(window.location.href)
    const scope = 'snsapi_userinfo'
    const state = 'STATE'
    
    const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appId}&redirect_uri=${redirectUri}&response_type=code&scope=${scope}&state=${state}#wechat_redirect`
    
    window.location.href = authUrl
  })
}

// 初始化微信JS-SDK
export const initWechatSDK = async (config) => {
  return new Promise((resolve, reject) => {
    if (typeof wx === 'undefined') {
      reject(new Error('微信JS-SDK未加载'))
      return
    }

    wx.config({
      debug: false,
      appId: config.appId,
      timestamp: config.timestamp,
      nonceStr: config.nonceStr,
      signature: config.signature,
      jsApiList: [
        'checkJsApi',
        'onMenuShareTimeline',
        'onMenuShareAppMessage',
        'onMenuShareQQ',
        'onMenuShareWeibo',
        'onMenuShareQZone',
        'hideMenuItems',
        'showMenuItems',
        'hideAllNonBaseMenuItem',
        'showAllNonBaseMenuItem',
        'translateVoice',
        'startRecord',
        'stopRecord',
        'onVoiceRecordEnd',
        'playVoice',
        'pauseVoice',
        'stopVoice',
        'onVoicePlayEnd',
        'uploadVoice',
        'downloadVoice',
        'chooseImage',
        'previewImage',
        'uploadImage',
        'downloadImage',
        'getNetworkType',
        'openLocation',
        'getLocation',
        'hideOptionMenu',
        'showOptionMenu',
        'closeWindow',
        'scanQRCode',
        'chooseWXPay',
        'openProductSpecificView',
        'addCard',
        'chooseCard',
        'openCard'
      ]
    })
    
    wx.ready(() => {
      console.log('微信JS-SDK初始化成功')
      resolve()
    })
    
    wx.error((res) => {
      console.error('微信JS-SDK初始化失败:', res)
      reject(res)
    })
  })
}

// 微信分享
export const shareToWechat = (shareData) => {
  if (typeof wx === 'undefined') {
    console.error('微信JS-SDK未加载')
    return
  }

  // 分享到朋友圈
  wx.onMenuShareTimeline({
    title: shareData.title,
    link: shareData.link,
    imgUrl: shareData.imgUrl,
    success: () => {
      console.log('分享到朋友圈成功')
    },
    cancel: () => {
      console.log('取消分享到朋友圈')
    }
  })

  // 分享给朋友
  wx.onMenuShareAppMessage({
    title: shareData.title,
    desc: shareData.desc,
    link: shareData.link,
    imgUrl: shareData.imgUrl,
    success: () => {
      console.log('分享给朋友成功')
    },
    cancel: () => {
      console.log('取消分享给朋友')
    }
  })
}

// 获取微信用户位置
export const getWechatLocation = () => {
  return new Promise((resolve, reject) => {
    if (typeof wx === 'undefined') {
      reject(new Error('微信JS-SDK未加载'))
      return
    }

    wx.getLocation({
      type: 'wgs84',
      success: (res) => {
        resolve({
          latitude: res.latitude,
          longitude: res.longitude,
          speed: res.speed,
          accuracy: res.accuracy
        })
      },
      fail: (err) => {
        reject(err)
      }
    })
  })
}

// 微信支付
export const wechatPay = (payData) => {
  return new Promise((resolve, reject) => {
    if (typeof wx === 'undefined') {
      reject(new Error('微信JS-SDK未加载'))
      return
    }

    wx.chooseWXPay({
      timestamp: payData.timestamp,
      nonceStr: payData.nonceStr,
      package: payData.package,
      signType: payData.signType,
      paySign: payData.paySign,
      success: (res) => {
        resolve(res)
      },
      fail: (err) => {
        reject(err)
      }
    })
  })
}
