package service

import (
	"datarp/database"
	"datarp/logger"
	"datarp/models"
	"fmt"
	"strconv"
)

func ValidateUIDs(oldUID, newUID string) error {
	// 转换字符串UID为整数
	oldUIDInt, err := strconv.Atoi(oldUID)
	if err != nil {
		logger.Error("老UID格式错误: " + err.Error())
		return fmt.Errorf("老UID格式错误")
	}

	newUIDInt, err := strconv.Atoi(newUID)
	if err != nil {
		logger.Error("新UID格式错误: " + err.Error())
		return fmt.Errorf("新UID格式错误")
	}

	var oldPlayer, newPlayer models.PlayerUID
	
	if err := database.DB.Where("uid = ?", oldUIDInt).First(&oldPlayer).Error; err != nil {
		logger.Error("查询老UID失败: " + err.Error())
		return fmt.Errorf("老UID不存在")
	}

	if err := database.DB.Where("uid = ?", newUIDInt).First(&newPlayer).Error; err != nil {
		logger.Error("查询新UID失败: " + err.Error())
		return fmt.Errorf("新UID不存在")
	}

	if oldPlayer.AccountUID != newPlayer.AccountUID {
		logger.Error(fmt.Sprintf("账号UID不一致: 老号=%s, 新号=%s", oldPlayer.AccountUID, newPlayer.AccountUID))
		return fmt.Errorf("账号UID不一致")
	}

	return nil
}

func MigratePlayerData(oldUID, newUID string) error {
	// 转换字符串UID为整数
	oldUIDInt, _ := strconv.Atoi(oldUID)
	newUIDInt, _ := strconv.Atoi(newUID)

	// 获取最后一位数字
	oldLastDigit := oldUID[len(oldUID)-1:]
	newLastDigit := newUID[len(newUID)-1:]

	// 获取老号数据
	var oldData models.PlayerData
	oldTableName := fmt.Sprintf("t_player_data_%s", oldLastDigit)
	if err := database.DB.Table(oldTableName).Where("uid = ?", oldUIDInt).First(&oldData).Error; err != nil {
		logger.Error("获取老号数据失败: " + err.Error())
		return fmt.Errorf("获取老号数据失败")
	}

	// 更新新号数据
	newTableName := fmt.Sprintf("t_player_data_%s", newLastDigit)
	if err := database.DB.Table(newTableName).Where("uid = ?", newUIDInt).Updates(map[string]interface{}{
		"level":     oldData.Level,
		"exp":       oldData.Exp,
		"bin_data":  oldData.BinData,
		"json_data": oldData.JsonData,
	}).Error; err != nil {
		logger.Error("更新新号数据失败: " + err.Error())
		return err
	}

	// 删除block数据
	blockTableName := fmt.Sprintf("t_block_data_%s", newLastDigit)
	if err := database.DB.Table(blockTableName).Where("uid = ?", newUIDInt).Delete(nil).Error; err != nil {
		logger.Error("删除block数据失败: " + err.Error())
		return err
	}

	logger.Info(fmt.Sprintf("成功迁移存档: %d -> %d", oldUIDInt, newUIDInt))
	return nil
}