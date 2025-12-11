package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
)

// AddConfigDictReq 配置字典请求结构体
// 使用 validate 标签定义验证规则，label 标签用于错误消息中的字段名
type AddConfigDictReq struct {
	ParentId  int64  `json:"parentId" validate:"number,gte=0" label:"字典集id"`
	Name      string `json:"name" validate:"min=2,max=50" label:"名称"`
	Type      int64  `json:"type" validate:"number,gte=1,lte=12" label:"类型"`
	UniqueKey string `json:"uniqueKey" validate:"min=2,max=50" label:"标识"`
	Value     string `json:"value" validate:"max=2048" label:"字典项值"`
	OrderNum  int64  `json:"orderNum" validate:"gte=0,lte=9999" label:"排序"`
	Remark    string `json:"remark" validate:"max=200" label:"备注"`
	Status    int64  `json:"status" validate:"number,gte=0,lte=1" label:"状态"`
}

// 响应结构体
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 创建验证器并配置中文翻译
func newValidator() (*validator.Validate, ut.Translator, error) {
	validate := validator.New()
	
	// 注册标签名函数，使用 label 标签作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		if name == "" {
			return fld.Name
		}
		return name
	})

	// 创建中文翻译器
	zhTranslator := zh.New()
	uni := ut.New(zhTranslator, zhTranslator)
	trans, _ := uni.GetTranslator("zh")
	
	// 注册默认的中文翻译
	err := translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, nil, err
	}

	return validate, trans, nil
}

// 验证请求并返回错误信息
func validateRequest(req interface{}, validate *validator.Validate, trans ut.Translator) error {
	err := validate.Struct(req)
	if err != nil {
		// 将验证错误转换为中文错误消息
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationError := range validationErrors {
				return fmt.Errorf(validationError.Translate(trans))
			}
		}
		return err
	}
	return nil
}

// HTTP 处理函数：添加配置字典
func addConfigDictHandler(w http.ResponseWriter, r *http.Request) {
	// 只接受 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req AddConfigDictReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Code:    400,
			Message: fmt.Sprintf("参数解析失败: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 创建验证器
	validate, trans, err := newValidator()
	if err != nil {
		response := Response{
			Code:    500,
			Message: fmt.Sprintf("验证器初始化失败: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证请求
	if err := validateRequest(&req, validate, trans); err != nil {
		response := Response{
			Code:    400,
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证通过，返回成功响应
	response := Response{
		Code:    200,
		Message: "验证成功",
		Data:    req,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 健康检查接口
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Code:    200,
		Message: "服务运行正常",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// 注册路由
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/config-dict/add", addConfigDictHandler)

	// 启动服务器
	port := ":8080"
	fmt.Printf("服务器启动在端口 %s\n", port)
	fmt.Println("测试接口:")
	fmt.Println("  GET  http://localhost:8080/health")
	fmt.Println("  POST http://localhost:8080/api/config-dict/add")
	fmt.Println("\n示例请求体:")
	fmt.Println(`{
  "parentId": 0,
  "name": "测试字典",
  "type": 1,
  "uniqueKey": "test_key",
  "value": "测试值",
  "orderNum": 1,
  "remark": "这是备注",
  "status": 1
}`)

	log.Fatal(http.ListenAndServe(port, nil))
}

