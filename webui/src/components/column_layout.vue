<template>
    <el-container>
        <el-aside width="200px">
            <el-menu class="el-menu-vertical-demo" @open="handleOpen" @close="handleClose" :router="true"
                :default-active="activeMenu" :default-openeds="['online_doc']">
                <!-- <el-sub-menu index="kbase">
                    <template #title>
                        <el-icon>
                            <Files />
                        </el-icon>
                        <span>知识库</span>
                    </template>
</el-sub-menu> -->
                <el-sub-menu index="online_doc">
                    <template #title>
                        <el-icon>
                            <Collection />
                        </el-icon>
                        <span>在线文档</span>
                    </template>
                    <el-menu-item index="new" :route="{ path: '/doc/new' }">
                        <el-icon>
                            <Plus />
                        </el-icon>
                        <span>创建文档</span>
                    </el-menu-item>
                    <el-menu-item index="list" :route="{ path: '/doc/list' }">
                        <el-icon><Files /></el-icon>
                        <span>文档列表</span>
                    </el-menu-item>
                </el-sub-menu>



                <el-menu-item index="file_list" :route="{ path: '/file/list' }">
                    <el-icon>
                        <DocumentCopy />
                    </el-icon>
                    <span>文件列表</span>
                </el-menu-item>
            </el-menu>
        </el-aside>
        <el-main>
            <router-view></router-view>
        </el-main>
    </el-container>
    <el-dialog></el-dialog>
</template>

<script lang="js" setup>
import { ref } from 'vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
const { Files, Collection, Plus, DocumentCopy } = ElementPlusIconsVue
const activeMenu = ref('')
</script>

<style lang="css" scoped></style>