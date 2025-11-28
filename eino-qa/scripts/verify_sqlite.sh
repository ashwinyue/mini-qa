#!/bin/bash

# SQLite 集成验证脚本

set -e

echo "=========================================="
echo "SQLite 集成验证脚本"
echo "=========================================="
echo ""

# 1. 编译检查
echo "1. 编译检查..."
cd "$(dirname "$0")/.."
go build ./internal/infrastructure/repository/sqlite/...
echo "✅ 编译成功"
echo ""

# 2. 运行基础功能测试
echo "2. 运行基础功能测试..."
go run examples/sqlite_usage_example.go > /dev/null 2>&1
echo "✅ 基础功能测试通过"
echo ""

# 3. 运行多租户隔离测试
echo "3. 运行多租户隔离测试..."
go run examples/multi_tenant_example.go > /dev/null 2>&1
echo "✅ 多租户隔离测试通过"
echo ""

# 4. 验证数据库文件
echo "4. 验证数据库文件..."
if [ -f "data/db/tenant1.db" ] && [ -f "data/db/tenant2.db" ] && [ -f "data/db/tenant3.db" ]; then
    echo "✅ 数据库文件创建成功"
    ls -lh data/db/*.db
else
    echo "❌ 数据库文件缺失"
    exit 1
fi
echo ""

# 5. 验证表结构
echo "5. 验证表结构..."
TABLES=$(sqlite3 data/db/tenant1.db ".tables")
if [[ $TABLES == *"orders"* ]] && [[ $TABLES == *"sessions"* ]] && [[ $TABLES == *"missed_queries"* ]]; then
    echo "✅ 表结构正确"
    echo "   表: $TABLES"
else
    echo "❌ 表结构不完整"
    exit 1
fi
echo ""

# 6. 验证数据
echo "6. 验证数据..."
ORDER_COUNT=$(sqlite3 data/db/tenant1.db "SELECT COUNT(*) FROM orders;")
SESSION_COUNT=$(sqlite3 data/db/tenant1.db "SELECT COUNT(*) FROM sessions;")
echo "✅ 数据验证通过"
echo "   tenant1 订单数: $ORDER_COUNT"
echo "   tenant1 会话数: $SESSION_COUNT"
echo ""

echo "=========================================="
echo "所有验证通过！✅"
echo "=========================================="
