import { createRouter, createWebHistory } from "vue-router";
import { routes } from "./routes";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(_to, from, savedPosition) {
    return new Promise((resolve) => {
      setTimeout(
        () => {
          if (savedPosition) {
            return resolve({ ...savedPosition });
          }
        },
        from.name ? 1000 : 500
      );
    });
  },
});

export default router;
