package chromedriver

import (
	"cat_ben/src/db"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"time"
)

func GetTokenAndSave() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover:", r)
		}
	}()
	service, err := selenium.NewChromeDriverService("c:/chromedriver/chromedriver.exe", 2345)
	if err != nil {
		log.Fatalf("启动 ChromeDriver 失败: %v", err)
	}
	defer service.Stop()

	// 2. 配置浏览器选项
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Path: "", // 留空表示使用默认 Chrome 路径
		Args: []string{
			//"--headless", // 无头模式
			"--no-sandbox",
			"--disable-gpu",
			//"--window-size=1280,800",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, "http://localhost:2345/wd/hub")
	if err != nil {
		log.Fatalf("创建 WebDriver 失败: %v", err)
	}
	defer wd.Quit()

	// 4. 导航到网页
	if err := wd.Get("https://www.xueqiu.com"); err != nil {
		log.Fatalf("导航失败: %v", err)
	}
	time.Sleep(5 * time.Second)
	cookies, err := wd.GetCookies()
	if err != nil {
		log.Fatalf("获取 cookie 失败: %v", err)
	}
	cookieStr := ""
	for i, cookie := range cookies {
		cookieStr += cookie.Name
		cookieStr += "="
		cookieStr += cookie.Value
		cookieStr += "; "
		fmt.Printf("Cookie #%d:\n", i+1)
		fmt.Printf("  Name:   %s\n", cookie.Name)
		fmt.Printf("  Value:  %s\n", cookie.Value)
		fmt.Printf("  Domain: %s\n", cookie.Domain)
		fmt.Printf("  Path:   %s\n", cookie.Path)
		fmt.Printf("  Expiry: %v\n", cookie.Expiry)
		fmt.Printf("  Secure: %t\n", cookie.Secure)
		fmt.Println("-----------------------")
	}
	fmt.Println(cookieStr)
	db.UpdateValue("xueqiu_token", cookieStr)

}
