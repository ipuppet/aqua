/* eslint-disable no-console */

class LocalData {
  get<T = any>(name: string, defaultValue: T | null = null): T | null {
    try {
      const objData = localStorage.getItem(name + "_obj")
      if (objData) {
        const data = JSON.parse(objData)
        const now = new Date().getTime()
        if (data.life !== undefined && now > data.life) {
          this.remove(name)
          return defaultValue
        }
        return data.obj as T
      }
      const strData = localStorage.getItem(name + "_str")
      if (strData) {
        return strData as any as T
      }
      return defaultValue
    } catch (error) {
      console.error(`Error retrieving from localStorage: ${name}`, error)
      return defaultValue
    }
  }

  set<T = any>(name: string, value: T, lifeSeconds?: number): boolean {
    try {
      if (typeof value === "object" && value !== null) {
        const data: any = { obj: value }
        if (lifeSeconds !== undefined) {
          data.life = new Date().getTime() + lifeSeconds * 1000
        }
        localStorage.setItem(name + "_obj", JSON.stringify(data))
      } else {
        localStorage.setItem(name + "_str", String(value))
      }
      return true
    } catch (error) {
      console.error(`Error setting to localStorage: ${name}`, error)
      return false
    }
  }

  remove(name: string): boolean {
    try {
      localStorage.removeItem(name + "_obj")
      localStorage.removeItem(name + "_str")
      return true
    } catch (error) {
      console.error(`Error removing from localStorage: ${name}`, error)
      return false
    }
  }
}

class SessionData {
  get<T = any>(name: string, defaultValue: T | null = null): T | null {
    try {
      const objData = sessionStorage.getItem(name + "_obj")
      if (objData) {
        return JSON.parse(objData) as T
      }
      const strData = sessionStorage.getItem(name + "_str")
      if (strData) {
        return strData as any as T
      }
      return defaultValue
    } catch (error) {
      console.error(`Error retrieving from sessionStorage: ${name}`, error)
      return defaultValue
    }
  }

  set<T = any>(name: string, value: T): boolean {
    try {
      if (typeof value === "object" && value !== null) {
        sessionStorage.setItem(name + "_obj", JSON.stringify(value))
      } else {
        sessionStorage.setItem(name + "_str", String(value))
      }
      return true
    } catch (error) {
      console.error(`Error setting to sessionStorage: ${name}`, error)
      return false
    }
  }

  remove(name: string): boolean {
    try {
      sessionStorage.removeItem(name + "_obj")
      sessionStorage.removeItem(name + "_str")
      return true
    } catch (error) {
      console.error(`Error removing from sessionStorage: ${name}`, error)
      return false
    }
  }
}

export const localData = new LocalData()
export const sessionData = new SessionData()
