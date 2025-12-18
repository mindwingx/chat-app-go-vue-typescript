import { ApiClient } from "./client/axios";

const api = new ApiClient(`${import.meta.env.VITE_API_URL}`)

export const handshake = {
  async check(): Promise<boolean> {
    try {
      await api.get("/handshake")
      return true
    } catch(err) {
      console.error("server unavailable")
      return false
    }
  }
}
