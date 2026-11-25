/**
 * 本地和会话数据存储模块
 *
 * 提供类型安全的存储接口，支持：
 * - LocalData: 使用 localStorage 的持久化存储（支持过期时间）
 * - SessionData: 使用 sessionStorage 的会话存储
 */

/* eslint-disable no-console */

/**
 * 本地存储类 - 支持过期时间
 *
 * 使用 JSON 序列化，支持对象和基本类型
 * 支持可选的生命周期设置（单位：秒）
 */
class LocalData {
  /**
   * 获取存储的值
   * @param name 键名
   * @param defaultValue 如果不存在或已过期时的默认值
   * @returns 存储的值或默认值
   */
  get<T = any>(name: string, defaultValue: T | null = null): T | null {
    try {
      // 尝试获取对象存储
      const objData = localStorage.getItem(name + "_obj")
      if (objData) {
        const data = JSON.parse(objData)
        const now = new Date().getTime()

        // 检查是否过期
        if (data.life !== undefined && now > data.life) {
          this.remove(name)
          return defaultValue
        }

        return data.obj as T
      }

      // 尝试获取字符串存储
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

  /**
   * 设置存储的值
   * @param name 键名
   * @param value 要存储的值
   * @param lifeSeconds 生命周期（秒），不设置则永久存储
   * @returns 是否设置成功
   */
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

  /**
   * 删除存储的值
   * @param name 键名
   * @returns 是否删除成功
   */
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

  /**
   * 清空所有本地存储
   * @returns 是否清空成功
   */
  clear(): boolean {
    try {
      localStorage.clear()
      return true
    } catch (error) {
      console.error("Error clearing localStorage", error)
      return false
    }
  }
}

/**
 * 会话存储类 - 不支持过期时间
 *
 * 在当前浏览器标签页关闭时自动清空
 * 使用 JSON 序列化，支持对象和基本类型
 */
class SessionData {
  /**
   * 获取存储的值
   * @param name 键名
   * @param defaultValue 如果不存在时的默认值
   * @returns 存储的值或默认值
   */
  get<T = any>(name: string, defaultValue: T | null = null): T | null {
    try {
      // 尝试获取对象存储
      const objData = sessionStorage.getItem(name + "_obj")
      if (objData) {
        return JSON.parse(objData) as T
      }

      // 尝试获取字符串存储
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

  /**
   * 设置存储的值
   * @param name 键名
   * @param value 要存储的值
   * @returns 是否设置成功
   */
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

  /**
   * 删除存储的值
   * @param name 键名
   * @returns 是否删除成功
   */
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

  /**
   * 清空所有会话存储
   * @returns 是否清空成功
   */
  clear(): boolean {
    try {
      sessionStorage.clear()
      return true
    } catch (error) {
      console.error("Error clearing sessionStorage", error)
      return false
    }
  }
}

// 单例实例
export const localData = new LocalData()
export const sessionData = new SessionData()

export default {
  localData,
  sessionData
}
