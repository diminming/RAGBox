<template>
    <el-container>
        <el-header>
            <div>我的知识库</div>
        </el-header>
        <el-container>
            <el-aside width="200px">
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
            </el-aside>
            <el-main>
                <router-view></router-view>
            </el-main>
        </el-container>
    </el-container>
</template>
<script lang="js" setup>
import { ref } from 'vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
const { ChatDotRound, CirclePlus, ChatLineRound, Files, Collection } = ElementPlusIconsVue
const activeMenu = ref('file_list')
</script>