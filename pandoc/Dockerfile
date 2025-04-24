# 使用 Python 3.9 基础镜像
FROM python:3.9

# 更新包列表并安装 Pandoc
RUN apt-get update && apt-get install -y pandoc

# 设置工作目录
WORKDIR /app

# 复制当前目录内容到工作目录
COPY . /app

# 安装依赖
RUN pip install --no-cache-dir flask

# 暴露端口
EXPOSE 80

# 运行应用
CMD ["python", "app.py"]
