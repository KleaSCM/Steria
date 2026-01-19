/**
 * Electron Main Process.
 *
 * Handles the lifecycle of the application windows and system integration.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */

import { app, BrowserWindow } from 'electron';
import path from 'path';
import { fileURLToPath } from 'url';
import { createRequire } from 'module';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const require = createRequire(import.meta.url);

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) {
    app.quit();
}

const CreateWindow = () => {
    // Create the browser window.
    const MainWindow = new BrowserWindow({
        width: 1200,
        height: 800,
        webPreferences: {
            preload: path.join(__dirname, 'preload.js'),
            nodeIntegration: true,
            contextIsolation: false, // For now, to simplify local file access if needed
        },
        titleBarStyle: 'hidden', // Modern look
        titleBarOverlay: {
            color: '#00000000',
            symbolColor: '#ffffff',
            height: 30
        },
        vibrancy: 'fullscreen-ui', // MacOS glass effect
        backgroundMaterial: 'acrylic', // Windows glass effect
        backgroundColor: '#00000000', // Transparent for glass effects
    });

    // Load the index.html of the app.
    if (process.env.VITE_DEV_SERVER_URL) {
        MainWindow.loadURL(process.env.VITE_DEV_SERVER_URL);
    } else if (!app.isPackaged) {
        // Fallback for dev mode without plugin
        MainWindow.loadURL('http://localhost:5173');
    } else {
        MainWindow.loadFile(path.join(__dirname, '../dist/index.html'));
    }

    // Open the DevTools.
    // MainWindow.webContents.openDevTools();
};

// IPC Handlers
import { ipcMain } from 'electron';
import fs from 'fs';
import { execFile } from 'child_process';
import crypto from 'crypto';

const IdentityPath = path.join(app.getPath('userData'), 'identity.json');

const HashPassword = (password: string, salt: string): Promise<string> => {
    return new Promise((resolve, reject) => {
        crypto.scrypt(password, salt, 64, (err, derivedKey) => {
            if (err) reject(err);
            resolve(derivedKey.toString('hex'));
        });
    });
};

ipcMain.handle('get-identity', async () => {
    // Return null if we want to force login every time app starts
    // Or return partial identity if we implemented "remember me"
    // For now, let's require login on launch for security
    return null;
});

ipcMain.handle('create-account', async (_, data) => {
    try {
        const Salt = crypto.randomBytes(16).toString('hex');
        const HashedPassword = await HashPassword(data.Password, Salt);

        const NewIdentity = {
            Name: data.Name,
            Email: data.Email,
            HashedPassword,
            Salt
        };

        fs.writeFileSync(IdentityPath, JSON.stringify(NewIdentity));
        return true;
    } catch (error) {
        return false;
    }
});

ipcMain.handle('login', async (_, data) => {
    try {
        if (!fs.existsSync(IdentityPath)) return null;

        const StoredIdentity = JSON.parse(fs.readFileSync(IdentityPath, 'utf-8'));

        // Check identifier (Name or Email)
        const IsMatch = StoredIdentity.Name === data.Identifier ||
            (StoredIdentity.Email && StoredIdentity.Email === data.Identifier);

        if (!IsMatch) return null;

        const InputHash = await HashPassword(data.Password, StoredIdentity.Salt);

        if (InputHash === StoredIdentity.HashedPassword) {
            return { Name: StoredIdentity.Name, Email: StoredIdentity.Email };
        }

        return null;
    } catch (error) {
        return null;
    }
});

ipcMain.handle('run-steria', async (_, args: string[]) => {
    return new Promise((resolve, reject) => {
        // Assume 'steria' binary is in the same directory as the app or in PATH
        // For development, we point to the compiled go binary
        const SteriaPath = app.isPackaged
            ? path.join(process.resourcesPath, 'steria')
            : path.join(__dirname, '../../steria'); // Adjust depending on where 'steria' binary is

        execFile(SteriaPath, args, (error, stdout, stderr) => {
            if (error) {
                resolve({ success: false, output: stderr || error.message });
            } else {
                resolve({ success: true, output: stdout });
            }
        });
    });
});

ipcMain.handle('get-projects', async () => {
    // Re-using run-steria logic logic implicitly or explicitly
    const SteriaPath = app.isPackaged
        ? path.join(process.resourcesPath, 'steria')
        : path.join(__dirname, '../../steria');

    return new Promise((resolve) => {
        execFile(SteriaPath, ['projects', 'list', '--json'], (error, stdout) => {
            if (error) {
                resolve({});
            } else {
                try {
                    resolve(JSON.parse(stdout));
                } catch {
                    resolve({});
                }
            }
        });
    });
});

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
app.on('ready', CreateWindow);

// Quit when all windows are closed, except on macOS.
app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});

app.on('activate', () => {
    // On OS X it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) {
        CreateWindow();
    }
});
