<template>
    <div class="zmos-gva-container-wrap">
        <div class="zmos-gva-search-box">
            <el-form ref="elSearchFormRef" :inline="true" :model="searchForm" class="zmos-gva-search-form-inline" @keyup.enter="onSearch">
                {{- range $i, $v := .VueFields }}
                    <el-form-item label="{{ $v.Label }}">
                        <el-input v-model="searchForm.{{ $v.Key }}" placeholder="" clearable />
                    </el-form-item>
                {{- end }}
                <el-form-item>
                    <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
                    <el-button icon="refresh" @click="onSearchReset">重置</el-button>
                </el-form-item>
            </el-form>
        </div>
        <div class="zmos-gva-table-box">
            <div class="zmos-gva-btn-list">
                <el-button type="primary" icon="plus" @click="openCreateFormDialog">新增</el-button>
                <el-popover v-model:visible="batchDeleteVisible" :disabled="!multipleSelection.length" placement="top" width="160">
                    <p>确定要删除吗？</p>
                    <div style="text-align: right; margin-top: 8px;">
                        <el-button @click="batchDeleteVisible = false">取消</el-button>
                        <el-button type="primary" @click="onBatchDelete">确定</el-button>
                    </div>
                    <template #reference>
                        <el-button icon="delete" type="danger" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="batchDeleteVisible = true">
                            批量删除
                        </el-button>
                    </template>
                </el-popover>
            </div>
            <el-table
                ref="multipleTable"
                style="width: 100%;height: 680px;"
                border
                tooltip-effect="dark"
                :data="tableData"
                row-key="{{ .PrimaryKeyJson }}"
                @selection-change="handleSelectionChange"
            >
                <el-table-column type="selection" width="55" />
                {{- range $i, $v := .VueFields }}
                <el-table-column prop="{{ $v.Key }}" label="{{ $v.Label }}" align="center"></el-table-column>
                {{- end -}}
                <el-table-column align="left" label="操作" width="220">
                    <template #default="scope">
                        <el-button link icon="View" @click="getDetailAndShowDetailDialog(scope.row)">查看</el-button>
                        <el-button type="primary" link icon="edit" class="table-button" @click="getDetailAndShowUpdateFormDialog(scope.row)">编辑</el-button>
                        <el-button type="danger" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
            <div class="zmos-gva-pagination">
                <el-pagination
                    layout="total, sizes, prev, pager, next, jumper"
                    :current-page="page_num"
                    :page-size="page_size"
                    :page-sizes="[10, 15, 20, 30, 50, 100]"
                    :total="total"
                    @current-change="handleCurrentChange"
                    @size-change="handleSizeChange"
                />
            </div>
        </div>
        <el-dialog v-model="formDialogVisible" :before-close="closeFormDialog" :title="type==='create'?'添加':'修改'" destroy-on-close>
            <el-scrollbar height="500px">
                <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="120px" style="padding-top: 16px">
                    {{- range $i, $v := .VueFields }}
                        <el-form-item label="{{ $v.Label }}" prop="{{ $v.Key }}">
                            <el-input v-model="formData.{{ $v.Key }}" placeholder="" clearable></el-input>
                        </el-form-item>
                    {{- end -}}
                </el-form>
            </el-scrollbar>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="closeFormDialog">取 消</el-button>
                    <el-button type="primary" @click="onSubmit">确 定</el-button>
                </div>
            </template>
        </el-dialog>

        <el-dialog v-model="detailDialogVisible" style="width: 800px" lock-scroll :before-close="closeDetailDialog" title="查看详情" destroy-on-close>
            <el-scrollbar height="550px">
                <el-descriptions :column="1" border>
                    {{- range $i, $v := .VueFields }}
                        <el-descriptions-item label="{{ $v.Label }}">
                            {{"{{"}} formData.{{ $v.Key }} {{"}}"}}
                        </el-descriptions-item>
                    {{- end -}}
                </el-descriptions>
            </el-scrollbar>
        </el-dialog>
    </div>
</template>
<script setup lang="ts" name="{{ .UnderlineName }}">
// 全量引入格式化工具 请按需保留
import {ElMessage, ElMessageBox} from 'element-plus'
import {ref, reactive, onMounted} from 'vue'

// 引入api和工具方法
import {create, deleteByIds, updateById, getDetailById, getPageList} from '@/service/public/{{ .UnderlineName }}'

// 获取预备使用字典，若无用可删除
const otherData = reactive({

})
const getOtherData = async () => {

}

onMounted(async () => {
    await getOtherData()
    await getTableData()
})
// =========== 表格控制部分 ===========
const page_num = ref(1)
const page_size = ref(10)
const total = ref(0)
const tableData = ref([])
const searchForm = ref({})
const elSearchFormRef = ref()

// 搜索重置
const onSearchReset = () => {
    searchForm.value = {}
    getTableData()
}

// 搜索
const onSearch = () => {
    elSearchFormRef.value?.validate(async (valid) => {
        if (!valid) return
        page_num.value = 1
        page_size.value = 10
        getTableData()
    })
}

// 修改页面容量
const handleSizeChange = (val) => {
    page_size.value = val
    getTableData()
}

// 修改页码
const handleCurrentChange = (val) => {
    page_num.value = val
    getTableData()
}

// 查询，根据所有条件
const getTableData = async () => {
    const params = {page_num: page_num.value, page_size: page_size.value, ...searchForm.value}
    const res = await getPageList(params)
    if (res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
    }
}

// ============== 表格控制部分结束 ===============
// ============== 批量删除开始 ===============
// 多选数据
const multipleSelection = ref([])
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 批量删除控制标记
const batchDeleteVisible = ref(false)

// 批量删除
const onBatchDelete = async () => {
    const ids = []
    if (multipleSelection.value.length === 0) {
        ElMessage({
            type: 'warning',
            message: '请选择要删除的数据'
        })
        return
    }
    multipleSelection.value &&
    multipleSelection.value.map(item => {
        ids.push(item.{{ .PrimaryKeyJson }})
    })
    const res = await deleteByIds({ {{- .PrimaryKeyJson }}: ids.join(",")})
    if (res.code === 0) {
        ElMessage({
            type: 'success',
            message: '删除成功'
        })
        if (tableData.value.length === ids.length && page_num.value > 1) {
            page_num.value--
        }
        batchDeleteVisible.value = false
        getTableData()
    }
}
// ============== 批量删除结束 ===============
// ============== 表单开始 ===============
const formData = ref({
    {{- range $i, $v := .VueFields }}
        {{ $v.Key }}: '',
    {{- end }}
})
const elFormRef = ref()
// 验证规则
const rule = reactive({
    {{- range $i, $v := .VueFields }}
        {{ $v.Key }}: [
            {
                required: true,
                message: '{{ $v.Label }}不能为空',
                trigger: ['input', 'blur'],
            },
        ],
    {{- end }}
})
// 弹窗控制标记
const formDialogVisible = ref(false)
// 弹窗行为控制标记（弹窗内部增create还是改update）
const type = ref('')


// 打开新增弹窗
const openCreateFormDialog = () => {
    type.value = 'create'
    formDialogVisible.value = true
}

// 打开更新弹窗
const getDetailAndShowUpdateFormDialog = async (row) => {
    const res = await getDetailById({ {{- .PrimaryKeyJson }}: row.{{ .PrimaryKeyJson -}} })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        formDialogVisible.value = true
    }
}

// 弹窗内确定提交
const onSubmit = async () => {
    elFormRef.value?.validate(async (valid) => {
        if (!valid) return
        let res
        switch (type.value) {
            case 'create':
                res = await create(formData.value)
                break
            case 'update':
                res = await updateById(formData.value)
                break
        }
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: '操作成功'
            })
            closeFormDialog()
            getTableData()
        }
    })
}

// 关闭弹窗
const closeFormDialog = () => {
    formDialogVisible.value = false
    formData.value = {
        {{- range $i, $v := .VueFields }}
        {{ $v.Key }}: '',
        {{- end }}
    }
}

// 查看详情控制标记
const detailDialogVisible = ref(false)

// 打开详情
const getDetailAndShowDetailDialog = async (row) => {
    const res = await getDetailById({ {{- .PrimaryKeyJson }}: row.{{ .PrimaryKeyJson -}} })
    if (res.code === 0) {
        formData.value = res.data
        detailDialogVisible.value = true
    }
}

// 关闭详情弹窗
const closeDetailDialog = () => {
    detailDialogVisible.value = false
    formData.value = {
        {{- range $i, $v := .VueFields }}
        {{ $v.Key }}: '',
        {{- end }}
    }
}

// 删除行
const deleteRow = async (row) => {
    try {
        const temp = await ElMessageBox.confirm('确定要删除吗?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
        })
        if (temp) {
            await deleteByIds({ {{- .PrimaryKeyJson }}: row.{{ .PrimaryKeyJson -}} })
            await getTableData()
        }
    } catch (e) {
    }
}
// ============== 表单结束 ===============
</script>
