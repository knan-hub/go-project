# 使用官方的 Golang 基础镜像
FROM golang:1.23 AS build

# 设置工作目录
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和 go.sum 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目源代码到容器中
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-project .

# 使用轻量的 alpine 安装 git 后的镜像上传到腾讯云镜像仓库作为基础镜像
# 此处可根据实际情况补充使用新基础镜像的逻辑

# 向量数据库
RUN mkdir -p /git_temp
RUN mkdir -p /knowledgeBase_temp

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件到最终镜像
COPY --from=build /app/go-project .
COPY --from=build /app/config/settings.yaml .
# 可根据实际需求决定是否保留下面这行
# COPY --from=build /app .

# 设置环境变量
ENV PORT=80

# 暴露端口
EXPOSE $PORT

# 启动应用程序
CMD ["./go-project", "./settings.yaml"]
