import { useUserStore } from "@/stores/user";
import type { NavigationGuard, RouteRecordRaw } from "vue-router";

const checkAuth: NavigationGuard = () => {
  const user = useUserStore();

  if (user.loggedIn) return true;
  else return { path: "Home" };
};

export const routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "Home",
    component: () => import("@/views/HomePage.vue"),
    props: ({ query }) => ({ code: query.verificationCode }),
    alias: "/verify",
  },
  {
    path: "/search",
    component: () => import("@/views/ResultsPage.vue"),
    props: ({ query }) => ({ query: query.query }),
    children: [
      {
        path: "",
        name: "Results",
        component: () => import("@/components/SearchResults.vue"),
        props: ({ query }) => ({ query: query.query }),
      },
    ],
  },
  {
    path: "/wardrobe",
    name: "Wardrobe",
    component: () => import("@/views/UserWardrobe.vue"),
    beforeEnter: checkAuth,
  },
  {
    path: "/profile",
    name: "Profile",
    component: () => import("@/views/UserSettings.vue"),
    beforeEnter: checkAuth,
  },
  {
    path: "/:pathMatch(.*)",
    redirect: "/",
  },
];
