<template>
    <div class="btn_bar">
        <el-button type="primary" size="small" @click="showUploadDrawer = true">上传</el-button>
    </div>
    <div class="file-list-container">
        <el-table :data="file_lst" size="small">

            <el-table-column prop="id" label="ID" width="80"></el-table-column>

            <el-table-column prop="origin_filename" label="文件名"></el-table-column>

            <el-table-column prop="mime_type" label="文件类型" width="180">
                <template #default="{ row }">
                    <span>{{ row.mime_type }}</span>
                </template>
            </el-table-column>

            <el-table-column prop="update_timestamp" label="更新时间" width="180">
                <template #default="{ row }">
                    <span>{{ new Date(row.update_timestamp * 1000).toLocaleString() }}</span>
                </template>
            </el-table-column>

            <el-table-column prop="create_timestamp" label="创建时间" width="180">
                <template #default="{ row }">
                    <span>{{ new Date(row.create_timestamp * 1000).toLocaleString() }}</span>
                </template>
            </el-table-column>

            <el-table-column label="操作" width="180">
                <template #default="{ row }">
                    <span>
                        <el-button size="small" link type="primary" @click="lstSegments(row)">查看分片</el-button>
                    </span>
                </template>
            </el-table-column>

        </el-table>
    </div>
    <el-drawer v-model="showUploadDrawer" title="文件上传" :with-header="true">
        <el-upload class="upload-demo" drag action="/api_v1/upload" name="file" directory multiple
            :on-change="handleChange">
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
                拖拽文件至此处或<em>点击</em>完成上传
            </div>
        </el-upload>
    </el-drawer>
    <el-dialog v-model="dialogFormVisible" title="查看分片" width="70%">
        <div class="segment-dialog-body">
            <el-table :data="segment_lst" size="small">
                <el-table-column prop="idx" label="序号" width="80"></el-table-column>
                <el-table-column prop="content" label="分片内容"></el-table-column>
            </el-table>
        </div>
    </el-dialog>
</template>

<script lang="js" setup>
const showUploadDrawer = ref(false)
const dialogFormVisible = ref(false)
const segment_lst = ref([])

import { ref, onMounted, reactive } from 'vue'
import { UploadFilled } from '@element-plus/icons-vue'
import { File } from '@/views/file/file.js'

const file_lst = ref([])
const pagination = reactive({
    page: 1,
    page_size: 10,
    total: 0
})

const get_file_lst = async () => {
    File.getLstByPage(pagination.page, pagination.page_size).then((res) => {
        file_lst.value = res.data
        pagination.total = res.total
    }).catch((err) => {
        console.error(err)
    })
}

const lstSegments = (row) => {
    new File(row.id).lstSegments().then((res) => {
        segment_lst.value = res.data.map((item, idx) => ({ idx: idx, content: item }))
        dialogFormVisible.value = true
    }).catch((err) => {
        console.error(err)
    })
}

onMounted(() => {
    get_file_lst()
})

</script>

<style lang="css">
.file-list-container {
    margin-top: 1rem;
    width: 100%;
}

.segment-dialog-body {
    max-height: 70vh;
    overflow-y: auto;
}
</style>