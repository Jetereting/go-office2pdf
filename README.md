# word to pdf 接口工具

# 直接使用
```shell
docker run --name office2pdf -p 3000:3000 -d jetereting/office2pdf
curl http://localhost:3000/convert?fileSrc=FILE_SRC
```

# 修改
## 构建二进制
```shell
GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o office2pdf main.go
```

## 构建运行 docker 服务

```shell
docker build . -t jetereting/office2pdf:latest

docker run --name office2pdf -p 3000:3000 -d office2pdf

docker logs -f office2pdf

curl http://localhost:3000/convert?fileSrc=FILE_SRC
#http://localhost:3000/convert?fileSrc=http://qn.eiyou.ga/tmp/%E6%B8%B8%E6%88%8F%E5%90%8D%E7%A7%B0%E6%B8%B8%E6%88%8F%E8%BD%AF%E4%BB%B6%E8%AF%B4%E6%98%8E%E4%B9%A6.doc
```
