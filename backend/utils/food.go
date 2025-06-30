package utils

import (
	"fmt"
	"main/global"
	"main/models"
)

func DeleteFoodByID(foodid int) error {
	if foodid <= 0 {
		return fmt.Errorf("invalid food ID: %d", foodid)
	}

	if err := global.DB.Where("id = ?", foodid).Delete(&models.Food{}).Error; err != nil {
		return fmt.Errorf("failed to delete food with ID %d: %v", foodid, err)
	}

	return nil
}
