package main

import (
	"clashconfig/api"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func DownLoadTemplate(url string, path string) {
	log.Printf("Rule template URL: %s", url)
	log.Println("Start downloading the rules template")
	resp, err := http.Get(url)
	if nil != err {
		log.Fatalf("Rule template download failed, please manually download save as [%s]\n", path)
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if nil != err || resp.StatusCode != http.StatusOK {
		log.Fatalf("Rule template download failed, please manually download save as [%s]\n", path)
	}
	ioutil.WriteFile(path, s, 0777)
	log.Printf("Rules template download complete. [%s]\n", path)
}
func main() {
	var listenAddr string
	var listenPort string
	var h bool
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&listenAddr, "l", "0.0.0.0", "Listen address")
	flag.StringVar(&listenPort, "p", "5050", "Listen Port")
	flag.Parse()

	if h {
		flag.Usage()
		return
	}
	_, err := os.Stat("ConnersHua.yaml")
	if err != nil && os.IsNotExist(err) {
		DownLoadTemplate("https://raw.githubusercontent.com/ConnersHua/Profiles/master/Clash/Pro.yaml", "ConnersHua.yaml")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/v2ray2clash", api.V2ray2Clash)
	router.GET("/ssr2clashr", api.SSR2ClashR)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", listenAddr, listenPort),
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
