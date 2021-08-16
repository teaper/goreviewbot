package rss

import (
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type Item struct {
	Title string `xml:"title"`
	Dc string `xml:"dc"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
	Guid string `xml:"guid"`
	Link string `xml:"link"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items []Item `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel `xml:"channel"`
}


/**
<item>
      <title>🎉 Go 1.17 Release Candidate 1 is released!

🏃‍♂️ Run it in dev! Run it in prod! File bugs! https://golang.org/issue/new

📢 Announcement: https://groups.google.com/g/golang-announce/c/gJE7OtHlRbM/m/21x8zAR-AAAJ

⬇️ Download: https://golang.org/dl/#go1.17rc1

#golang</title>
      <dc:creator>@golang</dc:creator>
      <description><![CDATA[<p>🎉 Go 1.17 Release Candidate 1 is released!

🏃‍♂️ Run it in dev! Run it in prod! File bugs! <a href="https://golang.org/issue/new">golang.org/issue/new</a>

📢 Announcement: <a href="https://groups.google.com/g/golang-announce/c/gJE7OtHlRbM/m/21x8zAR-AAAJ">groups.google.com/g/golang-a…</a>

⬇️ Download: <a href="https://golang.org/dl/#go1.17rc1">golang.org/dl/#go1.17rc1</a>

<a href="https://nitter.net/search?q=%23golang">#golang</a></p>
<imgs src="https://nitter.net/pic/media%2FE6NAesDXsAAF1uK.png" style="max-width:250px;" />]]></description>
      <pubDate>Tue, 13 Jul 2021 20:29:39 GMT</pubDate>
      <guid>https://nitter.net/golang/status/1415045781233545218#m</guid>
      <link>https://nitter.net/golang/status/1415045781233545218#m</link>
    </item>
 */

func GetRssPage(clientURL string,pubDate *string) string  {
	//curl 下载 rss 文件命令（-k 绕开 TLS 验证）：curl -k  https://nitter.net/golang/rss > rss
	//bt, err := ioutil.ReadFile("rss/rss") //从文件拿数据的方式
	//if err != nil {
	//	log.Fatal(err)
	//}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //绕过 TLS 验证
	}
	client := &http.Client{Transport: tr} //重写 http 的值
	log.Println("cfg.Rss.ClientURL ==>",clientURL)
	res, err := client.Get(clientURL)
	if err != nil {
		log.Println(err)
	}
	bt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	rss := RSS{}
	err = xml.Unmarshal(bt, &rss)
	if err != nil {
		log.Println(err)
	}
	log.Println(rss.Channel.Items[0].Title)

	if rss.Channel.Items[0].PubDate != *pubDate {
		*pubDate = rss.Channel.Items[0].PubDate
		return rss.Channel.Items[0].Title
	}
	log.Println("新文章的更新时间：",*pubDate)
	return ""
}
