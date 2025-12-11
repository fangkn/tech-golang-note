

gloang 可以用以下的库，用于对协议中一些字段的检测，可以减少对字段的非法检测。

```go
"github.com/go-playground/validator/v10"
translations "github.com/go-playground/validator/v10/translations/zh"
```

可以参考 [https://godoc.book.murphyyi.com/advanced-go-programming-book/ch5-web/ch5-04-validator.html](https://godoc.book.murphyyi.com/advanced-go-programming-book/ch5-web/ch5-04-validator.html)


协议定义：

``` go
type AddConfigDictReq struct {
    ParentId  int64  `json:"parentId"   validate:"number,gte=0"         label:"字典集id"`
    Name      string `json:"name"       validate:"min=2,max=50"         label:"名称"`
    Type      int64  `json:"type"       validate:"number,gte=1,lte=12"  label:"类型"`
    UniqueKey string `json:"uniqueKey"  validate:"min=2,max=50"         label:"标识"`
    Value     string `json:"value"      validate:"max=2048"             label:"字典项值"`
    OrderNum  int64  `json:"orderNum"   validate:"gte=0,lte=9999"       label:"排序"`
    Remark    string `json:"remark"     validate:"max=200"              label:"备注"`
    Status    int64  `json:"status"     validate:"number,gte=0,lte=1"   label:"状态"`

}
```

可在中间件中检测：
```go
  
func AddConfigDictHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.AddConfigDictReq
        if err := httpx.Parse(r, &req); err != nil {
            httpx.Error(w, errorx.NewHandlerError(errorx.ParamErrorCode, err.Error()))
            return
        }

        validate := validator.New()
        validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
            name := fld.Tag.Get("label")
            return name
        })

        trans, _ := ut.New(zh.New()).GetTranslator("zh")
        validateErr := translations.RegisterDefaultTranslations(validate, trans)
        if validateErr = validate.StructCtx(r.Context(), req); validateErr != nil {
            for _, err := range validateErr.(validator.ValidationErrors) {
                httpx.Error(w, errorx.NewHandlerError(errorx.ParamErrorCode, errors.New(err.Translate(trans)).Error()))
                return
            }
        }

        l := dict.NewAddConfigDictLogic(r.Context(), svcCtx)
        err := l.AddConfigDict(&req)
        if err != nil {
            httpx.Error(w, err)
            return
        }

        response.Response(w, nil, err)
    }
}

```


用例：source/validator 

结果如下： 

```sh

curl 'http://localhost:8080/api/config-dict/add' -d '{
  "parentId": 0,
  "name": "测试字典",
  "type": 1,
  "uniqueKey": "test_key",
  "value": "测试值",
  "orderNum": 1,
  "remark": "这是备注",
  "status": 1
}'
{"code":200,"message":"验证成功","data":{"parentId":0,"name":"测试字典","type":1,"uniqueKey":"test_key","value":"测试值","orderNum":1,"remark":"这是备注","status":1}}

curl 'http://localhost:8080/api/config-dict/add' -d '{
  "parentId": 0,
  "name": "测试字典",
  "type": 100,
  "uniqueKey": "test_key",
  "value": "测试值",
  "orderNum": 1,
  "remark": "这是备注",
  "status": 1
}'
{"code":400,"message":"类型必须小于或等于12"}
```