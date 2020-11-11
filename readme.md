# Cross Devices Service
Web platform for files and words transfer service, just like `ftp`.


## manual
Run `make` to start the service.  

## Download and Upload directory
Change the value of line `10` of `controllers/file.go`
```golang
const DownloadPath = "./DIR_NAME"
```
