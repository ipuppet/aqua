import { ACCESS_TOKEN_KEY, setCookie } from "@aqua/http"
import { ElMessageBox } from "element-plus"

export function redirectToLogin() {
  let domain = import.meta.env.VITE_API ?? ""
  domain = domain.replace("api", "passport")
  const loginUrl = "https://" + domain + "?from=" + encodeURIComponent(window.location.href)
  // 本地开发时使用浏览器 prompt 录入 token，避免跳转外链
  if (loginUrl.includes("localhost")) {
    ElMessageBox.prompt("请输入 Access Token", "Access Token", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      inputPlaceholder: "请输入 Access Token"
    })
      .then(({ value }: { value: string }) => {
        if (value) {
          setCookie(ACCESS_TOKEN_KEY, value.trim(), 2)
        }
      })
      .catch(() => {})
    return
  }
  window.location.href = loginUrl
}
