package server

import (
	"github.com/labstack/echo"
	"bytes"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"gitlab.com/chrislewispac/rmd-server/models"
	"gopkg.in/guregu/null.v3"
	"io"
	//"os"
	//"path/filepath"
	//"strings"
)

const (
	msgSuccessOnPostFile = "Successfully uploaded a file"
	msgSuccessOnGetFileList = "Successfully fetched list of files"
	msgSuccessOnDeleteFile = "Successfully deleted a file"
	//msgSuccessOnPostProfileImage = "Successfully uploaded a profile image"
	//msgSuccessOnPostCv = "Successfully uploaded a CV file"
	//msgErrorFileTypeNotAccepted = "Unacceptable file extension"
)

type PostFileRequest struct {
	MimeType	string	`json:"mime_type"`
	FileName	string	`json:"file_name"`
	Content		[]byte	`json:"content"`
}

type PostFileResponse struct {
	FileName	string	`json:"file_name"`
	FilePath	string	`json:"file_path"`
	FileSize	int64	`json:"file_size`
}

// Download a file from Seaweedfs
func (s *Server) GetFile(c echo.Context) (err error) {
	//uid := userID(c)
	fid := c.Param("fid")
	//TODO: check user's privilege to access the file
	//
	publicUrl, _, err := s.fs.GetUrl(fid)

	if(err != nil) {
		return c.String(http.StatusOK, "Failed to get file URL: " + err.Error())
	}
	// Get the data
	resp, err := http.Get(publicUrl)
	if err != nil {
		return c.String(http.StatusOK, "Failed to get file from Seaweedfs server: " + err.Error())
	}
	defer resp.Body.Close()

	// Writer the body to file
	// TODO: append content-type from file registry
	//c.Response().Header().Set("Content-Type", "application")
	io.Copy(c.Response().Writer, resp.Body)

	return nil
}

func (s *Server) GetUserFileList(c echo.Context) (err error) {
	res := Models.NewResponse()
	//uid := userID(c)

	//
	res.Msg = msgSuccessOnGetFileList

	return c.JSON(http.StatusOK, res)
}

// Upload a file to Seaweedfs
// Accepts application/json content-type, content value is base64 encoded byte stream
func (s *Server) PostFile(c echo.Context) (err error) {
	res := Models.NewResponse()
	//uid := userID(c)
	//
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		handleErr(err)
	}
	//
	fileReq := PostFileRequest{}
	err = json.Unmarshal(b, &fileReq)
	if err != nil {
		handleErr(err)
	}

	// assign upload to get fid, and do the upload to the seaweedfs
	fid, size, err := s.fs.AssignUpload(fileReq.FileName, fileReq.MimeType, bytes.NewReader(fileReq.Content));

	if(err != nil) {
		res.Msg = errMsgExists
		res.Error = null.StringFrom(err.Error())
		return c.JSON(http.StatusOK, res)
	}

	//TODO: register user's file to the file permission registry

	// create response
	fileResp := PostFileResponse{
		FileName: fileReq.FileName,
		FilePath: "/auth/file/" + fid,
		FileSize: size,
	}
	// format response
	data, _ := json.Marshal(fileResp)
	res.Msg = msgSuccessOnPostFile
	res.Data = data

	return c.JSON(http.StatusOK, res)
}

// Delete a file to Seaweedfs
func (s *Server) DeleteFile(c echo.Context) (err error) {
	res := Models.NewResponse()

	//uid := userID(c)
	fid := c.Param("fid")
	//
	err = s.fs.Delete(fid, 1)
	if(err != nil) {
		res.Msg = errMsgExists
		res.Error = null.StringFrom(err.Error())
		return c.JSON(http.StatusOK, res)
	}

	//TODO: remove the file entry from file registry

	// format response
	res.Msg = msgSuccessOnDeleteFile

	return c.JSON(http.StatusOK, res)
}
//
//// Upload a profile image
//// Accepts x-www-form-urlencoded encoding
//func (s* Server) PostProfileImage(c echo.Context) (err error) {
//	res := Models.NewResponse()
//	uid := userID(c)
//
//	// Source file
//	file, err := c.FormFile("file")
//	if err != nil {
//		res.Msg = errMsgExists
//		res.Error = null.StringFrom(err.Error())
//		return c.JSON(http.StatusOK, res)
//	}
//	file.Header.Get()
//	src, err := file.Open()
//	if err != nil {
//		res.Msg = errMsgExists
//		res.Error = null.StringFrom(err.Error())
//		return c.JSON(http.StatusOK, res)
//	}
//	defer src.Close()
//	// check file extension
//	srcExt := filepath.Ext(file.Filename)
//	if !(strings.EqualFold(srcExt, "jpeg") || strings.EqualFold(srcExt, "jpg") || strings.EqualFold(srcExt, "png") || strings.EqualFold(srcExt, "gif")){
//		res.Msg = errMsgExists
//		res.Error = msgErrorFileTypeNotAccepted
//		return c.JSON(http.StatusOK, res)
//	}
//	// Destination
//	dstFolder := "./profile_imgs/"
//	dstFileName := uid + filepath.Ext(file.Filename)
//	dst, err := os.Create(dstFolder + dstFileName)
//	if err != nil {
//		res.Msg = errMsgExists
//		res.Error = null.StringFrom(err.Error())
//		return c.JSON(http.StatusOK, res)
//	}
//
//	if _, err = io.Copy(dst, src); err != nil {
//		res.Msg = errMsgExists
//		res.Error = null.StringFrom(err.Error())
//		return c.JSON(http.StatusOK, res)
//	}
//	// TODO: save information of uploaded cv to DB
//
//	res.Msg = msgSuccessOnPostProfileImage
//	return c.JSON(http.StatusOK, res)
//}
//
//func (s* Server) GetProfileImage(c echo.Context) (err error) {
//	uid := c.Param("uid")
//
//	dstFolder := "./profile_imgs/"
//	// TODO: retrieve information of profile image from DB
//	dstFileName := uid + ".jpg"
//
//	http.ServeFile(c.Response().Writer, c.Request(), dstFolder + dstFileName)
//	return nil;
//}
//
//// Upload a CV file
//// Accepts x-www-form-urlencoded encoding
//func (s* Server) PostCvFile(c echo.Context) (err error) {
//	res := Models.NewResponse()
//	uid := userID(c)
//
//	// Source file
//	file, err := c.FormFile("file")
//	if err != nil {
//		res.Msg = errMsgExists
//		res.Error = null.StringFrom(err.Error())
//		return c.JSON(http.StatusOK, res)
//	}
//	file.Header.Get()
//	src, err := file.Open()
//	if err != nil {
//		res.Msg = errMsgExists
//		res.Error = null.StringFrom(err.Error())
//		return c.JSON(http.StatusOK, res)
//	}
//	defer src.Close()
//	// check file extension
//	srcExt := filepath.Ext(file.Filename)
//	if !(strings.EqualFold(srcExt, "pdf") || strings.EqualFold(srcExt, "jpeg") || strings.EqualFold(srcExt, "jpg") || strings.EqualFold(srcExt, "png") || strings.EqualFold(srcExt, "gif")){
//		res.Msg = errMsgExists
//		res.Error = msgErrorFileTypeNotAccepted
//		return c.JSON(http.StatusOK, res)
//	}
//	// Destination
//	dstFolder := "./cv/"
//	dstFileName := uid + filepath.Ext(file.Filename)
//	dst, err := os.Create(dstFolder + dstFileName)
//	if err != nil {
//		res.Msg = errMsgExists
//		res.Error = null.StringFrom(err.Error())
//		return c.JSON(http.StatusOK, res)
//	}
//
//	if _, err = io.Copy(dst, src); err != nil {
//		res.Msg = errMsgExists
//		res.Error = null.StringFrom(err.Error())
//		return c.JSON(http.StatusOK, res)
//	}
//
//	// TODO: save information of uploaded cv to DB
//
//	res.Msg = msgSuccessOnPostCv
//	return c.JSON(http.StatusOK, res)
//}
//
//
//func (s* Server) GetCvFile(c echo.Context) (err error) {
//	uid := c.Param("uid")
//
//	dstFolder := "./cv/"
//	// TODO: retrieve information of profile image from DB
//	dstFileName := uid + ".jpg"
//
//	http.ServeFile(c.Response().Writer, c.Request(), dstFolder + dstFileName)
//	return nil;
//}