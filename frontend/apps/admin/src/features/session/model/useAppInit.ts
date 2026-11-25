import { useUserStore } from "#/shared/stores/user"
import { useMenuStore } from "#/shared/stores/menu"
import { sessionData } from "#/shared/lib/storage"
import { client, buildApiUrl } from "#/shared/lib/api"

export function useAppInit() {
  const userStore = useUserStore()
  const menuStore = useMenuStore()

  function applyUserInfo(info: any) {
    if (!info) return
    userStore.setUserInfo(info)
    sessionData.set("user.info", info)
  }

  function getInfo() {
    const setStructureAction = (structure: any) => {
      userStore.setUserStructure(structure)
    }
    const setMenusAction = (menus: any) => {
      menuStore.setMenus(menus)
    }
    const menus = sessionData.get("menu.menus", []) ?? []
    const structure = sessionData.get("user.structure")
    if (!structure || !menus.length) {
      client
        .GET("/api/users/me/menu", { baseUrl: buildApiUrl("user") })
        .then(({ data: res, error }) => {
          if (!error && res) {
            const menuData = res.menus ?? []
            sessionData.set("menu.menus", menuData as any)
            setMenusAction(menuData as any)
          }
        })
        .catch(() => {})

      client
        .GET("/api/users/me/structure", { baseUrl: buildApiUrl("user") })
        .then(({ data: res, error }) => {
          if (!error && res) {
            const structureData = res.structure ?? {}
            sessionData.set("user.structure", structureData as any)
            setStructureAction(structureData as any)
          }
        })
        .catch(() => {})
    } else {
      setMenusAction(menus)
      setStructureAction(structure)
    }
  }

  function init() {
    const cachedUserInfo = sessionData.get("user.info")
    if (cachedUserInfo) {
      userStore.setUserInfo(cachedUserInfo)
      userStore.setLoginState(true)
    }
    client
      .GET("/api/users/me", { baseUrl: buildApiUrl("user") })
      .then(({ data: res, error }) => {
        if (!error && res) {
          // GetSelf directly returns User
          applyUserInfo(res as any)
          userStore.setLoginState(true)
          getInfo()
        }
      })
      .catch(() => {})
  }

  return { init }
}
