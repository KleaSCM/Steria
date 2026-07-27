<script lang="ts">
    /**
     * Navigation Bar Component.
     *
     * Provides top-level navigation for the application.
     *
     * Author: KleaSCM
     * Email: KleaSCM@gmail.com
     */

    import { Router, Routes, type Route } from "../router.svelte";

    const NavItems: Route[] = [Routes.Home, Routes.Projects, Routes.Settings];
</script>

<nav class="glass-panel">
    {#each NavItems as item}
        <button
            class:active={Router.Current === item}
            onclick={() => Router.Navigate(item)}
        >
            {item}
        </button>
    {/each}
</nav>

<style>
    nav {
        position: fixed;
        top: 1rem;
        left: 50%;
        transform: translateX(-50%);
        padding: 0.5rem;
        display: flex;
        gap: 0.5rem;
        z-index: 100;
    }

    button {
        background: transparent;
        border: 1px solid transparent;
        color: var(--text-dim);
        padding: 0.5rem 1.5rem;
        cursor: pointer;
        transition: all 0.2s ease;
        border-radius: 8px;
    }

    button:hover {
        color: var(--text-main);
        background: rgba(255, 255, 255, 0.05);
    }

    button.active {
        background: rgba(255, 126, 182, 0.2);
        color: var(--primary);
        border-color: rgba(255, 126, 182, 0.3);
        box-shadow: 0 0 10px rgba(255, 126, 182, 0.2);
    }
</style>
