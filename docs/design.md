# li - A simple text editor

li (pronounced 'lee') is a simple text editor based on [hibiken's mini](https://github.com/hibiken/mini) which was inspired from [antirez's kilo](http://antirez.com/news/108) editor.

The philosophy behind li is: simplicity, and hackability. Though not hackable in the same way as emacs or vim is hackable. This is more for developing and small, efficient editor specifically for you. There is a list of non-goals below containing many features that are part of Vim which are not part of the goals of this editor. Think of it as a hackable nano.

This both a small customisable editor, as well as a sandbox for your ideas. One interesting thing I've put is the ability to silently recomiple the program in the background. Since there is no runtime readable files, this is the next best thing for it.

There are no runtime readable files, it is designed to be easy to modify basic variables in the `config.go` file and recompiled.

You are also strongly encouraged to write your own keymaps here as well directly in the source code. By default it used Vim bindings.

The true features of the editor are perhaps too simply for the end user, the real features are the core editor runtime.

- Keyword and comment syntax highlighting
- File open/close
- Render text
- Status line

However it is shipped with a default Vim like configuration which contains:

- Normal/Insert mode
- Undo/Redo
- Basic Vim commands e.g `d`,`c`, `hjkl`

There are some features that I do not intend on ever implementing as they are against the general philophy.

- Multiple windows
- Tabs

Though of course this is your editor and you are free to do what you like ;)

Things I am open to however but are not on an immediate TODO are

- LSP support
- Advanced Syntax highlighting

# Architecture

li has a very small kernel. It handles IO, signals, rendering, errors. This part is not intended to be modified by users as it is really the most basic core. Located in the `core` package.

On top of this the kernel's API is not exactly safe. Think of it like C, you are more than welcome to try to index out of bounds. If you want to have a safer API, use the SDK's API.

On top of the core is a basic SDK for the the editor. It provides common functions such as line manipulation, chunking changes into undo and redos.

Then there is the keymaps in the `config.go` file. Here you will find *all* the keymaps for the editor as well as certain variables.

# Core

The core is a minimal kernel for the editor.

When the user presses a key, it is passed to the `core.Input` channel.
When the program receives a signal, it is passed to the `core.Signal` channel.
Aware of filetypes and for changing their configuration, but not the actual configuration

Handles rendering including basic syntax highlighting

Create a new editor with `NewEditor(in io.Reader, out io.Writer)`
which returns a struct with the following fields

// Channel to read keys
KeyChan() chan-> Key

// Channel to read signals
SigChan() chan-> Signal

// The line rows
Rows() []Row

# Callbacks

There are several callbacks

- File open

# Design methodology

API-centric design. When brainstorming the design of the API, imagine yourself as the hacker, customising the software. What is the simplest and most intuitive API that you would like provided?
Once you have done that, put your maintainer hat back on and ask is this feasible. Experiment and modify later.

I originally wanted to give the users full access to the rows so that they can modify them themselves, but I think this is better abstracted away with a getter since maybe later we will change the internal storage of the rows into a rope data structure which is more common for editors.
