<script lang="ts">
  /**
   * Main Application Component.
   *
   * Serves as the root component for the application, organizing the main layout
   * and navigation structure.
   *
   * Author: KleaSCM
   * Email: KleaSCM@gmail.com
   */

  import MeteorShower from "./lib/components/MeteorShower.svelte";
  import Navbar from "./lib/components/Navbar.svelte";
  import { Router, Routes } from "./lib/router.svelte";
  import { IdentityStore } from "./lib/stores/identity.svelte";
  import { onMount } from "svelte";

  // Views
  import Home from "./lib/pages/Home.svelte";
  import Projects from "./lib/pages/Projects.svelte";
  import Settings from "./lib/pages/Settings.svelte";
  import "./app.css";

  onMount(async () => {
    await IdentityStore.Initialize();
  });
</script>

<main>
  <MeteorShower />
  <Navbar />

  <div class="content-wrapper">
    {#if Router.Current === Routes.Home}
      <Home />
    {:else if Router.Current === Routes.Projects}
      <Projects />
    {:else if Router.Current === Routes.Settings}
      <Settings />
    {/if}
  </div>
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    overflow: hidden;
  }

  .content-wrapper {
    flex: 1;
    overflow-y: auto;
    position: relative;
    z-index: 1;
  }
</style>
