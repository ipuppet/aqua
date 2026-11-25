import { defineStore } from "pinia"
import type { UserInfo, UserStructure } from "#/shared/types"

export const useUserStore = defineStore("user", {
  state: () => ({
    isLogin: false,
    info: { avatar: "", nickname: "" } as UserInfo,
    structure: {} as UserStructure
  }),
  actions: {
    setUserInfo(info: UserInfo) {
      this.info = info
    },
    setUserStructure(structure: UserStructure) {
      this.structure = structure
    },
    setLoginState(isLogin: boolean) {
      this.isLogin = isLogin
    }
  }
})
