###
```shell
# 打镜像
docker build -t isc-envoy-control-service:1.0.0 .
# 镜像保存
docker save isc-envoy-control-service:1.0.0 -o isc-envoy-control-service.tar
# 镜像文件上传
scp -v isc-envoy-control-service.tar root@10.30.30.78:/root/zhouzy/isc-envoy-control-service.tar
# 开发环境中载入
docker load -i /root/zhouzy/isc-envoy-control-service.tar
```

