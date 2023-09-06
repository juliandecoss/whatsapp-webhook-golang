package services

import (
	"encoding/json"
	"fmt"
)

func Logger(data map[string]string) {
	loggerData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error al formatear JSON:", err)
	}
	fmt.Print(string(loggerData))
}
