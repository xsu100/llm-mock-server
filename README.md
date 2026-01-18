# LLM Mock Server

[English](#english) | [中文](#中文)

---

## English

A lightweight, high-performance LLM Mock Server built in Go. It simulates OpenAI-compatible APIs (Text Chat & Image Generation) to help developers save costs and speed up development/testing of LLM-based applications.

### 🚀 Features

- **OpenAI & Vertex AI Compatible**: Drop-in replacement for OpenAI endpoints and Vertex AI Gemini `generateContent` API.
- **Dynamic Rules**: Manage mock responses via a simple Admin API. Support for:
    - **Fixed Match**: Match specific keywords.
    - **Regex Match**: Match complex patterns in the prompt.
    - **Random**: Return random responses from a pool.
- **Persistence**: Rules are stored in a local SQLite database.
- **Lightweight**: Zero external dependencies (CGO-free SQLite).

### 🛠️ Quick Start

1. **Build**:
   ```bash
   go build -o mock-server main.go
   ```

2. **Run**:
   ```bash
   ./mock-server
   ```

3. **Configure your Client**:
   Change your OpenAI client `base_url` to `http://localhost:8080/v1`.

### 📚 API Reference

#### OpenAI Endpoints
- `POST /v1/chat/completions`: Chat completions mock.
- `POST /v1/images/generations`: Image generation mock.
- `POST /v1/videos/generations`: Video generation mock.

#### Vertex AI Endpoints
- `POST /v1/projects/:project/locations/:location/publishers/google/models/:model:generateContent`

### 🧪 Testing

For a comprehensive list of curl commands to test all endpoints, please refer to [test_commands.md](./test_commands.md).

### 📖 Vertex AI SDK Usage

To use this server with Vertex AI SDKs, you need to override the default API endpoint.

#### Python Example
```python
import vertexai
from vertexai.generative_models import GenerativeModel

# Initialize with the mock server endpoint
vertexai.init(
    project="your-project-id",
    location="us-central1",
    api_endpoint="http://localhost:8080"
)

model = GenerativeModel("gemini-1.5-pro")
response = model.generate_content("Hello Gemini!")
print(response.text)
```

#### Go Example
```go
import (
    "context"
    "google.golang.org/api/option"
    "github.com/google/generative-ai-go/genai"
)

ctx := context.Background()
client, _ := genai.NewClient(ctx, 
    option.WithEndpoint("http://localhost:8080/v1"),
    option.WithoutAuthentication(),
)
model := client.GenerativeModel("gemini-1.5-pro")
resp, _ := model.GenerateContent(ctx, genai.Text("Hello!"))
```

#### Admin API
- `GET /admin/rules`: List all mock rules.
- `POST /admin/rules`: Create a new rule.
- `DELETE /admin/rules/:id`: Delete a rule.

#### Admin API Examples

**Create a Chat Rule (Keyword Match)**:
```bash
curl -X POST http://localhost:8080/admin/rules \
-H "Content-Type: application/json" \
-d '{
  "name": "Joke Rule",
  "type": "fixed",
  "pattern": "joke",
  "response": "{\"choices\":[{\"message\":{\"content\":\"Why did the AI cross the road? To reach the other side of the mock!\"}}]}",
  "enabled": true
}'
```

**Create a Vertex AI Rule (Regex Match)**:
```bash
curl -X POST http://localhost:8080/admin/rules \
-H "Content-Type: application/json" \
-d '{
  "name": "Vertex Greeting",
  "type": "regex",
  "pattern": "hello.*gemini",
  "response": "{\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"Greetings from Vertex Mock!\"}]}}]}",
  "enabled": true
}'
```

**List Rules**:
```bash
curl http://localhost:8080/admin/rules
```

---

## 中文

一个基于 Go 语言开发的轻量级、高性能大模型 Mock 服务器。它模拟了 OpenAI 兼容的 API（文本对话和图片生成），旨在帮助开发者在开发大模型应用时节省昂贵的 API 调用成本，并加快测试速度。

### 🚀 功能特性

- **OpenAI 与 Vertex AI 协议兼容**：无缝替换 OpenAI 接口及 Vertex AI Gemini `generateContent` 接口。
- **动态规则配置**：通过简单的管理 API 实时管理 Mock 响应，支持：
    - **固定匹配**：匹配特定关键词触发。
    - **正则匹配**：基于正则表达式匹配复杂的提示词。
    - **随机返回**：从规则池中随机返回结果。
- **本地持久化**：规则存储在本地 SQLite 数据库中。
- **轻量化**：无外部重型依赖（使用纯 Go 实现的 SQLite）。

### 🛠️ 快速开始

1. **编译**:
   ```bash
   go build -o mock-server main.go
   ```

2. **运行**:
   ```bash
   ./mock-server
   ```

3. **配置客户端**:
   将你的 OpenAI 客户端 `base_url` 修改为 `http://localhost:8080/v1`。

### 📚 API 参考

#### OpenAI 兼容接口
- `POST /v1/chat/completions`：聊天文本生成。
- `POST /v1/images/generations`：图片生成。
- `POST /v1/videos/generations`：视频生成。

#### Vertex AI 兼容接口
- `POST /v1/projects/:project/locations/:location/publishers/google/models/:model:generateContent`

### 🧪 测试指南

有关测试所有接口的完整 curl 命令列表，请参阅 [test_commands.md](./test_commands.md)。

### 📖 Vertex AI SDK 使用指南

要在 Vertex AI SDK 中使用此 Mock 服务器，您需要覆盖默认的 API 端点。

#### Python 示例
```python
import vertexai
from vertexai.generative_models import GenerativeModel

# 初始化并指向 Mock Server 地址
vertexai.init(
    project="your-project-id",
    location="us-central1",
    api_endpoint="http://localhost:8080"
)

model = GenerativeModel("gemini-1.5-pro")
response = model.generate_content("Hello Gemini!")
print(response.text)
```

#### Go 示例
```go
import (
    "context"
    "google.golang.org/api/option"
    "github.com/google/generative-ai-go/genai"
)

ctx := context.Background()
client, _ := genai.NewClient(ctx, 
    option.WithEndpoint("http://localhost:8080/v1"),
    option.WithoutAuthentication(),
)
model := client.GenerativeModel("gemini-1.5-pro")
resp, _ := model.GenerateContent(ctx, genai.Text("Hello!"))
```

#### 管理接口 (Admin)
- `GET /admin/rules`：查看所有规则。
- `POST /admin/rules`：创建新规则。
- `DELETE /admin/rules/:id`：删除指定规则。

#### 管理接口示例

**创建聊天规则 (关键词匹配)**:
```bash
curl -X POST http://localhost:8080/admin/rules \
-H "Content-Type: application/json" \
-d '{
  "name": "笑话规则",
  "type": "fixed",
  "pattern": "笑话",
  "response": "{\"choices\":[{\"message\":{\"content\":\"为什么大模型会过马路？为了去 Mock 服务器那头！\"}}]}",
  "enabled": true
}'
```

**创建 Vertex AI 规则 (正则匹配)**:
```bash
curl -X POST http://localhost:8080/admin/rules \
-H "Content-Type: application/json" \
-d '{
  "name": "Vertex 问候",
  "type": "regex",
  "pattern": "你好.*gemini",
  "response": "{\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"来自 Vertex Mock 的问候！\"}]}}]}",
  "enabled": true
}'
```

**获取所有规则**:
```bash
curl http://localhost:8080/admin/rules
```
