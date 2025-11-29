# 部署指南

## 概述

本指南介绍如何在不同环境中部署 Eino QA System。

## 目录

- [开发环境部署](#开发环境部署)
- [生产环境部署](#生产环境部署)
- [Docker 部署](#docker-部署)
- [Kubernetes 部署](#kubernetes-部署)
- [性能优化](#性能优化)
- [监控和日志](#监控和日志)
- [故障排查](#故障排查)

---

## 开发环境部署

### 前置要求

- Go 1.23+
- Docker 和 Docker Compose
- DashScope API Key

### 快速部署

1. **克隆项目**

```bash
git clone <repository-url>
cd eino-qa
```

2. **配置环境变量**

```bash
cp .env.example .env
# 编辑 .env 文件，填入 API Key
```

3. **启动 Milvus**

```bash
docker-compose -f docker-compose.milvus.yml up -d
```

4. **运行服务**

```bash
make run
```


---

## 生产环境部署

### 系统要求

**硬件要求**:
- CPU: 4 核心以上
- 内存: 8GB 以上
- 磁盘: 50GB 以上 SSD

**软件要求**:
- Linux (Ubuntu 20.04+ / CentOS 7+)
- Go 1.23+
- Milvus 2.4+
- Systemd (用于服务管理)

### 部署步骤

#### 1. 准备服务器

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装必要工具
sudo apt install -y git make curl

# 安装 Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. 部署 Milvus

**使用 Docker Compose**:

```bash
# 下载 Milvus docker-compose 配置
wget https://github.com/milvus-io/milvus/releases/download/v2.4.0/milvus-standalone-docker-compose.yml

# 启动 Milvus
docker-compose -f milvus-standalone-docker-compose.yml up -d

# 验证 Milvus 运行状态
docker ps | grep milvus
```

**生产环境建议**:
- 使用 Milvus 集群模式以提高可用性
- 配置持久化存储
- 启用认证和 TLS

#### 3. 部署应用

```bash
# 创建应用目录
sudo mkdir -p /opt/eino-qa
cd /opt/eino-qa

# 克隆代码
git clone <repository-url> .

# 编译应用
make build-prod

# 创建数据目录
sudo mkdir -p /var/lib/eino-qa/db
sudo mkdir -p /var/log/eino-qa
```


#### 4. 配置应用

创建生产配置文件 `/opt/eino-qa/config/config.prod.yaml`:

```yaml
server:
  port: 8080
  mode: release

dashscope:
  api_key: ${DASHSCOPE_API_KEY}
  chat_model: qwen-plus
  embed_model: text-embedding-v2
  embedding_dimension: 1536
  max_retries: 3
  timeout: 30s

milvus:
  host: localhost
  port: 19530
  username: ${MILVUS_USERNAME}
  password: ${MILVUS_PASSWORD}
  timeout: 10s

database:
  base_path: /var/lib/eino-qa/db

logging:
  level: info
  format: json
  output: file
  file_path: /var/log/eino-qa/app.log

security:
  api_keys:
    - ${API_KEY_1}
    - ${API_KEY_2}
```

创建环境变量文件 `/opt/eino-qa/.env.prod`:

```bash
DASHSCOPE_API_KEY=your_production_api_key
MILVUS_USERNAME=your_milvus_username
MILVUS_PASSWORD=your_milvus_password
API_KEY_1=your_production_api_key_1
API_KEY_2=your_production_api_key_2
```

#### 5. 配置 Systemd 服务

创建服务文件 `/etc/systemd/system/eino-qa.service`:

```ini
[Unit]
Description=Eino QA System
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
User=eino-qa
Group=eino-qa
WorkingDirectory=/opt/eino-qa
EnvironmentFile=/opt/eino-qa/.env.prod
ExecStart=/opt/eino-qa/bin/server -config /opt/eino-qa/config/config.prod.yaml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/eino-qa /var/log/eino-qa

[Install]
WantedBy=multi-user.target
```

创建服务用户：

```bash
sudo useradd -r -s /bin/false eino-qa
sudo chown -R eino-qa:eino-qa /opt/eino-qa
sudo chown -R eino-qa:eino-qa /var/lib/eino-qa
sudo chown -R eino-qa:eino-qa /var/log/eino-qa
```

启动服务：

```bash
# 重新加载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start eino-qa

# 设置开机自启
sudo systemctl enable eino-qa

# 查看服务状态
sudo systemctl status eino-qa

# 查看日志
sudo journalctl -u eino-qa -f
```


#### 6. 配置 Nginx 反向代理

安装 Nginx：

```bash
sudo apt install -y nginx
```

创建 Nginx 配置 `/etc/nginx/sites-available/eino-qa`:

```nginx
upstream eino_qa {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name your-domain.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 证书配置
    ssl_certificate /etc/ssl/certs/your-cert.pem;
    ssl_certificate_key /etc/ssl/private/your-key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 日志
    access_log /var/log/nginx/eino-qa-access.log;
    error_log /var/log/nginx/eino-qa-error.log;

    # 请求大小限制
    client_max_body_size 10M;

    # 超时设置
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;

    location / {
        proxy_pass http://eino_qa;
        proxy_http_version 1.1;
        
        # 请求头
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持（用于流式响应）
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 健康检查端点
    location /health {
        proxy_pass http://eino_qa;
        access_log off;
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/eino-qa /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## Docker 部署

### 创建 Dockerfile

创建 `Dockerfile`:

```dockerfile
# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git make

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN make build-prod

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1000 eino && \
    adduser -D -u 1000 -G eino eino

# 复制编译后的二进制文件
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/config /app/config

# 创建数据目录
RUN mkdir -p /app/data/db /app/logs && \
    chown -R eino:eino /app

USER eino

EXPOSE 8080

CMD ["/app/server"]
```


### 创建 Docker Compose 配置

创建 `docker-compose.yml`:

```yaml
version: '3.8'

services:
  # Milvus 向量数据库
  etcd:
    image: quay.io/coreos/etcd:v3.5.5
    environment:
      - ETCD_AUTO_COMPACTION_MODE=revision
      - ETCD_AUTO_COMPACTION_RETENTION=1000
      - ETCD_QUOTA_BACKEND_BYTES=4294967296
      - ETCD_SNAPSHOT_COUNT=50000
    volumes:
      - etcd_data:/etcd
    command: etcd -advertise-client-urls=http://127.0.0.1:2379 -listen-client-urls http://0.0.0.0:2379 --data-dir /etcd

  minio:
    image: minio/minio:RELEASE.2023-03-20T20-16-18Z
    environment:
      MINIO_ACCESS_KEY: minioadmin
      MINIO_SECRET_KEY: minioadmin
    volumes:
      - minio_data:/minio_data
    command: minio server /minio_data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  milvus:
    image: milvusdb/milvus:v2.4.0
    command: ["milvus", "run", "standalone"]
    environment:
      ETCD_ENDPOINTS: etcd:2379
      MINIO_ADDRESS: minio:9000
    volumes:
      - milvus_data:/var/lib/milvus
    ports:
      - "19530:19530"
      - "9091:9091"
    depends_on:
      - etcd
      - minio

  # Eino QA 应用
  eino-qa:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DASHSCOPE_API_KEY=${DASHSCOPE_API_KEY}
      - API_KEY_1=${API_KEY_1}
      - API_KEY_2=${API_KEY_2}
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs
      - ./config:/app/config
    depends_on:
      - milvus
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

volumes:
  etcd_data:
  minio_data:
  milvus_data:
```

### 构建和运行

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f eino-qa

# 停止服务
docker-compose down

# 清理数据（谨慎使用）
docker-compose down -v
```

---

## Kubernetes 部署

### 创建 Kubernetes 资源

#### 1. ConfigMap

创建 `k8s/configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: eino-qa-config
  namespace: default
data:
  config.yaml: |
    server:
      port: 8080
      mode: release
    
    dashscope:
      api_key: ${DASHSCOPE_API_KEY}
      chat_model: qwen-plus
      embed_model: text-embedding-v2
      embedding_dimension: 1536
    
    milvus:
      host: milvus-service
      port: 19530
    
    database:
      base_path: /data/db
    
    logging:
      level: info
      format: json
      output: stdout
```


#### 2. Secret

创建 `k8s/secret.yaml`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: eino-qa-secret
  namespace: default
type: Opaque
stringData:
  DASHSCOPE_API_KEY: "your_dashscope_api_key"
  API_KEY_1: "your_api_key_1"
  API_KEY_2: "your_api_key_2"
```

#### 3. Deployment

创建 `k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: eino-qa
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: eino-qa
  template:
    metadata:
      labels:
        app: eino-qa
    spec:
      containers:
      - name: eino-qa
        image: your-registry/eino-qa:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: DASHSCOPE_API_KEY
          valueFrom:
            secretKeyRef:
              name: eino-qa-secret
              key: DASHSCOPE_API_KEY
        - name: API_KEY_1
          valueFrom:
            secretKeyRef:
              name: eino-qa-secret
              key: API_KEY_1
        - name: API_KEY_2
          valueFrom:
            secretKeyRef:
              name: eino-qa-secret
              key: API_KEY_2
        volumeMounts:
        - name: config
          mountPath: /app/config
        - name: data
          mountPath: /app/data
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
      volumes:
      - name: config
        configMap:
          name: eino-qa-config
      - name: data
        persistentVolumeClaim:
          claimName: eino-qa-data
```

#### 4. Service

创建 `k8s/service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: eino-qa-service
  namespace: default
spec:
  selector:
    app: eino-qa
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

#### 5. Ingress

创建 `k8s/ingress.yaml`:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: eino-qa-ingress
  namespace: default
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - your-domain.com
    secretName: eino-qa-tls
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: eino-qa-service
            port:
              number: 80
```

#### 6. PersistentVolumeClaim

创建 `k8s/pvc.yaml`:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: eino-qa-data
  namespace: default
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 50Gi
  storageClassName: standard
```

### 部署到 Kubernetes

```bash
# 创建命名空间
kubectl create namespace eino-qa

# 应用配置
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/pvc.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml

# 查看部署状态
kubectl get pods -n eino-qa
kubectl get svc -n eino-qa
kubectl get ingress -n eino-qa

# 查看日志
kubectl logs -f deployment/eino-qa -n eino-qa

# 扩容
kubectl scale deployment eino-qa --replicas=5 -n eino-qa
```


---

## 性能优化

### 1. 应用层优化

**连接池配置**:

```yaml
milvus:
  max_connections: 100
  max_idle_connections: 20
  connection_timeout: 10s

database:
  max_open_conns: 50
  max_idle_conns: 10
  conn_max_lifetime: 1h
```

**缓存配置**:

```yaml
cache:
  enabled: true
  type: redis  # memory, redis
  redis:
    host: localhost
    port: 6379
    db: 0
  ttl: 1h
```

**并发控制**:

```yaml
server:
  max_concurrent_requests: 1000
  request_timeout: 30s
  read_timeout: 10s
  write_timeout: 10s
```

### 2. Milvus 优化

**索引优化**:

```go
// 使用 HNSW 索引提高检索性能
indexParams := map[string]interface{}{
    "index_type": "HNSW",
    "metric_type": "L2",
    "params": map[string]interface{}{
        "M": 16,
        "efConstruction": 256,
    },
}
```

**搜索参数优化**:

```go
searchParams := map[string]interface{}{
    "ef": 64,  // 搜索时的候选数量
}
```

### 3. 数据库优化

**SQLite 优化**:

```sql
-- 启用 WAL 模式
PRAGMA journal_mode=WAL;

-- 增加缓存大小
PRAGMA cache_size=-64000;  -- 64MB

-- 启用内存映射
PRAGMA mmap_size=268435456;  -- 256MB
```

### 4. 系统层优化

**文件描述符限制**:

```bash
# 编辑 /etc/security/limits.conf
eino-qa soft nofile 65536
eino-qa hard nofile 65536
```

**内核参数优化**:

```bash
# 编辑 /etc/sysctl.conf
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048
net.ipv4.tcp_tw_reuse = 1
```

---

## 监控和日志

### 1. 日志收集

**使用 Filebeat 收集日志**:

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/eino-qa/app.log
  json.keys_under_root: true
  json.add_error_key: true

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "eino-qa-%{+yyyy.MM.dd}"
```

### 2. 指标监控

**Prometheus 配置**:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'eino-qa'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

**Grafana 仪表板**:

关键指标：
- 请求 QPS
- 响应时间 (P50, P95, P99)
- 错误率
- Milvus 查询延迟
- 数据库连接数

### 3. 告警配置

**Prometheus 告警规则**:

```yaml
groups:
- name: eino-qa
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }} requests/sec"

  - alert: HighResponseTime
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High response time detected"
      description: "P95 response time is {{ $value }} seconds"
```

---

## 故障排查

### 常见问题

#### 1. 服务无法启动

**症状**: 服务启动失败或立即退出

**排查步骤**:

```bash
# 查看服务状态
sudo systemctl status eino-qa

# 查看日志
sudo journalctl -u eino-qa -n 100

# 检查配置文件
./bin/server -config config/config.yaml -validate

# 检查端口占用
sudo lsof -i :8080
```

**常见原因**:
- 配置文件错误
- 端口被占用
- 权限不足
- 依赖服务未启动

#### 2. Milvus 连接失败

**症状**: 日志显示 "failed to connect to milvus"

**排查步骤**:

```bash
# 检查 Milvus 状态
docker ps | grep milvus

# 测试连接
telnet localhost 19530

# 查看 Milvus 日志
docker logs milvus-standalone
```

**解决方法**:
- 确保 Milvus 正在运行
- 检查网络连接
- 验证认证信息

#### 3. 内存泄漏

**症状**: 内存使用持续增长

**排查步骤**:

```bash
# 查看内存使用
ps aux | grep server

# 使用 pprof 分析
go tool pprof http://localhost:8080/debug/pprof/heap

# 查看 goroutine 泄漏
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

#### 4. 性能下降

**症状**: 响应时间变长

**排查步骤**:

```bash
# 查看系统负载
top
htop

# 查看磁盘 I/O
iostat -x 1

# 查看网络连接
netstat -an | grep 8080

# 分析慢查询
# 查看应用日志中的慢请求
```

### 日志级别调整

动态调整日志级别（无需重启）:

```bash
# 设置为 debug 级别
curl -X POST http://localhost:8080/admin/log-level \
  -H "X-API-Key: admin_key" \
  -d '{"level": "debug"}'

# 恢复为 info 级别
curl -X POST http://localhost:8080/admin/log-level \
  -H "X-API-Key: admin_key" \
  -d '{"level": "info"}'
```

---

## 备份和恢复

### 数据备份

#### 1. SQLite 数据库备份

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backup/eino-qa"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份所有租户数据库
for db in /var/lib/eino-qa/db/*.db; do
    tenant=$(basename $db .db)
    sqlite3 $db ".backup '$BACKUP_DIR/${tenant}_${DATE}.db'"
done

# 压缩备份
tar -czf $BACKUP_DIR/backup_${DATE}.tar.gz $BACKUP_DIR/*_${DATE}.db
rm $BACKUP_DIR/*_${DATE}.db

# 删除 7 天前的备份
find $BACKUP_DIR -name "backup_*.tar.gz" -mtime +7 -delete
```

#### 2. Milvus 数据备份

```bash
# 使用 Milvus 备份工具
docker exec milvus-standalone \
  milvus-backup create \
  --backup-name backup_$(date +%Y%m%d)
```

### 数据恢复

#### 1. SQLite 恢复

```bash
# 解压备份
tar -xzf backup_20241129_120000.tar.gz

# 恢复数据库
cp tenant1_20241129_120000.db /var/lib/eino-qa/db/tenant1.db

# 重启服务
sudo systemctl restart eino-qa
```

#### 2. Milvus 恢复

```bash
# 恢复 Milvus 数据
docker exec milvus-standalone \
  milvus-backup restore \
  --backup-name backup_20241129
```

---

## 安全加固

### 1. 网络安全

**防火墙配置**:

```bash
# 只允许必要的端口
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw deny 8080/tcp  # 应用端口不对外暴露
sudo ufw enable
```

### 2. 应用安全

**启用 HTTPS**:

```yaml
server:
  tls:
    enabled: true
    cert_file: /etc/ssl/certs/server.crt
    key_file: /etc/ssl/private/server.key
```

**API Key 轮换**:

```bash
# 定期更新 API Key
# 1. 生成新的 API Key
# 2. 更新配置文件
# 3. 重启服务
# 4. 通知客户端更新
```

### 3. 数据安全

**敏感数据加密**:

```yaml
security:
  encryption:
    enabled: true
    key_file: /etc/eino-qa/encryption.key
```

---

## 更新和升级

### 滚动更新

```bash
# 1. 备份数据
./scripts/backup.sh

# 2. 拉取新代码
git pull origin main

# 3. 编译新版本
make build-prod

# 4. 重启服务
sudo systemctl restart eino-qa

# 5. 验证服务
curl http://localhost:8080/health
```

### Kubernetes 滚动更新

```bash
# 更新镜像
kubectl set image deployment/eino-qa \
  eino-qa=your-registry/eino-qa:v1.1.0 \
  -n eino-qa

# 查看更新状态
kubectl rollout status deployment/eino-qa -n eino-qa

# 回滚（如果需要）
kubectl rollout undo deployment/eino-qa -n eino-qa
```

---

## 附录

### A. 系统要求清单

| 组件 | 最低要求 | 推荐配置 |
|------|---------|---------|
| CPU | 2 核心 | 4 核心+ |
| 内存 | 4GB | 8GB+ |
| 磁盘 | 20GB | 50GB+ SSD |
| 网络 | 100Mbps | 1Gbps |

### B. 端口列表

| 端口 | 服务 | 说明 |
|------|------|------|
| 8080 | Eino QA | 应用服务 |
| 19530 | Milvus | gRPC 端口 |
| 9091 | Milvus | HTTP 端口 |
| 2379 | etcd | 客户端端口 |
| 9000 | MinIO | API 端口 |

### C. 环境变量参考

| 变量名 | 必填 | 说明 |
|--------|------|------|
| DASHSCOPE_API_KEY | 是 | DashScope API Key |
| API_KEY_1 | 否 | 向量管理 API Key |
| MILVUS_HOST | 否 | Milvus 主机地址 |
| MILVUS_PORT | 否 | Milvus 端口 |

---

## 获取帮助

如有部署问题，请：
1. 查看日志文件
2. 参考故障排查章节
3. 查看项目文档
4. 提交 Issue

---

**文档版本**: v1.0.0  
**最后更新**: 2024-11-29
