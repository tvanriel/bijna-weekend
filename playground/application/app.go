package application

import (
	"io"
	"os"
	"strconv"

	"github.com/labstack/echo"
	bijnaweekend "github.com/tvanriel/bijna-weekend"
	"go.uber.org/zap"
)

func Run() { 
        e := echo.New()

        log, _ := zap.NewDevelopment()

        e.GET("image", func(c echo.Context) error {
                gc, err := bijnaweekend.PastelGradient()
                
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }
                tc, err := bijnaweekend.PastelTextColor()
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }

                mascotteX, err := strconv.ParseFloat(c.QueryParam("mascotte_x"), 64)
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }
                mascotteY, err := strconv.ParseFloat(c.QueryParam("mascotte_y"), 64)
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }
                m, err := bijnaweekend.NewMascotte(mascotteX, mascotteY, "assets/" + c.QueryParam("image"))
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }
                taglineX, err := strconv.ParseFloat(c.QueryParam("tagline_x"), 64)
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }
                taglineY, err := strconv.ParseFloat(c.QueryParam("tagline_y"), 64)
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }

                t, err := bijnaweekend.NewTagline(taglineX, taglineY, "assets/" + c.QueryParam("font"), 30, c.QueryParam("text"), tc)
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }

                err = bijnaweekend.BijnaWeekend(
                        &bijnaweekend.Config{
                                Height: m.Image.Bounds().Dy(),
                                Width: m.Image.Bounds().Dx(),
                                Mascotte: m,
                                Tagline: t,
                                Palette: gc,
                                Writer: c.Response().Writer,
                        },
                )
                if err != nil {
                        log.Error("err", zap.Error(err))
                        return err
                }
                return nil
        })

        
        e.GET("/", func(ctx echo.Context) error {
                f, err := os.Open(
                        "index.html",
                )
                if err != nil {
                        return err
                }
                ctx.Response().Header().Set("Content-Type", "text/html")
                io.Copy(ctx.Response().Writer, f)
                return nil
        })
        

        e.Start(":8090")
}
