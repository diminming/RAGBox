import { createRouter, createWebHistory } from 'vue-router'
import { ElMessageBox } from 'element-plus'

const routes = [
	{
		path: '/',
		name: 'index',
		component: () => import('@/components/common_layout.vue'),
		redirect: '/dashboard',
		children: [
			{
				path: 'dashboard',
				name: 'dashboard',
				component: () => import('@/views/dashboard/index.vue'),
				meta: {
					title: '仪表盘',
					requiresAuth: true
				}
			},
			{
				path: 'chat',
				name: 'chat-component',
				redirect: '/chat/new',
				component: () => import('@/components/chat_layout.vue'),
				children: [
					{
						path: "new",
						name: "new-chat",
						component: () => import('@/views/chat/new_chat.vue'),
						meta: {
							title: '新建对话',
							requiresAuth: true
						}
					},
					{
						path: "history",
						name: "history-chat",
						component: () => import('@/views/chat/chat_lst.vue'),
						meta: {
							title: '历史对话',
							requiresAuth: true
						}
					}
				]
			},
			{
				path: 'doc',
				name: 'doc-component',
				redirect: '/doc/online',
				component: () => import('@/components/column_layout.vue'),
				children: [
					{
						path: "online",
						name: "online-doc",
						component: () => import('@/views/doc/online_doc.vue'),
						meta: {
							title: '在线文档',
							requiresAuth: true
						}
					}
				]
			},
			{
				path: 'file',
				name: 'file-component',
				redirect: '/file/list',
				component: () => import('@/components/column_layout.vue'),
				children: [
					{
						path: "list",
						name: "file-list",
						component: () => import('@/views/file/file_lst.vue'),
						meta: {
							title: '文件列表',
							requiresAuth: true
						}
					}
				]
			},
			{
				path: 'user',
				name: 'user-list',
				component: () => import('@/views/user/user_lst.vue'),
				meta: {
					title: '用户列表',
					requiresAuth: true
				}
			},
		]
	},
	{
		path: '/login',
		name: 'login',
		component: () => import('@/views/login/index.vue'),
		meta: {
			requiresAuth: false
		},
	},

	{
		path: '/:pathMatch(.*)*',
		name: 'not-found',
		component: () => import('@/views/err_404.vue'),
		meta: {
			title: '页面未找到',
			requiresAuth: true
		}
	}
]

const router = createRouter({
	history: createWebHistory(),
	routes
})

router.beforeEach((to, from, next) => {
	// 示例：检查是否登录
	const isLogin = localStorage.getItem('token')
	if (to.meta.requiresAuth && !isLogin) {
		ElMessageBox.alert('未登录或登录已过期，请重新登录', '未登录', {
			// if you want to disable its autofocus
			// autofocus: false,
			confirmButtonText: '重新登陆',
			callback: (action) => {
				if (action === 'confirm') {
					localStorage.removeItem('token');
					window.location.href = '/login';
				}
			},
		})
		next('/login')
	} else {
		next()
	}
})

export default router
