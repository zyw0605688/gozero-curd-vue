package curd

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/logic/curd"
	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/svc"
	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/types"
)

func CurdCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CurdCreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := curd.NewCurdCreateLogic(r.Context(), svcCtx)
		resp, err := l.CurdCreate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
