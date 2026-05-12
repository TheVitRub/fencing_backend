package http

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// allowedExts — расширения файлов, разрешённые к загрузке.
// Расширены до распространённых форматов изображений (jpg/png/webp/gif/svg).
var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
	".svg":  true,
}

// UploadConfig — параметры загрузки, прокидываются из конфигурации приложения.
type UploadConfig struct {
	Dir       string // путь к директории на диске
	URLPrefix string // префикс URL для доступа клиента
	MaxSizeMB int    // максимальный размер одного файла в МБ
}

// uploadHandler принимает один файл (поле "file" в multipart-форме)
// и сохраняет его в директорию из конфига под уникальным именем.
// Возвращает { "url": "/uploads/<random>.<ext>" }.
//
// Требует JWT (роут зарегистрирован в защищённой группе /api/admin).
func (h *Handler) uploadHandler(c *gin.Context) {
	if h.upload.Dir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "загрузка файлов не настроена"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "файл не передан: " + err.Error()})
		return
	}

	// Проверка размера
	maxBytes := int64(h.upload.MaxSizeMB) * 1024 * 1024
	if maxBytes > 0 && file.Size > maxBytes {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "файл слишком большой (максимум " + itoa(h.upload.MaxSizeMB) + " МБ)",
		})
		return
	}

	// Проверка расширения
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "недопустимое расширение файла (разрешены: jpg, png, webp, gif, svg)",
		})
		return
	}

	// Генерируем уникальное имя — 16 случайных байт в hex
	name, err := randomFileName(ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка генерации имени файла"})
		return
	}

	dst := filepath.Join(h.upload.Dir, name)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		h.logger.Error("не удалось сохранить файл",
			h.logger.F("error", err),
			h.logger.F("dst", dst),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить файл"})
		return
	}

	url := strings.TrimRight(h.upload.URLPrefix, "/") + "/" + name
	h.logger.Info("файл загружен",
		h.logger.F("name", name),
		h.logger.F("size", file.Size),
	)
	c.JSON(http.StatusOK, gin.H{
		"url":  url,
		"name": file.Filename,
		"size": file.Size,
	})
}

// randomFileName генерирует имя вида <16-byte-hex><ext>.
func randomFileName(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b) + ext, nil
}

// itoa — мини-itoa без аллокаций для подстановки в сообщения об ошибках.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
