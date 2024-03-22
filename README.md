# okx_test
okx api测试学习

```bash
go mod tidy
go run main.go

# 访问127.0.0.1:8080使用
```

- tips

```golang
	client := &http.Client{
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(&url.URL{Host: "localhost:7890", Scheme: "http"}), // 设置Clash代理的地址和端口
		//},
	}// 使用代理就取消注释
```
