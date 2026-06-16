package tasklog

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/diagnosis"
	"github.com/gocronx-team/gocron/internal/modules/i18n"
	"github.com/gocronx-team/gocron/internal/modules/llm"
	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/routers/base"
)

// Diagnose 调用 LLM 对某条任务执行日志做失败归因与修复建议。
func Diagnose(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		base.RespondError(c, i18n.T(c, "param_error"))
		return
	}

	logModel := new(models.TaskLog)
	if err := logModel.Find(id); err != nil {
		base.RespondError(c, i18n.T(c, "log_not_found"))
		return
	}
	if strings.TrimSpace(logModel.Result) == "" {
		base.RespondError(c, i18n.T(c, "log_no_result"))
		return
	}

	client, err := llm.FromSettings()
	if err != nil {
		base.RespondError(c, i18n.T(c, "llm_not_configured"))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), diagnosis.Timeout)
	defer cancel()

	d, err := diagnosis.Diagnose(ctx, client, logModel, i18n.GetLocale(c) == i18n.EnUS)
	if err != nil {
		logger.Errorf("日志诊断#调用LLM失败#日志ID-%d#%s", id, err)
		base.RespondError(c, i18n.T(c, "llm_call_failed"))
		return
	}

	base.RespondSuccess(c, i18n.T(c, "operation_success"), gin.H{
		"root_cause":  d.RootCause,
		"suggestions": d.Suggestions,
	})
}
