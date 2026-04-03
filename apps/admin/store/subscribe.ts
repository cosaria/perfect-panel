import { create } from "zustand";
import { getSubscribeList } from "@/services/admin-api/sdk.gen";
import type { SubscribeItem } from "@/services/admin-api/types.gen";

interface SubscribeState {
  // Data
  subscribes: SubscribeItem[];

  // Loading states
  loading: boolean;
  loaded: boolean;

  // Actions
  fetchSubscribes: () => Promise<void>;

  // Getters
  getSubscribeName: (subscribeId?: number) => string;
  getSubscribeById: (subscribeId: number) => SubscribeItem | undefined;
}

export const useSubscribeStore = create<SubscribeState>((set, get) => ({
  // Initial state
  subscribes: [],
  loading: false,
  loaded: false,

  // Actions
  fetchSubscribes: async () => {
    if (get().loading) return;

    set({ loading: true });
    try {
      const { data } = await getSubscribeList({ query: { page: 1, size: 999999999 } });
      set({
        subscribes: data?.list || [],
        loaded: true,
      });
    } catch (_error) {
      // Handle error silently
      set({ loaded: true });
    } finally {
      set({ loading: false });
    }
  },

  // Getters
  getSubscribeName: (subscribeId?: number) => {
    if (!subscribeId) return "--";
    const subscribe = get().subscribes.find((s) => s.id === subscribeId);
    return subscribe?.name ?? `Subscribe ${subscribeId}`;
  },

  getSubscribeById: (subscribeId: number) => {
    return get().subscribes.find((s) => s.id === subscribeId);
  },
}));

export const useSubscribe = () => {
  const store = useSubscribeStore();

  // Auto-fetch subscribes
  if (!store.loaded && !store.loading) {
    store.fetchSubscribes();
  }

  return {
    subscribes: store.subscribes,
    loading: store.loading,
    loaded: store.loaded,
    fetchSubscribes: store.fetchSubscribes,
    getSubscribeName: store.getSubscribeName,
    getSubscribeById: store.getSubscribeById,
  };
};

export default useSubscribeStore;
