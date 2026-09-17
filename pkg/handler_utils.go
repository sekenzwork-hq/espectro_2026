package pkg

func EmptyMessage(message string) map[string]any {
	return map[string]any{"status": 406, "message": message}
}
