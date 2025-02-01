package handlers

import (
	"datarp/logger"
	"datarp/service"
	"encoding/json"
	"net/http"
)

type MigrateRequest struct {
	OldUID string `json:"oldUid"`
	NewUID string `json:"newUid"`
}

type MigrateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func HandleMigrate(w http.ResponseWriter, r *http.Request) {
	// 添加通用的 CORS 头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	// 处理 OPTIONS 预检请求
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req MigrateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("解析请求失败: " + err.Error())
		sendResponse(w, false, "请求格式错误")
		return
	}

	// 验证UID
	if err := service.ValidateUIDs(req.OldUID, req.NewUID); err != nil {
		sendResponse(w, false, err.Error())
		return
	}

	// 执行迁移
	if err := service.MigratePlayerData(req.OldUID, req.NewUID); err != nil {
		sendResponse(w, false, "迁移失败："+err.Error())
		return
	}

	sendResponse(w, true, "迁移成功")
}

func sendResponse(w http.ResponseWriter, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MigrateResponse{
		Success: success,
		Message: message,
	})
} 