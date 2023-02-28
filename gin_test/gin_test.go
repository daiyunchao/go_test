package gin_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGinStart(t *testing.T) {
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}

// 获取Get参数
// http://localhost:8001/user/tom/18
func TestHandleGetParams(t *testing.T) {
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/user/:name/:age", func(c *gin.Context) {
		name := c.Param("name")
		age := c.Param("age")
		c.String(http.StatusOK, fmt.Sprint("Hello MyName is "+name+", age is "+age))
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// Url参数
// http://localhost:8001/user?name=tom&age=19
func TestUrlParams(t *testing.T) {
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/user", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		c.String(http.StatusOK, fmt.Sprint("Hello MyName is "+name+", age is "+age))
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// 路由组
// http://localhost:8001/v1/user?name=zhangsan&age=19
// http://localhost:8001/v2/user?name=zhangsan&age=19
func TestRouteGroup(t *testing.T) {
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	v1 := r.Group("v1")
	v2 := r.Group("v2")
	v1.GET("/user", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		c.String(http.StatusOK, fmt.Sprint("V1 MyName is "+name+", age is "+age))
	})
	v2.GET("/user", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		c.String(http.StatusOK, fmt.Sprint("V2 MyName is "+name+", age is "+age))
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// 404
func Test404Info(t *testing.T) {
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/user", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		c.String(http.StatusOK, fmt.Sprint("Hello MyName is "+name+", age is "+age))
	})
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprint("404"))
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// 绑定 Json转结构体
type User struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Nickname string `json:"nickname"`
}

// 将请求data转换成结构体
// curl http://127.0.0.1:8001/user -X POST -d "{\"name\":\"Tom\",\"age\":18,\"nickname\":\"TomCat\"}"
func TestJson2Struct(t *testing.T) {
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.POST("/user", func(c *gin.Context) {
		var user User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			return
		}

		c.String(http.StatusOK, fmt.Sprint("Hello MyName is "+user.Name+", age is "+strconv.Itoa(user.Age)))
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// 各种响应方式
func TestResponse(t *testing.T) {
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response

	//JSON
	r.GET("/json", func(c *gin.Context) {
		user := User{
			Name:     "TOM",
			Age:      19,
			Nickname: "Cat",
		}
		c.JSON(200, gin.H{"data": user, "status": 200})
	})

	//STRUCT
	r.GET("/struct", func(c *gin.Context) {
		user := User{
			Name:     "TOM",
			Age:      19,
			Nickname: "Cat",
		}
		c.JSON(200, user)
	})

	//ProtoBuff
	r.GET("/proto", func(c *gin.Context) {
		label := "Label"
		data := &protoexample.Test{
			Label: &label,
			Reps:  []int64{int64(1), int64(2)},
		}
		c.ProtoBuf(200, data)
	})

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// 模板渲染
func TestTmpl(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLGlob("tmpl/*")
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "后端构建的标题", "ce": "后端构建的内容"})
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// 静态路由
// http://localhost:8001/tmpl/
func TestStatic(t *testing.T) {
	r := gin.Default()
	r.Static("tmpl", "./tmpl")
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

func BeforeGlobalMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now().String()
		c.Set("st", st)
	}
}

// 全局中间件
func TestGlobalMiddleWare(t *testing.T) {
	r := gin.Default()
	r.Use(BeforeGlobalMiddleWare())
	r.GET("/user", func(c *gin.Context) {
		st, ok := c.Get("st")
		if ok {
			c.String(200, st.(string))
		} else {
			c.String(200, "Not Found")
		}
	})
	r.Static("tmpl", "./tmpl")
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

func BeforeGlobalMiddleWare2() gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()
		c.Set("st", st)
		c.Next()
		newSt, _ := c.Get("st")
		diffTime := time.Since(newSt.(time.Time))
		fmt.Println("执行API消耗时间为:", diffTime)
	}
}

// Next
func TestNextMiddleWare(t *testing.T) {
	r := gin.Default()
	r.Use(BeforeGlobalMiddleWare2())
	r.GET("/user", func(c *gin.Context) {
		_, ok := c.Get("st")
		time.Sleep(time.Second)
		if ok {
			c.String(200, "Hello")
		} else {
			c.String(200, "Not Found")
		}
	})
	r.Static("tmpl", "./tmpl")
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// 局部中间件
func TestMethodMiddleWare(t *testing.T) {
	r := gin.Default()
	r.GET("/has", BeforeGlobalMiddleWare2(), func(c *gin.Context) {
		_, ok := c.Get("st")
		time.Sleep(time.Second)
		if ok {
			c.String(200, "Hello")
		} else {
			c.String(200, "Not Found")
		}
	})
	r.GET("/no", func(c *gin.Context) {
		time.Sleep(time.Second)
		c.String(200, "Hello")
	})
	r.Static("tmpl", "./tmpl")
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

// 测试设置cookie
func TestCookie(t *testing.T) {
	r := gin.Default()
	r.GET("/user", func(c *gin.Context) {
		cookie, err := c.Cookie("loginName")
		if err != nil {
			c.SetCookie("loginName", "TOM", 60, "/", "localhost", false, true)
		} else {
			fmt.Println("cookieLoginName: ", cookie)
		}

	})
	r.Static("tmpl", "./tmpl")
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8001")
}

type Order struct {
	OrderId   int    `json:"orderId" binding:"required,gt=10"`
	OrderName string `json:"orderName" binding:"required"`
	OrderTime string `json:"orderTime" validate:"checkTime"`
}

func checkTime(fl validator.FieldLevel) bool {
	if fl.Field().String() == "1" {
		return true
	}
	return false
}
func TestVerifyParamStruct(t *testing.T) {
	r := gin.Default()
	validate := validator.New()
	r.POST("/order", func(c *gin.Context) {
		err := validate.RegisterValidation("checkTime", checkTime)
		if err != nil {
			return
		}
		var order Order
		err = c.ShouldBindJSON(&order)
		if err != nil {
			c.String(500, "参数错误", err)
			return
		}
		err = validate.Struct(order)
		if err != nil {
			c.String(500, "参数错误", err)
			return
		}
		c.JSON(200, order)
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8002")
}
