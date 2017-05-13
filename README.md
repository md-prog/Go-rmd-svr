# Recruit Me Docs - API Server - Dev


[![build status](https://gitlab.com/chrislewispac/rmd-server/badges/dev/build.svg)](https://gitlab.com/chrislewispac/rmd-server/commits/dev)

[![coverage report](https://gitlab.com/chrislewispac/rmd-server/badges/dev/coverage.svg)](https://gitlab.com/chrislewispac/rmd-server/commits/dev)

# Please do all of your work on dev branch or features/x branch!

To Run Application
-------------------
```
git clone ...
cd rmd-server
glide install
touch config.yml //<-- put config contents here, must be valid yaml (check spaces)
go get github.com/codegangsta/gin
gin run main.go
```


Testing
---------

#### To run tests continually with web ui:

```
$GOPATH/bin/goconvey
```

#### To run tests in terminal only

```
go test -v $(go list ./... | grep -v /vendor/)
```

#### Api design

#### [CLICK HERE FOR API SPECS](https://docs.google.com/document/d/11sxHwO3Ti7Ea-Vq0No4MMBfrbOgJaEUZTVjglronlfI/edit?usp=sharing)

We only return status 200. Please review code examples for createXErrorResposne and createXSuccessResponse methods.

```
func createAuthErrorResponse(user Models.User, errMsg string) *Models.Res {
	user.Password = ""
	user.ID = ""
	anon := struct {
		User Models.User `json:"user"`
	}{
		user,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = errMsgExists
	res.Data = data
	res.Error = null.StringFrom(errMsg)
	return res
}

func createAuthSuccessResponse(user Models.User, successMsg string) *Models.Res {
	user.Password = ""
	anon := struct {
		User Models.User `json:"user"`
	}{
		user,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = successMsg
	res.Data = data
	return res
}

```

Please review the file [Auth Test](https://gitlab.com/chrislewispac/rmd-server/blob/dev/server/auth_test.go) for examples on how to test the API routes given our structure.

### Dependencies

[Glide - Dependency Management](https://github.com/Masterminds/glide)

[Gin - Livereload](https://github.com/codegangsta/gin)

[Plivo - SMS](https://github.com/micrypt/go-plivo)

[Goconvey - Testing](https://github.com/smartystreets/goconvey)

[Gabs - Json parsing](https://github.com/Jeffail/gabs)

[Configor - Configurations](https://github.com/jinzhu/configor)

[Echo - Routing/Framework](https://github.com/labstack/echo)

[SQLX/pq - Postgres intereactions and map to struct](https://github.com/jmoiron/sqlx)

[Redis - Sessions](https://github.com/go-redis/redis)

[SparkPost - Email Handling](https://github.com/SparkPost/gosparkpost)

[Tesseract - OCR Engine](https://github.com/otiai10/gosseract)

[Logging with Logrus](https://github.com/Sirupsen/logrus)

[JWT](https://github.com/dgrijalva/jwt-go)

[Redis on Server](https://github.com/sameersbn/docker-redis)

[Seaweedfs on Server](https://github.com/chrislusf/seaweedfs)

[Weedo - Seaweedfs client](https://github.com/ginuerzh/weedo)