# Golang_Scraping
Golang/Selenium(ChromeDriver)

## condition
you must be installed latest "ChromeDriver".
https://chromedriver.chromium.org/downloads

## attension
> Web data scraping and crawling aren’t illegal by themselves, but it is important to be ethical while doing it.
> Don’t tread onto other people’s sites without being considerate. Respect the rules of their site. Consider reading over their Terms of Service, read the robots.txt file.
> If you suspect a site is preventing you from crawling, consider contacting the webmaster and asking permission to crawl their site. Don’t burn out their bandwidth–try using a slower crawl rate (like 1 request per 10-15 seconds). Don’t publish any content you find that was not intended to be published.

https://www.import.io/post/6-misunderstandings-about-web-scraping/

## TODO
1. <code>go mod init</code>
2. <code>mkdir tmp/</code>
3. <code>vi config.ini</code>
```:config.ini
[web]
# for Basic('user:password@')
url = https://user:password@example.com

[login]
username = your nice username
password = your nice password

[chatwork]
cwToken = your nice token
cwURL   = https://xxxxxxxxxxx
cw2Me   = xxxxxxxxxx
```
4. <code>go run main.go</code>

## other
### crontab
```zsh
% cd
% wget https://www8.cao.go.jp/chosei/shukujitsu/syukujitsu.csv
% crontab -e
12 18 * * 1-5 bash -c "sleep $((RANDOM \% 1800))s"; grep `date "+\%Y/\%-m/\%-d"`, syukujitsu.csv > /dev/null || cd ~/Desktop/go_kintai; bash -l -c 'go run ~/go_chatwork/main.go'
```