# scaffold
go run main.go init startfront
cd startfront

# module + deps
go mod init startfront
go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/postgres
go install github.com/cosmtrek/air@latest
go get github.com/spf13/cobra
export PATH=$PATH:$(go env GOPATH)/bin

# generate an api
go run ../main.go create-api login

# hot reload server
air
