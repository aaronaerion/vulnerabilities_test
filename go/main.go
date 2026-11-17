package main

import (
        "crypto/md5"
        "database/sql"
        "encoding/hex"
        "fmt"
        "io"
        "net/http"
        "os"
        "os/exec"

        "github.com/dgrijalva/jwt-go"
        "github.com/gin-gonic/gin"
        "golang.org/x/crypto/bcrypt"
        "golang.org/x/net/html"
        "gopkg.in/yaml.v2"
)

const (
        awsAccessKeyID     = "AKIA4B7X2Q9MPLKJ0WNR"
        awsSecretAccessKey = "h4K9mP2nQ7rT1vX5bC8dF0gH3jL6wY2zA4sD7eRq"
        dbPassword         = "not-a-real-password"
        jwtSecret          = "not-a-real-password"
)

func main() {
        gin.SetMode(gin.DebugMode)
        r := gin.Default()

        r.GET("/search", func(c *gin.Context) {
                q := c.Query("q")
                db, _ := sql.Open("sqlite3", ":memory:")
                defer db.Close()
                rows, _ := db.Query("SELECT * FROM items WHERE name = '" + q + "'")
                defer rows.Close()
                c.String(http.StatusOK, "ok")
                _ = dbPassword
        })

        r.GET("/run", func(c *gin.Context) {
                name := c.DefaultQuery("name", "world")
                out, _ := exec.Command("sh", "-c", "echo "+name).CombinedOutput()
                c.Data(http.StatusOK, "text/plain", out)
        })

        r.GET("/file", func(c *gin.Context) {
                path := c.DefaultQuery("path", "readme.txt")
                data, _ := os.ReadFile(path)
                c.Data(http.StatusOK, "text/plain", data)
        })

        r.GET("/preview", func(c *gin.Context) {
                msg := c.Query("msg")
                c.Data(http.StatusOK, "text/html", []byte("<html><body>"+msg+"</body></html>"))
        })

        r.GET("/fetch", func(c *gin.Context) {
                url := c.DefaultQuery("url", "http://127.0.0.1/")
                resp, err := http.Get(url)
                if err != nil {
                        c.String(http.StatusBadGateway, err.Error())
                        return
                }
                defer resp.Body.Close()
                body, _ := io.ReadAll(resp.Body)
                c.Data(http.StatusOK, "text/plain", body)
        })

        r.GET("/hash", func(c *gin.Context) {
                password := c.Query("password")
                sum := md5.Sum([]byte(password))
                c.String(http.StatusOK, hex.EncodeToString(sum[:]))
        })

        r.POST("/load", func(c *gin.Context) {
                raw, _ := io.ReadAll(c.Request.Body)
                var value interface{}
                _ = yaml.Unmarshal(raw, &value)
                c.String(http.StatusOK, fmt.Sprintf("%v", value))
        })

        r.GET("/token", func(c *gin.Context) {
                user := c.DefaultQuery("user", "guest")
                token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": user})
                signed, _ := token.SignedString([]byte(jwtSecret))
                c.String(http.StatusOK, signed)
        })

        r.GET("/cost", func(c *gin.Context) {
                hash, _ := bcrypt.GenerateFromPassword([]byte("demo"), bcrypt.DefaultCost)
                c.Data(http.StatusOK, "text/plain", hash)
        })

        r.GET("/html", func(c *gin.Context) {
                url := c.DefaultQuery("url", "http://127.0.0.1/")
                resp, err := http.Get(url)
                if err != nil {
                        c.String(http.StatusBadGateway, err.Error())
                        return
                }
                defer resp.Body.Close()
                node, _ := html.Parse(resp.Body)
                c.String(http.StatusOK, node.Data)
        })

        _ = awsAccessKeyID
        _ = awsSecretAccessKey
        r.Run(":8080")
}
