package profiledto

import "literatur/models"

type ProfileResponse struct {
	ID      int                         `json:"id" gorm:"primary_key:auto_increment"`
	User    models.UsersProfileResponse `json:"user"`
	UserID  int                         `json:"user_id"`
	Phone   string                      `json:"phone" gorm:"type: varchar(255)"`
	Gender  string                      `json:"gender" gorm:"type: varchar(255)"`
	Address string                      `json:"address" gorm:"type: text"`
}
