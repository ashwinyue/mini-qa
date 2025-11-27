/**
 * 系统状态 Store
 * 
 * 管理系统配置、模型和租户信息
 */
import { create } from 'zustand'

/**
 * 创建系统状态 Store
 */
const useSystemStore = create((set) => ({
    // 状态
    config: null,        // 系统配置
    isLoading: false,    // 加载状态
    error: null,         // 错误信息

    // 设置配置
    setConfig: (config) => set({ config }),

    // 设置加载状态
    setLoading: (loading) => set({ isLoading: loading }),

    // 设置错误
    setError: (error) => set({ error }),

    // 更新当前模型
    updateCurrentModel: (model) =>
        set((state) => ({
            config: state.config ? { ...state.config, currentModel: model } : null,
        })),

    // 设置租户 ID
    setTenantId: (tenantId) => {
        localStorage.setItem('tenantId', tenantId)
        set((state) => ({
            config: state.config ? { ...state.config, tenantId } : null,
        }))
    },

    // 获取租户 ID
    getTenantId: () => {
        return localStorage.getItem('tenantId') || 'default'
    },
}))

export default useSystemStore
