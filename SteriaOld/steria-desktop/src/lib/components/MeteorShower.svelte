<script lang="ts">
    /**
     * Meteor Shower Background Component.
     *
     * Renders an animated meteor shower effect for background ambiance.
     *
     * Author: KleaSCM
     * Email: KleaSCM@gmail.com
     */

    import { onMount } from "svelte";

    interface Meteor {
        Id: number;
        X: number;
        Y: number;
        Length: number;
        Speed: number;
        Delay: number;
    }

    let Meteors: Meteor[] = [];
    const MeteorCount = 20;

    onMount(() => {
        // Initialize meteors
        Meteors = Array.from({ length: MeteorCount }, (_, i) => ({
            Id: i,
            X: Math.random() * 100, // vw
            Y: Math.random() * 100, // vh
            Length: Math.random() * 80 + 20,
            Speed: Math.random() * 1 + 0.5,
            Delay: Math.random() * 5,
        }));
    });
</script>

<div class="meteor-shower">
    {#each Meteors as meteor (meteor.Id)}
        <div
            class="meteor"
            style="
				left: {meteor.X}vw;
				top: {meteor.Y}vh;
				width: {meteor.Length}px;
				animation-duration: {meteor.Speed}s;
				animation-delay: {meteor.Delay}s;
			"
        ></div>
    {/each}
</div>

<style>
    .meteor-shower {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        overflow: hidden;
        z-index: -1;
        pointer-events: none;
        background: radial-gradient(
            ellipse at bottom,
            #1b2735 0%,
            #090a0f 100%
        );
    }

    .meteor {
        position: absolute;
        height: 2px;
        background: linear-gradient(
            90deg,
            rgba(255, 255, 255, 0),
            rgba(255, 255, 255, 0.8)
        );
        border-radius: 999px;
        transform: rotate(-45deg);
        opacity: 0;
        animation-name: shower;
        animation-timing-function: linear;
        animation-iteration-count: infinite;
        box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
    }

    @keyframes shower {
        0% {
            transform: rotate(-45deg) translateX(0);
            opacity: 0;
        }
        10% {
            opacity: 1;
        }
        70% {
            opacity: 1;
        }
        100% {
            transform: rotate(-45deg) translateX(-100vh); /* Move diagonally */
            opacity: 0;
        }
    }
</style>
