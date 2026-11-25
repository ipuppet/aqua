import { createRouter, createWebHashHistory } from "vue-router"
import type { RouteRecordRaw } from "vue-router"

const LoginPage = () => import("#/pages/login/LoginPage.vue")
const RegisterPage = () => import("#/pages/register/RegisterPage.vue")
const ForgotPasswordPage = () => import("#/pages/forgot-password/ForgotPasswordPage.vue")
const ResetPasswordPage = () => import("#/pages/reset-password/ResetPasswordPage.vue")

export const RouteName = {
  Login: "Login",
  Register: "Register",
  ForgotPassword: "ForgotPassword",
  ResetPassword: "ResetPassword"
}

const routes: Array<RouteRecordRaw> = [
  { path: "/", name: RouteName.Login, component: LoginPage },
  { path: "/register", name: RouteName.Register, component: RegisterPage },
  { path: "/forgot-password", name: RouteName.ForgotPassword, component: ForgotPasswordPage },
  { path: "/reset-password", name: RouteName.ResetPassword, component: ResetPasswordPage }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
