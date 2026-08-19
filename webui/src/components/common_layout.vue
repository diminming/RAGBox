<template>
    <el-container>
        <el-header>
            <div class="header-item logo">
                <div class="logo-text" @click="goto('/dashboard')">RAGBox</div>
            </div>
            <div class="header-item">
                <span class="txt" @click="goto('/chat')">对话</span>
            </div>
            <div class="header-item">
                <span class="txt" @click="goto('/doc')">知识库</span>
            </div>
            <div class="header-actions">
                <div class="header-item">
                    <el-icon :size="'1.3rem'">
                        <Setting />
                    </el-icon>
                </div>
                <div class="header-item">
                    <el-avatar> {{ username }} </el-avatar>
                </div>
            </div>
        </el-header>
        <el-container>
            <!-- <el-aside width="200px">
                <el-menu class="el-menu-vertical-demo" @open="handleOpen" @close="handleClose" :router="true"
                    :default-active="activeMenu" :default-openeds="['chat', 'kbase']">
                    <el-sub-menu index="chat">
                        <template #title>
                            <el-icon>
                                <ChatDotRound />
                            </el-icon>
                            <span>对话</span>
                        </template>
                        <el-menu-item index="new_chat" :route="{ path: '/file/new_chat' }">
                            <el-icon>
                                <CirclePlus />
                            </el-icon>
                            <span>新建对话</span>
                        </el-menu-item>
                        <el-menu-item index="history_chat" :route="{ path: '/chat/history' }">
                            <el-icon>
                                <ChatLineRound />
                            </el-icon>
                            <span>历史对话</span>
                        </el-menu-item>
                    </el-sub-menu>

                    <el-sub-menu index="kbase">
                        <template #title>
                            <el-icon>
                                <Files />
                            </el-icon>
                            <span>知识库</span>
                        </template>

                        <el-menu-item index="online_doc" :route="{ path: '/doc/online' }">
                            <el-icon>
                                <Collection />
                            </el-icon>
                            <span>在线文档</span>
                        </el-menu-item>

                        <el-menu-item index="file_list" :route="{ path: '/file/list' }">
                            <el-icon>
                                <Collection />
                            </el-icon>
                            <span>文件列表</span>
                        </el-menu-item>

                    </el-sub-menu>

                </el-menu>
            </el-aside> -->
            <el-main>
                <router-view></router-view>
            </el-main>
        </el-container>
    </el-container>
    <el-dialog></el-dialog>
</template>
<script lang="js" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
const { ChatDotRound, CirclePlus, ChatLineRound, Files, Collection, Setting } = ElementPlusIconsVue
const activeMenu = ref('file_list')
const username = ref('admin')

const router = useRouter()
const goto = (path) => {
    router.push(path)
}

</script>
<style lang="css" scoped>
div.header-item {
    /* display: inline-block; */
    width: 7rem;
    text-align: center;
    font-weight: bold;
    font-size: 1rem;    
}

div.header-item>span:hover {
    /* background-color: #2adcff; */
    cursor: pointer;
}

div.logo {
    font-size: 1.5rem;
    font-weight: bold;
    width: 7rem;
    /* color: #409EFF; */
}

.menu-item {
    font-weight: bold;
    width: 7rem;
}

.el-header {
    display: flex;
    align-items: center;
    box-sizing: border-box;
}

.header-actions {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-left: auto;
}

.header-actions .header-item {
    width: auto;
    margin-left: 2rem;
}
</style>