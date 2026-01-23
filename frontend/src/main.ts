import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import i18n from "./i18n";
import "./styles/global.scss";

// 根据 UA 粗略判断设备类型（用于路由分流）
const isMobileUA = /Mobi|Android|iPhone|iPad|iPod|Windows Phone/i.test(navigator.userAgent);
document.documentElement.dataset.uaDevice = isMobileUA ? "mobile" : "desktop";

// 统一计算可视区域尺寸，写入 CSS 变量供移动端布局使用
const updateViewportVars = () => {
  const visualWidth = window.visualViewport?.width ?? 0;
  const visualHeight = window.visualViewport?.height ?? 0;
  const fallbackWidth =
    document.documentElement.clientWidth || window.innerWidth || window.screen?.width || 0;
  const fallbackHeight =
    document.documentElement.clientHeight || window.innerHeight || window.screen?.height || 0;
  const width = visualWidth > 0 ? visualWidth : fallbackWidth;
  const height = visualHeight > 0 ? visualHeight : fallbackHeight;
  const safeWidth = Math.max(width, 1);
  const safeHeight = Math.max(height, 1);
  const root = document.documentElement;
  root.style.setProperty("--viewport-w", `${safeWidth}px`);
  root.style.setProperty("--viewport-h", `${safeHeight}px`);
  root.style.setProperty("--mobile-unit", `${Math.max(1, Math.min(safeWidth, safeHeight) / 100)}px`);
};

updateViewportVars();
window.addEventListener("resize", updateViewportVars);
window.addEventListener("orientationchange", updateViewportVars);
window.addEventListener("load", updateViewportVars);
window.visualViewport?.addEventListener("resize", updateViewportVars);

// 挂载应用（Pinia + Router + i18n）
createApp(App).use(createPinia()).use(router).use(i18n).mount("#app");
