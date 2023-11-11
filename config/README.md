# 配置文件说明

文件名: `config.yaml`

## 文件内容示例:

```yaml
# cSpell: disable

server: # 服务器设置
  port: 8005 # 服务器端口
  encrypt_salt: "awa" # SH265 加密盐
  request_auth: "qwq" # 接口请求鉴权 (Authorization)
  jwt_encrypt: "awa" # JWT 加密串 (Token)
  jwt_issuer: "server" # JWT 签发人 (Token)
  admin_auth: "qwq" # 管理员接口鉴权 (Admin)

crawler: # 爬虫设置
  browser_socket_url: "ws://127.0.0.1:7999/devtools/browser/bc02fd66-c01b-4a3b-85fa-a3cd912a49a3" # 生产环境浏览器 websocket 链接
  fetch_interval: 30 # 数据抓取间隔 (分钟)
  max_parallel: 10 # 最大浏览器页面并行数
```
