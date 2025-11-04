# 密码生成
当前不直接存储密码，而是存储密码的bcrypt加密值，该值可通过以下方式计算：
```
./CrossDevicesService-{platform} genpass <password>
```

## Example
```
go run main.go genpass hello
```
输出 `$2a$10$gdt6CfQlR/OdTEWxWkV2Ve0yZiD3wBGKk5a65dc0T8bPdzsfK/8Yq`

# Set environmental variable
```
# Download and Upload directory 
export LOCAL_DIR_PATH="/path" 
```

# 启动服务
修改`conf/app.conf`的`httpport`以设置端口，启动服务命令：
```
./CrossDevicesService-{platform}
```
