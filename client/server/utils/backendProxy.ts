function sanitizeHeaders(headers: Record<string, string | string[] | undefined>) {
  const nextHeaders = new Headers()

  for (const [key, value] of Object.entries(headers)) {
    const normalized = key.toLowerCase()
    if (["host", "connection", "content-length"].includes(normalized)) {
      continue
    }

    if (Array.isArray(value)) {
      nextHeaders.set(key, value.join(","))
      continue
    }

    if (value) {
      nextHeaders.set(key, value)
    }
  }

  return nextHeaders
}

function buildTargetUrl(baseUrl: string, path: string, query: Record<string, unknown>) {
  const normalizedBase = baseUrl.replace(/\/$/, "")
  const normalizedPath = path ? `/${path.replace(/^\//, "")}` : ""
  const url = new URL(`${normalizedBase}${normalizedPath}`)

  for (const [key, value] of Object.entries(query)) {
    if (Array.isArray(value)) {
      for (const item of value) {
        url.searchParams.append(key, String(item))
      }
      continue
    }

    if (value !== undefined) {
      url.searchParams.set(key, String(value))
    }
  }

  return url.toString()
}

export async function proxyBackendRequest(event: H3Event, path: string) {
  const runtimeConfig = useRuntimeConfig(event)
  const headers = getHeaders(event)
  const overrideBaseUrl = headers["x-backend-base-url"]
  const targetUrl = buildTargetUrl(
    typeof overrideBaseUrl === "string" && overrideBaseUrl.trim()
      ? overrideBaseUrl.trim()
      : runtimeConfig.backendBaseUrl,
    path,
    getQuery(event)
  )

  const method = getMethod(event)
  const requestHeaders = sanitizeHeaders(headers)
  requestHeaders.delete("x-backend-base-url")

  const body = ["GET", "HEAD"].includes(method) ? undefined : await readRawBody(event, false)
  const upstreamResponse = await fetch(targetUrl, {
    method,
    headers: requestHeaders,
    body: body ?? undefined
  })

  setResponseStatus(event, upstreamResponse.status, upstreamResponse.statusText)

  const contentType = upstreamResponse.headers.get("content-type")
  if (contentType) {
    setHeader(event, "content-type", contentType)
  }

  const upstreamBody = await upstreamResponse.text()
  if (!upstreamBody) {
    return ""
  }

  return upstreamBody
}
