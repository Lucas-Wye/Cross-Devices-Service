package models

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Permission struct {
	Path  string `json:"path"`
	Read  bool   `json:"read"`
	Write bool   `json:"write"`
}

type User struct {
	Username     string       `json:"username"`
	PasswordHash string       `json:"passwordHash"`
	Role         string       `json:"role"`
	Permissions  []Permission `json:"permissions"`
}

type UserConfig struct {
	Users []User `json:"users"`
}

var (
	userConfig *UserConfig
	once       sync.Once
	configLock = new(sync.RWMutex)
)

func loadUsers() {
	file, err := os.Open("conf/users.json")
	if err != nil {
		logs.Error("Failed to open users.json: %v", err)
		userConfig = &UserConfig{}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := UserConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		logs.Error("Failed to decode users.json: %v", err)
		userConfig = &UserConfig{}
		return
	}
	userConfig = &config
}

func GetUserConfig() *UserConfig {
	once.Do(loadUsers)
	configLock.RLock()
	defer configLock.RUnlock()
	return userConfig
}

func GetUser(username string) (*User, bool) {
	config := GetUserConfig()
	for i, user := range config.Users {
		if user.Username == username {
			// 如果是admin，自动拥有所有文件夹的读写权限
			if user.Role == "admin" {
				// 获取所有文件夹名
				folders, err := os.ReadDir(GetLocalDirPath())
				if err == nil {
					permMap := make(map[string]bool)
					for _, p := range user.Permissions {
						permMap[p.Path] = true
					}
					for _, f := range folders {
						if f.IsDir() {
							name := f.Name()
							if !permMap[name] {
								user.Permissions = append(user.Permissions, Permission{
									Path:  name,
									Read:  true,
									Write: true,
								})
							}
						}
					}
					// 更新内存中的用户权限
					config.Users[i] = user
				}
			}
			return &user, true
		}
	}
	return nil, false
}

func SaveUserConfig(config *UserConfig) error {
	configLock.Lock()
	defer configLock.Unlock()

	file, err := os.Create("conf/users.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(config); err != nil {
		return err
	}
	userConfig = config
	return nil
}
