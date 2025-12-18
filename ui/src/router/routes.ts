import {createRouter, createWebHistory} from "vue-router";
import Join from "./../views/Join.vue";
import Chat from "./../views/Chat.vue";
import { joinGuard } from "./middleware.ts";
import NotFound from "./../views/NotFound.vue";

const routes = [
    {path: "/", component: Join},
    {path: "/chat", component: Chat},
    {path: "/:pathMatch(.*)*", component: NotFound}
]

const router = createRouter({
    history:createWebHistory(),
    routes,
})

router.beforeEach(joinGuard)

export { router }