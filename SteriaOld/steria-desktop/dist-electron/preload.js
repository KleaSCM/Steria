/**
 * Electron Preload Script.
 *
 * Exposes protected APIs to the renderer process.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
// See the Electron documentation for details on how to use preload scripts:
// https://www.electronjs.org/docs/latest/tutorial/process-model#preload-scripts
window.addEventListener('DOMContentLoaded', () => {
    const ReplaceText = (selector, text) => {
        const Element = document.getElementById(selector);
        if (Element) {
            Element.innerText = text;
        }
    };
    for (const type of ['chrome', 'node', 'electron']) {
        ReplaceText(`${type}-version`, process.versions[type] ?? 'unknown');
    }
});
import { contextBridge, ipcRenderer } from 'electron';
contextBridge.exposeInMainWorld('steria', {
    GetIdentity: () => ipcRenderer.invoke('get-identity'),
    CreateAccount: (data) => ipcRenderer.invoke('create-account', data),
    Login: (data) => ipcRenderer.invoke('login', data),
    GetProjects: () => ipcRenderer.invoke('get-projects'),
    Run: (args) => ipcRenderer.invoke('run-steria', args),
});
