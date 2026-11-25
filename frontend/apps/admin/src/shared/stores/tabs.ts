import { defineStore } from "pinia"
import type { RouteLocationNormalized } from "vue-router"

export const useTabsStore = defineStore("tabs", {
  state: () => {
    let visitedViews: RouteLocationNormalized[] = []
    let cachedViews: string[] = []
    try {
      const savedRoutes = localStorage.getItem("TABS_ROUTES")
      const savedCached = localStorage.getItem("TABS_CACHED")
      if (savedRoutes) visitedViews = JSON.parse(savedRoutes)
      if (savedCached) cachedViews = JSON.parse(savedCached)
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error("Failed to load tabs from localStorage", e)
    }
    return {
      visitedViews,
      cachedViews
    }
  },
  actions: {
    addView(view: RouteLocationNormalized) {
      this.addVisitedView(view)
      this.addCachedView(view)
    },
    addVisitedView(view: RouteLocationNormalized) {
      if (this.visitedViews.some((v: RouteLocationNormalized) => v.path === view.path)) return
      // limit to 10
      if (this.visitedViews.length >= 10) {
        this.visitedViews.shift()
      }
      this.visitedViews.push(
        Object.assign({}, view, {
          title: (view.meta.title as string) ?? "no-name"
        })
      )
      this.saveToLocalStorage()
    },
    addCachedView(view: RouteLocationNormalized) {
      if (this.cachedViews.includes(view.name as string)) return
      if (!view.meta.noCache) {
        this.cachedViews.push(view.name as string)
      }
      this.saveToLocalStorage()
    },
    delView(view: RouteLocationNormalized) {
      return new Promise(resolve => {
        this.delVisitedView(view)
        this.delCachedView(view)
        this.saveToLocalStorage()
        resolve({
          visitedViews: [...this.visitedViews],
          cachedViews: [...this.cachedViews]
        })
      })
    },
    delVisitedView(view: RouteLocationNormalized) {
      for (const [i, v] of this.visitedViews.entries()) {
        if (v.path === view.path) {
          this.visitedViews.splice(i, 1)
          break
        }
      }
      this.saveToLocalStorage()
    },
    delCachedView(view: RouteLocationNormalized) {
      const index = this.cachedViews.indexOf(view.name as string)
      index > -1 && this.cachedViews.splice(index, 1)
      this.saveToLocalStorage()
    },
    delOthersViews(view: RouteLocationNormalized) {
      return new Promise(resolve => {
        this.delOthersVisitedViews(view)
        this.delOthersCachedViews(view)
        this.saveToLocalStorage()
        resolve({
          visitedViews: [...this.visitedViews],
          cachedViews: [...this.cachedViews]
        })
      })
    },
    delOthersVisitedViews(view: RouteLocationNormalized) {
      this.visitedViews = this.visitedViews.filter((v: RouteLocationNormalized) => {
        return (v.meta.affix as boolean) || v.path === view.path
      })
      this.saveToLocalStorage()
    },
    delOthersCachedViews(view: RouteLocationNormalized) {
      const index = this.cachedViews.indexOf(view.name as string)
      if (index > -1) {
        this.cachedViews = this.cachedViews.slice(index, index + 1)
      } else {
        this.cachedViews = []
      }
      this.saveToLocalStorage()
    },
    delAllViews() {
      return new Promise(resolve => {
        this.delAllVisitedViews()
        this.delAllCachedViews()
        this.saveToLocalStorage()
        resolve({
          visitedViews: [...this.visitedViews],
          cachedViews: [...this.cachedViews]
        })
      })
    },
    delAllVisitedViews() {
      // keep affix tags
      const affixTags = this.visitedViews.filter((tag: RouteLocationNormalized) => tag.meta.affix)
      this.visitedViews = affixTags
      this.saveToLocalStorage()
    },
    delAllCachedViews() {
      this.cachedViews = []
      this.saveToLocalStorage()
    },
    updateVisitedView(view: RouteLocationNormalized) {
      for (let v of this.visitedViews) {
        if (v.path === view.path) {
          v = Object.assign(v, view)
          break
        }
      }
      this.saveToLocalStorage()
    },
    saveToLocalStorage() {
      const visitedViewsState = this.visitedViews.map((v: any) => ({
        name: v.name,
        path: v.path,
        fullPath: v.fullPath,
        meta: v.meta,
        params: v.params,
        query: v.query
      }))
      localStorage.setItem("TABS_ROUTES", JSON.stringify(visitedViewsState))
      localStorage.setItem("TABS_CACHED", JSON.stringify(this.cachedViews))
    }
  }
})
