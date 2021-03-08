package controllers

import (
	"github.com/Xhofe/alist/server/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

// list request bean
type ListReq struct {
	Dir      string `json:"dir" binding:"required"`
	Password string `json:"password"`
}

// handle list request
func List(c *gin.Context) {
	var list ListReq
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(200, MetaResponse(400, "Bad Request:"+err.Error()))
		return
	}
	log.Debugf("list:%+v", list)
	// find folder model
	dir, name := filepath.Split(list.Dir)
	file, err := models.GetFileByDirAndName(dir, name)
	if err != nil {
		// folder model not exist
		if file == nil {
			c.JSON(200, MetaResponse(404, "folder not found."))
			return
		}
		c.JSON(200, MetaResponse(500, err.Error()))
		return
	}
	// check password
	if file.Password != "" && file.Password != list.Password {
		if list.Password == "" {
			c.JSON(200, MetaResponse(401, "need password."))
		} else {
			c.JSON(200, MetaResponse(401, "wrong password."))
		}
		return
	}
	files, err := models.GetFilesByDir(list.Dir + "/")
	if err != nil {
		c.JSON(200, MetaResponse(500, err.Error()))
		return
	}
	// delete password
	for i, _ := range *files {
		(*files)[i].Password = ""
	}
	c.JSON(200, DataResponse(files))
}