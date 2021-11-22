package forms

type SmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Type   string `json:"type" form:"type" binding:"required,oneof=1 2"`
}
