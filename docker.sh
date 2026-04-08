#!/bin/bash

# 构建并运行 vGate 容器
echo "开始构建镜像..."
docker build --pull -t vgate_v1.1 .

if [ $? -eq 0 ]; then
    echo "构建成功，启动容器..."
    
    # 停止并删除已存在的同名容器
    docker rm -f vgate-app 2>/dev/null
    
    # 运行新容器doc
    docker run -d \
        --name vgate-app \
        --restart unless-stopped \
        -p 8080:8080 \
        -p 6789:6789 \
        vgate_v1.1
    
    if [ $? -eq 0 ]; then
        echo "容器启动成功！"
        echo "容器名称: vgate-app"
        echo "端口映射: 8080->8080, 6789->6789"
        echo "查看日志: docker logs -f vgate-app"
    else
        echo "容器启动失败！"
        exit 1
    fi
else
    echo "镜像构建失败！"
    exit 1
fi