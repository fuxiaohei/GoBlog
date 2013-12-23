package model

var (
	CategoryM *CategoryModel
	TagM *TagModel
	UserM *UserModel
	SessionM *SessionModel
)

// init models
func Init() {
	CategoryM = NewCategoryModel()
	TagM = NewTagModel()
	UserM = NewUserModel()
	SessionM = NewSessionModel()
}
