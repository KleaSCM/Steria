/**
 * Simple Router using Svelte 5 Runes.
 *
 * Manages application state and view navigation.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */

export const Routes = {
    Home: 'Home',
    Projects: 'Projects',
    Settings: 'Settings',
} as const;

export type Route = keyof typeof Routes;

let CurrentRoute = $state<Route>(Routes.Home);

export const Router = {
    get Current() {
        return CurrentRoute;
    },
    Navigate(route: Route) {
        CurrentRoute = route;
    }
};
