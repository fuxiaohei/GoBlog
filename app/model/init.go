package model

var (
	ArticleM *ArticleModel
	CategoryM *CategoryModel
	TagM *TagModel
	UserM *UserModel
	SessionM *SessionModel
)

// init models
func Init() {
	ArticleM = NewArticleModel()
	CategoryM = NewCategoryModel()
	TagM = NewTagModel()
	UserM = NewUserModel()
	SessionM = NewSessionModel()
}
