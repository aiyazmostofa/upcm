import { writable } from "svelte/store"

const user = writable({ ID: "", token: "", authLevel: "" })

export default { user }