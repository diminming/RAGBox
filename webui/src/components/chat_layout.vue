<template>
    <el-container>
        <el-aside width="200px">
            <el-menu class="el-menu-vertical-demo" @open="handleOpen" @close="handleClose" :router="true"
                :default-active="activeMenu" :default-openeds="['chat', 'kbase']">
                <el-menu-item index="new_chat" :route="{ path: '/chat/new' }">
                    <el-icon>
                        <CirclePlus />
                    </el-icon>
                    <span>新建对话</span>
                </el-menu-item>
                <!-- <el-menu-item index="history_chat" :route="{ path: '/chat/history' }">
                    <el-icon>
                        <ChatLineRound />
                    </el-icon>
                    <span>历史对话</span>
                </el-menu-item> -->
                <el-sub-menu index="chat">
                    <template #title>
                        <el-icon>
                            <ChatDotRound />
                        </el-icon>
                        <span>历史对话</span>
                    </template>
                    <template v-for="item in historyChatList" :key="item.id">
                        <el-menu-item :index="`history_chat_${item.id}`" :route="{ path: `/chat/history/${item.id}` }">
                            <span>{{ item.title }}</span>
                        </el-menu-item>
                    </template>
                </el-sub-menu>
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
const { ChatDotRound, CirclePlus } = ElementPlusIconsVue
const activeMenu = ref('new_chat')
const historyChatList = ref([
    { id: 1, title: '对话1' },
    { id: 2, title: '对话2' },
    { id: 3, title: '对话3' }
])

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