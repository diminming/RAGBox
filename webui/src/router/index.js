import { createRouter, createWebHistory } from 'vue-router'

const routes = [
	{
		path: '/',
        name: 'index',
		component: () => import('@/views/dashboard/index.vue'),
		meta: {
			title: '门户首页'
		}
	},
	{
		path: '/login',
		name: 'login',
		component: () => import('@/views/login/index.vue'),
	},
	{
		path: '/chat',
		name: 'chat-component',
		redirect: '/chat/new_chat',
		component: () => import('@/components/common_layout.vue'),
		children: [
			{
				path: "new_chat",
				name: "new-chat",
				component: () => import('@/views/chat/new_chat.vue'),
				meta: {
					title: '新建对话'
				}
			},
			{
				path: "history",
				name: "history-chat",
				component: () => import('@/views/chat/chat_lst.vue'),
				meta: {
					title: '历史对话'
				}
			}
		]
	},
	{
		path: '/doc',
		name: 'doc-component',
		redirect: '/doc/online',
		component: () => import('@/components/common_layout.vue'),
		children: [
			{
				path: "/doc/online",
				name: "online-doc",
				component: () => import('@/views/doc/online_doc.vue'),
				meta: {
					title: '在线文档'
				}
			}
		]
	},
	{
		path: '/file',
		name: 'file-component',
		redirect: '/file/list',
		component: () => import('@/components/common_layout.vue'),
		children: [
			{
				path: "/file/new_chat",
				name: "new-chat",
				component: () => import('@/views/chat/new_chat.vue'),
				meta: {
					title: '新建对话'
				}
			},
			{
				path: "/file/chat_list",
				name: "chat-list",
				component: () => import('@/views/chat/chat_lst.vue'),
				meta: {
					title: '对话列表'
				}
			},
			{
				path: "/file/list",
				name: "file-list",
				component: () => import('@/views/file/file_lst.vue'),
				meta: {
					title: '文件列表'
				}
			}
		]
	},
	{
		path: '/user',
		name: 'user-list',
		component: () => import('@/views/user/user_lst.vue'),
		meta: {
			title: '用户列表'
		}
	},
	{
		path: '/:pathMatch(.*)*',
        name: 'not-found',
		component: () => import('@/views/err_404.vue'),
		meta: {
			title: '页面未找到'
		}
	}
]

const router = createRouter({
	history: createWebHistory(),
	routes
})

export default router
