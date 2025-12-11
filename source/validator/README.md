# Validator Demo

这是一个使用 `github.com/go-playground/validator/v10` 进行字段验证的示例项目。

## 功能说明

- 使用 validator 库对结构体字段进行验证
- 支持中文错误消息翻译
- 提供 HTTP API 接口演示验证功能

## 验证规则示例

在 `AddConfigDictReq` 结构体中定义了以下验证规则：

- `parentId`: 数字类型，大于等于 0
- `name`: 字符串，长度在 2-50 之间
- `type`: 数字类型，范围在 1-12 之间
- `uniqueKey`: 字符串，长度在 2-50 之间
- `value`: 字符串，最大长度 2048
- `orderNum`: 数字类型，范围在 0-9999 之间
- `remark`: 字符串，最大长度 200
- `status`: 数字类型，范围在 0-1 之间

## 使用方法

### 1. 安装依赖

```bash
make deps
# 或
go mod download
go mod tidy
```

### 2. 运行程序

```bash
make run
# 或
go run main.go
```

### 3. 测试接口

#### 健康检查
```bash
curl http://localhost:8080/health
```

#### 添加配置字典（验证成功示例）
```bash
curl -X POST http://localhost:8080/api/config-dict/add \
  -H "Content-Type: application/json" \
  -d '{
    "parentId": 0,
    "name": "测试字典",
    "type": 1,
    "uniqueKey": "test_key",
    "value": "测试值",
    "orderNum": 1,
    "remark": "这是备注",
    "status": 1
  }'
```

#### 验证失败示例（name 字段太短）
```bash
curl -X POST http://localhost:8080/api/config-dict/add \
  -H "Content-Type: application/json" \
  -d '{
    "parentId": 0,
    "name": "a",
    "type": 1,
    "uniqueKey": "test_key",
    "value": "测试值",
    "orderNum": 1,
    "remark": "这是备注",
    "status": 1
  }'
```

预期返回中文错误消息：`名称长度必须至少为 2 个字符`

## 核心代码说明

1. **创建验证器**：使用 `validator.New()` 创建验证器实例
2. **注册标签名函数**：使用 `RegisterTagNameFunc` 将 `label` 标签作为字段名，用于错误消息
3. **配置中文翻译**：使用 `universal-translator` 和中文翻译包，将验证错误转换为中文
4. **验证结构体**：使用 `validate.Struct()` 对结构体进行验证

## 参考文档

- [validator 官方文档](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [中文翻译包](https://pkg.go.dev/github.com/go-playground/validator/v10/translations/zh)

