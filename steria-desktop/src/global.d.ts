/**
 * Global Type Definitions.
 * 
 * Declares the global `steria` API exposed by Electron.
 * 
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */

export interface SteriaAPI {
    GetIdentity: () => Promise<Identity | null>; // Returns basic info if session exists
    CreateAccount: (data: CreateAccountData) => Promise<boolean>;
    Login: (data: LoginData) => Promise<Identity | null>;
    GetProjects: () => Promise<Record<string, string>>;
    Run: (args: string[]) => Promise<{ success: boolean; output: string }>;
}

export interface CreateAccountData {
    Name: string;
    Password: string;
    Email?: string;
}

export interface LoginData {
    Identifier: string; // Name or Email
    Password: string;
}

export interface Identity {
    Name: string;
    Email?: string;
}

declare global {
    interface Window {
        steria: SteriaAPI;
    }
}
