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
	for _, user := range config.Users {
		if user.Username == username {
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
