# Carbon Report Month 模块 - HTTP API 文档

## 概述

碳报告月报（Carbon Report Month）模块用于管理每月碳排放相关数据的统计报告，包括周转量、基准值、能耗、碳减排量等核心指标。

---

## 基础信息

- **基础路径**: `/api/carbon-report-month`
- **认证方式**: Token 认证（所有接口均需登录）
- **数据格式**: JSON

---

## 数据模型

### CarbonReportMonth 碳报告月报

```json
{
  "id": 1234567890,
  "collection_date": "2026-03-01 00:00:00",
  "turnover": 45000.50,
  "baseline": 36000.00,
  "energyConsumption": 8500.25,
  "carbonReduction": 9000.50,
  "createBy": 1001,
  "createTime": "2026-03-20 10:30:00",
  "updateBy": 1001,
  "updateTime": "2026-03-20 15:45:00"
}
```

**字段说明**：
- `id`: 记录 ID（雪花算法生成）
- `collection_date`: 数据采集日期（月份）
- `turnover`: 周转量（单位：km·人次）
- `baseline`: 基准值（单位：kg CO₂）
- `energyConsumption`: 能耗（人工填写）
- `carbonReduction`: 碳减排量（单位：kg CO₂）
- `createBy/createTime`: 创建人 ID/创建时间
- `updateBy/updateTime`: 更新人 ID/更新时间

---

## API 接口

### 1. 创建碳报告月报

**请求**：
```http
POST /api/carbon-report-month/carbonReportMonth
Content-Type: application/json
Authorization: Bearer <token>
```

**请求体**：
```json
{
  "collection_date": "2026-03-01 00:00:00",
  "turnover": 45000.50,
  "baseline": 36000.00,
  "energyConsumption": 8500.25,
  "carbonReduction": 9000.50
}
```

**响应**：
```json
{
  "code": 200,
  "data": {
    "id": 1234567890
  },
  "message": "success"
}
```

**说明**：
- 自动记录当前登录用户为创建人
- `collection_date` 格式：`yyyy-MM-dd HH:mm:ss`

---

### 2. 更新碳报告月报

**请求**：
```http
PUT /api/carbon-report-month/carbonReportMonth
Content-Type: application/json
Authorization: Bearer <token>
```

**请求体**：
```json
{
  "id": 1234567890
}
```

**响应**：
```json
{
  "code": 200,
  "data": null,
  "message": "success"
}
```

**说明**：
- 目前仅支持更新操作标记（实际业务字段可扩展）
- 自动记录当前登录用户为更新人

---

### 3. 删除碳报告月报

**请求**：
```http
DELETE /api/carbon-report-month/carbonReportMonth?id=1234567890
Authorization: Bearer <token>
```

**参数**：
- `id` (query, required): 记录 ID

**响应**：
```json
{
  "code": 200,
  "data": null,
  "message": "success"
}
```

**说明**：
- 需要提供创建人 ID 进行权限验证

---

### 4. 根据 ID 查询碳报告月报

**请求**：
```http
GET /api/carbon-report-month/carbonReportMonth?id=1234567890
Authorization: Bearer <token>
```

**参数**：
- `id` (query, required): 记录 ID

**响应**：
```json
{
  "code": 200,
  "data": {
    "id": 1234567890,
    "collection_date": "2026-03-01 00:00:00",
    "turnover": 45000.50,
    "baseline": 36000.00,
    "energyConsumption": 8500.25,
    "carbonReduction": 9000.50,
    "createBy": 1001,
    "createTime": "2026-03-20 10:30:00",
    "updateBy": 1001,
    "updateTime": "2026-03-20 15:45:00"
  },
  "message": "success"
}
```

---

### 5. 分页查询碳报告月报

**请求**：
```http
GET /api/carbon-report-month/carbonReportMonths/page?pageNum=1&pageSize=10&startDate=2026-01-01&endDate=2026-12-31&sortBy=collection_date&sortOrder=desc
Authorization: Bearer <token>
```

**查询参数**：
- `pageNum` (query, optional): 页码，默认 1
- `pageSize` (query, optional): 每页数量，默认 10
- `startDate` (query, optional): 开始日期（格式：yyyy-MM-dd）
- `endDate` (query, optional): 结束日期（格式：yyyy-MM-dd）
- `sortBy` (query, optional): 排序字段，默认 `collection_date`
- `sortOrder` (query, optional): 排序方式，`asc` 或 `desc`

**响应**：
```json
{
  "code": 200,
  "data": {
    "data": [
      {
        "id": 1234567890,
        "collection_date": "2026-03-01 00:00:00",
        "turnover": 45000.50,
        "baseline": 36000.00,
        "energyConsumption": 8500.25,
        "carbonReduction": 9000.50,
        "createBy": 1001,
        "createTime": "2026-03-20 10:30:00",
        "updateBy": 1001,
        "updateTime": "2026-03-20 15:45:00"
      }
    ],
    "total": 12,
    "pageNum": 1,
    "pageSize": 10,
    "totalPages": 2
  },
  "message": "success"
}
```

**说明**：
- 支持按日期范围筛选
- 支持自定义排序字段和顺序
- 返回分页元数据（总记录数、总页数等）

---

## 错误码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（Token 无效或过期） |
| 500 | 服务器内部错误 |

**错误响应示例**：
```json
{
  "code": 400,
  "message": "invalid id",
  "data": null
}
```

---

## cURL 示例

### 创建碳报告月报
```bash
curl -X POST http://localhost:8080/api/carbon-report-month/carbonReportMonth \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "collection_date": "2026-03-01 00:00:00",
    "turnover": 45000.50,
    "baseline": 36000.00,
    "energyConsumption": 8500.25,
    "carbonReduction": 9000.50
  }'
```

### 分页查询
```bash
curl -X GET "http://localhost:8080/api/carbon-report-month/carbonReportMonths/page?pageNum=1&pageSize=10" \
  -H "Authorization: Bearer <your_token>"
```

### 根据 ID 查询
```bash
curl -X GET "http://localhost:8080/api/carbon-report-month/carbonReportMonth?id=1234567890" \
  -H "Authorization: Bearer <your_token>"
```

---

## 路由结构

```
/api/carbon-report-month/
├── carbonReportMonth          # 单条记录操作
│   ├── POST                   # 创建
│   ├── PUT                    # 更新
│   ├── DELETE                 # 删除
│   └── GET                    # 按 ID 查询
│
└── carbonReportMonths/        # 列表查询
    └── page                   # 分页查询
```

---

## 注意事项

1. **时间格式**：所有日期时间字段使用 `yyyy-MM-dd HH:mm:ss` 格式
2. **数值精度**：金额和排放量字段保留 2-4 位小数
3. **权限控制**：删除操作需要验证创建人权限
4. **并发安全**：使用雪花算法生成唯一 ID，支持高并发场景
5. **分页默认值**：未指定分页参数时，默认返回第 1 页，每页 10 条记录
6. **能耗字段**：`energyConsumption` 为人工填写字段，用于记录月度能耗数据
