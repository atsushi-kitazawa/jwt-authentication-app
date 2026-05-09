<script setup lang="ts">
type FieldSpec = {
  key: string
  label: string
  placeholder: string
  type?: "text" | "email" | "password"
}

type ActionSpec = {
  id: string
  title: string
  method: "GET" | "POST" | "PUT" | "DELETE"
  path: string
  description: string
  requiresAuth: boolean
  fields: FieldSpec[]
}

type ResponseState = {
  title: string
  status: string
  body: string
  tone: "idle" | "success" | "error"
}

const DEFAULT_ACTION_ID = "health"

const config = useRuntimeConfig()

const actions: ActionSpec[] = [
  {
    id: "health",
    title: "Health Check",
    method: "GET",
    path: "/health",
    description: "認証なしで API の生存確認を行います。",
    requiresAuth: false,
    fields: []
  },
  {
    id: "login",
    title: "Login",
    method: "POST",
    path: "/login",
    description: "name と password を送り、JWT を取得して保持します。",
    requiresAuth: false,
    fields: [
      { key: "name", label: "Name", placeholder: "alice" },
      { key: "password", label: "Password", placeholder: "secret", type: "password" }
    ]
  },
  {
    id: "root",
    title: "Protected Root",
    method: "GET",
    path: "/",
    description: "JWT を付けて保護されたトップAPIを確認します。",
    requiresAuth: true,
    fields: []
  },
  {
    id: "create-user",
    title: "Create User",
    method: "POST",
    path: "/users",
    description: "新規ユーザーを作成します。認証は不要です。",
    requiresAuth: false,
    fields: [
      { key: "name", label: "Name", placeholder: "alice" },
      { key: "email", label: "Email", placeholder: "alice@example.com", type: "email" },
      { key: "password", label: "Password", placeholder: "secret", type: "password" }
    ]
  },
  {
    id: "list-users",
    title: "List Users",
    method: "GET",
    path: "/users",
    description: "登録済みユーザー一覧を取得します。",
    requiresAuth: true,
    fields: []
  },
  {
    id: "get-user",
    title: "Get User",
    method: "GET",
    path: "/users/{id}",
    description: "指定した ID のユーザーを取得します。",
    requiresAuth: true,
    fields: [
      { key: "id", label: "User ID", placeholder: "1" }
    ]
  },
  {
    id: "get-user-by-name",
    title: "Get User By Name",
    method: "GET",
    path: "/users/name/{name}",
    description: "指定した name のユーザーを取得します。",
    requiresAuth: true,
    fields: [
      { key: "name", label: "Name", placeholder: "alice" }
    ]
  },
  {
    id: "update-user",
    title: "Update User",
    method: "PUT",
    path: "/users/{id}",
    description: "ユーザー情報を更新します。password も必須です。",
    requiresAuth: true,
    fields: [
      { key: "id", label: "User ID", placeholder: "1" },
      { key: "name", label: "Name", placeholder: "alice-updated" },
      { key: "email", label: "Email", placeholder: "alice@example.com", type: "email" },
      { key: "password", label: "Password", placeholder: "new-secret", type: "password" }
    ]
  },
  {
    id: "delete-user",
    title: "Delete User",
    method: "DELETE",
    path: "/users/{id}",
    description: "ユーザーを論理削除します。",
    requiresAuth: true,
    fields: [
      { key: "id", label: "User ID", placeholder: "1" }
    ]
  }
]

const defaultAction =
  actions.find((action) => action.id === DEFAULT_ACTION_ID) ??
  actions[0]

if (!defaultAction) {
  throw new Error(`Default action '${DEFAULT_ACTION_ID}' is not configured.`)
}

const selectedActionId = ref(DEFAULT_ACTION_ID)
const backendBaseUrl = ref("http://127.0.0.1:8080")
const token = ref("")
const pending = ref(false)
const response = ref<ResponseState>({
  title: "Ready",
  status: "画面右下にレスポンスが表示されます。",
  body: "最初は Health Check か Create User から始めるのがおすすめです。",
  tone: "idle"
})

const formValues = reactive<Record<string, Record<string, string>>>(
  Object.fromEntries(actions.map((action) => [action.id, {}]))
)

const selectedAction = computed<ActionSpec>(() => {
  return actions.find((action) => action.id === selectedActionId.value) ?? defaultAction
})

function getActionForm(actionId: string) {
  if (!formValues[actionId]) {
    formValues[actionId] = {}
  }

  return formValues[actionId]
}

const selectedFormValues = computed(() => getActionForm(selectedAction.value.id))

const tokenPreview = computed(() => {
  if (!token.value) {
    return "JWTはまだありません。"
  }

  if (token.value.length <= 48) {
    return token.value
  }

  return `${token.value.slice(0, 24)}...${token.value.slice(-16)}`
})

const canSubmit = computed(() => {
  if (pending.value) {
    return false
  }

  if (selectedAction.value.requiresAuth && !token.value.trim()) {
    return false
  }

  return selectedAction.value.fields.every((field) => {
    return (selectedFormValues.value[field.key] || "").trim() !== ""
  })
})

onMounted(() => {
  const storedBaseUrl = localStorage.getItem("jwt-client-base-url")
  const storedToken = localStorage.getItem("jwt-client-token")
  const storedForms = localStorage.getItem("jwt-client-forms")

  if (storedBaseUrl) {
    backendBaseUrl.value = storedBaseUrl
  }
  if (storedToken) {
    token.value = storedToken
  }
  if (storedForms) {
    const parsed = JSON.parse(storedForms) as Record<string, Record<string, string>>
    for (const action of actions) {
      formValues[action.id] = { ...getActionForm(action.id), ...(parsed[action.id] || {}) }
    }
  }
})

watch(backendBaseUrl, (value) => {
  if (import.meta.client) {
    localStorage.setItem("jwt-client-base-url", value)
  }
})

watch(token, (value) => {
  if (import.meta.client) {
    localStorage.setItem("jwt-client-token", value)
  }
})

watch(
  formValues,
  (value) => {
    if (import.meta.client) {
      localStorage.setItem("jwt-client-forms", JSON.stringify(value))
    }
  },
  { deep: true }
)

function selectAction(actionId: string) {
  selectedActionId.value = actionId
}

function resetToken() {
  token.value = ""
  response.value = {
    title: "Token Cleared",
    status: "保存中のJWTを削除しました。",
    body: "認証が必要なAPIを試す場合は、もう一度 Login を実行してください。",
    tone: "idle"
  }
}

function prettyPrintBody(raw: string) {
  if (!raw) {
    return "(no content)"
  }

  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}

function buildBackendPath(action: ActionSpec) {
  const actionForm = getActionForm(action.id)
  let path = action.path
  for (const field of action.fields) {
    const value = (actionForm[field.key] || "").trim()
    path = path.replace(`{${field.key}}`, encodeURIComponent(value))
  }

  return path
}

function buildRequestBody(action: ActionSpec) {
  const actionForm = getActionForm(action.id)

  if (action.id === "login") {
    return {
      name: actionForm.name?.trim() || "",
      password: actionForm.password?.trim() || ""
    }
  }

  if (action.id === "create-user" || action.id === "update-user") {
    return {
      name: actionForm.name?.trim() || "",
      email: actionForm.email?.trim() || "",
      password: actionForm.password?.trim() || ""
    }
  }

  return undefined
}

async function runAction(action = selectedAction.value) {
  if (!backendBaseUrl.value.trim()) {
    response.value = {
      title: "Validation Error",
      status: "Backend Base URL は必須です。",
      body: "画面上部の Backend Base URL を入力してください。",
      tone: "error"
    }
    return
  }

  if (action.requiresAuth && !token.value.trim()) {
    response.value = {
      title: "Validation Error",
      status: "JWT token is required.",
      body: "このAPIは認証必須です。先に Login を成功させてください。",
      tone: "error"
    }
    return
  }

  const actionForm = getActionForm(action.id)
  const missingField = action.fields.find((field) => {
    return !(actionForm[field.key] || "").trim()
  })
  if (missingField) {
    response.value = {
      title: "Validation Error",
      status: `${missingField.label} is required.`,
      body: "不足している入力項目を埋めてから、もう一度実行してください。",
      tone: "error"
    }
    return
  }

  pending.value = true

  const startedAt = performance.now()
  const backendPath = buildBackendPath(action)
  const proxyBase = config.public.apiProxyBase || "/api/backend"
  const proxyUrl = backendPath === "/" ? proxyBase : `${proxyBase}${backendPath}`
  const requestBody = buildRequestBody(action)

  try {
    const rawResponse = await fetch(proxyUrl, {
      method: action.method,
      headers: {
        "content-type": requestBody ? "application/json" : "",
        "x-backend-base-url": backendBaseUrl.value.trim(),
        ...(token.value.trim() ? { authorization: `Bearer ${token.value.trim()}` } : {})
      },
      body: requestBody ? JSON.stringify(requestBody) : undefined
    })

    const rawText = await rawResponse.text()
    const elapsed = Math.round(performance.now() - startedAt)

    response.value = {
      title: `${action.method} ${action.path}`,
      status: `${rawResponse.status} ${rawResponse.statusText} in ${elapsed}ms`,
      body: prettyPrintBody(rawText),
      tone: rawResponse.ok ? "success" : "error"
    }

    if (action.id === "login" && rawResponse.ok) {
      try {
        const parsed = JSON.parse(rawText) as { token?: string }
        token.value = parsed.token || ""
      } catch {
        token.value = ""
      }
    }
  } catch (error) {
    response.value = {
      title: "Request Failed",
      status: "Nuxt proxy request failed.",
      body: error instanceof Error ? error.message : "Unknown error",
      tone: "error"
    }
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <div class="page-shell">
    <div class="ambient ambient-left" />
    <div class="ambient ambient-right" />

    <main class="app-frame">
      <section class="hero-card">
        <div>
          <p class="eyebrow">JWT Authentication Client</p>
          <h1>Nuxt で API を触れる本番向けクライアント</h1>
          <p class="hero-copy">
            ログイン、JWT保持、各エンドポイントの実行、JSONレスポンス確認までをブラウザでまとめて扱えます。
          </p>
        </div>

        <div class="hero-controls">
          <label class="field">
            <span>Backend Base URL</span>
            <input v-model="backendBaseUrl" type="url" placeholder="http://127.0.0.1:8080">
          </label>

          <div class="session-card">
            <div class="session-head">
              <span>Session</span>
              <button type="button" class="ghost-button" @click="resetToken">
                Clear Token
              </button>
            </div>
            <p :class="token ? 'session-state ok' : 'session-state muted'">
              {{ token ? "Authenticated" : "Not logged in" }}
            </p>
            <pre class="token-box">{{ tokenPreview }}</pre>
          </div>
        </div>
      </section>

      <section class="workspace">
        <aside class="menu-panel">
          <div class="panel-header">
            <h2>API Menu</h2>
            <p>実行したい API を選択</p>
          </div>

          <button
            v-for="action in actions"
            :key="action.id"
            type="button"
            class="api-button"
            :class="{ active: action.id === selectedAction.id }"
            @click="selectAction(action.id)"
          >
            <span class="api-method">{{ action.method }}</span>
            <span class="api-text">
              <strong>{{ action.title }}</strong>
              <small>{{ action.path }}</small>
            </span>
          </button>
        </aside>

        <section class="detail-stack">
          <article class="detail-card">
            <div class="panel-header">
              <h2>{{ selectedAction.title }}</h2>
              <p>{{ selectedAction.description }}</p>
            </div>

            <div class="pill-row">
              <span class="route-pill">{{ selectedAction.method }} {{ selectedAction.path }}</span>
              <span :class="selectedAction.requiresAuth ? 'auth-pill locked' : 'auth-pill open'">
                {{ selectedAction.requiresAuth ? "JWT required" : "Public endpoint" }}
              </span>
            </div>

            <div class="form-grid">
              <label
                v-for="field in selectedAction.fields"
                :key="`${selectedAction.id}-${field.key}`"
                class="field"
              >
                <span>{{ field.label }}</span>
                <input
                  v-model="selectedFormValues[field.key]"
                  :type="field.type || 'text'"
                  :placeholder="field.placeholder"
                >
              </label>
            </div>

            <div class="action-row">
              <button type="button" class="primary-button" :disabled="!canSubmit" @click="runAction()">
                {{ pending ? "Running..." : "Run Request" }}
              </button>
              <p class="helper-copy">
                {{ selectedAction.requiresAuth ? "Login 成功後に実行できます。" : "この API は未ログインでも試せます。" }}
              </p>
            </div>
          </article>

          <article class="response-card">
            <div class="panel-header">
              <h2>{{ response.title }}</h2>
              <p :class="['response-status', response.tone]">
                {{ response.status }}
              </p>
            </div>

            <pre class="response-box">{{ response.body }}</pre>
          </article>
        </section>
      </section>
    </main>
  </div>
</template>
