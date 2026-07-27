/**
 * Steria (STERIA) — 重ね合わせバージョン管理システムよ。
 *
 * Steria is a version control system built on the physics concept of
 * **quantum superposition**. All tracked versions of a file coexist
 * simultaneously as content-addressed hashes. She does not hide,
 * does not collapse — she just watches, and remembers everything.
 *
 * DESIGN PHILOSOPHY:
 * Steria is out of your way. You do not think about Steria.
 * She watches you. She is too shy to ask questions.
 * The only time she asks is when you run `steria choose`.
 *
 * Every `steria done` is signed with your identity and carries a
 * message — like pressing Ctrl+S in a text document, but for
 * your whole project. Save, sync, done. You work, she watches.
 *
 * Steria is NOT Git. Steria does not use commits-as-deltas,
 * branches, checkouts, or any Git-native concept. Steria is a
 * complete reimagining of version control. Comparing her to Git
 * is actively unhelpful.
 *
 * STORAGE:
 * .steria/
 *   ├── objects/          — Content-addressable hash store
 *   │   └── aa/
 *   │       └── bbcc...   — Object file named by hash prefix
 *   ├── index             — Superposition state manifest
 *   ├── config            — Remote URL + local identity
 *   └── head              — Pointer to latest done state
 *
 * COMMAND OVERVIEW:
 * ┌──────────┬──────────────────────────────────────────────┐
 * │ Command  │ What it does                                  │
 * ├──────────┼──────────────────────────────────────────────┤
 * │ init     │ Create a repo on the remote server            │
 * │ watch    │ Create .steria, start watching                │
 * │ done     │ Save state + sign + auto-sync to remote       │
 * │ choose   │ Collapse superposition — pick versions        │
 * │ clone    │ Fetch entire remote superposition locally     │
 * └──────────┴──────────────────────────────────────────────┘
 *
 * WORKFLOW:
 *   # Local-only (single user)
 *   steria watch → (work) → steria done → (work) → steria done → ...
 *
 *   # Starting fresh with remote
 *   steria init MyProject
 *   steria watch
 *   (work) → steria done → (work) → steria done → ...
 *
 *   # Joining a team project
 *   steria clone MyProject
 *   steria choose
 *   (work) → steria done → (work) → steria done → ...
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */

# Steria VCS

**Superposition Version Control** — she watches *everything*.

Steriaは量子力学の**重ね合わせ**をベースにしたバージョン管理システムよ。
全バージョンがハッシュとして同時に存在してるの。
Gitみたいに「過去に戻る」んじゃなくて、「全部が今ここにある」って感じ 

Steriaはあなたの作業を見守ってるだけで、邪魔はしないの。

---

## 核となる考え方

Steria models file versioning as **quantum superposition**.

- **All versions coexist** — Every `steria done` hashes the current state
  and stores it in the content-addressable object store. Previous versions
  are not replaced — they remain in superposition forever.
- **Observation collapses the state** — `steria choose` is the act of
  observation. You see every version side-by-side. You pick. You resolve.
  The superposition collapses into a single chosen state.
- **No checkout** — You never "check out" a revision. The past is already
  here, sitting in `.steria/objects/` and visible in your directory.
- **No branches** — Branching implies forking a timeline. Steria does not
  have timelines. She has superposition. All versions, one space.

Steriaはファイルのバージョン管理を**量子の重ね合わせ**としてモデル化してるの。
`steria done`するたびに状態がハッシュとして保存されて、前のは上書きされずに
重ね合わせの中に残り続けるの。
`steria choose`が「観測」の役割 — 全部のバージョンから選ぶってわけ 

---

## ユーザー

Each user who interacts with Steria has a **name** — set once, stored
forever. You never type your name again. Every `steria done` is
automatically signed with it.

- Set on first use — Steria asks for your name
- Stored globally in `~/.steria/config`
- Overridable per-project in `.steria/config`
- Every done records your name automatically — you only write the message

```
$ steria done "fixed the collision detection bug"
  → KleaSCM: "fixed the collision detection bug"

$ steria done "added the new character sprite"
  → 百合子: "added the new character sprite"
```

Identity is how `steria choose` shows you who did what:

```
path/to/file.go
  a1b2c3d4...  KleaSCM: "initial implementation"
  d4e5f6a7...  百合子: "refactored the parser"
  9f8e7d6c...  Kaori: "added edge case handling"
```

初めてSteriaを使う時に名前を設定すると、その後は自動で署名されるの。
`steria done`ではメッセージだけ書けばいいの 
`steria choose`で「誰が何をしたか」が一目で分かるようになってるの。

---

## Commands — コマンド

### `steria watch [Directory]`

Creates `.steria` in the specified directory and begins watching.
**This is how you start tracking a project. That's it.**

- Creates `.steria/` directory structure
- Starts watching for file modifications, additions, deletions
- No flags. No config. Just watch.
- Future: `.SteriaIgnore` for files to exclude

```
$ cd CuteProject
$ steria watch
 Steria is watching CuteProject...
```

Steriaがプロジェクトを見張り始めるの 
`.steria`フォルダを作って、そこからは全部監視よ。

### `steria done "message"`

Saves the current state into the superposition, signs it automatically
with your identity, and syncs to remote if configured.

**This is the 90% command.** You use this constantly.
It is Ctrl+S for version control — save, sign, sync, done.

- Hashes every tracked file (SHA-256)
- Stores objects in `.steria/objects/`
- Signs with your pre-configured identity — you only write the message
- Updates `.steria/index` with new superposition entries + identity
- Updates `.steria/head`
- If remote configured: sends objects + index delta to server
- Previous states remain untouched — superposition grows

```
$ steria done "fixed the collision detection bug"
KleaSCM: "fixed the collision detection bug"
 done
```

### `steria choose [File...]`

The **collapse** operation. Steria finally speaks.
She opens a diff viewer TUI showing every version of every file —
labelled by who made it, what they said, and when — and asks you
to pick the winner.

- Terminal diff viewer TUI with navigation
- Each version labelled: `Identity: "message"` + hash + timestamp
- Navigate versions, view diffs against current, pick the winner
- Conflicts resolved explicitly — Steria never auto-merges
- Output: single resolved state per file in the working directory
- Main use: after `steria clone` to collapse the downloaded superposition

**差分を見ながら、どのバージョンを採用するか選べるTUIになってるの。**

```
$ steria choose
┌─── Steria Choose ───────────────────────────────────────┐
│ path/to/file.go                                          │
│                                                          │
│ [1] a1b2c3d4  KleaSCM: "initial implementation"         │
│     │  +func Parse(input string) Result {                │
│     │  +    return Result{Value: input}                  │
│     │  +}                                                │
│                                                          │
│ [2] d4e5f6a7  百合子: "refactored the parser"            │
│     │  -func Parse(input string) Result {                │
│     │  -    return Result{Value: input}                  │
│     │  +func Parse(input string) (Result, error) {       │
│     │  +    if input == "" {                             │
│     │  +        return Result{}, errors.New("empty")     │
│     │  +    }                                            │
│     │  +    return Result{Value: input}, nil             │
│     │  +}                                                │
│                                                          │
│ [3] 9f8e7d6c  Kaori: "added edge case handling"          │
│     │  ...                                               │
│                                                          │
│ Pick a version [1-3]: █                                  │
└──────────────────────────────────────────────────────────┘
```

### `steria clone [RepoName]`

Fetches the entire superposition from a remote server and materialises it
locally. Every version of every file is present side-by-side.

- Downloads the complete object store from the remote
- Reconstructs the full index locally
- Configures remote in `.steria/config` for future `done` syncs
- All versions visible immediately — nothing hidden
- After: `steria choose` to collapse to a working state

```
$ steria clone MyProject
 Steria is cloning MyProject...
 done — all versions are here
$ steria choose
```

リモートサーバーから重ね合わせ全体を取ってくるの。
全部のバージョンがそのままローカルに現れるから、すぐに全部見えるわ 
その後は`steria choose`で収縮させるの。

### `steria init [RepoName]`

Creates a new repository on the remote server **and** sets up local
`.steria` with the remote URL configured. One command, everything ready.

- Registers the repo name on the server
- Creates `.steria/` locally with remote URL in config
- After init: just `steria watch` to start tracking
- Future `steria done` calls auto-sync to this remote

```
$ steria init MyProject
Remote: steria://server/MyProject
.steria created — ready to watch
```

---

## アーキテクチャ

### Storage Model

Content-addressable object store with identity-signed entries.

```
Object = SHA-256(FileContent) → FileContent
```

- Objects stored in `.steria/objects/XX/YYYYYY...`
- Index maps file paths to superposition of hash pointers
- Each index entry includes `Identity: "message"` metadata
- Head points to the latest done state

### Server Model

Central server model. `steria done` auto-syncs.

- Server stores the canonical superposition for each repo
- Single-user: local binary is its own server
- Team: Steria daemon on a LAN/internet-accessible machine
- Protocol: minimal HTTP-based API

```
┌─────────┐  clone/done  ┌──────────┐
│  Team   │ ◄──────────► │  Steria  │
│  Member │               │  Server  │
│  A      │               │          │
└─────────┘               │  .steria │
                          │  objects │
┌─────────┐               └──────────┘
│  Team   │
│  Member │
│  B      │
└─────────┘
```

### Index Format

Each index entry maps a file path to its superposition of signed states:

```
path/to/file.go
  a1b2c3d4...  KleaSCM: "initial implementation"
  d4e5f6a7...  Yuka: "refactored the parser"
  9f8e7d6c...  Kaori: "added edge case handling"

path/to/image.png
  ffeeddcc...  百合子: "added new character sprite"
```

Each done appends new hashes to the superposition.
The identity is part of the record — you always know who did what.

---

## ライフサイクル

### Starting a new project (local only)

```
steria watch
  → .steria created, watching...

(work work work)

steria done "initial project setup"
  → KleaSCM: "initial project setup"
  → saved locally 
```

### Starting a new project (with remote team)

```
steria init MyProject
  → remote repo created

steria watch
  → .steria created, watching...

(work work work)

steria done "initial project setup"
  → KleaSCM: "initial project setup"
  → saved locally + synced to remote 
```

### Joining a team project

```
steria clone MyProject
  → all versions downloaded

steria choose
  → collapse to working state

(work work work)

steria done "fixed the ui layout"
  → KleaSCM: "fixed the ui layout"
  → saved locally + synced to remote 
```

---

## Future Considerations — これから

- `.SteriaIgnore` — glob-based ignore file
- `steria log` — view done history
- `steria status` — show changes since last done
- `steria diff` — compare versions within superposition
- Authentication + authorisation for server
- Partial clone — subset of superposition
- Server daemon mode (`steria serve`)
- `steria config` — manage identity and remote settings

---
