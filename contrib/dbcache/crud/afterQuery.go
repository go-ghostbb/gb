package crud

import (
	"fmt"
	"gorm.io/gorm"
)

func (h *Handler) afterQuery(db *gorm.DB) {
	fmt.Println("after query")
	fmt.Println("after query")
}
