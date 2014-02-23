package plugin

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/fuxiaohei/GoBlog/app/model"
	"github.com/fuxiaohei/GoInk"
	"html/template"
	"net/mail"
	"net/smtp"
	"path"
	"strings"
)

var (
	// email template original from http://ben-lab.com/tech/1789.html
	emailReplyTemplate = `<table style="width: 99.8%;height:99.8% "><tbody><tr><td style="background:#FAFAFA">
    <div style="background-color:white;border-top:2px solid #0079BC;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
        <h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;"><span style="color: #0079bc;font-weight: bold;">&gt; </span>您在<a style="text-decoration:none;color: #e64346;" href="{{.link}}"> {{.site}} </a>上的留言有回复啦！</h2>
        <div style="padding:0 12px 0 12px;margin-top:18px">
            <p><strong>{{.author_p}}</strong> 同学，您曾在文章《{{.title}}》上发表评论:</p>
            <p style="background-color: #f5f5f5;border: 0px solid #DDD;padding: 10px 15px;margin:18px 0">{{.text_p}}</p>
            <p><strong>{{.author}}</strong> 给您的回复如下:</p>
            <p style="background-color: #f5f5f5;border: 0px solid #DDD;padding: 10px 15px;margin:18px 0">{{.text}}</p>
            <p>您可以点击 <a style="text-decoration:none; color:#0079BC" href="{{.permalink}}">查看回复的完整內容 </a>，欢迎再次光临 <a style="text-decoration:none; color:#0079bc" href="{{.link}}">{{.site}} </a>。</p>
        </div>
    </div>
</td></tr></tbody></table>`
	emailCreatedTemplate = `<table style="width: 99.8%;height:99.8% "><tbody><tr><td style="background:#FAFAFA">
    <div style="background-color:white;border-top:2px solid #0079BC;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
        <h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;"><span style="color: #0079bc;font-weight: bold;">&gt; </span>{{.author}} 在<a style="text-decoration:none;color: #e64346;" href="{{.link}}"> {{.site}} </a>上的文章 《{{.title}}》 发表了评论！</h2>
        <div style="padding:0 12px 0 12px;margin-top:18px">
            <p><strong>{{.author}}</strong> 同学，在文章《{{.title}}》上发表评论:</p>
            <p style="background-color: #f5f5f5;border: 0px solid #DDD;padding: 10px 15px;margin:18px 0">{{.text}}</p>
            <p>您可以点击 <a style="text-decoration:none; color:#0079BC" href="{{.permalink}}">查看完整內容 </a>，欢迎再次光临 <a style="text-decoration:none; color:#0079bc" href="{{.link}}">{{.site}} </a>。</p>
        </div>
    </div>
</td></tr></tbody></table>`
)

// Email plugin instance, contains activation status, handler status, settings and templates.
type EmailPlugin struct {
	isActive            bool
	isHandlerRegistered bool
	settings            map[string]string
	templates           map[string]*template.Template
}

func init() {
	EmailPlugin := new(EmailPlugin)
	EmailPlugin.isActive = false
	EmailPlugin.isHandlerRegistered = false
	EmailPlugin.settings = make(map[string]string)
	EmailPlugin.templates = make(map[string]*template.Template)
	// load templates
	t1, e := template.New("reply").Parse(emailReplyTemplate)
	if e != nil {
		panic(e)
	}
	EmailPlugin.templates["reply"] = t1
	t2, e2 := template.New("created").Parse(emailCreatedTemplate)
	if e2 != nil {
		panic(e)
	}
	EmailPlugin.templates["created"] = t2
	register(EmailPlugin)
}

// Name returns email plugin name
func (p *EmailPlugin) Name() string {
	return "邮件提醒"
}

// Key returns email plugin unique key.
func (p *EmailPlugin) Key() string {
	return "email_notify_plugin"
}

// Desc returns email plugin descripition.
func (p *EmailPlugin) Desc() string {
	return "评论及回复等邮件提醒"
}

// ToStorage creates plugin storage map for plugin list json.
func (p *EmailPlugin) ToStorage() map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = p.Name()
	m["description"] = p.Desc()
	m["is_activate"] = p.isActive
	return m
}

// Activate activates email plugin.
// It will change activate status to on.
// If the handlers are not registered, it will add handlers to router table.
func (p *EmailPlugin) Activate() {
	if p.isHandlerRegistered {
		p.isActive = true
		return
	}
	if model.Storage.Has("plugin/" + p.Key()) {
		model.Storage.Get("plugin/"+p.Key(), &p.settings)
	}
	// email middleware handler
	fn := func(context *GoInk.Context) {
		context.On("comment_created", func(co interface{}) {
			if !p.isActive {
				return
			}
			go func() {
				if len(p.settings["smtp_host"]) < 1 || len(p.settings["smtp_email_user"]) < 1 || len(p.settings["smtp_email_password"]) < 1 {
					println("no setting in smtp email plugin")
					return
				}
				p.sendEmail(co.(*model.Comment), true)
			}()
		})
		context.On("comment_reply", func(co interface{}) {
			if !p.isActive {
				return
			}
			go func() {
				if len(p.settings["smtp_host"]) < 1 || len(p.settings["smtp_email_user"]) < 1 || len(p.settings["smtp_email_password"]) < 1 {
					fmt.Println("no setting in smtp email plugin")
					return
				}
				p.sendEmail(co.(*model.Comment), false)
			}()
		})
	}
	Handler("comment_email_notify", fn, false)
	p.SetSetting(p.settings)
	p.isActive = true
	p.isHandlerRegistered = true
}

// Deactivate closes this email plugin.
func (p *EmailPlugin) Deactivate() {
	p.isActive = false
}

// IsActive returns plugin activation status.
func (p *EmailPlugin) IsActive() bool {
	return p.isActive
}

// Version returns plugin version string.
func (p *EmailPlugin) Version() string {
	return "0.1.5"
}

// HasSetting returns bool of setting form existing.
func (p *EmailPlugin) HasSetting() bool {
	return true
}

// Form returns the form template string for email plugin.
func (p *EmailPlugin) Form() string {
	return `<h4 class="title">
        SMTP 邮箱提醒设置
    </h4>
    <p class="item">
        <label for="smtp-addr">SMTP服务器</label>
        <input class="ipt" id="smtp-addr" type="text" name="smtp_host" placeholder="SMTP服务器地址，如 smtp.gmail.com:465" value="` + p.settings["smtp_host"] + `"/>
    </p>
    <p class="item">
        <label for="smtp-email">SMTP邮箱</label>
        <input class="ipt" id="smtp-email" type="email" name="smtp_email_user" placeholder="使用SMTP的邮箱" value="` + p.settings["smtp_email_user"] + `"/>
    </p>
    <p class="item">
        <label for="smtp-password">邮箱密码</label>
        <input class="ipt" id="smtp-password" type="text" name="smtp_email_password" placeholder="邮箱密码" value="` + p.settings["smtp_email_password"] + `"/>
    </p>
    <p class="submit item">
        <label>&nbsp;</label>
        <button class="btn btn-blue">保存设置</button>
        <span class="tip">暂不支持验证，请自行保证正确</span>
    </p>`
}

// SetSettings saves plugin settings to json.
func (p *EmailPlugin) SetSetting(settings map[string]string) {
	p.settings = settings
	model.Storage.Set("plugin/"+p.Key(), p.settings)
}

func (p *EmailPlugin) sendEmail(co *model.Comment, isCreate bool) {
	var (
		tpl     *template.Template
		buff    bytes.Buffer
		pco     *model.Comment
		content *model.Content
		err     error
		user    *model.User
		title   string
		from    mail.Address
		to      mail.Address
	)

	// get article or page content instance
	content = model.GetContentById(co.Cid)
	if content == nil {
		println("error content getting in email plugin")
		return
	}
	// email for notify new commment creation
	if co.Pid < 1 {
		tpl = p.templates["created"]

		err = tpl.Execute(&buff, map[string]interface{}{
			"link":      model.GetSetting("site_url"),
			"site":      model.GetSetting("site_title"),
			"author":    co.Author,
			"text":      template.HTML(co.Content),
			"title":     content.Title,
			"permalink": path.Join(model.GetSetting("site_url"), content.Link()),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		user = model.GetUsersByRole("ADMIN")[0]
		from = mail.Address{"no-reply@" + model.GetSetting("site_url"), p.settings["smtp_email_user"]}
		to = mail.Address{user.Nick, user.Email}
		title = co.Author + "在您的网站发表新评论"
		p.sendSmtp(from, to, title, buff.Bytes())
		return
	}
	// send mail for the author of comment replying to
	pco = model.GetCommentById(co.Pid)
	tpl = p.templates["reply"]
	err = tpl.Execute(&buff, map[string]interface{}{
		"link":      model.GetSetting("site_url"),
		"site":      model.GetSetting("site_title"),
		"author_p":  pco.Author,
		"text_p":    template.HTML(pco.Content),
		"author":    co.Author,
		"text":      template.HTML(co.Content),
		"title":     content.Title,
		"permalink": path.Join(model.GetSetting("site_url"), content.Link()),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	user = model.GetUsersByRole("ADMIN")[0]

	from = mail.Address{pco.Author + "@" + model.GetSetting("site_url"), p.settings["smtp_email_user"]}
	to = mail.Address{pco.Author, pco.Email}
	title = "您的评论有了回复"
	p.sendSmtp(from, to, title, buff.Bytes())
	// send email to notice admin new comment creation
	// this comment is a reply comment
	if isCreate {
		go func() {
			tpl = p.templates["created"]

			err = tpl.Execute(&buff, map[string]interface{}{
				"link":      model.GetSetting("site_url"),
				"site":      model.GetSetting("site_title"),
				"author":    co.Author,
				"text":      template.HTML("回复" + pco.Author + ":<br/>" + co.Content),
				"title":     content.Title,
				"permalink": path.Join(model.GetSetting("site_url"), content.Link()),
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			user = model.GetUsersByRole("ADMIN")[0]
			from = mail.Address{"no-reply@" + model.GetSetting("site_url"), p.settings["smtp_email_user"]}
			to = mail.Address{user.Nick, user.Email}
			title = co.Author + "在您的网站发表新评论"
			p.sendSmtp(from, to, title, buff.Bytes())
		}()
	}
}

// send email via smtp
func (p *EmailPlugin) sendSmtp(from mail.Address, to mail.Address, title string, body []byte) {
	// set email vars
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	host := p.settings["smtp_host"]
	email := p.settings["smtp_email_user"]

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(title)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString(body)

	auth := smtp.PlainAuth(
		"",
		email,
		p.settings["smtp_email_password"],
		strings.Split(host, ":")[0],
	)

	err := smtp.SendMail(
		host,
		auth,
		email,
		[]string{to.Address},
		[]byte(message),
	)
	println("send to email ", to.String(), fmt.Sprint(err))
}
