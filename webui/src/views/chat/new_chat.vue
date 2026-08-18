<template>
    <div style="height: 70vh;" class="markdown-body markdown-body-light" v-html="renderedHtml">

    </div>
    <div style="position: fixed; bottom: 1rem; left: 50%; transform: translateX(-50%); width: 100%; max-width: 800px;">
        <el-input v-model="content" style="position: relative; width: 80%; left: 50%; transform: translateX(-50%); " :rows="7" type="textarea" placeholder="请输入您的问题" />
        <div style="margin-top: 1rem; text-align: right; width: 60%">
            <div style="display: inline-block;">
                <el-checkbox-button v-model="enableKbase" size="large">
                    <template #default>
                        <el-icon>
                            <Search />
                        </el-icon>
                    </template>
                </el-checkbox-button>
            </div>

            <div style="display: inline-block; margin-left: 2rem">
                <el-button @click="sendMessage" type="primary" :icon="Upload"></el-button>
            </div>
        </div>
    </div>


</template>

<script lang="js" setup>
import { Search, Upload } from '@element-plus/icons-vue'
import { ref, computed } from 'vue'
import request from '@/utils/request.js'
import { marked } from 'marked'
// import 'github-markdown-css/github-markdown.css'
import 'github-markdown-css/github-markdown-light.css'
import { ElLoading } from 'element-plus'


const content = ref('')
const enableKbase = ref(false)
const response = ref('')

const renderedHtml = computed(() => {
    return marked.parse(response.value)
})

function sendMessage() {

    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)',
    })

    request.post('/chat', {
        content: content.value,
        enableKbase: enableKbase.value
    }).then(res => {
        response.value = res.data
        loading.close()
    })
}



</script>

<style scoped>
/* 限制展示区域的布局（可选） */
.markdown-body {
    box-sizing: border-box;
    min-width: 200px;
    max-width: 980px;
    margin: 0 auto;
    padding: 45px;
    overflow: auto;
}

/* 适配移动端 */
@media (max-width: 767px) {
    .markdown-body {
        padding: 15px;
    }
}
</style>