import { onBeforeUnmount, onMounted, ref } from "vue";

type Theme = "light" | "dark";

// storageKey 主题持久化 key
const storageKey = "blog-theme";
const theme = ref<Theme>("light");
let hasManual = false;
let media: MediaQueryList | null = null;
let initialized = false;

export function useTheme() {
  // applyTheme 统一写入 state 与 DOM data-theme
  function applyTheme(next: Theme) {
    theme.value = next;
    document.documentElement.dataset.theme = next;
  }

  // handleSystemChange 跟随系统主题变化（仅在未手动切换时生效）
  function handleSystemChange(event: MediaQueryListEvent) {
    if (hasManual) return;
    applyTheme(event.matches ? "dark" : "light");
  }

  // initTheme 读取本地配置，或跟随系统主题
  function initTheme() {
    if (initialized) return;
    const saved = localStorage.getItem(storageKey) as Theme | null;
    if (saved === "light" || saved === "dark") {
      hasManual = true;
      applyTheme(saved);
      initialized = true;
      return;
    }
    media = window.matchMedia("(prefers-color-scheme: dark)");
    applyTheme(media.matches ? "dark" : "light");
    media.addEventListener("change", handleSystemChange);
    initialized = true;
  }

  function toggleTheme() {
    // 手动切换主题并持久化
    const next = theme.value === "dark" ? "light" : "dark";
    hasManual = true;
    localStorage.setItem(storageKey, next);
    applyTheme(next);
  }

  onMounted(() => {
    initTheme();
  });

  onBeforeUnmount(() => {
    media?.removeEventListener("change", handleSystemChange);
  });

  return {
    theme,
    toggleTheme,
  };
}
