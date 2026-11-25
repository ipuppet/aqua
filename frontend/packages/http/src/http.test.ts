import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import {
  ACCESS_TOKEN_KEY,
  REFRESH_TOKEN_KEY,
  buildApiUrl,
  clearAuthTokens,
  configureCookies,
  createAquaClient,
  delCookie,
  getCookie,
  setCookie
} from "./index"

function installCookieJar(hostname = "admin.example.com") {
  const jar = new Map<string, string>()
  const writes: string[] = []
  Object.defineProperty(globalThis, "document", {
    configurable: true,
    value: {
      get cookie() {
        return Array.from(jar.entries())
          .map(([k, v]) => `${k}=${v}`)
          .join("; ")
      },
      set cookie(raw: string) {
        writes.push(raw)
        const [pair] = raw.split(";")
        const eq = pair.indexOf("=")
        const name = pair.slice(0, eq).trim()
        const value = pair.slice(eq + 1).trim()
        if (raw.includes("expires=Thu, 01 Jan 1970")) jar.delete(name)
        else jar.set(name, value)
      }
    }
  })
  Object.defineProperty(globalThis, "window", {
    configurable: true,
    value: { location: { hostname } }
  })
  return { jar, writes }
}

function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" }
  })
}

describe("cookies", () => {
  afterEach(() => configureCookies({ domain: "" }))

  it("writes domain + Secure; SameSite=Lax on non-localhost", () => {
    const { writes } = installCookieJar("admin.example.com")
    configureCookies({ domain: ".example.com" })
    setCookie("k", "v", 1)
    expect(writes[0]).toContain("domain=.example.com")
    expect(writes[0]).toContain("Secure")
    expect(writes[0]).toContain("SameSite=Lax")
  })

  it("uses host-only cookies when domain is empty on localhost", () => {
    const { writes } = installCookieJar("localhost")
    configureCookies({ domain: "" })
    setCookie("k", "v", 1)
    expect(writes[0]).not.toContain("domain=")
    expect(writes[0]).toContain("SameSite=Lax")
    expect(writes[0]).not.toContain("Secure")
  })

  it("round-trips get/set/del and clearAuthTokens", () => {
    installCookieJar()
    configureCookies({ domain: "" })
    setCookie(ACCESS_TOKEN_KEY, "a")
    setCookie(REFRESH_TOKEN_KEY, "r")
    expect(getCookie(ACCESS_TOKEN_KEY)).toBe("a")
    delCookie(ACCESS_TOKEN_KEY)
    expect(getCookie(ACCESS_TOKEN_KEY)).toBeNull()
    setCookie(ACCESS_TOKEN_KEY, "a2")
    clearAuthTokens()
    expect(getCookie(ACCESS_TOKEN_KEY)).toBeNull()
    expect(getCookie(REFRESH_TOKEN_KEY)).toBeNull()
  })
})

describe("buildApiUrl", () => {
  it("maps local hosts to per-service ports", () => {
    expect(buildApiUrl("admin", "localhost")).toBe("http://localhost:8082")
    expect(buildApiUrl("passport", "127.0.0.1")).toBe("http://127.0.0.1:8081")
    expect(buildApiUrl("blog", "192.168.1.8")).toBe("http://192.168.1.8:8084")
  })

  it("uses the gateway path in production", () => {
    expect(buildApiUrl("admin", "api.example.com")).toBe("https://api.example.com/admin")
  })

  it("inserts eo to match the current hostname", () => {
    expect(buildApiUrl("blog", "api.example.com", "www.eo.example.com")).toBe("https://api.eo.example.com/blog")
  })
})

describe("createAquaClient", () => {
  beforeEach(() => {
    installCookieJar()
    configureCookies({ domain: "" })
    vi.restoreAllMocks()
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it("authMode=none attaches the token and does not refresh on 401", async () => {
    setCookie(ACCESS_TOKEN_KEY, "at")
    setCookie(REFRESH_TOKEN_KEY, "rt")

    const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
      const request = input instanceof Request ? input : new Request(String(input))
      expect(request.headers.get("Authorization")).toBe("Bearer at")
      return jsonResponse(401, { message: "expired" })
    })
    vi.stubGlobal("fetch", fetchMock)

    const client = createAquaClient({
      baseUrl: "https://api.example/blog",
      authMode: "none"
    })
    const { response, error } = await client.GET("/api/about")
    expect(response.status).toBe(401)
    expect(error).toBeDefined()
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })

  it("authMode=refresh single-flights concurrent 401s into one refresh", async () => {
    setCookie(REFRESH_TOKEN_KEY, "rt")
    setCookie(ACCESS_TOKEN_KEY, "old")

    let refreshCalls = 0
    const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
      const request = input instanceof Request ? input : new Request(String(input))
      const url = request.url

      if (url.endsWith("/api/token")) {
        refreshCalls += 1
        await new Promise(r => setTimeout(r, 30))
        return jsonResponse(200, { accessToken: "new-at", refreshToken: "new-rt" })
      }

      if (request.headers.get("Authorization") === "Bearer new-at") {
        return jsonResponse(200, { articles: [] })
      }
      return jsonResponse(401, { message: "expired" })
    })
    vi.stubGlobal("fetch", fetchMock)

    const client = createAquaClient({
      baseUrl: "https://api.example/admin",
      passportApiUrl: "https://api.example/passport",
      authMode: "refresh"
    })

    const [a, b] = await Promise.all([client.GET("/api/about"), client.GET("/api/tags")])
    expect(a.error).toBeUndefined()
    expect(b.error).toBeUndefined()
    expect(refreshCalls).toBe(1)
    expect(getCookie(ACCESS_TOKEN_KEY)).toBe("new-at")
  })

  it("authMode=refresh notifies onAuthFailure once when refresh fails", async () => {
    setCookie(REFRESH_TOKEN_KEY, "rt")

    const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
      const request = input instanceof Request ? input : new Request(String(input))
      if (request.url.endsWith("/api/token")) {
        return jsonResponse(500, { message: "refresh boom" })
      }
      return jsonResponse(401, { message: "expired" })
    })
    vi.stubGlobal("fetch", fetchMock)

    const onAuthFailure = vi.fn()
    const client = createAquaClient({
      baseUrl: "https://api.example/admin",
      passportApiUrl: "https://api.example/passport",
      authMode: "refresh",
      onAuthFailure
    })

    await client.GET("/api/about")
    await client.GET("/api/about")
    expect(onAuthFailure).toHaveBeenCalledTimes(1)
    expect(onAuthFailure.mock.calls[0]?.[0]).toContain("refresh boom")
    expect(getCookie(REFRESH_TOKEN_KEY)).toBeNull()
  })

  it("does not attempt refresh for the token endpoint itself", async () => {
    setCookie(REFRESH_TOKEN_KEY, "rt")

    const fetchMock = vi.fn(async () => jsonResponse(401, { message: "bad token" }))
    vi.stubGlobal("fetch", fetchMock)

    const client = createAquaClient({
      baseUrl: "https://api.example/passport",
      passportApiUrl: "https://api.example/passport",
      authMode: "refresh"
    })

    const { response } = await client.POST("/api/token", {
      body: { refreshToken: "rt" }
    })
    expect(response.status).toBe(401)
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })
})
