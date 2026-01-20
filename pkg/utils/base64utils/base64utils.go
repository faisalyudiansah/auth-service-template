package base64utils

import "fmt"

func ConvertToBase64URI(base64 string) string {
	return fmt.Sprintf("data:image/png;base64,%v", base64)
}
