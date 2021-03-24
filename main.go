package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"gopkg.in/ini.v1"
)

const (
	// TimeLayout 出力時間のフォーマット
	TimeLayout string = "2006-01-02 15:04:05"
	// DateLayout 日付フォーマット
	DateLayout string = "2006/01/02"
)

// ScrapingList スクレイピング用のstruct
type ScrapingList struct {
	URL           string
	UserName      string
	Password      string
	ChatworkToken string
	ChatworkURL   string
	Chatwork2Me   string
}

// Scraping リスト
var Scraping ScrapingList
var comment string = `Subject：
xxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxx
`

func init() {
	config, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("config.iniファイルの読み込みに失敗しました:%v", err)
	}
	Scraping = ScrapingList{
		URL:           config.Section("web").Key("url").MustString(""),
		UserName:      config.Section("login").Key("username").MustString(""),
		Password:      config.Section("login").Key("password").MustString(""),
		ChatworkToken: config.Section("chatwork").Key("cwToken").MustString(""),
		ChatworkURL:   config.Section("chatwork").Key("cwURL").MustString(""),
		Chatwork2Me:   config.Section("chatwork").Key("cw2Me").MustString(""),
	}
}

// ChatWorkMessagePost チャットを通知する
func ChatWorkMessagePost(message string) {
	values := url.Values{}
	values.Add("body", message)

	req, err := http.NewRequest(
		"POST",
		Scraping.ChatworkURL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		log.Fatalf("リクエスト失敗しました。:%v", err)
	}

	// コンテントタイプを設定
	req.Header.Set("X-ChatWorkToken", Scraping.ChatworkToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// メッセージ投稿
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("メッセージ送信失敗しました。:%v", err)
	}
	defer resp.Body.Close()

	// メッセージ投稿結果表示
	if resp.StatusCode == http.StatusOK {
		fmt.Println("メッセージ投稿に成功しました。")
	} else {
		fmt.Println("メッセージ投稿に失敗しました。")
	}

	fmt.Println("メッセージ投稿ステータスコード=[",
		resp.StatusCode,
		"]レスポンス内容=[",
		resp.Status,
		"]",
	)
}

func main() {
	fmt.Println("処理開始...")

	// ChromeDriver
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--window-size=1280,800",
		}),
		agouti.Debug,
	)

	err := driver.Start()
	if err != nil {
		errLog("Chromeドライバ起動に失敗しました。")
		log.Fatalln(err)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		errLog("Chromeでページを生成できませんでした。")
		log.Fatalln(err)
	}

	err = page.Navigate(Scraping.URL)
	if err != nil {
		errLog("該当のページに遷移できませんでした。")
		log.Fatalln(err)
	}

	// ログイン処理
	page.FindByXPath(`//*[@id="username"]`).Fill(Scraping.UserName)
	page.FindByXPath(`//*[@id="password"]`).Fill(Scraping.Password)
	page.FindByName(`login`).Click()

	// HTML取得
	getAll, err := page.HTML()
	if err != nil {
		errLog("遷移後のHTML(一覧)取得失敗しました。")
		log.Fatalln(err)
	}
	readerGetAll := strings.NewReader(getAll)
	contentDom, _ := goquery.NewDocumentFromReader(readerGetAll)

	seeText := contentDom.Find("#header > h1 > span.current-project").Text()
	if seeText != "xxxxxx" {
		message := "ログインに失敗しました。"
		errLog(message)
		log.Fatalln(message)
	}

	err = page.FindByXPath(`//*[@id="content"]/form[2]/div/table/thead/tr/th[7]/a`).Click()
	if err != nil {
		errLog("要素をクリックできませんでした。")
		log.Fatalln(err)
	}

	// HTML取得
	getAll, err = page.HTML()
	if err != nil {
		errLog("遷移後のHTML(一覧)取得失敗しました。")
		log.Fatalln(err)
	}
	readerGetAll = strings.NewReader(getAll)
	contentDom, _ = goquery.NewDocumentFromReader(readerGetAll)

	// 行数取得
	rawCnt := len(contentDom.Find("#content > form:nth-child(5) > div > table > tbody > tr").Nodes)

	// 今日の日付を選択
	todayStr := time.Now().Format(DateLayout)
	for i := 1; i <= rawCnt; i++ {
		targetDay := contentDom.Find(`tbody > tr:nth-child(` + strconv.Itoa(i) + `) > td:nth-child(10)`).Text()
		if todayStr == targetDay {
			err = page.FindByXPath(`//form[2]/div/table/tbody/tr[` + strconv.Itoa(i) + `]/td[7]/a`).Click()
			if err != nil {
				errLog("要素をクリックできませんでした。")
				log.Fatalln(err)
			}
			break
		}
	}

	// HTML取得
	getDetail, err := page.HTML()
	if err != nil {
		errLog("遷移後のHTML(詳細ページ)取得失敗しました。")
		log.Fatalln(err)
	}
	readerGetDetail := strings.NewReader(getDetail)
	contentDomDetail, _ := goquery.NewDocumentFromReader(readerGetDetail)

	seeText = contentDomDetail.Find(`#content > div.issue.tracker-21.status-1.priority-2.priority-default.details > div.attributes > div:nth-child(2) > div:nth-child(2) > div > div.label > span`).Text()
	if seeText != "SYS追加項目" {
		message := "詳細ページへ遷移できませんでした。"
		errLog(message)
		log.Fatalln(message)
	}
	// コメントを入力
	page.FindByXPath(`//*[@id="content"]/div[1]/a[1]`).Click()
	page.FindByXPath(`//*[@id="issue_notes"]`).Fill(comment)
	page.FindByName(`commit`).Click()

	// 2秒まつ
	time.Sleep(time.Second * 2)

	var dir = "./tmp/" + time.Now().Format(strings.Replace(DateLayout, "/", "-", -1))

	// ディレクトリ作成
	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("ディレクトリ作成します...")
			err = os.Mkdir(dir, 0755)
			if err != nil {
				fmt.Println("ディレクトリ作成失敗しました...")
				panic(err)
			}
		}
	}

	// スクリーンショットをとる
	page.Screenshot(dir + "/input_" + time.Now().Format(TimeLayout) + ".png")

	postMessage := Scraping.Chatwork2Me + time.Now().Format(DateLayout) + "\n記入した内容：\n[code]" + comment + "[/code]"
	ChatWorkMessagePost(postMessage)

	fmt.Println("正常に処理が終了しました。処理を終了します。")
}

func errLog(message string) {
	ChatWorkMessagePost(Scraping.Chatwork2Me + time.Now().Format(DateLayout) + message)
}
