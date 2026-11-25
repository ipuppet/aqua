package handler

import (
	"crypto/rsa"
	"strings"

	"github.com/gin-gonic/gin"

	"aqua/pkg/apperr"
	"aqua/pkg/rbac"
	"aqua/pkg/transport/middleware"
	"aqua/pkg/transport/response"
	"aqua/services/admin/internal/httpbind"
	"aqua/services/admin/logic"
)

// FileHandler exposes file upload/delete APIs.
type FileHandler struct {
	fileService *logic.FileService
	publicKey   *rsa.PublicKey
	rbacReader  rbac.Reader
}

// NewFileHandler constructs the handler.
func NewFileHandler(fileService *logic.FileService, rbacReader rbac.Reader, publicKey *rsa.PublicKey) *FileHandler {
	if fileService == nil {
		return nil
	}
	return &FileHandler{
		fileService: fileService,
		publicKey:   publicKey,
		rbacReader:  rbacReader,
	}
}

// Register binds the file routes.
func (h *FileHandler) Register(e *gin.Engine) {
	if h == nil {
		return
	}
	r := e.Group("/api/file")
	r.Use(middleware.Authenticate(h.publicKey))
	r.Use(middleware.RequireAnyPermission(h.rbacReader, "blog.manager"))
	r.GET("/image/:module", h.fsListImages)
	r.POST("/image/:module", h.fsUploadImage)
	r.DELETE("/image/*path", h.fsDeleteImage)
}

// FSRequestURI captures upload route params.
type FSRequestURI struct {
	Module string `uri:"module" binding:"required"`
}

// FSDeleteURI captures delete route params.
type FSDeleteURI struct {
	Path string `uri:"path" binding:"required"` // e.g. blog/2025-02-27/xxx.png
}

// 文件处理函数
func (h *FileHandler) fsUploadImage(c *gin.Context) {
	var uriParam FSRequestURI
	if !httpbind.BindURI(c, &uriParam) {
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		response.Fail(c, apperr.BadRequest("image is required"))
		return
	}

	src, err := file.Open()
	if err != nil {
		response.OKOrFail(c, err)
		return
	}
	defer func() { _ = src.Close() }()

	info, err := h.fileService.UploadImage(c.Request.Context(), uriParam.Module, file.Filename, src)
	response.OKWithDataOrFail(c, info, err)
}

func (h *FileHandler) fsDeleteImage(c *gin.Context) {
	var uriParam FSDeleteURI
	if !httpbind.BindURI(c, &uriParam) {
		return
	}
	path := strings.TrimPrefix(uriParam.Path, "/")
	if path == "" {
		response.Fail(c, apperr.BadRequest("path is required"))
		return
	}

	err := h.fileService.DeleteImage(c.Request.Context(), path)
	response.OKOrFail(c, err)
}

func (h *FileHandler) fsListImages(c *gin.Context) {
	var uriParam FSRequestURI
	if !httpbind.BindURI(c, &uriParam) {
		return
	}

	refresh := c.Query("refresh") == "true"
	images, err := h.fileService.ListImages(c.Request.Context(), uriParam.Module, refresh)
	response.OKWithDataOrFail(c, images, err)
}
