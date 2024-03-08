func (l *{{ .Name }}ListLogic) {{ .Name }}List(req *types.{{ .Name }}ListRequest) (resp *types.{{ .Name }}ListResponse, err error) {
    // 返回参数
    resp = &types.{{ .Name }}ListResponse{}
    list := []model.{{ .Name }}{}
    total := int64(0)

    // 构造查询条件
    querySql := l.svcCtx.PublicDb.Model(&model.{{ .Name }}{})
    countSql := l.svcCtx.PublicDb.Model(&model.{{ .Name }}{})

    // 根据参数，增加过滤条件，有需要自行修改，如
    {{- range $k, $v := .VueFields }}
    if req.{{ $v.Name }} != nil {
        querySql.Where("{{ $v.Key }} = ?", req.{{ $v.Name }})
        countSql.Where("{{ $v.Key }} = ?", req.{{ $v.Name }})
    }
    {{- end }}

    // 分页
    if req.PageNum > 0 && req.PageSize > 0 {
        offset := (req.PageNum - 1) * req.PageSize
        querySql.Offset(offset).Limit(req.PageSize)
    }

    // 查询
    err = querySql.Find(&list).Error
    err = countSql.Count(&total).Error

    // 处理返回列表
    for _, item := range list {
    temp := &types.{{ .Name }}UpdateResponse{}
    err = copier.Copy(temp, item)
    resp.List = append(resp.List, temp)
    }
    resp.Total = total

    return
}