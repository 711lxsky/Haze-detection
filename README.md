# 🌡️ 雾霾检测工具

本着对 **B测** 认真负责的态度，三人共同从 0 到 1 开发了一款小而美的轻量级 **雾霾检测工具** 🎉！

现已部署于公网环境 [weather.hxzzz.asia](https://weather.hxzzz.asia)

项目采用前后端分离架构：
- 前端使用 **Vue** + **ElementUI** 💻
- 后端采用 **Gin** + **Gorm** 框架 🚀
- 数据库为 **MySQL** 🗂️
- 支持 **Docker** 容器化部署 📦

---

## 📂 项目整体结构

```shell
.
├── core  # 📚 后端核心模块
│   ├── config # ⚙️ 配置文件
│   │   ├── global_config.go
│   │   └── qweather_config.go
│   ├── constant # 🔧 常量定义
│   │   ├── common_error.go
│   │   └── err_info.go
│   ├── docker # 🐳 Docker容器化相关文件
│   │   ├── docker-compose.yml
│   │   ├── Dockerfile
│   │   ├── haze-detection-core-1.0.0.tar
│   │   ├── init.sql
│   ├── go.mod # 📦 项目依赖
│   ├── go.sum
│   ├── handler # 📩 请求处理
│   │   ├── qweather_request.go # 针对和风天气API调用的封装处理
│   │   ├── response.go # 响应体封装
│   │   ├── weather_city.go
│   │   ├── weather_lon_lat.go # 根据经纬度查询天气
│   │   └── weather_pos.go # 根据位置查询天气
│   ├── main
│   │   ├── main.go # 🏁 项目入口文件
│   │   ├── mock_weather.json # 模拟数据
│   │   └── router.go # 路由配置
│   ├── model # 📊 数据模型
│   │   └── query_record.go
│   ├── request # 📥 请求体封装
│   │   ├── query_position.go
│   │   └── query_weather_lon_lat.go
│   ├── sql
│   │   └── weather_district_id.csv
│   └── util # 🛠️ 工具类
│       ├── db_init.go
│       └── gin_init.go
├── README.md # 📜 项目描述
└── web
    ├── auto-imports.d.ts
    ├── components.d.ts
    ├── index.html
    ├── package.json
    ├── package-lock.json
    ├── public
    │   └── vite.svg
    ├── README.md
    ├── src
    │   ├── App.vue
    │   ├── assets
    │   │   ├── cloud.svg
    │   │   ├── leaf.svg
    │   │   ├── search.svg
    │   │   └── vue.svg
    │   ├── main.ts
    │   ├── style.css
    │   ├── ts
    │   │   ├── api.ts
    │   │   └── renderCharts.ts
    │   ├── views
    │   │   └── HomeView.vue
    │   └── vite-env.d.ts
    ├── tsconfig.app.json
    ├── tsconfig.json
    ├── tsconfig.node.json
    └── vite.config.ts
```

---

## 🚀 项目运行

### 🐳 容器化部署 (推荐)

后端提供服务镜像文件 `haze-detection-core-1.0.0.tar`，可一键导入！MySQL 镜像直接拉取即可 🎉。

相关指令:
```shell
cd /core/docker
docker load -i haze-detection-core-1.0.0.tar
docker-compose up -d
```

前端部署 🖥️
相关指令
```shell
cd  /web
# 替换 ./web/src/views/HomeView.vue 中的两个后端地址
docker build . -t haze-detection-front
docker run -d -p 3002:80 haze-detection-front
```
> 即可在本地3002端口打开
---

### 🏠 本地部署

克隆本项目到本地：
```shell
git clone https://github.com/711lxsky/Haze-detection.git
```

#### 📚 针对 `core` 模块
1. 进入 `core` 目录，执行以下命令安装依赖：  
   ```shell
   go mod tidy
   ```
2. 运行项目：  
   ```shell
   go run ./main/main.go
   ```

#### 🖥️ 针对 `web` 模块

```shell
cd ./web
npm install
npm run dev # 调试模式
```
---

希望你喜欢这个项目！如果有任何问题或建议，请随时提交 Issue 或 Pull Request 😊。
