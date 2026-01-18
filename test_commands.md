# LLM Mock Server Test Commands / 测试命令指南

[English](#english) | [中文](#中文)

---

<a name="english"></a>
## English

Use these commands to verify that your Mock Server is working correctly.

> [!NOTE]
> Make sure the server is running on `http://localhost:8080` before executing these commands.

### 1. OpenAI Compatible Endpoints

#### Chat Completion (Default)
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "Hello!"}]}'
```

#### Image Generation (Default)
```bash
curl -X POST http://localhost:8080/v1/images/generations \
-H "Content-Type: application/json" \
-d '{"prompt": "A futuristic city", "n": 1, "size": "1024x1024"}'
```

#### Video Generation (Default)
```bash
curl -X POST http://localhost:8080/v1/videos/generations \
-H "Content-Type: application/json" \
-d '{"prompt": "A waterfall in the forest"}'
```

### 2. Vertex AI (Gemini) Endpoint

#### Generate Content (Default)
```bash
curl -X POST http://localhost:8080/v1/projects/my-project/locations/us-central1/publishers/google/models/gemini-1.5-pro:generateContent \
-H "Content-Type: application/json" \
-d '{"contents": [{"role": "user", "parts": [{"text": "Explain quantum physics."}]}]}'
```

### 3. Admin API & Rule Testing

#### Step 3.1: Create a Rule
```bash
curl -X POST http://localhost:8080/admin/rules \
-H "Content-Type: application/json" \
-d '{
  "name": "Weather Rule",
  "type": "fixed",
  "pattern": "weather",
  "response": "{\"choices\":[{\"message\":{\"content\":\"The weather is perfectly sunny in the mock world!\"}}]}",
  "enabled": true
}'
```

#### Step 3.2: Test the Rule (OpenAI)
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "How is the weather?"}]}'
```

---

<a name="中文"></a>
## 中文

使用以下命令验证您的 Mock 服务器是否正常运行。

> [!NOTE]
> 在执行这些命令之前，请确保服务器运行在 `http://localhost:8080`。

### 1. OpenAI 兼容接口测试

#### 文本对话 (默认)
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "你好！"}]}'
```

#### 图片生成 (默认)
```bash
curl -X POST http://localhost:8080/v1/images/generations \
-H "Content-Type: application/json" \
-d '{"prompt": "未来的城市", "n": 1, "size": "1024x1024"}'
```

#### 视频生成 (默认)
```bash
curl -X POST http://localhost:8080/v1/videos/generations \
-H "Content-Type: application/json" \
-d '{"prompt": "森林里的瀑布"}'
```

### 2. Vertex AI (Gemini) 接口测试

#### 内容生成 (默认)
```bash
curl -X POST http://localhost:8080/v1/projects/my-project/locations/us-central1/publishers/google/models/gemini-1.5-pro:generateContent \
-H "Content-Type: application/json" \
-d '{"contents": [{"role": "user", "parts": [{"text": "解释一下量子物理。"}]}]}'
```

### 3. 管理接口与规则测试

#### 步骤 3.1: 创建一个规则
```bash
curl -X POST http://localhost:8080/admin/rules \
-H "Content-Type: application/json" \
-d '{
  "name": "天气规则",
  "type": "fixed",
  "pattern": "天气",
  "response": "{\"choices\":[{\"message\":{\"content\":\"Mock 世界里天气晴朗！\"}}]}",
  "enabled": true
}'
```

#### 步骤 3.2: 测试该规则 (OpenAI)
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "天气怎么样？"}]}'
```
