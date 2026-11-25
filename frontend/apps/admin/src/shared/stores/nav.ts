import { defineStore } from "pinia"
import type { NavButton, NavState } from "#/shared/types"

type NavButtonInput = Partial<NavButton> & { click?: NavButton["click"] }

function initButton(button: NavButtonInput): NavButton {
  return {
    text: button.text ?? "button",
    click: button.click ?? (() => {}),
    type: button.type ?? "",
    icon: button.icon ?? "",
    size: button.size ?? "default"
  }
}

export const useNavStore = defineStore("nav", {
  state: (): NavState => ({
    buttons: [],
    showSearch: false,
    searchKeyword: "",
    onSearch: null,
    onClearSearch: null
  }),
  actions: {
    setNavButtons(buttons: NavButtonInput[]) {
      this.buttons = []
      for (const button of buttons) {
        this.buttons.push(initButton(button))
      }
    },
    setNavButton(button: NavButtonInput) {
      this.buttons = []
      this.buttons.push(initButton(button))
    },
    clearNavButtons() {
      this.buttons = []
    },
    setSearch(show: boolean, onSearch?: (keyword: string) => void, onClearSearch?: () => void) {
      this.showSearch = show
      this.onSearch = onSearch ?? null
      this.onClearSearch = onClearSearch ?? null
      if (!show) {
        this.searchKeyword = ""
      }
    },
    clearSearch() {
      this.searchKeyword = ""
      if (this.onClearSearch) {
        this.onClearSearch()
      }
    },
    triggerSearch() {
      if (this.onSearch) {
        this.onSearch(this.searchKeyword)
      }
    }
  }
})
