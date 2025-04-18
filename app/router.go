package app

import (
	"github.com/gin-gonic/gin"
	"github.com/nahidhasan98/remind-name/bot"
	"github.com/nahidhasan98/remind-name/feedback"
	"github.com/nahidhasan98/remind-name/subscription"
	"github.com/nahidhasan98/remind-name/web"
)

func (app *App) RegisterRoute() {
	app.LoadHTMLGlob("view/*")
	app.Static("/assets", "./assets")

	app.GET("/", web.Index)

	subcriptionHandler := subscription.NewHandler()
	app.POST("/subscription", subcriptionHandler.AddSubscription)

	feedbackHandler := feedback.NewHandler()
	app.POST("/feedback", feedbackHandler.SaveFeedback)

	// telegram bot webhook
	app.POST("/teleweb", gin.WrapH(bot.Telegram_Webhook))

	// auto pull build restart
	// apbrRG := app.Group("/apbr")
	// apbrRG.GET("/", nil)
	// apbrRG.POST("/pull", nil)
	// apbrRG.POST("/build", nil)
	// apbrRG.POST("/restart", nil)

	// api group
	// api := app.Group(fmt.Sprintf("/api/%s", config.API_VERSION))

	// platform details
	// pdRG := api.Group("/platform")
	// pdRG.GET("/:platform", controller.GetPlatformDetails)
}
