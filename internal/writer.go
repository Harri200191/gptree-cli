package internal

import (
    "os"
)

func WriteToFile(filename, content string) error {
    return os.WriteFile(filename, []byte(content), 0644)
}
