package model

import "github.com/fuxiaohei/GoBlog/app"

var (
	ArticleM  *ArticleModel
	CategoryM *CategoryModel
	TagM      *TagModel
	UserM     *UserModel
	SessionM  *SessionModel
	CommentM  *CommentModel
	SettingM *SettingModel
)

// init models
func Init() {
	ArticleM = NewArticleModel()
	CategoryM = NewCategoryModel()
	TagM = NewTagModel()
	UserM = NewUserModel()
	SessionM = NewSessionModel()
	CommentM = NewCommentModel()
	SettingM = NewSettingModel()

	// do some more
	app.Ink.View.NewFunc("Setting", SettingM.GetItem)
}
