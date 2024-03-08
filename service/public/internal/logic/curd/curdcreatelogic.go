package curd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-cmd/cmd"
	"github.com/zyw0605688/gozero-curd-vue/service/public/model"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/svc"
	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurdCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurdCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurdCreateLogic {
	return &CurdCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurdCreateLogic) CurdCreate(req *types.CurdCreateRequest) (resp *types.Empty, err error) {
	// 获取结构体
	var modelStruct interface{}
	for k, v := range model.ModelList {
		if k == req.ModelName {
			modelStruct = v
		}
	}
	m := reflect.TypeOf(modelStruct)
	// 大驼峰名字
	name := m.Name()
	// 主键名字
	primaryKeyName := m.Field(0).Name
	primaryKeyJson := m.Field(0).Tag.Get("json")
	underlineName := GetUnderlineWord(name)
	createContent := ""
	deleteContentRequest := ""
	deleteContentResponse := ""
	updateContent := ""
	listContent := ""
	// 前端字段kv
	vueFields := []map[string]string{}

	// create需要的字段
	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `json:"` + field.Tag.Get("json") + `"`
		createContent += (item + " `" + tag + "`" + "\n")
	}
	// delete Request需要的字段
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `form:"` + field.Tag.Get("json") + `"`
		deleteContentRequest += (item + " `" + tag + "`" + "\n")
	}
	// delete Response需要的字段
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `json:"` + field.Tag.Get("json") + `"`
		deleteContentResponse += (item + " `" + tag + "`" + "\n")
	}
	// update及列表返回等需要的字段
	for i := 0; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `json:"` + field.Tag.Get("json") + `"`
		updateContent += (item + " `" + tag + "`" + "\n")
	}
	// list request需要的字段
	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v *%v`, field.Name, field.Type)
		tag := `form:"` + field.Tag.Get("json") + ",optional" + `"`
		listContent += (item + " `" + tag + "`" + "\n")
	}
	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		key := field.Tag.Get("json")
		item := field.Tag.Get("gorm")
		label := getColumnFromGormTag(item)
		fmt.Println(label)
		vueFields = append(vueFields, map[string]string{"Key": key, "Label": label, "Name": field.Name})
	}

	CreateRequest := getCreateRequest(name, createContent)
	CreateResponse := getCreateResponse(name, updateContent)
	DeleteRequest := getDeleteRequest(name, deleteContentRequest)
	DeleteResponse := getDeleteResponse(name, deleteContentResponse)
	UpdateRequest := getUpdateRequest(name, updateContent)
	UpdateResponse := getUpdateResponse(name, updateContent)
	DetailRequest := getDetailRequest(name, deleteContentRequest)
	DetailResponse := getDetailResponse(name, updateContent)
	ListRequest := getListRequest(name, listContent)
	ListResponse := getListResponse(name)
	ServerContent := getServerContent(name, underlineName)

	res := ""
	res += CreateRequest + "\n"
	res += CreateResponse + "\n"
	res += DeleteRequest + "\n"
	res += DeleteResponse + "\n"
	res += UpdateRequest + "\n"
	res += UpdateResponse + "\n"
	res += DetailRequest + "\n"
	res += DetailResponse + "\n"
	res += ListRequest + "\n"
	res += ListResponse + "\n"
	res += ServerContent + "\n"

	// 生成.api文件
	err = createApiFile(underlineName, res)
	// 将.api文件加入总的api文件
	err = addApiFile(underlineName)
	// 运行goctl命令生成代码
	err = goCtlGenFile()
	// 编辑logic文件
	err = editLogicFile(name, underlineName, primaryKeyName, primaryKeyJson, vueFields)
	// 生成前端文件
	err = genWebApiFile(underlineName)
	err = genWebVueFile(name, underlineName, primaryKeyJson, vueFields)
	resp = &types.Empty{}
	return
}

func getCreateRequest(name, content string) string {
	return fmt.Sprintf(`type %vCreateRequest {
		%v
	}`, name, content)
}
func getCreateResponse(name, content string) string {
	return fmt.Sprintf(`type %vCreateResponse {
		%v
	}`, name, content)
}
func getDeleteRequest(name, content string) string {
	return fmt.Sprintf(`type %vDeleteRequest {
		%v
	}`, name, content)
}
func getDeleteResponse(name, content string) string {
	return fmt.Sprintf(`type %vDeleteResponse {
		%v
	}`, name, content)
}

func getUpdateRequest(name, content string) string {
	return fmt.Sprintf(`type %vUpdateRequest {
		%v
	}`, name, content)
}
func getUpdateResponse(name, content string) string {
	return fmt.Sprintf(`type %vUpdateResponse {
		%v
	}`, name, content)
}

func getDetailRequest(name, content string) string {
	return fmt.Sprintf(`type %vDetailRequest {
		%v
	}`, name, content)
}
func getDetailResponse(name, content string) string {
	return fmt.Sprintf(`type %vDetailResponse {
		%v
	}`, name, content)
}

// 列表
func getListRequest(name, listContent string) string {
	tag_page_size := "`" + "form:" + `"page_size"` + "`"
	tag_page_num := "`" + "form:" + `"page_num"` + "`"
	return fmt.Sprintf(`type %vListRequest {
		PageSize int  %v
		PageNum int  %v
		%v
	}`, name, tag_page_size, tag_page_num, listContent)
}
func getListResponse(name string) string {
	tag_list := "`" + "json:" + `"list"` + "`"
	tag_total := "`" + "json:" + `"total"` + "`"
	return fmt.Sprintf(`type %vListResponse {
		List  []*%vUpdateResponse %v
		Total int64   %v                         
	}`, name, name, tag_list, tag_total)
}

func getServerContent(name, underlineName string) string {

	return fmt.Sprintf(`

@server(
	group: %v
	prefix: /api/v1/%v
)

service public-api {
	@handler %vCreate
	post /(%vCreateRequest) returns (%vCreateResponse)
	
	@handler %vDelete
	delete /(%vDeleteRequest) returns (%vDeleteResponse)
	
	@handler %vUpdate
	put /(%vUpdateRequest) returns (%vUpdateResponse)
	
	@handler %vDetail
	get /detail(%vDetailRequest) returns (%vDetailResponse)

	@handler %vList
	get /list(%vListRequest) returns (%vListResponse)
}

`, underlineName, underlineName, name, name, name, name, name, name, name, name, name, name, name, name, name, name, name)
}

func createApiFile(underlineName, res string) error {
	// 生成存放的文件路径
	fileName := underlineName + ".api"
	projectWd, _ := os.Getwd()
	filePath := filepath.Join(projectWd, "../../sources/public", fileName)
	// 如果存在，就删除文件
	_, err := os.Stat(filePath)
	if err == nil {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	// 写入文件
	err = ioutil.WriteFile(filePath, []byte(res), 0644)
	// 格式化文件
	cmdArgs := []string{"api", "format", "-dir", filePath}
	c := cmd.NewCmd("goctl", cmdArgs...)
	<-c.Start()
	return nil
}

func addApiFile(underlineName string) error {
	projectWd, _ := os.Getwd()
	filePath := filepath.Join(projectWd, "../../sources/public/index.api")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644) // 打开文件并设置写入模式
	if err != nil {
		fmt.Println("打开index.api文件失败：", err)
		return err
	}
	defer file.Close() // 关闭文件

	// 向文件中追加内容
	content := fmt.Sprintf(`%vimport "%v.api"`, "\n", underlineName)
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("向index.api文件中追加内容失败：", err)
		return err
	}
	return nil
}

func goCtlGenFile() error {
	// goctl生成文件
	cmdArgs := []string{"api", "go", "-api", "../../sources/public/index.api", "-dir", ".", "--home", "../../../template"}
	c := cmd.NewCmd("goctl", cmdArgs...)
	<-c.Start()
	return nil
}

func editLogicFile(name, underlineName, primaryKeyName, primaryKeyJson string, vueFields []map[string]string) error {

	// 新增逻辑
	createLogic := fmt.Sprintf(`

func (l *%vCreateLogic) %vCreate(req *types.%vCreateRequest) (resp *types.%vCreateResponse, err error) {
	modelParams := &model.%v{}
	err = copier.Copy(modelParams, req)
	if err != nil {
		logx.Error("复制入参失败", err)
		return nil, err
	}
	err = l.svcCtx.PublicDb.Save(&modelParams).Error
	if err != nil {
		logx.Error("写入数据库失败", err)
		return nil, err
	}

	resp = &types.%vCreateResponse{}
	err = copier.Copy(resp, modelParams)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}

	return
}

`, name, name, name, name, name, name)

	// 删除逻辑
	deleteLogic := fmt.Sprintf(`

func (l *%vDeleteLogic) %vDelete(req *types.%vDeleteRequest) (resp *types.%vDeleteResponse, err error) {
	ids := strings.Split(req.%v, ",")
	err = l.svcCtx.PublicDb.Where("%v in ?", ids).Delete(&model.%v{}).Error
	resp = &types.%vDeleteResponse{
		%v: req.%v,
	}
	return
}

`, name, name, name, name, primaryKeyName, primaryKeyJson, name, name, primaryKeyName, primaryKeyName)

	// 修改逻辑
	updateLogic := fmt.Sprintf(`

func (l *%vUpdateLogic) %vUpdate(req *types.%vUpdateRequest) (resp *types.%vUpdateResponse, err error) {
	item := &model.%v{}
	l.svcCtx.PublicDb.Where("%v = ?", req.%v).First(&item)
	err = copier.Copy(item, req)
	if err != nil {
		logx.Error("复制参数失败", err)
		return nil, err
	}
	err = l.svcCtx.PublicDb.Save(&item).Error
	resp = &types.%vUpdateResponse{}
	err = copier.Copy(resp, item)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}

	return
}

`, name, name, name, name, name, primaryKeyJson, primaryKeyName, name)

	// 详情逻辑
	detailLogic := fmt.Sprintf(`

func (l *%vDetailLogic) %vDetail(req *types.%vDetailRequest) (resp *types.%vDetailResponse, err error) {
	resp = &types.%vDetailResponse{}
	item := model.%v{}
	err = l.svcCtx.PublicDb.Where("%v = ?", req.%v).First(&item).Error
	err = copier.Copy(resp, item)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}
	return
}

`, name, name, name, name, name, name, primaryKeyJson, primaryKeyName)

	// 列表逻辑
	listLogic, _ := getListLogic(name, vueFields)

	// 生成存放的文件路径
	fileName := strings.ToLower(name)
	projectWd, _ := os.Getwd()
	createLogicFile := filepath.Join(projectWd, "./internal/logic/"+underlineName+"/", fileName+"createlogic.go")
	deleteLogicFile := filepath.Join(projectWd, "./internal/logic/"+underlineName+"/", fileName+"deletelogic.go")
	updateLogicFile := filepath.Join(projectWd, "./internal/logic/"+underlineName+"/", fileName+"updatelogic.go")
	detailLogicFile := filepath.Join(projectWd, "./internal/logic/"+underlineName+"/", fileName+"detaillogic.go")
	listLogicFile := filepath.Join(projectWd, "./internal/logic/"+underlineName+"/", fileName+"listlogic.go")

	fileList := map[string]string{
		createLogicFile: createLogic,
		deleteLogicFile: deleteLogic,
		updateLogicFile: updateLogic,
		detailLogicFile: detailLogic,
		listLogicFile:   listLogic,
	}
	for k, v := range fileList {
		// 向文件中追加内容
		itemFileRaw, err := os.OpenFile(k, os.O_APPEND|os.O_WRONLY, 0644) // 打开文件并设置写入模式
		if err != nil {
			fmt.Println("打开index.api文件失败：", err)
			return err
		}
		defer itemFileRaw.Close() // 关闭文件
		_, err = itemFileRaw.WriteString(fmt.Sprintf("\n%v\n", v))
		if err != nil {
			fmt.Println("向index.api文件中追加内容失败：", err)
			return err
		}
	}

	return nil
}

func getListLogic(name string, vueFields []map[string]string) (string, error) {
	projectWd, _ := os.Getwd()
	file := filepath.Join(projectWd, "../../sources/template/list.tpl")
	tpl, err := template.ParseFiles(file)
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return "", err
	}
	var tplBuffer bytes.Buffer
	data := map[string]interface{}{
		"Name":      name,
		"VueFields": vueFields,
	}
	err = tpl.Execute(&tplBuffer, data)
	if err != nil {
		return "", err
	}
	return tplBuffer.String(), nil
}

func genWebApiFile(underlineName string) error {
	projectWd, _ := os.Getwd()
	fileDir := filepath.Join(projectWd, "../../")
	file := filepath.Join(projectWd, "../../sources/template/api.tpl")
	tpl, err := template.ParseFiles(file)

	if err != nil {
		fmt.Println("create template failed, err:", err)
		return err
	}
	apiFile, err := os.Create(fileDir + "./sources/temp/" + underlineName + ".ts")
	defer apiFile.Close()
	err = tpl.Execute(apiFile, underlineName)
	if err != nil {
		return err
	}
	return nil
}
func genWebVueFile(name, underlineName, primaryKeyJson string, vueFields interface{}) error {
	projectWd, _ := os.Getwd()
	fileDir := filepath.Join(projectWd, "../../")
	filePath := filepath.Join(projectWd, "../../sources/template/table.tpl")
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return err
	}
	file, err := os.Create(fileDir + "./sources/temp/" + underlineName + ".vue")
	defer file.Close()
	data := map[string]interface{}{
		"Name":           name,
		"UnderlineName":  underlineName,
		"PrimaryKeyJson": primaryKeyJson,
		"VueFields":      vueFields,
	}
	err = tpl.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}
func getColumnFromGormTag(tag string) string {
	parts := strings.Split(tag, ";") // 根据分号分割tag
	for _, part := range parts {
		fmt.Println(1, part)
		if strings.Contains(part, "comment:") {
			column := strings.TrimPrefix(part, "comment:")
			fmt.Println(2, column)
			return column
		}
	}
	return ""
}

// 从大驼峰，转下划线
func GetUnderlineWord(str string) (resp string) {
	var matchNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	str = matchNonAlphaNumeric.ReplaceAllString(str, "_")     //非常规字符转化为 _
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}") //拆分出连续大写
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")  //拆分单词
	return strings.ToLower(snake)
}
