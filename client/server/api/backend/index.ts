import { proxyBackendRequest } from "../../utils/backendProxy"

export default defineEventHandler(async (event) => {
  return proxyBackendRequest(event, "")
})
