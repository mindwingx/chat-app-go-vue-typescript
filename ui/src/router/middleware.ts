import type { NavigationGuardNext, RouteLocationNormalized } from "vue-router"
import { useUserStore } from "./../stores/user"

export function joinGuard(
    to: RouteLocationNormalized, 
    _from: RouteLocationNormalized, 
    next: NavigationGuardNext
) {
    const userStore = useUserStore()

    if (userStore.active === true){
        if (to.path === "/") {
            next("/chat")
        } else {
            next()
        }
    } else {
        if (to.path === "/chat") {
            next("/")
        } else {
            next()
        }
    }
}