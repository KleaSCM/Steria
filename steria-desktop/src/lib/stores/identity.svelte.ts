/**
 * Identity Store.
 * 
 * Manages the user's identity state using Svelte 5 Runes.
 * 
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */

import type { Identity } from '../../global';

let CurrentIdentity = $state<Identity | null>(null);

export const IdentityStore = {
    get Current() {
        return CurrentIdentity;
    },
    async Initialize() {
        // We can check if session exists, but we are enforcing login
        // So initialize does nothing, user must login
    },
    async CreateAccount(name: string, password: string, email?: string) {
        if (window.steria) {
            const Success = await window.steria.CreateAccount({ Name: name, Password: password, Email: email });
            if (Success) {
                // Auto login after create? Or make them login? 
                // Let's auto login
                await IdentityStore.Login(name, password);
            }
        }
    },
    async Login(identifier: string, password: string) {
        if (window.steria) {
            const Identity = await window.steria.Login({ Identifier: identifier, Password: password });
            if (Identity) {
                CurrentIdentity = Identity;
                return true;
            }
        }
        return false;
    }
};
