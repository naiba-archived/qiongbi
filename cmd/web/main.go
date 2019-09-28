package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/shopspring/decimal"
	"github.com/smartwalle/alipay"

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
	db.AutoMigrate(model.Trade{})
	pay = alipay.New(os.Getenv("AppID"), os.Getenv("PubKey"), os.Getenv("PriKey"), true)
}

func main() {
	engine := gin.Default()
	engine.Static("/asset", "resource/asset")
	engine.GET("/", func(c *gin.Context) {
		c.File("resource/template/index.html")
	})
	engine.POST("/notify", notify)
	engine.POST("/donate", donate)
	engine.GET("/return", onReturn)
	engine.Run(":8080")
}

type donateReq struct {
	Name   string `json:"name,omitempty" form:"name" binding:"required,max=12"`
	Email  string `json:"email,omitempty" form:"email" binding:"required,email"`
	Amount string `json:"amount,omitempty" form:"amount" binding:"required"`
	Note   string `json:"note,omitempty" form:"note" binding:"max=255"`
}

func onReturn(c *gin.Context) {
	c.Request.ParseForm()
	_, err := pay.VerifySign(c.Request.Form)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusOK, "/")
}

func donate(c *gin.Context) {
	var dr donateReq
	if err := c.ShouldBind(&dr); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	amt, err := decimal.NewFromString(dr.Amount)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if amt.LessThan(decimal.NewFromFloat(0.01)) {
		c.String(http.StatusBadRequest, "数额太小")
		return
	}

	var t model.Trade
	t.Name = dr.Name
	t.Email = dr.Email
	t.Amount = dr.Amount
	db.Create(t)

	var p = alipay.AliPayTradeWapPay{}
	p.NotifyURL = "https://" + os.Getenv("Domain") + "/notify"
	p.ReturnURL = "https://" + os.Getenv("Domain") + "/return"
	p.Subject = t.Name + "-捐助-" + t.Amount
	p.OutTradeNo = fmt.Sprintf("%d", t.ID)
	p.TotalAmount = t.Amount
	u, err := pay.TradeWapPay(p)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Redirect(http.StatusOK, u.String())
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
		t.Paid = true
		if err := db.Save(&t).Error; err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	pay.AckNotification(c.Writer)
}
