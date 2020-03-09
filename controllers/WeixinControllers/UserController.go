package WeixinControllers

import (
	"fmt"
	// "path/filepath"

	// "github.com/astaxie/beego"

	// "strings"

	// "time"

	// "os"

	"github.com/hawkwang/DocHub/helper"
	// "github.com/hawkwang/DocHub/helper/conv"
	"github.com/hawkwang/DocHub/models"
	// "github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type UserController struct {
	BaseController
}

func (this *UserController) Prepare() {
	this.BaseController.Prepare()
	this.Xsrf()
}

func (this *UserController) GetToken() {
	
	var (
		currenttoken string
		tokens []interface{}
	)
	currenttoken := this.XSRFToken()
	tokens = append(tokens, currenttoken)
	fmt.Println(tokens)
	this.ResponseJson(true, "获得token", tokens)
}

//用户登录
func (this *UserController) Login() {

	if this.IsLogin > 0 {
		this.Redirect("/user", 302)
		return
	}

	// GET 请求
	if this.Ctx.Request.Method == "GET" {
		this.Data["Seo"] = models.NewSeo().GetByPage("PC-Login", "会员登录", "会员登录", "会员登录", this.Sys.Site)
		this.Data["IsUser"] = true
		this.Data["PageId"] = "wenku-reg"
		this.TplName = "login.html"
		return
	}

	type Post struct {
		Email, Password string
	}

	var post struct {
		Email, Password string
	}

	this.ParseForm(&post)
	valid := validation.Validation{}
	res := valid.Email(post.Email, "Email")
	if !res.Ok {
		this.ResponseJson(false, "登录失败，邮箱格式不正确")
	}

	ModelUser := models.NewUser()
	users, rows, err := ModelUser.UserList(1, 1, "", "", "u.`email`=? and u.`password`=?", post.Email, helper.MD5Crypt(post.Password))
	if rows == 0 || err != nil {
		if err != nil {
			helper.Logger.Error(err.Error())
		}
		this.ResponseJson(false, "登录失败，邮箱或密码不正确")
	}

	user := users[0]
	this.IsLogin = helper.Interface2Int(user["Id"])

	if this.IsLogin > 0 {
		//查询用户有没有被封禁
		if info := ModelUser.UserInfo(this.IsLogin); info.Status == false { //被封禁了
			this.ResponseJson(false, "登录失败，您的账号已被管理员禁用")
		}
		this.BaseController.SetCookieLogin(this.IsLogin)
		this.ResponseJson(true, "登录成功")
	}
	this.ResponseJson(false, "登录失败，未知错误！")
}

//用户退出登录
func (this *UserController) Logout() {
	this.ResetCookie()
	if v, ok := this.Ctx.Request.Header["X-Requested-With"]; ok && v[0] == "XMLHttpRequest" {
		this.ResponseJson(true, "退出登录成功")
	}
	this.Redirect("/", 302)
}


// 检测用户是否已登录
func (this *UserController) CheckLogin() {
	if this.BaseController.IsLogin > 0 {
		this.ResponseJson(true, "已登录")
	}
	this.ResponseJson(false, "您当前处于未登录状态，请先登录")
}


