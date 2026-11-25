import { defineRule } from "vee-validate"
import { securityConfig } from "#/shared/config/security"

defineRule("required", (value: any) => {
  if (!value || !value.length) {
    return "此字段为必填项"
  }
  return true
})

defineRule("phone", (value: string) => {
  if (!value || !value.length) {
    return true
  }
  const phoneRegex = /^((13[0-9])|(14[579])|(15[^4])|(18[0-9])|(17[0135678]|18[0-9]|19[8]))\d{8}$/
  if (!phoneRegex.test(value)) {
    return "请输入有效的手机号码"
  }
  return true
})

defineRule("email", (value: string) => {
  if (!value || !value.length) {
    return true
  }
  const emailRegex = /[A-Za-z\d]+([-_.][A-Za-z\d]+)*@([A-Za-z\d]+[-.])+[A-Za-z\d]{2,4}/
  if (!emailRegex.test(value)) {
    return "请输入有效的邮箱地址"
  }
  return true
})

defineRule("password", (value: string) => {
  if (!value || !value.length) {
    return true
  }
  if (value.length < securityConfig.passwordMinLength) {
    return "密码长度至少为6位。"
  }
  if (value.length > securityConfig.passwordMaxLength) {
    return "密码长度不能超过50位。"
  }
  return true
})
