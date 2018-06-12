package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strings"
)

var (
	db orm.Ormer
)

//由于model这个名字叫 MovieInfo 那么操作的表其实 movie_info
type MovieInfo struct {
	Id                   int64
	Movie_id             string
	Movie_name           string
	Movie_pic            string
	Movie_director       string
	Movie_writer         string
	Movie_country        string
	Movie_language       string
	Movie_main_character string
	Movie_type           string
	Movie_on_time        string
	Movie_span           string
	Movie_grade          string
	_Create_time         string
}

func init() {
	orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句
	orm.RegisterDataBase("default", "mysql", "mysqlcli:12345678@tcp(10.68.7.20:3306)/go_data?charset=utf8", 30)
	orm.RegisterModel(new(MovieInfo))
	db = orm.NewOrm()
}

func AddMovie(movie_info *MovieInfo) (int64, error) {
	id, err := db.Insert(movie_info)
	return id, err
}

// 获取电影id.
func GetMovieID(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	reg := regexp.MustCompile(`<meta.*?url=http://m.douban.com/movie/subject/(.*?)/"/>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	// [[<meta name="mobile-agent" content="format=html5; url=http://m.douban.com/movie/subject/25827935/"/> 25827935]]
	return string(result[0][1])
}

// 获取电影名.
func GetMovieName(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	reg := regexp.MustCompile(`<span.*?property="v:itemreviewed">(.*)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return string(result[0][1])
}

// 获取电影图片
func GetMoviePicture(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	reg := regexp.MustCompile(`src="(.*?)" title`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return string(result[0][1])
}

// 获取导演名.
func GetMovieDirector(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	reg := regexp.MustCompile(`<a.*?rel="v:directedBy">(.*)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return string(result[0][1])
}

// 获取电影编剧.
func GetMovieWriter(movieHtml string) string {
	reg := regexp.MustCompile(`/celebrity/13.*?/">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	mainCharacters := ""
	for _, v := range result {
		mainCharacters += v[1] + "/"
	}

	return strings.Trim(mainCharacters, "/")
}

// 获取主演.
func GetMovieMainCharacters(movieHtml string) string {
	//if movieHtml == "" {
	//	return ""
	//}

	reg := regexp.MustCompile(`<a.*?rel="v:starring">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	mainCharacters := ""
	for _, v := range result {
		mainCharacters += v[1] + "/"
	}

	return strings.Trim(mainCharacters, "/")
}

// 上映日期
// <span.*?content="(.*?)">(.*?)</span>
func GetMovieOnTime(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	// 2016-09-14(中国大陆)/2016-10-27(香港)
	/*
	reg := regexp.MustCompile(`<span.*?content="(.*?)">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	mainCharacters := ""
	i := 0
	for _, v := range result{
		if i < 2 {
			mainCharacters += v[1] + "/"
		}
		i ++
	}

	return strings.Trim(mainCharacters, "/")
	*/

	reg := regexp.MustCompile(`<span property="v:initialReleaseDate" content="(.*?)\(.*?</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return result[0][1]
}

// 电影类型
func GetMovieType(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	reg := regexp.MustCompile(`<span.*?property="v:genre">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	mainCharacters := ""
	for _, v := range result {
		mainCharacters += v[1] + "/"

	}

	return strings.Trim(mainCharacters, "/")
}

// 电影评分
func GetMovieGrade(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	reg := regexp.MustCompile(`<strong.*?property="v:average">(.*?)</strong>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return string(result[0][1])
}

// 片长
func GetMovieSpan(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	reg := regexp.MustCompile(`<span property="v:runtime".*?content=.*?>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return string(result[0][1])
}

// 获取电影页面URL
func GetMovieUrls(movieHtml string) []string {
	reg := regexp.MustCompile(`<a.*?href="(https://movie.douban.com/.*?)"`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	var movieSets []string
	for _, v := range result {
		movieSets= append(movieSets, v[1])
	}
	return movieSets
}

