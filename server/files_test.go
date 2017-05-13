package server

import  (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"

	"gitlab.com/chrislewispac/rmd-server/models"
	"io/ioutil"
	"encoding/json"
)

func TestServer_GetUserFileList(t *testing.T) {
	return
}

func TestServer_CrudFile(t *testing.T) {
	var filePath string
	var fileSize int64
	Convey("when you POST /auth/file", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + "/auth/file"

		t.Log("Seaweedfs master url: " + ts.fs.Master().Url)

		Convey("with a valid file", func() {
			validFileReq := PostFileRequest{
				FileName: "test.json",
				MimeType: "application/json",
			}
			fileContent, err := ioutil.ReadFile("test_files/email_test.json")
			So(err, ShouldEqual, nil)

			validFileReq.Content = fileContent
			validStr, err := json.Marshal(validFileReq)
			So(err, ShouldEqual, nil)

			Convey("for an invalid user", func() {
				tStr := createInvalidTestUser(ts)

				req, err := http.NewRequest(echo.POST, url, bytes.NewReader(validStr))
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, body := doRequest(req)

				So(code, ShouldEqual, 200)
				So(body.Path("error").Data(), ShouldEqual, nil)
			})

			Convey("for a valid user", func() {

				_, tStr := createTestUser(ts)

				Convey("with a valid file, it should be uploaded", func() {
					req, err := http.NewRequest(echo.POST, url, bytes.NewReader(validStr))
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					code, res := doRequestRes(req)

					So(code, ShouldEqual, 200)
					So(res.Error.String, ShouldEqual, "")
					So(res.Msg, ShouldEqual, msgSuccessOnPostFile)
					So(res.Data, ShouldNotEqual, nil)
					var data PostFileResponse
					So(json.Unmarshal(res.Data, &data), ShouldEqual, nil)
					So(data.FileName, ShouldEqual, validFileReq.FileName)
					So(data.FileSize, ShouldNotEqual, 0)
					So(data.FilePath, ShouldNotEqual, nil)
					filePath = data.FilePath
					fileSize = data.FileSize

				})
			})
		})
	})


	Convey("when you GET /auth/file/:fid", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + filePath

		Convey("for an invalid user", func() {
			tStr := createInvalidTestUser(ts)

			Convey("you should get an error", func() {
				req, err := http.NewRequest(echo.GET, url+"123435", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))

				code, _ := doRequestRaw(req)

				So(code, ShouldEqual, 200)
			})
		})

		Convey("for a valid user", func() {
			_, tStr := createTestUser(ts)

			Convey("with an non-existent file id, you should get an error", func() {
				req, err := http.NewRequest(echo.GET, url+"123435", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))

				code, _ := doRequestRaw(req)

				So(code, ShouldEqual, 200)
				// TODO : check correct error message
			})

			Convey("with a valid fid", func() {
				Convey("and the matching user, you should get the file", func() {
					req, err := http.NewRequest(echo.GET, url, nil)
					So(err, ShouldEqual, nil)

					req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))

					code, body := doRequestRaw(req)

					So(code, ShouldEqual, 200)
					So(len(body), ShouldBeGreaterThanOrEqualTo, fileSize)
				})
			})
		})
	})
	Convey("when you DELETE /auth/file/:fid", t, func() {
		ts, addr, close := startTestServer(t)
		defer close()

		url := addr + filePath

		Convey("for an invalid user", func() {
			tStr := createInvalidTestUser(ts)

			Convey("you should get an error", func() {
				req, err := http.NewRequest(echo.DELETE, url, nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, _ := doRequest(req)

				So(code, ShouldEqual, 200)
				// TODO: check correct return message when authentication fails
				//
				// So(body.Path("error").Data(), ShouldEqual, errMsgSomethingWrong)
			})
		})

		Convey("for a valid user", func() {
			_, tStr := createTestUser(ts)

			Convey("but a non-existent fid, you should get an error", func() {
				req, err := http.NewRequest(echo.DELETE, addr + "/auth/file/12345", nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, res := doRequestRes(req)

				So(code, ShouldEqual, 200)
				So(res.Msg, ShouldEqual, errMsgExists)
			})

			Convey("and a valid fid, the file should be deleted", func() {
				req, err := http.NewRequest(echo.DELETE, url, nil)
				So(err, ShouldEqual, nil)

				req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", tStr))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				code, res := doRequestRes(req)

				So(code, ShouldEqual, 200)
				So(res.Msg, ShouldEqual, msgSuccessOnDeleteFile)
				So(res.Error, ShouldEqual, "")
				t.Log(res.Error)
			})
		})
	})
}

// creates a test user,
// sets it to the test server environment,
// and returns the token
func createTestUser(ts *Server) (string, string) {
	user := new(Models.User)
	err := ts.insertUser("test@test.com", nil, user)
	So(err, ShouldEqual, nil)

	token := generateToken(user.ID, 10*time.Minute, false)
	tStr, err := token.SignedString(SigningKey)
	So(err, ShouldEqual, nil)

	err = ts.storeToken(tStr, 10*time.Minute)
	So(err, ShouldEqual, nil)
	return user.ID, tStr
}

func createInvalidTestUser(ts *Server) string {
	token := generateToken("invalidUserId", 10*time.Minute, false)
	tStr, err := token.SignedString(SigningKey)
	So(err, ShouldEqual, nil)

	err = ts.storeToken(tStr, 10*time.Minute)
	So(err, ShouldEqual, nil)
	return tStr
}