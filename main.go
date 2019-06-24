package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
)

type CommandLine struct {
	UserName string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	App      string `form:"app" json:"app" binding:"required"`
	Path     string `form:"path" json:"path" `
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	cmd := router.Group("/cmd")
	cmd.POST("/input", InputCommand)

	if err := router.Run(); err != nil {
		log.Printf("%s", err)
	}
}

func InputCommand(c *gin.Context) {
	var cmd CommandLine
	if err := c.ShouldBind(&cmd); err != nil {
		log.Printf("%s", err)
	}
	if cmd.UserName == "admin" && cmd.Password == "admin" {
		switch cmd.App {
		case "mkdir":
			MakeDir(c, cmd.Path)
			break
		case "rm":
			DeleteFile(c, cmd.Path)
			break
		case "rm-r":
			DeleteDir(c, cmd.Path)
			break
		case "ls":
			ListFile(c, cmd.Path)
		default:
			Help(c)
			break
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Info": "沒有權限。"})
	}
}
func Help(c *gin.Context) {
	c.String(200,
		"[Create a dir]app:mkdir;path:/file/to/path \n"+
			"[Delete a file]app:rm;path:/file/to/path \n"+
			"[Delete a dir]app:rm-r;path:/file/to/path\n"+
			"[List current dir files]app:ls;path:/file/to/path")
}
func MakeDir(c *gin.Context, path string) {
	InpCMD := exec.Command("mkdir", path)
	_, err := InpCMD.Output()
	if err != nil {
		log.Printf("%s", err)
		c.JSON(200, gin.H{"MakeDir": err.Error()})
	} else {
		c.JSON(200, gin.H{"MakeDir": "成功"})
	}
}
func DeleteFile(c *gin.Context, path string) {
	InpCMD := exec.Command("rm", path)
	_, err := InpCMD.Output()
	if err != nil {
		log.Printf("%s", err)
		c.JSON(200, gin.H{"DeleteFile": err.Error()})
	} else {
		c.JSON(200, gin.H{"DeleteFile": "成功"})
	}
}

func DeleteDir(c *gin.Context, path string) {
	InpCMD := exec.Command("rm", "-r", path)
	_, err := InpCMD.Output()
	if err != nil {
		log.Printf("%s", err)
		c.JSON(200, gin.H{"DeleteDir": err.Error()})
	} else {
		c.JSON(200, gin.H{"DeleteDir": "成功"})
	}
}

func ListFile(c *gin.Context, path string) {
	InpCMD := exec.Command("ls", "-a", path)
	listFile, err := InpCMD.Output()
	if err != nil {
		log.Printf("%s", err)
		c.JSON(200, gin.H{"ListFile": err.Error()})
	} else {
		c.String(200, string(listFile))
	}
}
