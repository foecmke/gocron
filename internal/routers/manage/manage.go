package manage

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/modules/utils"
)

func Slack(c *gin.Context) {
	settingModel := new(models.Setting)
	slack, err := settingModel.Slack()
	jsonResp := utils.JsonResponse{}
	var result string
	if err != nil {
		logger.Error(err)
		result = jsonResp.Success(utils.SuccessContent, nil)
	} else {
		result = jsonResp.Success(utils.SuccessContent, slack)
	}
	c.String(http.StatusOK, result)
}

func UpdateSlack(c *gin.Context) {
	url := strings.TrimSpace(c.Query("url"))
	template := strings.TrimSpace(c.Query("template"))
	settingModel := new(models.Setting)
	err := settingModel.UpdateSlack(url, template)
	result := utils.JsonResponseByErr(err)
	c.String(http.StatusOK, result)
}

func CreateSlackChannel(c *gin.Context) {
	channel := strings.TrimSpace(c.Query("channel"))
	settingModel := new(models.Setting)
	var result string
	if settingModel.IsChannelExist(channel) {
		jsonResp := utils.JsonResponse{}
		result = jsonResp.CommonFailure("Channel已存在")
	} else {
		_, err := settingModel.CreateChannel(channel)
		result = utils.JsonResponseByErr(err)
	}
	c.String(http.StatusOK, result)
}

func RemoveSlackChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	settingModel := new(models.Setting)
	_, err := settingModel.RemoveChannel(id)
	result := utils.JsonResponseByErr(err)
	c.String(http.StatusOK, result)
}

// endregion

// region 邮件
func Mail(c *gin.Context) {
	settingModel := new(models.Setting)
	mail, err := settingModel.Mail()
	jsonResp := utils.JsonResponse{}
	var result string
	if err != nil {
		logger.Error(err)
		result = jsonResp.Success(utils.SuccessContent, nil)
	} else {
		result = jsonResp.Success("", mail)
	}
	c.String(http.StatusOK, result)
}

type MailServerForm struct {
	Host     string `binding:"Required;MaxSize(100)"`
	Port     int    `binding:"Required;Range(1-65535)"`
	User     string `binding:"Required;MaxSize(64);Email"`
	Password string `binding:"Required;MaxSize(64)"`
}

func UpdateMail(c *gin.Context) {
	var form MailServerForm
	if err := c.ShouldBindJSON(&form); err != nil {
		json := utils.JsonResponse{}
		result := json.CommonFailure("表单验证失败, 请检测输入")
		c.String(http.StatusOK, result)
		return
	}
	
	jsonByte, _ := json.Marshal(form)
	settingModel := new(models.Setting)
	template := strings.TrimSpace(c.Query("template"))
	err := settingModel.UpdateMail(string(jsonByte), template)
	result := utils.JsonResponseByErr(err)
	c.String(http.StatusOK, result)
}

func CreateMailUser(c *gin.Context) {
	username := strings.TrimSpace(c.Query("username"))
	email := strings.TrimSpace(c.Query("email"))
	settingModel := new(models.Setting)
	var result string
	if username == "" || email == "" {
		jsonResp := utils.JsonResponse{}
		result = jsonResp.CommonFailure("用户名、邮箱均不能为空")
	} else {
		_, err := settingModel.CreateMailUser(username, email)
		result = utils.JsonResponseByErr(err)
	}
	c.String(http.StatusOK, result)
}

func RemoveMailUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	settingModel := new(models.Setting)
	_, err := settingModel.RemoveMailUser(id)
	result := utils.JsonResponseByErr(err)
	c.String(http.StatusOK, result)
}

func WebHook(c *gin.Context) {
	settingModel := new(models.Setting)
	webHook, err := settingModel.Webhook()
	jsonResp := utils.JsonResponse{}
	var result string
	if err != nil {
		logger.Error(err)
		result = jsonResp.Success(utils.SuccessContent, nil)
	} else {
		result = jsonResp.Success("", webHook)
	}
	c.String(http.StatusOK, result)
}

func UpdateWebHook(c *gin.Context) {
	url := strings.TrimSpace(c.Query("url"))
	template := strings.TrimSpace(c.Query("template"))
	settingModel := new(models.Setting)
	err := settingModel.UpdateWebHook(url, template)
	result := utils.JsonResponseByErr(err)
	c.String(http.StatusOK, result)
}

// endregion
