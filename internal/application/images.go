package application

import "encoding/json"

// marshalImages превращает массив URL в JSON-строку для хранения в БД.
// nil или пустой массив сериализуются как "[]" — чтобы в БД не было NULL.
func marshalImages(images []string) string {
	if len(images) == 0 {
		return "[]"
	}
	b, err := json.Marshal(images)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// parseImages десериализует JSON-строку из БД в []string.
// При ошибке парсинга возвращает пустой слайс (не nil), чтобы клиент всегда получал массив.
func parseImages(s string) []string {
	if s == "" || s == "[]" {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return []string{}
	}
	return out
}

// pickCover выбирает «главную обложку» поста для предпросмотра в списке.
// Если задан явный image_url — он приоритетнее. Иначе берётся первое из массива.
func pickCover(explicit string, images []string) string {
	if explicit != "" {
		return explicit
	}
	if len(images) > 0 {
		return images[0]
	}
	return ""
}
