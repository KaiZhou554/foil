<!--
  GradientButton — animated gradient glow button.
  Copyright (c) 2026 Wendell47 (Wendel Araujo dos Santos)
  Licensed under the MIT License.
-->
<template>
  <button
    class="gradient-btn inline-flex items-center justify-center px-8 py-4
           border-0 relative overflow-hidden rounded-[10rem] font-bold cursor-pointer
           text-neutral-800 dark:text-neutral-200
           shadow-[0_0_7px_-5px_rgba(0,0,0,0.5)]
           hover:text-[rgb(33,0,85)] dark:hover:text-white
           active:scale-[0.97]
           transition-transform duration-[0.02s]"
  >
    <slot />
    <span class="hover-effect absolute inset-0 flex items-center justify-center z-1">
      <span
        class="block rounded-[40rem] w-40 h-40 opacity-50 blur-[20px]
               animate-[spin_3s_linear_infinite]
               bg-[linear-gradient(90deg,rgba(222,0,75,1)_0%,rgba(191,70,255,1)_49%,rgba(0,212,255,1)_100%)]
               group-hover:w-32 group-hover:h-32
               transition-[width,height] duration-400"
      />
    </span>
  </button>
</template>

<script setup lang="ts">
// GradientButton.vue — slot-based, no props needed beyond <slot>
</script>

<style scoped>
.gradient-btn:hover {
  background: rgb(193, 228, 248);
}

@media (prefers-color-scheme: dark) {
  .gradient-btn:hover {
    background: rgb(40, 40, 60);
  }
}
</style>
