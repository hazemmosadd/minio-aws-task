package engine

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	ss "test/StorageServices"
)

type router struct {
	storageServiceClient ss.StorageServiceClient
}

func NewRouter(storageServiceClient ss.StorageServiceClient) *router {
	return &router{storageServiceClient: storageServiceClient}
}

func (r *router) InstallFileHandler(engine *gin.Engine) {
	engine.POST("/upload", r.uploadFile)
	engine.GET("/signed/upload/:filename", r.getUploadSignedUrl)
	engine.GET("/signed/download/:filename", r.getDownloadSignedUrl)
	engine.GET("/download/base64/:filename", r.downloadFileBase64)
	engine.GET("/download/:filename", r.downloadFileWithUsingBuffers)
	engine.GET("/download/large/:filename", r.downloadFileLarge)

}

func (r *router) uploadFile(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	src, err := file.Open()
	// at this point we have the whole file in memory

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer src.Close()

	err = r.storageServiceClient.UploadObject("id-images", file.Filename, src, file) // this will stream the file

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})

}

func (r *router) getUploadSignedUrl(c *gin.Context) {

	url, err := r.storageServiceClient.GeneratePresignedURL("PUT", "id-images", c.Param("filename"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Use this url to upload the file", "url": url})
}

func (r *router) getDownloadSignedUrl(c *gin.Context) {
	url, err := r.storageServiceClient.GeneratePresignedURL("GET", "id-images", c.Param("filename"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Use this url to download the file", "url": url})
}

// / DOWNLOAD FILES IN DIFFERENT WAYS IN TERMS OF HOW THEY ARE SHOULD BE PROCESSED
func (r *router) downloadFileBase64(c *gin.Context) {

	bytes, err := r.storageServiceClient.DownloadObject("id-images", c.Param("filename"))
	// heer iam having the bytes in the memory
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File downloaded successfully", "data": bytes})

}

func (r *router) downloadFileWithUsingBuffers(c *gin.Context) {
	bytes, err := r.storageServiceClient.DownloadObject("id-images", c.Param("filename"))
	// here iam having the bytes in the memory
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/octet-stream") // to be downloaded as binaries
	// at the end of the day from the client the data will reach you as a stream
	// and you should decide weather to process this data chunck by chunk or wait for it to complete
	// and this will be a load on the memory
	_, err = c.Writer.Write(bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (r *router) downloadFileLarge(c *gin.Context) {
	stream, err := r.storageServiceClient.DownloadObjectStream("id-images", c.Param("filename"))
	_, err = io.Copy(c.Writer, stream)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
