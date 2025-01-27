package models

type PlayerUID struct {
    UID        int    `gorm:"column:uid"`
    AccountUID string `gorm:"column:account_uid"`
}

func (PlayerUID) TableName() string {
    return "t_player_uid"
}

type PlayerData struct {
    UID      int    `gorm:"column:uid"`
    Level    int    `gorm:"column:level"`
    Exp      int    `gorm:"column:exp"`
    BinData  []byte `gorm:"column:bin_data;type:mediumblob"`
}

func (PlayerData) TableName() string {
    return "" // 将在service层动态设置表名
} 