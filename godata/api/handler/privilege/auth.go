package privilege

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/id"
	"github.com/phoenix-agent-go/infra/jwt"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/usecase"
	"github.com/redis/go-redis/v9"
)

const (
	captchaCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	captchaLength  = 4
	captchaTTL     = 60 * time.Second
	captchaPrefix  = "captcha:"
)

type AuthHandler struct {
	svc        *service.PrivilegeService
	jwtManager *jwt.JWTManager
	rdb        *redis.Client
}

func NewAuthHandler(svc *service.PrivilegeService, jwtManager *jwt.JWTManager, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{svc: svc, jwtManager: jwtManager, rdb: rdb}
}

func (h *AuthHandler) Captcha(c *gin.Context) {
	code := generateCaptchaCode()
	key := captchaPrefix + strconv.FormatUint(id.MustGenerateID(), 10)

	if h.rdb != nil {
		_ = h.rdb.Set(c.Request.Context(), key, code, captchaTTL).Err()
	}

	img := generateCaptchaImage(code)
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	b64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	response.Success(c, model.CaptchaVO{CaptchaKey: key, Image: b64})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var dto model.LoginInfoDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	if h.svc == nil {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: 500}, "service not available")
		return
	}

	ip := c.ClientIP()
	userInfo, err := h.svc.Login(c.Request.Context(), dto, ip)
	if err != nil {
		if appErr, ok := err.(*usecase.AppError); ok {
			response.ErrorWithMsg(c, errcode.ErrCode{Code: appErr.Code}, appErr.Msg)
			return
		}
		response.Error(c, errcode.InternalError)
		return
	}

	userID, _ := strconv.ParseUint(userInfo.UserID, 10, 64)
	token, _ := h.jwtManager.GenerateToken(userID, userInfo.Username)
	userInfo.Token = token
	response.Success(c, userInfo)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, "退出成功")
}

func (h *AuthHandler) Menus(c *gin.Context) {
	if h.svc == nil {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: 500}, "service not available")
		return
	}
	userID, _ := c.Get("user_id")
	userIDStr := strconv.FormatUint(userID.(uint64), 10)

	menus, err := h.svc.GetUserMenus(c.Request.Context(), userIDStr)
	if err != nil {
		response.Error(c, errcode.InternalError)
		return
	}
	response.Success(c, menus)
}

func (h *AuthHandler) GetLoginUserInfo(c *gin.Context) {
	if h.svc == nil {
		response.ErrorWithMsg(c, errcode.ErrCode{Code: 500}, "service not available")
		return
	}
	userID, _ := c.Get("user_id")
	userIDStr := strconv.FormatUint(userID.(uint64), 10)

	user, err := h.svc.GetUserByID(c.Request.Context(), userIDStr)
	if err != nil {
		response.Error(c, errcode.Unauthorized)
		return
	}
	response.Success(c, user)
}

func generateCaptchaCode() string {
	b := make([]byte, captchaLength)
	for i := range b {
		b[i] = captchaCharset[rand.Intn(len(captchaCharset))]
	}
	return string(b)
}

func generateCaptchaImage(code string) image.Image {
	const width, height = 200, 80
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	bg := color.RGBA{R: 240, G: 240, B: 240, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bg)
		}
	}

	for i := 0; i < 6; i++ {
		lineColor := color.RGBA{
			R: uint8(rand.Intn(200)),
			G: uint8(rand.Intn(200)),
			B: uint8(rand.Intn(200)),
			A: 255,
		}
		x1, y1 := rand.Intn(width), rand.Intn(height)
		x2, y2 := rand.Intn(width), rand.Intn(height)
		drawLine(img, x1, y1, x2, y2, lineColor)
	}

	for i := 0; i < 100; i++ {
		dotColor := color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255,
		}
		img.Set(rand.Intn(width), rand.Intn(height), dotColor)
	}

	colors := []color.RGBA{
		{R: 220, G: 50, B: 50, A: 255},
		{R: 50, G: 120, B: 220, A: 255},
		{R: 50, G: 180, B: 50, A: 255},
		{R: 200, G: 150, B: 30, A: 255},
	}

	charWidth := (width - 40) / len(code)
	for i, ch := range code {
		charColor := colors[i%len(colors)]
		x := 20 + i*charWidth
		y := height/2 - 10 + rand.Intn(10) - 5
		drawChar(img, x, y, byte(ch), charColor)
	}

	return img
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}
	if steps == 0 {
		img.Set(x1, y1, c)
		return
	}
	xInc := float64(dx) / float64(steps)
	yInc := float64(dy) / float64(steps)
	x, y := float64(x1), float64(y1)
	for i := 0; i <= steps; i++ {
		img.Set(int(x), int(y), c)
		x += xInc
		y += yInc
	}
}

func drawChar(img *image.RGBA, startX, startY int, ch byte, c color.Color) {
	patterns := map[byte][][2]int{
		'A': {{0, 4}, {1, 3}, {1, 2}, {2, 1}, {2, 0}, {3, 0}, {3, 1}, {3, 2}, {3, 3}, {3, 4}, {4, 2}},
		'B': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 0}, {1, 2}, {1, 4}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
		'C': {{2, 0}, {1, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 4}, {2, 4}},
		'D': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 0}, {1, 4}, {2, 1}, {2, 2}, {2, 3}},
		'E': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 0}, {1, 2}, {1, 4}, {2, 0}, {2, 4}},
		'F': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 0}, {1, 2}, {2, 0}},
		'G': {{2, 0}, {1, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 4}, {2, 4}, {2, 3}, {2, 2}, {1, 2}},
		'H': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 2}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
		'J': {{1, 0}, {1, 1}, {1, 2}, {1, 3}, {0, 4}, {2, 4}},
		'K': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 2}, {2, 1}, {2, 3}},
		'L': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}},
		'M': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 1}, {2, 2}, {3, 1}, {4, 0}, {4, 1}, {4, 2}, {4, 3}, {4, 4}},
		'N': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 1}, {2, 2}, {3, 3}, {4, 0}, {4, 1}, {4, 2}, {4, 3}, {4, 4}},
		'P': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 0}, {1, 2}, {2, 0}, {2, 1}},
		'R': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 3}, {2, 4}},
		'S': {{2, 0}, {1, 0}, {0, 1}, {1, 2}, {2, 3}, {1, 4}, {0, 4}},
		'T': {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
		'U': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 4}, {2, 4}, {3, 0}, {3, 1}, {3, 2}, {3, 3}},
		'V': {{0, 0}, {0, 1}, {0, 2}, {1, 3}, {2, 4}, {3, 3}, {4, 0}, {4, 1}, {4, 2}},
		'W': {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 3}, {2, 2}, {3, 3}, {4, 0}, {4, 1}, {4, 2}, {4, 3}, {4, 4}},
		'X': {{0, 0}, {0, 4}, {1, 1}, {1, 3}, {2, 2}, {3, 1}, {3, 3}, {4, 0}, {4, 4}},
		'Y': {{0, 0}, {0, 1}, {1, 2}, {2, 2}, {2, 3}, {2, 4}, {3, 0}, {3, 1}, {4, 2}},
		'Z': {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {3, 1}, {2, 2}, {1, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
		'2': {{0, 0}, {1, 0}, {2, 0}, {2, 1}, {1, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}},
		'3': {{0, 0}, {1, 0}, {2, 0}, {2, 1}, {1, 2}, {2, 3}, {0, 4}, {1, 4}, {2, 4}},
		'4': {{0, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
		'5': {{0, 0}, {1, 0}, {2, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 3}, {0, 4}, {1, 4}},
		'6': {{1, 0}, {2, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 2}, {1, 4}, {2, 3}, {2, 4}},
		'7': {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {3, 1}, {2, 2}, {1, 3}, {1, 4}},
		'8': {{1, 0}, {2, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 2}, {1, 4}, {2, 1}, {2, 3}, {2, 4}},
		'9': {{1, 0}, {2, 0}, {0, 1}, {1, 2}, {2, 1}, {2, 2}, {2, 3}, {1, 4}, {2, 4}},
	}

	scale := 4
	if pattern, ok := patterns[ch]; ok {
		for _, p := range pattern {
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.Set(startX+p[0]*scale+dx, startY+p[1]*scale+dy, c)
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
