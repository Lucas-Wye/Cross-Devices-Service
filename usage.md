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

在Linux中，由于`$`符号的特殊性，设置环境变量时需要使用转移字符`\$`

# Set environmental variable
```
# Download and Upload directory 
export LOCAL_DIR_PATH="/path"

# admin user
export ADMIN_USERNAME="" 
export ADMIN_PASSWORD="" 

# user
export NORMAL_USERNAME=""
export NORMAL_PASSWORD="" 
```

# 启动服务
修改`conf/app.conf`的`httpport`以设置端口，启动服务命令：
```
./CrossDevicesService-{platform}
```
