<template>
	<main class="login-page">
		<svg class="network-background" viewBox="0 0 1440 900" preserveAspectRatio="xMidYMid slice" aria-hidden="true">
			<g class="network-lines">
				<line x1="-40" y1="180" x2="260" y2="350" />
				<line x1="260" y1="350" x2="470" y2="170" />
				<line x1="260" y1="350" x2="390" y2="650" />
				<line x1="390" y1="650" x2="700" y2="760" />
				<line x1="470" y1="170" x2="760" y2="260" />
				<line x1="760" y1="260" x2="1040" y2="100" />
				<line x1="760" y1="260" x2="930" y2="520" />
				<line x1="930" y1="520" x2="1240" y2="410" />
				<line x1="1040" y1="100" x2="1340" y2="220" />
				<line x1="1240" y1="410" x2="1480" y2="590" />
				<line x1="930" y1="520" x2="850" y2="850" />
			</g>
			<g class="network-nodes">
				<circle cx="260" cy="350" r="8" />
				<circle cx="470" cy="170" r="5" />
				<circle cx="390" cy="650" r="6" />
				<circle cx="760" cy="260" r="10" />
				<circle cx="1040" cy="100" r="6" />
				<circle cx="930" cy="520" r="8" />
				<circle cx="1240" cy="410" r="5" />
				<circle cx="700" cy="760" r="4" />
			</g>
		</svg>

		<section class="login-panel" aria-labelledby="login-title">
			<div class="brand-mark" aria-hidden="true"><span></span><span></span><span></span></div>
			<p class="eyebrow">RAGBOX WORKSPACE</p>
			<h1 id="login-title">欢迎回来</h1>
			<p class="subtitle">登录后继续管理你的知识库</p>

			<form class="login-form" @submit.prevent="handleSubmit">
				<label class="field-label" for="username">用户名</label>
				<div class="input-wrap">
					<el-icon class="input-icon"><User /></el-icon>
					<input
						id="username"
						v-model.trim="form.username"
						type="text"
						name="username"
						autocomplete="username"
						placeholder="请输入用户名"
						:disabled="isSubmitting"
					/>
				</div>

				<label class="field-label" for="password">密码</label>
				<div class="input-wrap">
					<el-icon class="input-icon"><Lock /></el-icon>
					<input
						id="password"
						v-model="form.password"
						type="password"
						name="password"
						autocomplete="current-password"
						placeholder="请输入密码"
						:disabled="isSubmitting"
					/>
				</div>

				<p v-if="errorMessage" class="error-message" role="alert">{{ errorMessage }}</p>
				<button class="submit-button" type="submit" :disabled="isSubmitting">
					<span>{{ isSubmitting ? '登录中...' : '登录' }}</span>
					<el-icon v-if="!isSubmitting"><ArrowRight /></el-icon>
				</button>
			</form>

			<p class="footer-note">安全连接 · 专注知识管理</p>
		</section>
	</main>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'

const form = reactive({
	username: '',
	password: ''
})
const isSubmitting = ref(false)
const errorMessage = ref('')

const handleSubmit = async () => {
	errorMessage.value = ''

	if (!form.username || !form.password) {
		errorMessage.value = '请输入用户名和密码'
		return
	}

	isSubmitting.value = true
	await new Promise((resolve) => window.setTimeout(resolve, 450))
	isSubmitting.value = false
	ElMessage.success('登录信息已提交')
}
</script>

<style scoped>
:global(*) {
	box-sizing: border-box;
}

:global(body) {
	margin: 0;
}

.login-page {
	position: relative;
	display: flex;
	min-height: 100vh;
	overflow: hidden;
	align-items: center;
	justify-content: flex-end;
	padding: 48px clamp(28px, 10vw, 160px) 48px 7vw;
	background: #f4f8f8;
	color: #173d42;
}

.login-page::before {
	position: absolute;
	inset: 0;
	background: linear-gradient(112deg, rgba(228, 242, 241, 0.88), rgba(250, 252, 249, 0.5) 58%, rgba(255, 255, 255, 0.8));
	content: '';
}

.network-background {
	position: absolute;
	inset: 0;
	width: 100%;
	height: 100%;
	opacity: 0.75;
}

.network-lines {
	fill: none;
	stroke: #9cc7c5;
	stroke-width: 1.4;
}

.network-nodes {
	fill: #75aaa7;
}

.network-nodes circle:nth-child(4) {
	fill: #e19c70;
}

.login-panel {
	position: relative;
	z-index: 1;
	width: min(100%, 390px);
	margin-right: clamp(0px, 2vw, 34px);
	padding: 42px 44px 34px;
	border: 1px solid rgba(255, 255, 255, 0.9);
	border-radius: 8px;
	background: rgba(255, 255, 255, 0.84);
	box-shadow: 0 24px 70px rgba(39, 78, 81, 0.12);
	backdrop-filter: blur(14px);
}

.brand-mark {
	display: flex;
	gap: 4px;
	align-items: flex-end;
	height: 24px;
	margin-bottom: 26px;
}

.brand-mark span {
	display: block;
	width: 7px;
	border-radius: 4px;
	background: #e19c70;
}

.brand-mark span:nth-child(1) { height: 13px; }
.brand-mark span:nth-child(2) { height: 20px; background: #4d9998; }
.brand-mark span:nth-child(3) { height: 16px; background: #82b9b4; }

.eyebrow {
	margin: 0 0 12px;
	color: #6a9694;
	font-size: 11px;
	font-weight: 700;
	letter-spacing: 1.6px;
}

h1 {
	margin: 0;
	color: #173d42;
	font-size: 32px;
	font-weight: 700;
	letter-spacing: 0;
}

.subtitle {
	margin: 10px 0 32px;
	color: #789091;
	font-size: 14px;
}

.login-form {
	display: grid;
	gap: 10px;
}

.field-label {
	margin-top: 8px;
	color: #436568;
	font-size: 13px;
	font-weight: 600;
}

.input-wrap {
	position: relative;
}

.input-icon {
	position: absolute;
	top: 50%;
	left: 14px;
	z-index: 1;
	color: #86aaa9;
	transform: translateY(-50%);
}

input {
	width: 100%;
	height: 48px;
	padding: 0 14px 0 42px;
	border: 1px solid #dce8e7;
	border-radius: 5px;
	outline: none;
	background: #fbfdfc;
	color: #173d42;
	font: inherit;
	font-size: 14px;
	transition: border-color 160ms ease, box-shadow 160ms ease;
}

input::placeholder { color: #a7b9b8; }
input:focus { border-color: #4d9998; box-shadow: 0 0 0 3px rgba(77, 153, 152, 0.12); }

.error-message {
	margin: 5px 0 0;
	color: #c66f58;
	font-size: 12px;
}

.submit-button {
	display: flex;
	width: 100%;
	height: 48px;
	align-items: center;
	justify-content: center;
	gap: 10px;
	margin-top: 16px;
	border: 0;
	border-radius: 5px;
	background: #286f70;
	color: white;
	cursor: pointer;
	font: inherit;
	font-size: 14px;
	font-weight: 700;
	transition: background 160ms ease, transform 160ms ease;
}

.submit-button:hover:not(:disabled) { background: #1f6061; transform: translateY(-1px); }
.submit-button:disabled { cursor: wait; opacity: 0.65; }

.footer-note {
	margin: 30px 0 0;
	color: #9aaead;
	font-size: 11px;
	text-align: center;
}

@media (max-width: 700px) {
	.login-page {
		justify-content: center;
		padding: 24px;
	}

	.login-panel {
		margin-right: 0;
		padding: 34px 28px 28px;
	}

	.network-background { opacity: 0.48; }
}
</style>