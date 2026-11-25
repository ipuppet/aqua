import { defineStore } from "pinia"
import type { MenuItem } from "#/shared/types"

export const useMenuStore = defineStore("menu", {
  state: () => ({
    active: true,
    loaded: false,
    menus: [] as MenuItem[]
  }),
  actions: {
    setMenuState(active: boolean) {
      this.active = active
    },
    setMenus(menus: MenuItem[]) {
      this.menus = menus
      this.loaded = true
    }
  }
})
