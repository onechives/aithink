import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/Home.vue";
import PostDetailView from "../views/PostDetail.vue";
import LoginView from "../views/Login.vue";
import RegisterView from "../views/Register.vue";
import AdminEditorView from "../views/AdminEditor.vue";
import AdminPostsView from "../views/AdminPosts.vue";
import AdminReviewView from "../views/AdminReview.vue";
import AdminUsersView from "../views/AdminUsers.vue";
import AdminNicknamesView from "../views/AdminNicknames.vue";
import MyPostsView from "../views/MyPosts.vue";
import AccountSettingsView from "../views/AccountSettings.vue";
import MessagesView from "../views/Messages.vue";
import MobileShell from "../components/MobileShell.vue";
import MobileHome from "../views/mobile/MobileHome.vue";
import MobilePostDetail from "../views/mobile/MobilePostDetail.vue";
import MobileEditor from "../views/mobile/MobileEditor.vue";
import MobileSettings from "../views/mobile/MobileSettings.vue";
import MobileMessages from "../views/mobile/MobileMessages.vue";
import MobileLogin from "../views/mobile/MobileLogin.vue";
import MobileRegister from "../views/mobile/MobileRegister.vue";
import MobileMyPosts from "../views/mobile/MobileMyPosts.vue";
import { useAuthStore } from "../stores/auth";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // PC 端路由
    { path: "/", name: "home", component: HomeView, meta: { titleKey: "nav.articles" } },
    { path: "/posts/:id", name: "post-detail", component: PostDetailView, meta: { titleKey: "post.detail" } },
    { path: "/login", name: "login", component: LoginView, meta: { titleKey: "login.title", hideMobileNav: true } },
    { path: "/register", name: "register", component: RegisterView, meta: { titleKey: "register.title", hideMobileNav: true } },
    { path: "/write/:id?", name: "editor", component: AdminEditorView, meta: { titleKey: "editor.newTitle" } },
    { path: "/me/posts", name: "my-posts", component: MyPostsView, meta: { titleKey: "myPosts.title" } },
    { path: "/me/settings", name: "account-settings", component: AccountSettingsView, meta: { titleKey: "settings.title" } },
    { path: "/me/messages", name: "messages", component: MessagesView, meta: { titleKey: "messages.title" } },
    { path: "/admin/posts", name: "admin-posts", component: AdminPostsView, meta: { titleKey: "admin.postManage" } },
    { path: "/admin/review/posts", name: "admin-review-posts", component: AdminReviewView, meta: { titleKey: "admin.postReview" } },
    { path: "/admin/review/users", name: "admin-review-users", component: AdminUsersView, meta: { titleKey: "admin.userReview" } },
    { path: "/admin/review/nicknames", name: "admin-review-nicknames", component: AdminNicknamesView, meta: { titleKey: "admin.nicknameReview" } },
    {
      path: "/m",
      component: MobileShell,
      children: [
        // 移动端独立路由
        { path: "", name: "m-home", component: MobileHome, meta: { titleKey: "nav.articles" } },
        { path: "posts/:id", name: "m-post-detail", component: MobilePostDetail, meta: { titleKey: "post.detail" } },
        { path: "write/:id?", name: "m-editor", component: MobileEditor, meta: { titleKey: "editor.newTitle" } },
        { path: "me/posts", name: "m-my-posts", component: MobileMyPosts, meta: { titleKey: "myPosts.title" } },
        { path: "settings", name: "m-settings", component: MobileSettings, meta: { titleKey: "settings.title" } },
        { path: "messages", name: "m-messages", component: MobileMessages, meta: { titleKey: "messages.title" } },
        { path: "login", name: "m-login", component: MobileLogin, meta: { titleKey: "login.title", hideMobileNav: true } },
        { path: "register", name: "m-register", component: MobileRegister, meta: { titleKey: "register.title", hideMobileNav: true } },
      ],
    },
  ],
});

router.beforeEach((to) => {
  const auth = useAuthStore();
  // 只有写作和消息需要强制登录
  const needsLogin = to.name === "editor" || to.name === "messages" || to.name === "m-editor" || to.name === "m-messages";
  if (needsLogin) {
    if (!auth.token) {
      const loginTarget = to.name?.toString().startsWith("m-") ? { name: "m-login" } : { name: "login" };
      return loginTarget;
    }
  }
  // 管理端路由必须是管理员
  if (to.name === "admin-posts" || to.name === "admin-review-posts" || to.name === "admin-review-users" || to.name === "admin-review-nicknames") {
    if (!auth.token || auth.role !== "admin") {
      return { name: "login" };
    }
  }

  const isMobileUA = document.documentElement.dataset.uaDevice === "mobile";
  const isMobileRoute = to.path === "/m" || to.path.startsWith("/m/");
  // 根据设备类型做路由重定向（移动端强制 /m 前缀）
  if (isMobileUA && !isMobileRoute && !to.path.startsWith("/admin")) {
    const targetPath = to.path === "/" ? "/m" : `/m${to.path}`;
    return { path: targetPath, query: to.query, hash: to.hash };
  }
  // PC 设备访问 /m 时自动回到 PC 路由
  if (!isMobileUA && isMobileRoute) {
    const targetPath = to.path.replace(/^\/m/, "") || "/";
    return { path: targetPath, query: to.query, hash: to.hash };
  }
});

export default router;
