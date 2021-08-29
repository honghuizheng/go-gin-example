package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/honghuizheng/go-gin-example/models"
	"github.com/honghuizheng/go-gin-example/pkg/e"
	"github.com/honghuizheng/go-gin-example/pkg/setting"
	"github.com/honghuizheng/go-gin-example/pkg/util"
	"github.com/unknwon/com"
	"log"
	"net/http"
)

//获取某个文章详情
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface {}
	if ! valid.HasErrors() {
		if models.ExistArticleById(id) {
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}
// 获取所有文章列表
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)

	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}

//新增文章
func AddArticle(c *gin.Context) {
	var (
		article = models.Article{}
		code int
	)
	err := c.Bind(&article)
	if err != nil {
		code = e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check,err := valid.Valid(&article)
	if err != nil {
		code = e.ERROR
	}
	if ! check {
		code = e.INVALID_PARAMS
	}


	articles := map[string]interface{}{
		"title":             article.Title,
		"desc":              article.Desc,
		"content":           article.Content,
		"created_by":        article.CreatedBy,
		"modified_by": 		 article.ModifiedBy,
		"tag_id":    		 article.TagID,
		"state": 			 article.State,
	}

	if models.ExistTagByID(article.TagID) {
		code = e.SUCCESS
		models.AddArticle(articles)
	} else {
		code = e.ERROR_EXIST_TAG
	}

	//name := c.Query("name")
	//state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	//createdBy := c.Query("created_by")
	//
	//valid := validation.Validation{}
	//valid.Required(name, "name").Message("名称不能为空")
	//valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	//valid.Required(createdBy, "created_by").Message("创建人不能为空")
	//valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	//valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	//code := e.INVALID_PARAMS
	//if ! valid.HasErrors() {
	//	if ! models.ExistTagByName(name) {
	//		code = e.SUCCESS
	//		models.AddTag(name, state, createdBy)
	//	} else {
	//		code = e.ERROR_EXIST_TAG
	//	}
	//}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}

//修改文章
func EditArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	var (
		article = models.Article{}
		code int
	)
	err := c.Bind(&article)
	if err != nil {
		code = e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check,err := valid.Valid(&article)
	if err != nil {
		code = e.ERROR
	}
	if ! check {
		code = e.INVALID_PARAMS
	}

	articles := map[string]interface{}{
		"title":             article.Title,
		"desc":              article.Desc,
		"content":           article.Content,
		"modified_by": 		 article.ModifiedBy,
		"tag_id":    		 article.TagID,
		"state": 			 article.State,
	}

	if models.ExistArticleById(id) {
		code = e.SUCCESS
		models.EditArticle(id,articles)
	} else {
		code = e.ERROR_NOT_EXIST_ARTICLE
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}

//删除文章标签
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistArticleById(id) {
			models.DeleteArticleById(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}
