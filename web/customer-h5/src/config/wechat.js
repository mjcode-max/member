// 微信配置
export const WECHAT_CONFIG = {
  // 微信AppID（可以从微信公众平台获取）
  // 注意：AppID是公开的，可以暴露在前端代码中
  // AppSecret必须保密，只能在后端使用
  appId: 'wx5c0484314f35d6d4',
  
  // 授权回调地址（可选）
  // 如果设置了此值，将使用此值作为redirect_uri，否则使用当前页面的URL
  // 注意：此域名必须在微信公众平台的"授权回调域名"中配置
  // 
  // 生产环境：https://yourdomain.com/login
  // 开发环境（内网穿透）：
  //   - 使用 ngrok: https://xxxxx.ngrok.io/login
  //   - 使用 natapp: https://xxxxx.natapp1.cc/login
  //   - 使用微信开发者工具：工具 -> 内网穿透 -> 设置公网域名
  // 
  // 内网开发说明：
  // 1. 微信授权需要重定向，内网地址无法被微信访问
  // 2. 必须使用内网穿透工具将内网地址映射到公网
  // 3. 在微信公众平台配置内网穿透的公网域名
  // 4. 在此处配置内网穿透的公网地址作为redirectUri
  redirectUri: '' // 留空则自动使用当前页面URL
}

