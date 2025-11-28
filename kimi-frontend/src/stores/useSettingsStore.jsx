/**
 * 设置状态 Store
 * 
 * 管理用户设置（主题、语言等）
 */
import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

/**
 * 创建设置状态 Store
 */
const useSettingsStore = create(
    persist(
        (set, get) => ({
            // 设置状态
            theme: 'light',              // 主题: 'light' | 'dark' | 'system'
            language: 'zh',              // 语言: 'zh' | 'en'
            autoExpandReferences: true,  // 自动展开参考网页

            // 设置主题
            setTheme: (theme) => {
                set({ theme })
                // 应用主题到 document
                applyTheme(theme)
            },

            // 设置语言
            setLanguage: (language) => {
                set({ language })
            },

            // 设置自动展开参考
            setAutoExpandReferences: (autoExpand) => {
                set({ autoExpandReferences: autoExpand })
            },

            // 重置所有设置
            resetSettings: () => {
                set({
                    theme: 'light',
                    language: 'zh',
                    autoExpandReferences: true,
                })
                applyTheme('light')
            },
        }),
        {
            name: 'kimi-settings-store',
            storage: createJSONStorage(() => localStorage),
            // 持久化所有设置
            partialize: (state) => ({
                theme: state.theme,
                language: state.language,
                autoExpandReferences: state.autoExpandReferences,
            }),
        }
    )
)

/**
 * 应用主题到 document
 */
function applyTheme(theme) {
    const root = document.documentElement
    
    if (theme === 'system') {
        // 检测系统主题
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        root.classList.toggle('dark', prefersDark)
    } else {
        root.classList.toggle('dark', theme === 'dark')
    }
}

// 初始化时应用保存的主题
if (typeof window !== 'undefined') {
    const savedTheme = localStorage.getItem('kimi-settings-store')
    if (savedTheme) {
        try {
            const { state } = JSON.parse(savedTheme)
            if (state?.theme) {
                applyTheme(state.theme)
            }
        } catch (e) {
            console.error('Failed to parse saved theme:', e)
        }
    }
}

export default useSettingsStore
