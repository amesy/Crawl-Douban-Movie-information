package controllers

import (
	"github.com/astaxie/beego"
	"crawl_movie/models"
	"github.com/astaxie/beego/httplib"
	"time"
)

type CrawlMovieController struct {
	beego.Controller
}

func (c *CrawlMovieController) CrawlMovie() {
	var movieInfo models.MovieInfo
	add := "10.68.7.20:6379"

	// 连接到Redis.
	models.ConnectRedis(add)

	// 爬虫入口URL.
	sURL := "https://movie.douban.com/subject/25827935/"
	models.PutinQueue(sURL)

	for {
		length := models.GetQueueLength()
		if length == 0 {
			break // 如果url队列为空，则退出当前循环。
		}

		surls := models.PopfromQueue()

		// 判断surls是否已被访问过.
		if models.IsVisit(surls) {
			continue
		}


		req := httplib.Get(surls)
		sMovieHtml, err := req.String()
		if err != nil {
			panic(err)
		}

		movieInfo.Movie_name = models.GetMovieName(sMovieHtml)
		// 记录电影信息.
		if movieInfo.Movie_name != "" {
			movieInfo.Movie_id 				= models.GetMovieID(sMovieHtml)
			movieInfo.Movie_pic 			= models.GetMoviePicture(sMovieHtml)
			movieInfo.Movie_director 		= models.GetMovieDirector(sMovieHtml)
			movieInfo.Movie_writer 			= models.GetMovieWriter(sMovieHtml)
			movieInfo.Movie_main_character 	= models.GetMovieMainCharacters(sMovieHtml)
			movieInfo.Movie_on_time 		= models.GetMovieOnTime(sMovieHtml)
			movieInfo.Movie_type			= models.GetMovieType(sMovieHtml)
			movieInfo.Movie_grade			= models.GetMovieGrade(sMovieHtml)
			movieInfo.Movie_span			= models.GetMovieSpan(sMovieHtml)

			models.AddMovie(&movieInfo)
		} else {
			continue
		}

		// 提取该页面所有链接.
		urls := models.GetMovieUrls(sMovieHtml)
		for _, url := range urls {
			models.PutinQueue(url)
			c.Ctx.WriteString("<br>" + url + "</br>")
		}

		// sUrl应当被记录到访问set中.
		models.AddToSet(surls)

		// 防止被目标网站察觉.
		time.Sleep(1 * time.Second)
	}

	c.Ctx.WriteString("<h1>End Of Crwal!</h1>")

}