import { defineStore } from "pinia"
import { ref } from "vue"
import { localData } from "#/shared/lib/storage"

export const useDarkModeStore = defineStore("darkMode", () => {
  const isDarkMode = ref(localData.get("darkMode") === "true" || localData.get("darkMode") === true)

  function setDarkMode(value: boolean) {
    isDarkMode.value = value
    localData.set("darkMode", value)
  }

  return {
    isDarkMode,
    setDarkMode
  }
})
