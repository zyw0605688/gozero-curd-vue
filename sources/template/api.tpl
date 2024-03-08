import http from '@/service/http'

export const create = async (data) => {
    return http.post("/public-api/{{ . }}", data)
}

export const deleteByIds = async (params) => {
    return http.delete("/public-api/{{ . }}", params)
}

export const updateById = async (data) => {
    return http.put("/public-api/{{ . }}", data)
}

export const getPageList = async (params) => {
    return http.get("/public-api/{{ . }}/list",params)
}

export const getDetailById = async (params) => {
    return http.get("/public-api/{{ . }}/detail", params)
}
