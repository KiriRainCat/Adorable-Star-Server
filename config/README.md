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
  proxy_port: 20172 # 代理服务器端口 (用于浏览器代理)
  fetch_interval: 30 # 数据抓取间隔 (分钟)
  max_parallel: 10 # 最大浏览器页面并行数

postgresql: # PostgreSQL 数据库设置
  dev_host: "127.0.0.1" # 开发环境连接地址
  host: "host" # 连接地址
  port: 5432 # 端口
  database: "name" # 数据库名
  user: "name" # 用户名
  password: "pwd" # 密码

redis: # Redis 数据库设置
  host: "host" # 连接地址
  port: 11813 # 端口
  password: "pwd" # 密码
  db: 0 # 默认数据库

smtp: # SMTP 邮件服务设置
  host: "smtp.host.com" # 连接地址
  port: 2023 # 端口
  key: "pwd" # 秘钥
  mail: "name@host.com" # 邮箱
```
