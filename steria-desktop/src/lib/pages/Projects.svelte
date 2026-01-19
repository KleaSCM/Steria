<script lang="ts">
    /**
     * Projects View.
     *
     * Lists local projects managed by Steria.
     *
     * Author: KleaSCM
     * Email: KleaSCM@gmail.com
     */

    import { onMount } from "svelte";

    let Projects = $state<Record<string, string>>({});
    let Loading = $state(true);

    onMount(async () => {
        if (window.steria) {
            Projects = await window.steria.GetProjects();
        }
        Loading = false;
    });

    const ProjectKeys = $derived(Object.keys(Projects).sort());
</script>

<div class="view-container">
    <h2>Projects</h2>
    <div class="glass-panel list">
        {#if Loading}
            <div class="empty-state">Loading...</div>
        {:else if ProjectKeys.length === 0}
            <div class="empty-state">
                No projects found. <br />
                <small
                    >Use `steria projects add` via CLI or stay tuned for UI
                    support!</small
                >
            </div>
        {:else}
            <div class="project-grid">
                {#each ProjectKeys as name}
                    <div class="project-card">
                        <h3>{name}</h3>
                        <p class="path">{Projects[name]}</p>
                        <div class="actions">
                            <button class="small">Open</button>
                            <button class="small">Settings</button>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<style>
    .view-container {
        padding: 4rem 2rem;
        width: 100%;
        max-width: 800px;
        margin: 0 auto;
    }

    .list {
        margin-top: 1rem;
        min-height: 200px;
        padding: 1.5rem;
    }

    .project-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
        gap: 1rem;
    }

    .project-card {
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        padding: 1rem;
        border-radius: 8px;
        transition: all 0.2s ease;
    }

    .project-card:hover {
        background: rgba(255, 255, 255, 0.1);
        border-color: var(--primary);
    }

    .project-card h3 {
        margin: 0 0 0.5rem 0;
        color: var(--text-main);
    }

    .path {
        font-size: 0.8rem;
        color: var(--text-dim);
        margin-bottom: 1rem;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .empty-state {
        color: var(--text-dim);
        font-style: italic;
        text-align: center;
        padding: 2rem;
    }

    .actions {
        display: flex;
        gap: 0.5rem;
    }

    button.small {
        font-size: 0.8rem;
        padding: 0.3em 0.8em;
    }
</style>
