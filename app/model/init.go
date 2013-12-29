package model

var (
	ArticleM *ArticleModel
	CategoryM *CategoryModel
	TagM *TagModel
	UserM *UserModel
	SessionM *SessionModel
	CommentM *CommentModel
)

// init models
func Init() {
	ArticleM = NewArticleModel()
	CategoryM = NewCategoryModel()
	TagM = NewTagModel()
	UserM = NewUserModel()
	SessionM = NewSessionModel()
	CommentM = NewCommentModel()
}
