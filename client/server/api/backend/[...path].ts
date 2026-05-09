import { proxyBackendRequest } from "../../utils/backendProxy"

export default defineEventHandler(async (event) => {
  const rawPath = event.context.params?.path
  const path = Array.isArray(rawPath) ? rawPath.join("/") : rawPath || ""
  return proxyBackendRequest(event, path)
})
