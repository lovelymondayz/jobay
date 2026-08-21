package handlers

import (
	"jobay/internal/models"
)

func nextID(items interface{}) int {
	switch v := items.(type) {
	case []models.Job:
		max := 0
		for _, item := range v {
			if item.ID > max {
				max = item.ID
			}
		}
		return max + 1
	case []models.Action:
		max := 0
		for _, item := range v {
			if item.ID > max {
				max = item.ID
			}
		}
		return max + 1
	}
	return 1
}
