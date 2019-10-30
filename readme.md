# FileService-go
web文件传输系统，可在服务器上实现类似ftp文件传输功能；采用beego框架开发  

## manual
Run `make` to start the service.  

### 设置下载和上传的文件夹
修改controllers/file.go的第10行的变量
```golang
const DownloadPath = "./upload"
```

