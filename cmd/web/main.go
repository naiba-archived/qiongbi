package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/shopspring/decimal"
	"github.com/smartwalle/alipay"

	"github.com/naiba/com"
	"github.com/naiba/qiongbi/model"
)

var db *gorm.DB
var pay *alipay.AliPay

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "data/qiongbi.db")
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	db.AutoMigrate(model.Trade{})
	pay = alipay.New(os.Getenv("AppID"), os.Getenv("PubKey"), os.Getenv("PriKey"), true)
	log.Println("load alipay conig", os.Getenv("AppID"), os.Getenv("PubKey"), os.Getenv("PriKey"))
}

func main() {
	engine := gin.Default()
	engine.SetFuncMap(template.FuncMap{
		"md5": com.MD5,
		"ft": func(t time.Time) string {
			return t.Format("2006-01-02 15:04")
		},
	})
	engine.Static("/asset", "resource/asset")
	engine.LoadHTMLGlob("resource/template/*")
	engine.GET("/", home)
	engine.POST("/notify", notify)
	engine.POST("/donate", donate)
	engine.GET("/return", onReturn)
	engine.Run(":8080")
}

func onReturn(c *gin.Context) {
	c.Request.ParseForm()
	_, err := pay.VerifySign(c.Request.Form)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/")
}

type donateReq struct {
	Name   string `json:"name,omitempty" form:"name" binding:"required,max=12"`
	Email  string `json:"email,omitempty" form:"email" binding:"required,email"`
	Amount string `json:"amount,omitempty" form:"amount" binding:"required"`
	Note   string `json:"note,omitempty" form:"note" binding:"max=255"`
}

func donate(c *gin.Context) {
	var dr donateReq
	if err := c.ShouldBind(&dr); err != nil {
		c.String(http.StatusBadRequest, "输入有误："+err.Error())
		return
	}

	amt, err := decimal.NewFromString(dr.Amount)
	if err != nil {
		c.String(http.StatusBadRequest, "金额有误："+err.Error())
		return
	}
	if amt.LessThan(decimal.NewFromFloat(0.01)) {
		c.String(http.StatusBadRequest, "金额太小")
		return
	}

	var t model.Trade
	t.Name = dr.Name
	t.Email = dr.Email
	t.Amount = amt.String()
	t.Note = dr.Note

	if err := db.Create(&t).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var p = alipay.AliPayTradePagePay{}
	p.NotifyURL = "https://" + os.Getenv("Domain") + "/notify"
	p.ReturnURL = "https://" + os.Getenv("Domain") + "/return"
	p.Subject = t.Name + "捐赠" + t.Amount
	p.OutTradeNo = fmt.Sprintf("%d", t.ID)
	p.TotalAmount = t.Amount
	u, err := pay.TradePagePay(p)
	if err != nil {
		c.String(http.StatusBadRequest, "网关错误："+err.Error())
		return
	}
	c.Redirect(http.StatusFound, u.String())
}

func notify(c *gin.Context) {
	n, err := pay.GetTradeNotification(c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var t model.Trade
	if err := db.First(&t, "id = ?", n.OutTradeNo).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if n.TradeStatus == alipay.K_TRADE_STATUS_TRADE_SUCCESS {
		if err := db.Model(&t).Update("paid", true).Error; err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	pay.AckNotification(c.Writer)
}

const pageSize = 8

type sumResult struct {
	Amt decimal.Decimal
}

func home(c *gin.Context) {
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	if page == 1 {
		page = 0
	}

	var totalPage, totalNum int
	db.Table("trades").Where("paid = ?", true).Count(&totalNum)
	totalPage = totalNum / pageSize

	var ts []model.Trade
	db.Where("paid = ?", true).Order("id DESC", true).Limit(pageSize).Offset(page * pageSize).Find(&ts)
	var all sumResult
	db.Table("trades").Select("sum(amount) as amt").Where("paid = ?", true).Scan(&all)

	now := time.Now().Format("2006年1月")

	if page == 0 {
		page = 1
	}
	if totalPage == 0 {
		totalPage = 1
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"now":         now,
		"totalNum":    totalNum,
		"sum":         all.Amt,
		"trades":      ts,
		"totalPage":   totalPage,
		"currentPage": page,
	})
}
