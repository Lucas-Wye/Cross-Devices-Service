package controllers

import (
	models "CrossDevicesService/models"

	"fmt"
	"github.com/astaxie/beego"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

var DownloadPath = models.GetLocalDirPath()

const CopyPastePath = "/copy"
const ServiceName = "Cross Devices Service"
const CopyPasteFile = "./logs/copy_paste.txt"

type FileController struct {
	beego.Controller
	user *models.User
}

func (this *FileController) Prepare() {
	sessionUser := this.GetSession("user")
	if sessionUser == nil {
		this.Redirect("/login", 302)
		this.StopRun()
	}

	// It's better to store just the username in session and refetch user data
	// to ensure permissions are always up-to-date.
	user, ok := models.GetUser(sessionUser.(*models.User).Username)
	if !ok {
		this.Redirect("/login", 302)
		this.StopRun()
	}
	this.user = user
}

func (c *FileController) GetUpload() {
	// 读取上传成功的提示参数
	if c.GetString("success") == "1" {
		c.Data["success"] = true
		if v := c.GetSession("last_filename"); v != nil {
			if s, ok := v.(string); ok {
				c.Data["filename"] = s
			}
			c.DelSession("last_filename")
		}
	}
	c.Data["user"] = c.user
	// 如果之前在同一请求中设置了错误信息，直接展示
	c.TplName = "upload.html"
}

func (this *FileController) Upload() {
	dir := this.GetString("dir")
	// fmt.Println("[55] dir =", dir)
	hasPermission := false
	for _, p := range this.user.Permissions {
		if (dir == "" && p.Path == "shared") || (dir == p.Path) {
			if p.Write {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		this.Data["error"] = "没有上传权限"
		this.TplName = "upload.html"
		return
	}

	f, h, err := this.GetFile("uploadfile") // 读取上传文件
	if err != nil {
		this.Data["error"] = "未选择文件或读取文件失败: " + err.Error()
		this.TplName = "upload.html"
		return
	}
	defer f.Close()

	targetDir := DownloadPath
	if dir != "" {
		targetDir = filepath.Join(DownloadPath, dir)
		// fmt.Println("[82] Target directory:", targetDir)
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			os.MkdirAll(targetDir, 0755)
		}
		if err != nil {
			this.Data["error"] = err.Error()
			this.TplName = "upload.html"
			return
		}
	}

	// 计算目标文件名（清洗、截断以避免超长/非法字符），若已存在则自动加后缀(1),(2)...
	// 先取 basename 以避免浏览器上传携带路径（IE/Edge 历史行为）
	origName := filepath.Base(h.Filename)
	fmt.Println("[105]", origName)
	cleanName := sanitizeFilename(origName, 180, false) // 预留空间给后缀
	fmt.Println("[107]", cleanName)
	finalName, destPath := ensureUniquePath(targetDir, cleanName, 180)
	beego.Info("Upload saving to:", destPath, "(name bytes:", len([]byte(finalName)), ")")

	if err := this.SaveToFile("uploadfile", destPath); err != nil {
		this.Data["error"] = "保存文件失败: " + err.Error()
		this.TplName = "upload.html"
		return
	}

	// 成功后重定向到上传页并附带成功提示（避免重复提交）
	// 注意：不要把超长文件名放进 URL，避免 Location 头过长导致报错
	// 将文件名暂存到 session，在 GetUpload 时读取并清理
	// 为安全起见限制存储的展示长度
	showName := truncateUTF8Bytes(finalName, 200)
	this.SetSession("last_filename", showName)
	redirectURL := "/upload?success=1"
	if dir != "" {
		redirectURL += "&dir=" + url.QueryEscape(dir)
	}
	this.Redirect(redirectURL, 302)
}

func (this *FileController) GetList() {
	var trees []models.FileNode
	for _, p := range this.user.Permissions {
		if p.Read {
			dirPath := filepath.Join(DownloadPath, p.Path)
			if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
				if root, err := models.GetFileTree(DownloadPath, p.Path); err == nil && root != nil {
					trees = append(trees, *root)
				}
			}
		}
	}
	this.Data["trees"] = trees
	if this.user.Role == "admin" {
		this.TplName = "download_admin.html"
	} else if this.user.Role == "user" {
		this.TplName = "download.html"
	}
}

func (this *FileController) DeleteFiles() {
	if this.user.Role != "admin" {
		this.Abort("403")
		return
	}

	filesToDelete := this.GetStrings("files")
	if len(filesToDelete) == 0 {
		this.Redirect("/download", 302)
		return
	}

	for _, fileRelPath := range filesToDelete {
		fullPath := filepath.Join(DownloadPath, fileRelPath)
		// Use RemoveAll to delete files or directories
		if err := os.RemoveAll(fullPath); err != nil {
			beego.Error("Failed to delete file:", fullPath, "Error:", err)
			// You might want to add some error feedback to the user
		}
	}

	this.Redirect("/download", 302)
}

func (this *FileController) Download() {
	// 说明：部分 Mac 上的 Chrome 会错误地将响应保存为当前页面的 HTML，常见原因：
	// 1. Content-Type 与文件扩展不匹配；2. Content-Disposition 中 filename 未正确加引号或编码；
	// 3. 额外写入了模板/文本导致混合响应。此实现改进了响应头设置并避免重复覆盖。
	dir := this.GetString("dir")
	filename := this.GetString("filename")

	// 基本权限校验
	hasPermission := false
	for _, p := range this.user.Permissions {
		if ((dir == "" && p.Path == "shared") || (dir == p.Path)) && p.Read {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		this.Abort("403")
		return
	}

	// 构造物理路径（防止出现重复斜杠）
	cleanDir := strings.TrimPrefix(strings.Trim(dir, "/"), "./")
	targetPath := filepath.Join(DownloadPath, cleanDir, filename)

	// 文件存在校验
	fi, err := os.Stat(targetPath)
	if err != nil || fi.IsDir() { // 目录不允许直接下载
		this.Abort("404")
		return
	}

	// 自动推断 MIME（zip 专门处理）
	ct := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".zip" {
		ct = "application/zip"
	} else if m := mime.TypeByExtension(ext); m != "" {
		ct = m
	}

	// RFC 6266/5987 兼容的文件名（处理空格、非 ASCII）；为兼容性限制长度
	headerName := truncateUTF8Bytes(filename, 150)
	asciiName := truncateUTF8Bytes(sanitizeFilename(headerName, 150, true), 150)
	escapedUTF8 := url.PathEscape(headerName)
	contentDisp := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", asciiName, escapedUTF8)

	o := this.Ctx.Output
	o.Header("Content-Type", ct)
	o.Header("Content-Disposition", contentDisp)
	o.Header("Content-Transfer-Encoding", "binary")
	o.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	o.Header("Pragma", "no-cache")
	o.Header("Expires", "0")

	// 使用框架的下载方法：内部会调用 http.ServeFile，不再额外写模板
	o.Download(targetPath)
	this.StopRun() // 确保后续不再渲染其他模板
}

// sanitizeFilename 清理文件名，使其对文件系统和 HTTP 头部安全。
// - 移除控制字符与平台非法字符（/ \ : * ? " < > |）。
// - 将非 ASCII 字符替换为下划线（如果 asciiOnly 为 true）。
// - 移除首尾的空格和点。
// - 截断到 maxBytes 长度，同时保留文件扩展名。
func sanitizeFilename(name string, maxBytes int, asciiOnly bool) string {
	// 移除控制字符和非法字符
	filtered := strings.Map(func(r rune) rune {
		// 非法字符列表
		if strings.ContainsRune(`/\\:*?"<>|`, r) {
			return '_'
		}
		// 根据参数决定是否替换非 ASCII 字符
		if asciiOnly && r > unicode.MaxASCII {
			return '_'
		}
		// 移除控制字符
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name)

	// 移除首尾的空格和点
	filtered = strings.Trim(filtered, ". ")
	if filtered == "" {
		filtered = "file"
	}

	// 分离扩展名
	ext := filepath.Ext(filtered)
	base := strings.TrimSuffix(filtered, ext)

	// 截断 base 部分以适应总长度限制
	extBytes := len([]byte(ext))
	allowBase := maxBytes - extBytes
	if allowBase < 1 {
		// 扩展名过长，则丢弃扩展名
		ext = ""
		allowBase = maxBytes
	}
	base = truncateUTF8Bytes(base, allowBase)

	if base == "" {
		base = "file"
		base = truncateUTF8Bytes(base, allowBase)
	}

	return base + ext
}

// ensureUniquePath: 确保文件名在 targetDir 下唯一，必要时追加 (1),(2)... 并确保总字节 <= maxBytes
func ensureUniquePath(targetDir, name string, maxBytes int) (finalName, destPath string) {
	tryName := name
	destPath = filepath.Join(targetDir, tryName)
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return tryName, destPath
	}

	base := strings.TrimSuffix(tryName, filepath.Ext(tryName))
	ext := filepath.Ext(tryName)
	for i := 1; ; i++ {
		suffix := fmt.Sprintf("(%d)", i)
		// 计算可用字节留给 base
		allowBase := maxBytes - len([]byte(ext)) - len([]byte(suffix))
		if allowBase < 1 {
			allowBase = 1
		}
		baseTrunc := truncateUTF8Bytes(base, allowBase)
		candidate := baseTrunc + suffix + ext
		candidatePath := filepath.Join(targetDir, candidate)
		if _, err := os.Stat(candidatePath); os.IsNotExist(err) {
			return candidate, candidatePath
		}
	}
}

// truncateUTF8Bytes 将字符串按 UTF-8 字节长度安全截断到不超过 maxBytes
func truncateUTF8Bytes(s string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(s) <= maxBytes {
		return s
	}
	// 逐个 rune 累计避免截断多字节
	var outLen int
	for i, r := range s {
		rl := utf8.RuneLen(r)
		if rl < 0 {
			rl = 1
		}
		if outLen+rl > maxBytes {
			return s[:i]
		}
		outLen += rl
	}
	return s
}

func (c *FileController) Paste() {
	// c.Data["input"] = models.ReadFromCopyPaste()
	c.Data["input"] = models.ReadFromFile(CopyPasteFile)
	c.TplName = "copy_paste.html"
}

func (c *FileController) Copy() {
	input := c.GetString("input")
	if input != "" {
		// models.Write2CopyPaste(input)
		models.WriteToFile(CopyPasteFile, "\n"+input)
	}
	c.Redirect(CopyPastePath, 302)
}
