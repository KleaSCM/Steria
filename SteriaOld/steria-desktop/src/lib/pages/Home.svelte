<script lang="ts">
    /**
     * Home View.
     *
     * The landing page of the application.
     *
     * Author: KleaSCM
     * Email: KleaSCM@gmail.com
     */

    import { IdentityStore } from "../stores/identity.svelte";

    let Mode = $state<"login" | "create">("create");
    let Name = $state("");
    let Email = $state("");
    let Password = $state("");
    let Loading = $state(false);

    const HandleSubmit = async () => {
        if (!Name && !Email) return;
        if (!Password) return;

        Loading = true;

        if (Mode === "create") {
            await IdentityStore.CreateAccount(Name, Password, Email);
        } else {
            // Name field serves as identifier
            await IdentityStore.Login(Name, Password);
        }

        Loading = false;
    };
</script>

<div class="view-container">
    <h1>Steria</h1>
    <p class="subtitle">Distributed Version Control System</p>

    <div class="glass-panel content">
        {#if IdentityStore.Current}
            <p>
                Welcome back, <span class="highlight"
                    >{IdentityStore.Current.Name}</span
                >! 🌸
            </p>
            <div class="actions">
                <button>New Project</button>
                <button>Sync All</button>
            </div>
        {:else}
            <h2>{Mode === "create" ? "Create Account" : "Welcome Back"}</h2>
            <p class="intro">
                {Mode === "create"
                    ? "Let's get you set up to share files."
                    : "Login to access your repository."}
            </p>

            <div class="form-group">
                <input
                    type="text"
                    placeholder={Mode === "create"
                        ? "Your Name"
                        : "Name or Email"}
                    bind:value={Name}
                />
                {#if Mode === "create"}
                    <input
                        type="email"
                        placeholder="Email (Optional)"
                        bind:value={Email}
                    />
                {/if}
                <input
                    type="password"
                    placeholder="Password"
                    bind:value={Password}
                />

                <button
                    onclick={HandleSubmit}
                    disabled={!Name || !Password || Loading}
                >
                    {Loading
                        ? "Processing..."
                        : Mode === "create"
                          ? "Get Started"
                          : "Login"}
                </button>
            </div>

            <p class="switch-mode">
                {Mode === "create"
                    ? "Already have an account?"
                    : "Don't have an account?"}
                <button
                    class="link"
                    onclick={() =>
                        (Mode = Mode === "create" ? "login" : "create")}
                >
                    {Mode === "create" ? "Login" : "Create Account"}
                </button>
            </p>
        {/if}
    </div>
</div>

<style>
    .view-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100%;
        gap: 1rem;
    }

    .subtitle {
        color: var(--text-dim);
        margin-top: -0.5rem;
    }

    .content {
        padding: 3rem;
        margin-top: 2rem;
        min-width: 400px;
        text-align: center;
    }

    .intro {
        color: var(--text-dim);
        margin-bottom: 2rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    input {
        background: rgba(0, 0, 0, 0.3);
        border: 1px solid rgba(255, 255, 255, 0.1);
        padding: 0.8em;
        border-radius: 8px;
        color: var(--text-main);
        outline: none;
        transition: border-color 0.2s;
    }

    input:focus {
        border-color: var(--primary);
    }

    .highlight {
        color: var(--primary);
        font-weight: bold;
    }

    .actions {
        display: flex;
        gap: 1rem;
        justify-content: center;
        margin-top: 1.5rem;
    }

    .switch-mode {
        margin-top: 1.5rem;
        font-size: 0.9em;
        color: var(--text-dim);
    }

    .link {
        background: none;
        border: none;
        color: var(--secondary);
        padding: 0;
        text-decoration: underline;
        cursor: pointer;
        display: inline;
        font-family: inherit;
        font-size: inherit;
    }
    .link:hover {
        color: var(--primary);
        box-shadow: none;
    }
</style>
