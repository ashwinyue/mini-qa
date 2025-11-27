import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

const useAuthStore = create(
    persist(
        (set) => ({
            user: null,
            token: null,
            loading: false,

            // Actions
            setUser: (user) => set({ user }),
            setToken: (token) => set({ token }),
            setLoading: (loading) => set({ loading }),

            login: (token, user) => {
                set({ token, user })
            },

            logout: () => {
                set({ token: null, user: null })
            },
        }),
        {
            name: 'auth-storage',
            storage: createJSONStorage(() => localStorage),
            partialize: (state) => ({ user: state.user, token: state.token }),
        }
    )
)

export default useAuthStore
